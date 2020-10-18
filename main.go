package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	_ "github.com/lib/pq"
)

func minFl(v1 float64, v2 float64) float64 {
	if v1 > v2 {
		return v2
	}

	return v1
}

func maxFl(v1 float64, v2 float64) float64 {
	if v1 < v2 {
		return v2
	}

	return v1
}

var store = map[int]float64{}
var names = map[int]string{}

const goal float64 = 200.0

var goalEnd = time.Date(2020, 12, 31, 23, 59, 59, 0, time.UTC)

func leftDays() int {
	return int(goalEnd.Sub(time.Now()).Hours() / 24)
}

func registerDistance(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	fromID := update.Message.From.ID

	dirtyParts := strings.Split(update.Message.Text, " ")
	parts := make([]string, 0, len(dirtyParts))

	for _, v := range dirtyParts {
		if v == "" {
			continue
		}

		parts = append(parts, v)
	}

	if len(parts) < 2 {
		return
	}

	numb, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return
	}

	err = insertData(db, fromID, numb, fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName))
	if err != nil {
		return
	}

	err = refreshData(db)
	if err != nil {
		return
	}

	currNum := store[fromID]

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, registerMessageDistanceMsg(numb, currNum, goal, leftDays()))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML

	bot.Send(msg)
}

func removeDistance(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	fromID := update.Message.From.ID

	dirtyParts := strings.Split(update.Message.Text, " ")
	parts := make([]string, 0, len(dirtyParts))

	for _, v := range dirtyParts {
		if v == "" {
			continue
		}

		parts = append(parts, v)
	}

	numb, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return
	}

	err = insertData(db, fromID, -numb, fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName))
	if err != nil {
		return
	}

	err = refreshData(db)
	if err != nil {
		return
	}

	currNum := store[fromID]

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, removeMessageDistanceMsg(numb, currNum, goal, leftDays()))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML

	bot.Send(msg)
}

func myDistance(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	fromID := update.Message.From.ID

	currDistance := 0.0
	if v, ok := store[fromID]; ok {
		currDistance = v
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, myMessageDistanceMsg(currDistance, goal, leftDays()))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML

	bot.Send(msg)
}

func distanceStats(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, statsMessageDistanceMsg(store, names, goal, leftDays()))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML

	bot.Send(msg)
}

func main() {
	var err error

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatalln("$PORT is required")
	}

	url, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatalln("$DATABASE_URL is required")
	}

	db, err := connect(url)
	if err != nil {
		log.Fatalf("Connection error: %s", err.Error())
	}

	err = refreshData(db)
	if err != nil {
		log.Fatalf("Refresh data error: %s", err.Error())
	}

	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatalln("$BOT_TOKEN is required")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 5

		updates, err := bot.GetUpdatesChan(u)
		if err != nil {
			panic(err)
		}

		for update := range updates {
			log.Println(update)

			if update.Message == nil || update.Message.From == nil { // ignore any non-Message Updates
				continue
			}

			switch {
			case strings.HasPrefix(update.Message.Text, "/add"):
				registerDistance(update, bot, db)
			case strings.HasPrefix(update.Message.Text, "/register_distance"):
				registerDistance(update, bot, db)
			case strings.HasPrefix(update.Message.Text, "/delete"):
				removeDistance(update, bot, db)
			case strings.HasPrefix(update.Message.Text, "/remove_distance"):
				removeDistance(update, bot, db)
			case strings.HasPrefix(update.Message.Text, "/my"):
				myDistance(update, bot)
			case strings.HasPrefix(update.Message.Text, "/my_distance"):
				myDistance(update, bot)
			case strings.HasPrefix(update.Message.Text, "/stats"):
				distanceStats(update, bot)
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		}
	}()

	http.HandleFunc("/", hello)
	http.ListenAndServe(":"+port, nil)
}

func hello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok"))
}
