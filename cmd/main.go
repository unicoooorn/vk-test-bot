package main

import (
	"hellobot/internal/ports/vk"
	"log"
)

func main() {
	bot, err := vk.NewVkBot(
		"your-token",
		12345678,
	)
	if err != nil {
		log.Fatal(err)
	}
	bot.PollAndServe()
}
