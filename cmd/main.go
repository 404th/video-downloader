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

	// Initialize Telegram bot with your token
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

	// Set up update configuration
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Get updates channel
	updates := bot.GetUpdatesChan(u)

	// Handle incoming messages
	for update := range updates {
		if update.Message == nil { // Ignore non-Message updates
			continue
		}

		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, config.WelcomeMessage)
			bot.Send(msg)
			continue
		}

		if utils.IsValidCategory(update.Message.Text) {
			// Check if the message contains a URL
			url := update.Message.Text
			if !strings.HasPrefix(url, "http") || (!strings.Contains(url, "instagram.com") && !strings.Contains(url, "x.com") && !strings.Contains(url, "twitter.com")) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Please send a valid Instagram or X URL.")
				bot.Send(msg)
				continue
			}

			// Notify user that processing has started
			processingMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏬ Downloading...")
			sentMsg, err := bot.Send(processingMsg)
			if err != nil {
				log.Printf("Failed to send processing message: %v", err)
				continue
			}

			// Download and send the video, then delete the processing message
			err = utils.DownloadAndSendVideo(bot, update.Message.Chat.ID, sentMsg.MessageID, url)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Something went wrong. Retry... /start")
				bot.Send(msg)
				continue
			}
		}
	}
}
