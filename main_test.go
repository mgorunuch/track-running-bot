package main

import (
	"log"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Test_images(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI("1319032110:AAHTa28_WGA6Ep_P-CIEvDPVGs22FVm-qtA")
	if err != nil {
		log.Panic(err)
	}

	buffer, err := drawChart(200, 100, []float64{5, 3, 1, 2, 4, 8, 1, 2, 0, 4, 3, 5})

	bts := tgbotapi.FileBytes{
		"name2.jpeg",
		buffer.Bytes(),
	}

	res, err := bot.UploadFile("sendPhoto", map[string]string{
		"chat_id": "73420519",
	}, "photo", bts)

	log.Printf("%+v", err)
	log.Printf("%+v", res)
}
