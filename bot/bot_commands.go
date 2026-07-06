package bot

import (
	"log"
	tele "gopkg.in/telebot.v3"
)

func RegisterBotCommands(b *tele.Bot, menu *tele.ReplyMarkup, serverMenu *tele.ReplyMarkup) {
	commands := []tele.Command{
		{Text: "start", Description: "open main menu"},
		{Text: "backup", Description: "trigger database backup"},
		{Text: "status", Description: "check ecosystem health"},
		{Text: "fullinfo", Description: "view server hardware info"},
		{Text: "speedtest", Description: "run api speed test"},
		{Text: "notes", Description: "open notes menu"},
		{Text: "add_journal", Description: "add a new journal note"},
		{Text: "add_idea", Description: "add a new idea note"},
		{Text: "add_task", Description: "add a new task note"},
		{Text: "view_notes", Description: "view all notes"},
		{Text: "view_journal", Description: "view journal notes"},
		{Text: "view_idea", Description: "view idea notes"},
		{Text: "view_task", Description: "view task notes"},
		{Text: "view_all_task", Description: "view all task notes directly"},
		{Text: "view_all_idea", Description: "view all idea notes directly"},
		{Text: "manage_notes", Description: "manage all notes"},
		{Text: "manage_journal", Description: "manage journal notes"},
		{Text: "manage_idea", Description: "manage idea notes"},
		{Text: "manage_task", Description: "manage task notes"},
	}

	// Set for Default Scope
	err := b.SetCommands(commands)
	if err != nil {
		log.Printf("Failed to set bot commands (default scope): %v", err)
	}

	b.Handle("/start", func(c tele.Context) error {
		msg := "🤖 Hi, I’m NwwOne\n\nThe central management system of the Nww Ecosystem.\n\nI'm here to assist with infrastructure monitoring, backups, project management, and automated operations across all connected services\n\nSelect an action below."
		return c.Send(msg, menu)
	})

	// Add Note Commands
	b.Handle("/add_journal", func(c tele.Context) error { return handleAddNoteCommand(c, "Journal") })
	b.Handle("/add_idea", func(c tele.Context) error { return handleAddNoteCommand(c, "Idea") })
	b.Handle("/add_task", func(c tele.Context) error { return handleAddNoteCommand(c, "Task") })

	// View Note Commands (Static)
	b.Handle("/view_notes", func(c tele.Context) error { return handleViewNoteCommand(c, "all") })
	b.Handle("/view_journal", func(c tele.Context) error { return handleViewNoteCommand(c, "journal") })
	b.Handle("/view_idea", func(c tele.Context) error { return handleViewNoteCommand(c, "idea") })
	b.Handle("/view_task", func(c tele.Context) error { return handleViewNoteCommand(c, "task") })

	// View All Direct Commands
	b.Handle("/view_all_task", func(c tele.Context) error { return handleViewAllNoteCommand(c, "task") })
	b.Handle("/view_all_idea", func(c tele.Context) error { return handleViewAllNoteCommand(c, "idea") })

	// Manage Note Commands (Dynamic)
	b.Handle("/manage_notes", func(c tele.Context) error { return handleManageNoteCommand(c, "all") })
	b.Handle("/manage_journal", func(c tele.Context) error { return handleManageNoteCommand(c, "journal") })
	b.Handle("/manage_idea", func(c tele.Context) error { return handleManageNoteCommand(c, "idea") })
	b.Handle("/manage_task", func(c tele.Context) error { return handleManageNoteCommand(c, "task") })

	// General Commands
	b.Handle("/backup", handleBackup)
	b.Handle("/status", handleStatus)
	b.Handle("/fullinfo", func(c tele.Context) error {
		if !checkAuth(c) {
			return c.Send("❌ Unauthorized access.")
		}
		return c.Send(GetFullInfo(), serverMenu)
	})
	b.Handle("/speedtest", handleSpeedTest)
	b.Handle("/notes", handleNotesMenu)
}