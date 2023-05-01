package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
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

const goal float64 = 700.0

var startDate = time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
var goalEnd = time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
var totalDays = math.Ceil(goalEnd.Sub(startDate).Hours() / 24)

func leftDays() int {
	return int(goalEnd.Sub(time.Now()).Hours() / 24)
}

func registerDistance(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	fmt.Print("Started: registerDistance")

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

	if currNum >= goal {
		animMsg := tgbotapi.NewStickerShare(update.Message.Chat.ID, "CAACAgIAAxkBAAIEB1_cuGhuVwkH421EKWjNt7pCWtbVAAJBAAOYv4AN5hxhbUCronYeBA")
		_, err = bot.Send(animMsg)
		if err != nil {
			log.Print(err)
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, registerMessageDistanceMsg(numb, currNum, goal, leftDays()))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML

	_, err = bot.Send(msg)
	if err != nil {
		log.Print(err)
	}
}

func removeDistance(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	fmt.Print("Started: removeDistance")

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

func myDistance(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	fmt.Print("Started: myDistance")

	fromID := update.Message.From.ID

	currDistance := 0.0
	if v, ok := store[fromID]; ok {
		currDistance = v
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, myMessageDistanceMsg(currDistance, goal, leftDays()))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML

	dt, err := getSingleUserData(db, fromID)
	if err != nil {
		bot.Send(msg)
		fmt.Printf("%+v", err)

		return
	}

	daysKm := runningDistanceToDayKm(dt, startDate, goalEnd)
	buffer, err := drawChart(uint(goal), uint(goalEnd.Sub(startDate).Hours()/24), [][]float64{daysKm})
	if err != nil {
		bot.Send(msg)
		fmt.Printf("%+v", err)

		return
	}

	bts := tgbotapi.FileBytes{
		Name:  time.Now().String() + ".png",
		Bytes: buffer.Bytes(),
	}

	res, err := bot.UploadFile("sendPhoto", map[string]string{
		"chat_id":             fmt.Sprint(msg.ChatID),
		"caption":             msg.Text,
		"parse_mode":          msg.ParseMode,
		"reply_to_message_id": fmt.Sprint(msg.ReplyToMessageID),
	}, "photo", bts)
	if err != nil {
		bot.Send(msg)
		fmt.Printf("%+v", err)

		return
	}

	fmt.Printf("%+v", res)
}

func avgPrediction(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	fmt.Print("Started: avgPrediction")

	fromID := update.Message.From.ID

	currDistance := 0.0
	if v, ok := store[fromID]; ok {
		currDistance = v
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, myMessageDistanceMsg(currDistance, goal, leftDays()))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML

	avgPerDay := currDistance / float64(365-leftDays())
	buffer, err := drawSuccessPredChard(int(goal), 365, avgPerDay, float64(365-leftDays()))
	if err != nil {
		bot.Send(msg)
		fmt.Printf("%+v", err)

		return
	} else {
		fmt.Print("ERROR: ", err)
	}

	bts := tgbotapi.FileBytes{
		Name:  time.Now().String() + ".png",
		Bytes: buffer.Bytes(),
	}

	res, err := bot.UploadFile("sendPhoto", map[string]string{
		"chat_id":             fmt.Sprint(msg.ChatID),
		"caption":             msg.Text,
		"parse_mode":          msg.ParseMode,
		"reply_to_message_id": fmt.Sprint(msg.ReplyToMessageID),
	}, "photo", bts)
	if err != nil {
		bot.Send(msg)
		fmt.Printf("%+v", err)

		return
	}

	fmt.Printf("%+v", res)
}

func distanceStats(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	fmt.Print("Started: distanceStats")

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, statsMessageDistanceMsg(store, names, goal, leftDays()))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML

	dat := make([][]float64, 0, len(store))
	for k := range store {
		dt, err := getSingleUserData(db, k)
		if err != nil {
			bot.Send(msg)
			fmt.Printf("%+v", err)

			return
		}

		daysKm := runningDistanceToDayKm(dt, startDate, goalEnd)
		dat = append(dat, daysKm)
	}

	buffer, err := drawChart(uint(goal), uint(goalEnd.Sub(startDate).Hours()/24), dat)
	if err != nil {
		bot.Send(msg)
		fmt.Printf("%+v", err)

		return
	}

	bts := tgbotapi.FileBytes{
		Name:  time.Now().String() + ".png",
		Bytes: buffer.Bytes(),
	}

	res, err := bot.UploadFile("sendPhoto", map[string]string{
		"chat_id":             fmt.Sprint(msg.ChatID),
		"caption":             msg.Text,
		"parse_mode":          msg.ParseMode,
		"reply_to_message_id": fmt.Sprint(msg.ReplyToMessageID),
	}, "photo", bts)
	if err != nil {
		bot.Send(msg)
		fmt.Printf("%+v", err)

		return
	}

	fmt.Printf("%+v", res)
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
			case strings.HasPrefix(update.Message.Text, "/my_distance"):
				myDistance(update, bot, db)
			case strings.HasPrefix(update.Message.Text, "/my_avgpred"):
				avgPrediction(update, bot, db)
			case strings.HasPrefix(update.Message.Text, "/my"):
				myDistance(update, bot, db)
			case strings.HasPrefix(update.Message.Text, "/stats"):
				distanceStats(update, bot, db)
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
