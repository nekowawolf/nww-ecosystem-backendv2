package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	tele "gopkg.in/telebot.v3"
)

var userNoteState = make(map[int64]string)

var (
	notesMenu = &tele.ReplyMarkup{}
	btnAddNote = notesMenu.Data("➕ Add Note", "btn_add_note")
	btnViewNotes = notesMenu.Data("📋 View Notes", "btn_view_notes")
	btnNotesBack = notesMenu.Data("🔙 Back to Web Tools", "btn_notes_back")
	btnBackToNotes = notesMenu.Data("⬅️ Back to Notes", "btn_back_to_notes_menu")

	addNoteCategoryMenu = &tele.ReplyMarkup{}
	btnCatJournal = addNoteCategoryMenu.Data("📓 Journal", "btn_cat_journal")
	btnCatIdea = addNoteCategoryMenu.Data("💡 Idea", "btn_cat_idea")
	btnCatTask = addNoteCategoryMenu.Data("✅ Task", "btn_cat_task")
	btnCatCancel = addNoteCategoryMenu.Data("❌ Cancel", "btn_cat_cancel")

	cancelNoteMenu = &tele.ReplyMarkup{}
	btnCancelInput = cancelNoteMenu.Data("❌ Cancel", "btn_cancel_note_input")
)

func init() {
	notesMenu.Inline(
		notesMenu.Row(btnAddNote, btnViewNotes),
		notesMenu.Row(btnNotesBack),
	)

	addNoteCategoryMenu.Inline(
		addNoteCategoryMenu.Row(btnCatJournal, btnCatIdea, btnCatTask),
		addNoteCategoryMenu.Row(btnCatCancel),
	)

	cancelNoteMenu.Inline(
		cancelNoteMenu.Row(btnCancelInput),
	)
}

func handleNotesMenu(c tele.Context) error {
	c.Respond()
	if !checkAuth(c) {
		return c.Send("❌ Unauthorized access.")
	}
	msg := "📝 *Notes Manager*\n\nManage your daily journals, ideas, and tasks here. What would you like to do?"
	return c.EditOrSend(msg, notesMenu, tele.ModeMarkdown)
}

func handleAddNoteMenu(c tele.Context) error {
	c.Respond()
	msg := "Please select a category for your new note:"
	return c.Edit(msg, addNoteCategoryMenu)
}

func handleCategorySelection(c tele.Context, category string) error {
	c.Respond()
	userNoteState[c.Chat().ID] = category
	msg := fmt.Sprintf("Category: *%s*\n\nPlease type your note and send it to me. (Or click Cancel)", category)
	return c.Edit(msg, cancelNoteMenu, tele.ModeMarkdown)
}

func handleNoteInput(c tele.Context) error {
	if !checkAuth(c) {
		return nil
	}

	category, ok := userNoteState[c.Chat().ID]
	if !ok {
		return nil
	}

	delete(userNoteState, c.Chat().ID)

	content := c.Message().Text
	if content == "" {
		return c.Send("❌ Error: Note content cannot be empty. Please try again from the Notes menu.", notesMenu)
	}

	lines := strings.Split(content, "\n")
	title := lines[0]
	if len(title) > 30 {
		title = title[:27] + "..."
	}

	note := models.Note{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Content:   content,
		Type:      strings.ToLower(category),
		CreatedAt: time.Now(),
	}

	_, err := config.Database.Collection("notes").InsertOne(context.Background(), note)
	if err != nil {
		log.Printf("Failed to insert note: %v", err)
		return c.Send("❌ Failed to save your note due to a database error.")
	}

	msg := fmt.Sprintf("✅ *Note successfully saved!*\n\n*Category:* %s\n*Date:* %s\n*Content:* %s", category, note.CreatedAt.Format("02 Jan 2006"), title)
	
	menu := &tele.ReplyMarkup{}
	btnAddAnother := menu.Data("➕ Add Another", "btn_add_note")
	btnBack := menu.Data("⬅️ Back to Notes", "btn_notes_back_main")
	menu.Inline(
		menu.Row(btnAddAnother, btnBack),
	)

	return c.Send(msg, menu, tele.ModeMarkdown)
}

func getMonthName(month int) string {
	months := []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	if month >= 1 && month <= 12 {
		return months[month]
	}
	return "Unknown"
}

func handleViewNotesYears(c tele.Context) error {
	c.Respond()

	pipeline := mongo.Pipeline{
		{{"$project", bson.D{{"year", bson.D{{"$year", "$created_at"}}}}}},
		{{"$group", bson.D{{"_id", "$year"}, {"count", bson.D{{"$sum", 1}}}}}},
		{{"$sort", bson.D{{"_id", -1}}}},
	}

	cursor, err := config.Database.Collection("notes").Aggregate(context.Background(), pipeline)
	if err != nil {
		return c.Send("❌ Failed to fetch data.")
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return c.Send("❌ Failed to process data.")
	}

	if len(results) == 0 {
		return c.Edit("📝 You don't have any notes yet.", notesMenu)
	}

	menu := &tele.ReplyMarkup{}
	var rows []tele.Row
	var buttons []tele.Btn

	for i, res := range results {
		
		year, ok := res["_id"].(int32)
		if !ok {
			continue
		}
		count, ok := res["count"].(int32)
		if !ok {
			continue
		}
		btnText := fmt.Sprintf("%d (%d Notes)", year, count)
		btn := menu.Data(btnText, "yr", fmt.Sprint(year))
		buttons = append(buttons, btn)

		if len(buttons) == 2 || i == len(results)-1 {
			rows = append(rows, menu.Row(buttons...))
			buttons = []tele.Btn{}
		}
	}
	
	rows = append(rows, menu.Row(btnBackToNotes))
	menu.Inline(rows...)

	msg := "📅 *Notes Archive*\n\nHere are the years you have notes. Please select a year:"
	return c.Edit(msg, menu, tele.ModeMarkdown)
}

func handleViewNotesMonths(c tele.Context) error {
	c.Respond()
	yearStr := c.Callback().Data
	year, err := strconv.Atoi(strings.TrimSpace(yearStr))
	if err != nil {
		return c.Send("❌ Invalid year.")
	}

	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

	matchStage := bson.D{{"$match", bson.D{{"created_at", bson.D{{"$gte", startDate}, {"$lt", endDate}}}}}}
	projectStage := bson.D{{"$project", bson.D{{"month", bson.D{{"$month", "$created_at"}}}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$month"}, {"count", bson.D{{"$sum", 1}}}}}}
	sortStage := bson.D{{"$sort", bson.D{{"_id", 1}}}}

	pipeline := mongo.Pipeline{matchStage, projectStage, groupStage, sortStage}

	cursor, err := config.Database.Collection("notes").Aggregate(context.Background(), pipeline)
	if err != nil {
		return c.Send("❌ Failed to fetch data.")
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return c.Send("❌ Failed to process data.")
	}

	menu := &tele.ReplyMarkup{}
	
	btnViewAll := menu.Data(fmt.Sprintf("📄 View All %d", year), "vyr", fmt.Sprint(year))
	var rows []tele.Row
	rows = append(rows, menu.Row(btnViewAll))

	var buttons []tele.Btn
	for i, res := range results {
		month, ok := res["_id"].(int32)
		if !ok {
			continue
		}
		count, ok := res["count"].(int32)
		if !ok {
			continue
		}
		btnText := fmt.Sprintf("%s (%d)", getMonthName(int(month)), count)
		dataStr := fmt.Sprintf("%d_%d", year, month)
		btn := menu.Data(btnText, "mo", dataStr)
		buttons = append(buttons, btn)

		if len(buttons) == 2 || i == len(results)-1 {
			rows = append(rows, menu.Row(buttons...))
			buttons = []tele.Btn{}
		}
	}
	
	btnBackToYears := menu.Data("🔙 Back to Years", "btn_view_notes")
	rows = append(rows, menu.Row(btnBackToYears))
	menu.Inline(rows...)

	msg := fmt.Sprintf("📅 *Year: %d*\nYou can view all notes in %d, or filter by available months:", year, year)
	return c.Edit(msg, menu, tele.ModeMarkdown)
}

func handleViewNotesDates(c tele.Context) error {
	c.Respond()
	dataStr := strings.TrimSpace(c.Callback().Data)
	parts := strings.Split(dataStr, "_")
	if len(parts) != 2 {
		return c.Send("❌ Invalid data.")
	}
	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	matchStage := bson.D{{"$match", bson.D{{"created_at", bson.D{{"$gte", startDate}, {"$lt", endDate}}}}}}
	projectStage := bson.D{{"$project", bson.D{{"day", bson.D{{"$dayOfMonth", "$created_at"}}}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$day"}, {"count", bson.D{{"$sum", 1}}}}}}
	sortStage := bson.D{{"$sort", bson.D{{"_id", 1}}}}

	pipeline := mongo.Pipeline{matchStage, projectStage, groupStage, sortStage}

	cursor, err := config.Database.Collection("notes").Aggregate(context.Background(), pipeline)
	if err != nil {
		return c.Send("❌ Failed to fetch data.")
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return c.Send("❌ Failed to process data.")
	}

	menu := &tele.ReplyMarkup{}
	
	btnViewAll := menu.Data(fmt.Sprintf("📄 View All %s %d", getMonthName(month), year), "vmo", dataStr)
	var rows []tele.Row
	rows = append(rows, menu.Row(btnViewAll))

	var buttons []tele.Btn
	for i, res := range results {
		day, ok := res["_id"].(int32)
		if !ok {
			continue
		}
		btnText := fmt.Sprintf("Day %02d", day)
		dateStr := fmt.Sprintf("%d_%d_%d", year, month, day)
		btn := menu.Data(btnText, "dt", dateStr)
		buttons = append(buttons, btn)

		if len(buttons) == 3 || i == len(results)-1 {
			rows = append(rows, menu.Row(buttons...))
			buttons = []tele.Btn{}
		}
	}
	
	btnBackToMonths := menu.Data("🔙 Back to Months", "yr", fmt.Sprint(year))
	rows = append(rows, menu.Row(btnBackToMonths))
	menu.Inline(rows...)

	msg := fmt.Sprintf("📅 *Month: %s %d*\nYou can view all notes this month, or select a specific date:", getMonthName(month), year)
	return c.Edit(msg, menu, tele.ModeMarkdown)
}

func handleViewNotesList(c tele.Context, mode string, dataStr string) error {
	c.Respond()
	
	var startDate, endDate time.Time
	var titleMsg string
	var backBtn tele.Btn

	menu := &tele.ReplyMarkup{}

	if mode == "year" {
		year, _ := strconv.Atoi(dataStr)
		startDate = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		titleMsg = fmt.Sprintf("📅 Notes in %d", year)
		backBtn = menu.Data("🔙 Back to Months", "yr", dataStr)
	} else if mode == "month" {
		parts := strings.Split(dataStr, "_")
		year, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])
		startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, 0)
		titleMsg = fmt.Sprintf("📅 Notes in %s %d", getMonthName(month), year)
		backBtn = menu.Data("🔙 Back to Dates", "mo", dataStr)
	} else if mode == "date" {
		parts := strings.Split(dataStr, "_")
		year, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])
		day, _ := strconv.Atoi(parts[2])
		startDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 0, 1)
		titleMsg = fmt.Sprintf("📅 Notes on %02d %s %d", day, getMonthName(month), year)
		backBtn = menu.Data("🔙 Back to Dates", "mo", fmt.Sprintf("%d_%d", year, month))
	}

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	opts := options.Find().SetSort(bson.D{{"created_at", 1}})
	cursor, err := config.Database.Collection("notes").Find(context.Background(), filter, opts)
	if err != nil {
		return c.Send("❌ Failed to fetch notes.")
	}
	defer cursor.Close(context.Background())

	var notes []models.Note
	if err = cursor.All(context.Background(), &notes); err != nil {
		return c.Send("❌ Failed to process notes.")
	}

	if len(notes) == 0 {
		msg := fmt.Sprintf("*%s*\n\nNo notes found.", titleMsg)
		menu.Inline(menu.Row(backBtn))
		return c.Edit(msg, menu, tele.ModeMarkdown)
	}

	msg := fmt.Sprintf("*%s*\n\n", titleMsg)
	for i, note := range notes {
		icon := "📓"
		if note.Type == "idea" {
			icon = "💡"
		} else if note.Type == "task" {
			icon = "✅"
		}
		
		timeStr := note.CreatedAt.Format("15:04")
		entry := fmt.Sprintf("%d. [%s] *%s* (%s)\n%s\n\n", i+1, icon, note.Title, timeStr, note.Content)
		
		if len(msg)+len(entry) > 4000 {
			msg += "\n... (truncated due to length)"
			break
		}
		msg += entry
	}

	menu.Inline(menu.Row(backBtn))
	return c.Edit(msg, menu, tele.ModeMarkdown)
}

func RegisterNotesHandlers(b *tele.Bot, webToolsMenu *tele.ReplyMarkup) {
	b.Handle("\fbtn_notes", handleNotesMenu)
	b.Handle(&btnNotesBack, func(c tele.Context) error {
		c.Respond()
		return c.Edit("🌐 Web Tools\n\nPlease select an action below:", webToolsMenu)
	})
	b.Handle(&btnBackToNotes, func(c tele.Context) error {
		return handleNotesMenu(c)
	})

	b.Handle(&btnAddNote, handleAddNoteMenu)
	b.Handle(&btnCatJournal, func(c tele.Context) error { return handleCategorySelection(c, "Journal") })
	b.Handle(&btnCatIdea, func(c tele.Context) error { return handleCategorySelection(c, "Idea") })
	b.Handle(&btnCatTask, func(c tele.Context) error { return handleCategorySelection(c, "Task") })
	
	b.Handle(&btnCatCancel, func(c tele.Context) error {
		c.Respond()
		delete(userNoteState, c.Chat().ID)
		return handleNotesMenu(c)
	})
	b.Handle(&btnCancelInput, func(c tele.Context) error {
		c.Respond()
		delete(userNoteState, c.Chat().ID)
		return handleNotesMenu(c)
	})

	b.Handle(&btnViewNotes, handleViewNotesYears)
	
	b.Handle("\fbtn_notes_back_main", func(c tele.Context) error {
		return handleNotesMenu(c)
	})

	// Dynamic Drill-down Handlers
	b.Handle(&tele.Btn{Unique: "yr"}, handleViewNotesMonths)
	b.Handle(&tele.Btn{Unique: "mo"}, handleViewNotesDates)
	b.Handle(&tele.Btn{Unique: "dt"}, func(c tele.Context) error { return handleViewNotesList(c, "date", c.Callback().Data) })
	b.Handle(&tele.Btn{Unique: "vyr"}, func(c tele.Context) error { return handleViewNotesList(c, "year", c.Callback().Data) })
	b.Handle(&tele.Btn{Unique: "vmo"}, func(c tele.Context) error { return handleViewNotesList(c, "month", c.Callback().Data) })
}

func CheckNotesText(c tele.Context) (bool, error) {
	_, ok := userNoteState[c.Chat().ID]
	if ok {
		return true, handleNoteInput(c)
	}
	return false, nil
}