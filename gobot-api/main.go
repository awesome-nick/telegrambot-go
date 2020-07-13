package main

import (
	"gobot/api"
	"log"
	"sync"
)

var wg sync.WaitGroup

func main() {

	wg.Add(1)
	log.Println("Starting http server in a separate thread...")
	go api.StartServer(&wg)

	wg.Wait()
}
