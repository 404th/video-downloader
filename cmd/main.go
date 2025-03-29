package main

import (
	"log"
	"strings"

	"github.com/404th/video-downloader/config"
	"github.com/404th/video-downloader/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Println("❌ Failed to initialize config: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("❌ Failed to initialize bot: %v", err)
	}

	switch cfg.BotMode {
	case "production":
		bot.Debug = false
	default:
		bot.Debug = true
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {
			// Send welcome message to user
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, config.WelcomeMessage)
			bot.Send(msg)

			// Get client's username or first name
			clientUsername := update.Message.From.UserName
			if clientUsername == "" {
				clientUsername = update.Message.From.FirstName
			}

			notificationMsgText := "New client started the bot: @" + clientUsername
			notificationMsg := tgbotapi.NewMessage(cfg.TelegramUsernameChatId, notificationMsgText)
			notificationMsg.ParseMode = "Markdown"

			if _, err := bot.Send(notificationMsg); err != nil {
				log.Printf("Failed to send message to: %v", err)
			}
			continue
		}

		if utils.IsValidCategory(update.Message.Text) {
			url := update.Message.Text
			if !strings.HasPrefix(url, "http") || (!strings.Contains(url, "instagram.com") && !strings.Contains(url, "x.com") && !strings.Contains(url, "twitter.com")) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Please send a valid Instagram or X URL.")
				bot.Send(msg)
				continue
			}

			processingMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏬ Downloading...")
			sentMsg, err := bot.Send(processingMsg)
			if err != nil {
				log.Printf("Failed to send processing message: %v", err)
				continue
			}

			err = utils.DownloadAndSendVideo(bot, update.Message.Chat.ID, sentMsg.MessageID, url)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Something went wrong. Retry... /start")
				bot.Send(msg)
				continue
			}
		}
	}
}
