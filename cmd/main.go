package main

import (
	"hellobot/internal/ports/vk"
	"log"
)

func main() {
	bot, err := vk.NewVkBot(
		"vk1.a.VBu8mKLhhnw9diHcr_WIISUCVjyI9GLsHxb-QywyMfF0dgKa9BU7XQ57xyUREHhDxbMmbU6r8qZsCuvv42JTixkNyFJXdu2WoVYx4Kv-6UGxgV6tPJJZKkXwxnglNAH9Ks1FumCX7rhkWlXhO_d5BTS1yTZTVuUCx0iWTx7dBDZTENLk16CttNzEjWYNi-uOJGjx3Bgllt2lVnSOqEekJw",
		"220417305",
	)
	if err != nil {
		log.Fatal(err)
	}
	bot.PollAndServe()

}
