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

	// Security Middleware
	b.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if !checkAuth(c) {
				return c.Send("❌ Unauthorized access.")
			}
			return next(c)
		}
	})

	menu := &tele.ReplyMarkup{}
	btnBackup := menu.Data("📦 Backup", "btn_backup")
	btnStatus := menu.Data("📊 Status", "btn_status")
	btnServer := menu.Data("🖥️ Server", "btn_server")
	btnTools := menu.Data("⚙️ Tools", "btn_tools")

	menu.Inline(
		menu.Row(btnBackup, btnStatus),
		menu.Row(btnServer, btnTools),
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

	toolsMenu := &tele.ReplyMarkup{}
	btnWebTools := toolsMenu.Data("🌐 Web Tools", "btn_web_tools")
	btnCryptoTools := toolsMenu.Data("🪙 Crypto Tools", "btn_crypto_tools")
	btnToolsBack := toolsMenu.Data("🔙 Back", "btn_tools_back")

	toolsMenu.Inline(
		toolsMenu.Row(btnWebTools, btnCryptoTools),
		toolsMenu.Row(btnToolsBack),
	)

	webToolsMenu := &tele.ReplyMarkup{}
	btnSpeedTest := webToolsMenu.Data("⚡ Test Speed API", "btn_speed_test")
	btnCDN := webToolsMenu.Data("🖼️ CDN GitHub", "btn_cdn_github")
	btnMissingImages := webToolsMenu.Data("🔍 Check Missing Images", "btn_missing_images")
	btnCheckInvalidLink := webToolsMenu.Data("🔗 Check Invalid Link", "btn_check_invalid_link")
	btnNotes := webToolsMenu.Data("📝 Notes", "btn_notes")
	btnWebToolsBack := webToolsMenu.Data("🔙 Back", "btn_web_tools_back")

	webToolsMenu.Inline(
		webToolsMenu.Row(btnSpeedTest),
		webToolsMenu.Row(btnCDN),
		webToolsMenu.Row(btnMissingImages),
		webToolsMenu.Row(btnCheckInvalidLink),
		webToolsMenu.Row(btnNotes),
		webToolsMenu.Row(btnWebToolsBack),
	)


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

	b.Handle(&btnTools, func(c tele.Context) error {
		c.Respond()
		return c.Edit("⚙️ Ecosystem Tools\n\nPlease select an action below:", toolsMenu)
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

	b.Handle(&btnToolsBack, func(c tele.Context) error {
		c.Respond()
		msg := "🤖 Hi, I’m NwwOne\n\nThe central management system of the Nww Ecosystem.\n\nI'm here to assist with infrastructure monitoring, backups, project management, and automated operations across all connected services\n\nSelect an action below."
		return c.Edit(msg, menu)
	})

	b.Handle(&btnWebTools, func(c tele.Context) error {
		c.Respond()
		return c.Edit("🌐 Web Tools\n\nPlease select an action below:", webToolsMenu)
	})

	b.Handle(&btnCryptoTools, func(c tele.Context) error {
		return c.Respond(&tele.CallbackResponse{Text: "in development", ShowAlert: true})
	})

	b.Handle(&btnWebToolsBack, func(c tele.Context) error {
		c.Respond()
		return c.Edit("⚙️ Tools\n\nPlease select an action below:", toolsMenu)
	})

	// Tools handlers
	b.Handle(&btnSpeedTest, handleSpeedTest)
	b.Handle(&btnCDN, handleCDNInit)
	b.Handle(&btnMissingImages, handleCheckMissingImages)
	b.Handle(&btnCheckInvalidLink, handleCheckInvalidLink)
	
	b.Handle(&tele.Btn{Unique: "exe_img_chk"}, handleExecuteImageCheck)
	b.Handle(&tele.Btn{Unique: "exe_lnk_chk"}, handleExecuteLinkCheck)
	b.Handle(&tele.Btn{Unique: "btn_cancel_chk"}, func(c tele.Context) error {
		c.Respond()
		return c.Edit("🌐 Web Tools\n\nPlease select an action below:", webToolsMenu)
	})
	
	RegisterNotesHandlers(b, webToolsMenu)

	// Catch text for Notes
	b.Handle(tele.OnText, func(c tele.Context) error {
		handled, err := CheckNotesText(c)
		if handled {
			return err
		}
		return nil
	})

	// Catch photos for CDN upload
	b.Handle(tele.OnPhoto, handlePhotoUpload)

	// Register all text commands
	RegisterBotCommands(b, menu, serverMenu)


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