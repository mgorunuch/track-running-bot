package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Test_images(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI("1319032110:AAHTa28_WGA6Ep_P-CIEvDPVGs22FVm-qtA")
	if err != nil {
		log.Panic(err)
	}

	buffer, err := drawChart(200, 100, [][]float64{{5, 3, 1, 2, 4, 8, 1, 2, 0, 4, 3, 5}})

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

func Test_123(t *testing.T) {
	res, err := drawSuccessPredChard(1000, 365, 3.65)

	fl, err := os.OpenFile("name.png", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	defer fl.Close()

	fl.Write(res.Bytes())

	fmt.Println(err)
}
