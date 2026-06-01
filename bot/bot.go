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

	menu := &tele.ReplyMarkup{}
	btnBackup := menu.Data("📦 Backup", "btn_backup")
	btnStatus := menu.Data("📊 Status", "btn_status")
	btnServer := menu.Data("🖥️ Server", "btn_server")

	menu.Inline(
		menu.Row(btnBackup, btnStatus),
		menu.Row(btnServer),
	)

	serverMenu := &tele.ReplyMarkup{}
	btnRAM := serverMenu.Data("💾 RAM", "btn_ram")
	btnCPU := serverMenu.Data("⚡ CPU", "btn_cpu")
	btnDisk := serverMenu.Data("💿 Disk", "btn_disk")
	btnDocker := serverMenu.Data("🐳 Docker", "btn_docker")
	btnNetwork := serverMenu.Data("🌐 Network", "btn_network")
	btnFullInfo := serverMenu.Data("📋 Full Info", "btn_full_info")
	btnBack := serverMenu.Data("🔙 Back", "btn_back")

	serverMenu.Inline(
		serverMenu.Row(btnRAM, btnCPU, btnDisk),
		serverMenu.Row(btnDocker, btnNetwork),
		serverMenu.Row(btnFullInfo),
		serverMenu.Row(btnBack),
	)

	b.Handle("/start", func(c tele.Context) error {
		msg := "🤖 Hi, I’m NwwOne\n\nThe central management system of the Nww Ecosystem.\n\nI'm here to assist with infrastructure monitoring, backups, project management, and automated operations across all connected services\n\nSelect an action below."
		return c.Send(msg, menu)
	})

	b.Handle(&btnBackup, func(c tele.Context) error {
		c.Respond()
		return handleBackup(c)
	})

	b.Handle(&btnStatus, func(c tele.Context) error {
		c.Respond()
		return handleStatus(c)
	})

	b.Handle(&btnServer, func(c tele.Context) error {
		c.Respond()
		return c.Edit("🖥️ Server Management\n\nSelect a metric to view:", serverMenu)
	})

	b.Handle(&btnRAM, func(c tele.Context) error {
		c.Respond()
		return c.Edit(GetRAMUsage(), serverMenu)
	})

	b.Handle(&btnCPU, func(c tele.Context) error {
		c.Respond()
		return c.Edit(GetCPUStatus(), serverMenu)
	})

	b.Handle(&btnDisk, func(c tele.Context) error {
		c.Respond()
		return c.Edit(GetDiskUsage(), serverMenu)
	})

	b.Handle(&btnDocker, func(c tele.Context) error {
		c.Respond()
		return c.Edit(GetDockerContainers(), serverMenu)
	})

	b.Handle(&btnNetwork, func(c tele.Context) error {
		c.Respond()
		return c.Edit(GetNetworkStats(), serverMenu)
	})

	b.Handle(&btnFullInfo, func(c tele.Context) error {
		c.Respond()
		return c.Edit(GetFullInfo(), serverMenu)
	})

	b.Handle(&btnBack, func(c tele.Context) error {
		c.Respond()
		msg := "🤖 Hi, I’m NwwOne\n\nThe central management system of the Nww Ecosystem.\n\nI'm here to assist with infrastructure monitoring, backups, project management, and automated operations across all connected services\n\nSelect an action below."
		return c.Edit(msg, menu)
	})

	// Commands
	b.Handle("/backup", handleBackup)
	b.Handle("/status", handleStatus)

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

func handleBackup(c tele.Context) error {
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
}

func handleStatus(c tele.Context) error {
	mongoStatus := "🔴 MongoDB Disconnected"
	if config.Database != nil {
		mongoStatus = "🟢 MongoDB Connected"
	}

	lastBackup := GetLastBackupDate()
	nextBackup := GetNextBackupTime()
	uptime := getUptime()

	statusMsg := fmt.Sprintf(`📊 App Ecosystem Status

	🟢 API Online
	%s
	🟢 Telegram Bot Running
	🟢 Docker Container Running

	📅 Next Backup: %s
	✅ Last Backup: %s

	⏱️ Uptime: %s`, mongoStatus, nextBackup, lastBackup, uptime)

	return c.Send(statusMsg)
}