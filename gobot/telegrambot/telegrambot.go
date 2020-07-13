package telegrambot

import (
	"gobot/database"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TelegramBot(wg *sync.WaitGroup) {

	// if err := godotenv.Load(os.ExpandEnv("$GOPATH/src/gobot/.env")); err != nil {
	// 	log.Fatalln("Error loading .env file")
	// }

	telegramAPItoken := os.Getenv("TELEGRAM_API_TOKEN")
	if telegramAPItoken == "" {
		log.Fatalln("Error loading telegram_api_token from .env-bot file!")
	}

	imgPath := os.Getenv("IMG_PATH")
	if imgPath == "" {
		imgPath = "/var/gobot/img/"
	}

	defer wg.Done()

	log.Println("TELEGRAM_API_TOKEN = " + telegramAPItoken)

	bot, err := tgbotapi.NewBotAPI(telegramAPItoken)
	if err != nil {
		log.Println("tgbotapi.NewBotAPI Error:")
		log.Fatalln(err)
	}

	isProduction := os.Getenv("PROD_ENV")
	if isProduction == "" {
		isProduction = "True"
	}

	if isProduction != "True" {
		bot.Debug = true
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println("bot.GetUpdatesChan Error:")
		log.Fatalln(err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if isProduction != "True" {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		}

		go recordRequest(update.Message)

		if update.Message.IsCommand() {
			go cmdHandler(imgPath, update.Message, bot)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}

		// if update.CallbackQuery != nil {
		// 	fmt.Print(update)

		// 	bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

		// 	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		// }
	}
}

func recordRequest(msg *tgbotapi.Message) {

	var req database.Request
	req.MessageId = uint64(msg.MessageID)
	req.UserId = uint64(msg.From.ID)
	req.FirstName = msg.From.FirstName
	req.Language = msg.From.LanguageCode
	req.ChatId = uint64(msg.Chat.ID)
	req.MessageText = msg.Text
	req.Date = uint64(msg.Date)

	req.Create()
}

func cmdHandler(imgPath string, u *tgbotapi.Message, b *tgbotapi.BotAPI) {
	c := u.Command()

	if c == "help" {
		go helpHandler(u, b)
		return
	}

	if i, ok := database.GetImgByCommand(c); ok {
		filename := imgPath + i.Filepath

		content, err := ioutil.ReadFile(filename)
		if err != nil {
			msg := tgbotapi.NewMessage(u.Chat.ID, "I don't know this command.")
			_, err := b.Send(msg)
			if err != nil {
				log.Println(err)
			}
		} else {
			bytes := tgbotapi.FileBytes{Name: i.Filepath, Bytes: content}
			msg := tgbotapi.NewPhotoUpload(u.Chat.ID, bytes)
			msg.Caption = i.ImageCaption
			_, err := b.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		message := database.GetBotMessageByCommand(c)
		msg := tgbotapi.NewMessage(u.Chat.ID, "")

		if message != nil {
			msg.Text = message.MessageText
		} else {
			msg.Text = "I don't know this command."
		}
		_, err := b.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func helpHandler(u *tgbotapi.Message, b *tgbotapi.BotAPI) {
	cs := database.GetAllImageCommands()
	ms := database.GetAllMessageCommands()
	s := "/" + strings.Join(*cs, " /") + " /" + strings.Join(*ms, " /")
	msg := tgbotapi.NewMessage(u.Chat.ID, "List of commands: "+s)
	_, err := b.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
