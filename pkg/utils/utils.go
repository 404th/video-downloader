package utils

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// DownloadAndSendVideo downloads the video with yt-dlp and sends it to Telegram
func DownloadAndSendVideo(bot *tgbotapi.BotAPI, chatID int64, processingMsgID int, url string) error {
	// Temporary file to store the video
	tempFile := filepath.Join(os.TempDir(), "video-"+ExtractVideoID(url)+".mp4")
	defer os.Remove(tempFile) // Clean up after sending

	// Configure yt-dlp command
	cmd := exec.Command(
		"yt-dlp",
		"--cookies", "./assets/instagram_cookies.txt",
		"-o", tempFile, // Save to file instead of stdout for Telegram
		"-f", "b", // Best format with video+audio
		url,
	)

	// Capture stderr for debugging
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Run the command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Log stderr
	go func() {
		errBytes, _ := io.ReadAll(stderr)
		if len(errBytes) > 0 {
			log.Printf("yt-dlp stderr: %s", errBytes)
		}
	}()

	// Wait for download to complete
	if err := cmd.Wait(); err != nil {
		return err
	}

	// Open the downloaded file
	file, err := os.Open(tempFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create Telegram video message
	video := tgbotapi.NewVideo(chatID, tgbotapi.FileReader{
		Name:   "video-" + ExtractVideoID(url) + ".mp4",
		Reader: file,
	})

	// Send the video
	_, err = bot.Send(video)
	if err != nil {
		return err
	}

	// Delete the "Downloading..." message after video is sent
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, processingMsgID)
	if _, err := bot.Request(deleteMsg); err != nil {
		log.Printf("Failed to delete processing message (ID: %d): %v", processingMsgID, err)
		return err
	}

	return nil
}

// ExtractVideoID extracts an ID from the URL
func ExtractVideoID(url string) string {
	parts := strings.Split(url, "/")
	for i, part := range parts {
		if (part == "reel" || part == "p" || part == "status") && i+1 < len(parts) {
			id := parts[i+1]
			if idx := strings.Index(id, "?"); idx != -1 {
				return id[:idx]
			}
			return id
		}
	}
	return "video"
}

// IsValidCategory determines if the message is a valid command or URL
func IsValidCategory(category string) bool {
	switch category {
	case "/start", "/stop":
		return false
	}
	return true
}
