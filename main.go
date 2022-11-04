package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/enescakir/emoji"
)

var helpMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Help", "help"),
	),
)

var generateJoke = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Joke", "joke"),
	),
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("GOLANGTZBOT"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		go singleUpdate(update, bot)
	}
}

func singleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {

	if update.CallbackQuery != nil {

		responseMsg := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "")

		receivedData := update.CallbackData()

		switch receivedData {
		case "help":
			responseMsg.Text = "Coming Soon"
			if _, err := bot.Request(responseMsg); err != nil {
				log.Panic(err)
			}
		case "joke":
			responseMsg.Text = Joke()
			responseMsg.CacheTime = 10
			if _, err := bot.Request(responseMsg); err != nil {
				log.Panic(err)
			}

		}
		return
	}

	if update.Message.NewChatMembers != nil {
		deleteRequest := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		if _, err := bot.Request(deleteRequest); err != nil {
			log.Panic(err)
		}

		member, _ := bot.GetChatMembersCount(tgbotapi.ChatMemberCountConfig{ChatConfig: update.FromChat().ChatConfig()})
		welcomeText := fmt.Sprintf("Karibu %v %v. You are member number %v", update.Message.From.FirstName, emoji.WavingHand, member)
		welcomeMsg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeText)
		welcomeMsg.ReplyMarkup = helpMenu
		sentMsg, err := bot.Send(welcomeMsg)

		if err != nil {
			log.Panic(err)
		}

		time.Sleep(60 * time.Second)
		deleteMsg := tgbotapi.NewDeleteMessage(sentMsg.Chat.ID, sentMsg.MessageID)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Panic(err)
		}

		return

	}

	if update.Message.LeftChatMember != nil {
		deleteRequest := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		if _, err := bot.Request(deleteRequest); err != nil {
			log.Panic(err)
		}
		goodbyeMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Another fallen soldier")

		if _, err := bot.Send(goodbyeMsg); err != nil {
			log.Panic(err)
		}
		return
	}

	if !update.Message.IsCommand() {
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "help":
		msg.Text = "So far I can only make Jokes. Use /joke"
		msg.ReplyMarkup = helpMenu
	case "hi":
		msg.Text = "Hello there! :)"
	case "status":
		msg.Text = "I'm incomplete :("
	case "joke":
		msg.Text = "Tap the button to get random programming jokes"
		msg.ReplyMarkup = generateJoke
	default:
		msg.Text = "I don't know that command"
	}

	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}

}
