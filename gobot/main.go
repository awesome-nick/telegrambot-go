package main

import (
	"gobot/telegrambot"
	"log"
	"sync"
)

var wg sync.WaitGroup

func main() {

	wg.Add(1)
	log.Println("Starting telegram bot in a separate thread...")
	go telegrambot.TelegramBot(&wg)

	wg.Wait()
}
