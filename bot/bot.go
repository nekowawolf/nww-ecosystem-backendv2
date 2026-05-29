package bot

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	tele "gopkg.in/telebot.v3"
)

var TelegramBot *tele.Bot

func InitBot() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Println("TELEGRAM_BOT_TOKEN is not set. Bot will not start.")
		return
	}

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("failed to start bot: %v", err)
		return
	}
	TelegramBot = b

	// Commands
	b.Handle("/backup", func(c tele.Context) error {
		err := c.Send("⏳  Starting backup...")
		if err != nil {
			return err
		}

		buf, filename, err := PerformBackup()
		if err != nil {
			return c.Send(fmt.Sprintf("❌ Backup failed: %v", err))
		}

		doc := &tele.Document{
			File:     tele.FromReader(bytes.NewReader(buf.Bytes())),
			FileName: filename,
		}

		return c.Send(doc)
	})

	b.Handle("/status", func(c tele.Context) error {
		mongoStatus := "🔴 MongoDB offline"
		if config.Database != nil {
			mongoStatus = "🟢 MongoDB online"
		}

		lastBackup := GetLastBackupDate()
		nextBackup := GetNextBackupTime()

		statusMsg := fmt.Sprintf("%s\n🟢 Heroku healthy\n📅 Next backup: %s\n✅ Last backup: %s",
			mongoStatus, nextBackup, lastBackup)

		return c.Send(statusMsg)
	})

	go b.Start()
	log.Println("Telegram bot is running")
}

func SendBackupArchive() {
	if TelegramBot == nil {
		log.Println("Telegram bot not initialized")
		return
	}

	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	if chatIDStr == "" {
		log.Println("TELEGRAM_CHAT_ID is not set")
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Printf("Invalid TELEGRAM_CHAT_ID: %v\n", err)
		return
	}

	chat := &tele.Chat{ID: chatID}

	buf, filename, err := PerformBackup()
	if err != nil {
		log.Printf("Automated backup failed: %v\n", err)
		TelegramBot.Send(chat, fmt.Sprintf("❌ Automated backup failed: %v", err))
		return
	}

	doc := &tele.Document{
		File:     tele.FromReader(bytes.NewReader(buf.Bytes())),
		FileName: filename,
	}

	_, err = TelegramBot.Send(chat, doc)
	if err != nil {
		log.Printf("Failed to send automated backup to chat: %v\n", err)
	} else {
		log.Println("Automated backup sent successfully")
	}
}
