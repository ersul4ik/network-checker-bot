package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var welcomeMessage = "Welcome to Network Checker Bot. Please send 'check' to start network checking "
var byeMessage = "Come again"
var beginCheckMessage = "Please write the URL that you want to check (www.google.com)"

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Start checking", "/check"),
		tgbotapi.NewInlineKeyboardButtonData("Exit", "/exit"),
	),
)

func getNetworkStatus(url string) string {
	pinger, err := ping.NewPinger(url)
	if err != nil {
		log.Println("\n PANIC")
		panic(err)
	}
	pinger.Count = 3
	if err := pinger.Run(); err != nil {
		log.Panic(err)
	}

	stats := pinger.Statistics()
	return fmt.Sprintf("Percentage of packets lost (%s): %f", url, stats.PacketLoss)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.CallbackQuery != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
			switch update.CallbackQuery.Data {
			case "/check":
				msg.Text = beginCheckMessage
				go func() {
					status := getNetworkStatus("www.google.com")
					msg.Text = status
					_, err := bot.Send(msg)
					if err != nil {
						log.Panic(err)
					}
				}()
			case "/exit":
				msg.Text = byeMessage
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}
			continue
		}

		if update.Message == nil { // ignore non-Message updates
			log.Printf("msg is empty")
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = welcomeMessage
				msg.ReplyMarkup = numericKeyboard
			default:
				msg.Text = "I don't know that command"
			}

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			continue
		}
	}
}
