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
var userNotesContext = make(map[int64][]primitive.ObjectID)
var userEditTarget = make(map[int64]primitive.ObjectID)

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
	title := ""
	body := ""

	if len(lines) > 1 {
		title = strings.TrimSpace(lines[0])
		body = strings.TrimSpace(strings.Join(lines[1:], "\n"))
	} else {
		title = category
		body = strings.TrimSpace(content)
	}

	note := models.Notes{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Content:   body,
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

func handleViewNotesCategory(c tele.Context) error {
	c.Respond()
	menu := &tele.ReplyMarkup{}
	btnAll := menu.Data("📑 All", "vyr", "all")
	btnJournal := menu.Data("📓 Journal", "vyr", "journal")
	btnIdea := menu.Data("💡 Idea", "vyr", "idea")
	btnTask := menu.Data("✅ Task", "vyr", "task")
	menu.Inline(menu.Row(btnAll, btnJournal), menu.Row(btnIdea, btnTask), menu.Row(btnBackToNotes))
	msg := "📅 *Notes Archive*\n\nPlease select a category to view your notes:"
	return c.Edit(msg, menu, tele.ModeMarkdown)
}

func handleViewNotesYears(c tele.Context, category string) error {
	c.Respond()
	filter := bson.M{}
	if category != "all" {
		filter["type"] = category
	}
	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$project", bson.D{{"year", bson.D{{"$year", "$created_at"}}}}}},
		{{"$group", bson.D{{"_id", "$year"}, {"count", bson.D{{"$sum", 1}}}}}},
		{{"$sort", bson.D{{"_id", -1}}}},
	}
	cursor, err := config.Database.Collection("notes").Aggregate(context.Background(), pipeline)
	if err != nil { return c.Send("❌ Failed to fetch data.") }
	defer cursor.Close(context.Background())
	var results []bson.M
	cursor.All(context.Background(), &results)

	if len(results) == 0 {
		return c.Edit("📝 No notes found for this category.", notesMenu)
	}

	menu := &tele.ReplyMarkup{}
	btnViewAll := menu.Data(fmt.Sprintf("📄 View All %s", strings.Title(category)), "vlist", category)
	var rows []tele.Row
	rows = append(rows, menu.Row(btnViewAll))

	var buttons []tele.Btn
	for i, res := range results {
		year, ok := res["_id"].(int32)
		if !ok { continue }
		count := res["count"].(int32)
		btnText := fmt.Sprintf("%d (%d Notes)", year, count)
		btn := menu.Data(btnText, "vmo", fmt.Sprintf("%s_%d", category, year))
		buttons = append(buttons, btn)
		if len(buttons) == 2 || i == len(results)-1 {
			rows = append(rows, menu.Row(buttons...))
			buttons = []tele.Btn{}
		}
	}
	rows = append(rows, menu.Row(menu.Data("🔙 Back to Categories", "btn_view_notes", "back")))
	menu.Inline(rows...)
	return c.Edit(fmt.Sprintf("📅 *%s Notes*\n\nPlease select a year:", strings.Title(category)), menu, tele.ModeMarkdown)
}

func handleViewNotesMonths(c tele.Context, data string) error {
	c.Respond()
	parts := strings.Split(data, "_")
	if len(parts) != 2 { return c.Send("❌ Error") }
	category := parts[0]
	year, _ := strconv.Atoi(parts[1])

	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	filter := bson.M{"created_at": bson.M{"$gte": startDate, "$lt": endDate}}
	if category != "all" { filter["type"] = category }

	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$project", bson.D{{"month", bson.D{{"$month", "$created_at"}}}}}},
		{{"$group", bson.D{{"_id", "$month"}, {"count", bson.D{{"$sum", 1}}}}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}
	cursor, _ := config.Database.Collection("notes").Aggregate(context.Background(), pipeline)
	defer cursor.Close(context.Background())
	var results []bson.M
	cursor.All(context.Background(), &results)

	menu := &tele.ReplyMarkup{}
	btnViewAll := menu.Data(fmt.Sprintf("📄 View All %d", year), "vlist", data)
	var rows []tele.Row
	rows = append(rows, menu.Row(btnViewAll))

	var buttons []tele.Btn
	for i, res := range results {
		month, ok := res["_id"].(int32)
		if !ok { continue }
		count := res["count"].(int32)
		btnText := fmt.Sprintf("%s (%d)", getMonthName(int(month)), count)
		btn := menu.Data(btnText, "vdt", fmt.Sprintf("%s_%d_%d", category, year, month))
		buttons = append(buttons, btn)
		if len(buttons) == 2 || i == len(results)-1 {
			rows = append(rows, menu.Row(buttons...))
			buttons = []tele.Btn{}
		}
	}
	rows = append(rows, menu.Row(menu.Data("🔙 Back to Years", "vyr", category)))
	menu.Inline(rows...)
	return c.Edit(fmt.Sprintf("📅 *%s Notes - %d*\n\nPlease select a month:", strings.Title(category), year), menu, tele.ModeMarkdown)
}

func handleViewNotesDates(c tele.Context, data string) error {
	c.Respond()
	parts := strings.Split(data, "_")
	if len(parts) != 3 { return c.Send("❌ Error") }
	category := parts[0]
	year, _ := strconv.Atoi(parts[1])
	month, _ := strconv.Atoi(parts[2])

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)
	filter := bson.M{"created_at": bson.M{"$gte": startDate, "$lt": endDate}}
	if category != "all" { filter["type"] = category }

	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$project", bson.D{{"day", bson.D{{"$dayOfMonth", "$created_at"}}}}}},
		{{"$group", bson.D{{"_id", "$day"}, {"count", bson.D{{"$sum", 1}}}}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}
	cursor, _ := config.Database.Collection("notes").Aggregate(context.Background(), pipeline)
	defer cursor.Close(context.Background())
	var results []bson.M
	cursor.All(context.Background(), &results)

	menu := &tele.ReplyMarkup{}
	btnViewAll := menu.Data(fmt.Sprintf("📄 View All %s %d", getMonthName(month), year), "vlist", data)
	var rows []tele.Row
	rows = append(rows, menu.Row(btnViewAll))

	var buttons []tele.Btn
	for i, res := range results {
		day, ok := res["_id"].(int32)
		if !ok { continue }
		btnText := fmt.Sprintf("Day %02d", day)
		btn := menu.Data(btnText, "vlist", fmt.Sprintf("%s_%d_%d_%d", category, year, month, day))
		buttons = append(buttons, btn)
		if len(buttons) == 3 || i == len(results)-1 {
			rows = append(rows, menu.Row(buttons...))
			buttons = []tele.Btn{}
		}
	}
	rows = append(rows, menu.Row(menu.Data("🔙 Back to Months", "vmo", fmt.Sprintf("%s_%d", category, year))))
	menu.Inline(rows...)
	return c.Edit(fmt.Sprintf("📅 *%s Notes - %s %d*\n\nPlease select a date:", strings.Title(category), getMonthName(month), year), menu, tele.ModeMarkdown)
}

func handleViewNotesStaticList(c tele.Context, data string) error {
	c.Respond()
	parts := strings.Split(data, "_")
	category := parts[0]
	filter := bson.M{}
	if category != "all" { filter["type"] = category }

	var titleMsg string
	if len(parts) == 1 {
		titleMsg = fmt.Sprintf("📑 *All %s Notes*", strings.Title(category))
		if category == "all" { titleMsg = "📑 *All Notes*" }
	} else if len(parts) == 2 {
		year, _ := strconv.Atoi(parts[1])
		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		filter["created_at"] = bson.M{"$gte": startDate, "$lt": endDate}
		titleMsg = fmt.Sprintf("📑 *%s Notes - %d*", strings.Title(category), year)
	} else if len(parts) == 3 {
		year, _ := strconv.Atoi(parts[1])
		month, _ := strconv.Atoi(parts[2])
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0)
		filter["created_at"] = bson.M{"$gte": startDate, "$lt": endDate}
		titleMsg = fmt.Sprintf("📑 *%s Notes - %s %d*", strings.Title(category), getMonthName(month), year)
	} else if len(parts) == 4 {
		year, _ := strconv.Atoi(parts[1])
		month, _ := strconv.Atoi(parts[2])
		day, _ := strconv.Atoi(parts[3])
		startDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 0, 1)
		filter["created_at"] = bson.M{"$gte": startDate, "$lt": endDate}
		titleMsg = fmt.Sprintf("📑 *%s Notes - %02d %s %d*", strings.Title(category), day, getMonthName(month), year)
	}

	opts := options.Find().SetSort(bson.D{{"created_at", 1}})
	cursor, _ := config.Database.Collection("notes").Find(context.Background(), filter, opts)
	defer cursor.Close(context.Background())
	var notes []models.Notes
	cursor.All(context.Background(), &notes)

	if len(notes) == 0 {
		return c.Send("📝 No notes found.")
	}

	var sb strings.Builder
	sb.WriteString(titleMsg + "\n════════════════\n\n")

	var currentIds []primitive.ObjectID
	wib := time.FixedZone("WIB", 7*3600)

	for i, note := range notes {
		currentIds = append(currentIds, note.ID)

		icon := "📝"
		if note.Type == "journal" { icon = "📓" }
		if note.Type == "idea" { icon = "💡" }
		if note.Type == "task" { icon = "✅" }

		timeStr := note.CreatedAt.In(wib).Format("02 Jan 06 15:04 WIB")
		noteBlock := fmt.Sprintf("%d. %s *%s* — %s\n> %s\n\n", i+1, icon, note.Title, timeStr, note.Content)
		if sb.Len()+len(noteBlock) > 3800 {
			sb.WriteString("_... (truncated)_\n")
			break
		}
		sb.WriteString(noteBlock)
	}

	userNotesContext[c.Chat().ID] = currentIds

	menu := &tele.ReplyMarkup{}
	btnDel := menu.Data("🗑️ Delete", "act_del")
	btnEdit := menu.Data("✏️ Edit", "act_edit")
	
	var backBtn tele.Btn
	if len(parts) == 1 { backBtn = menu.Data("🔙 Back", "btn_view_notes", "back") }
	if len(parts) == 2 { backBtn = menu.Data("🔙 Back", "vyr", category) }
	if len(parts) == 3 { backBtn = menu.Data("🔙 Back", "vmo", fmt.Sprintf("%s_%s", category, parts[1])) }
	if len(parts) == 4 { backBtn = menu.Data("🔙 Back", "vdt", fmt.Sprintf("%s_%s_%s", category, parts[1], parts[2])) }

	menu.Inline(menu.Row(btnEdit, btnDel), menu.Row(backBtn))
	return c.Edit(sb.String(), menu, tele.ModeMarkdown)
}

func handleActionDelete(c tele.Context) error {
	c.Respond()
	userNoteState[c.Chat().ID] = "delete_prompt"
	return c.Send("🗑️ Reply with the **number** of the note you want to delete (e.g., 1, 2, 3):", tele.ModeMarkdown)
}

func handleActionEdit(c tele.Context) error {
	c.Respond()
	userNoteState[c.Chat().ID] = "edit_prompt"
	return c.Send("✏️ Reply with the **number** of the note you want to edit:", tele.ModeMarkdown)
}

func handleDeleteInput(c tele.Context) error {
	delete(userNoteState, c.Chat().ID)
	text := strings.TrimSpace(c.Message().Text)
	idx, err := strconv.Atoi(text)
	if err != nil || idx < 1 {
		return c.Send("❌ Invalid number. Please try again from the Notes menu.", notesMenu)
	}
	
	ids, ok := userNotesContext[c.Chat().ID]
	if !ok || idx > len(ids) {
		return c.Send("❌ Note number not found in current context. Please load the list again.", notesMenu)
	}
	
	targetID := ids[idx-1]
	_, err = config.Database.Collection("notes").DeleteOne(context.Background(), bson.M{"_id": targetID})
	if err != nil {
		return c.Send("❌ Failed to delete note.", notesMenu)
	}
	
	return c.Send(fmt.Sprintf("✅ Note #%d has been deleted.", idx), notesMenu)
}

func handleEditInputPhase1(c tele.Context) error {
	delete(userNoteState, c.Chat().ID)
	text := strings.TrimSpace(c.Message().Text)
	idx, err := strconv.Atoi(text)
	if err != nil || idx < 1 {
		return c.Send("❌ Invalid number.", notesMenu)
	}
	ids, ok := userNotesContext[c.Chat().ID]
	if !ok || idx > len(ids) {
		return c.Send("❌ Note number not found.", notesMenu)
	}
	
	targetID := ids[idx-1]
	userEditTarget[c.Chat().ID] = targetID
	userNoteState[c.Chat().ID] = "edit_typing"
	
	return c.Send(fmt.Sprintf("✏️ Note #%d selected.\nNow reply with the new Note text (Title on first line, Content on next lines):", idx))
}

func handleEditInputPhase2(c tele.Context) error {
	delete(userNoteState, c.Chat().ID)
	targetID, ok := userEditTarget[c.Chat().ID]
	if !ok {
		return c.Send("❌ Error: Target lost.", notesMenu)
	}
	delete(userEditTarget, c.Chat().ID)
	
	content := c.Message().Text
	lines := strings.Split(content, "\n")
	title := ""
	body := ""

	if len(lines) > 1 {
		title = strings.TrimSpace(lines[0])
		body = strings.TrimSpace(strings.Join(lines[1:], "\n"))
	} else {
		title = "Updated Note"
		body = strings.TrimSpace(content)
	}
	
	update := bson.M{
		"$set": bson.M{
			"title": title,
			"content": body,
		},
	}
	_, err := config.Database.Collection("notes").UpdateOne(context.Background(), bson.M{"_id": targetID}, update)
	if err != nil {
		return c.Send("❌ Failed to update.", notesMenu)
	}
	return c.Send("✅ Note updated successfully!", notesMenu)
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

	b.Handle(&btnViewNotes, handleViewNotesCategory)
	b.Handle("\fbtn_view_notes", handleViewNotesCategory)
	
	b.Handle("\fbtn_notes_back_main", func(c tele.Context) error {
		return handleNotesMenu(c)
	})

	b.Handle(&tele.Btn{Unique: "vyr"}, func(c tele.Context) error { return handleViewNotesYears(c, c.Callback().Data) })
	b.Handle(&tele.Btn{Unique: "vmo"}, func(c tele.Context) error { return handleViewNotesMonths(c, c.Callback().Data) })
	b.Handle(&tele.Btn{Unique: "vdt"}, func(c tele.Context) error { return handleViewNotesDates(c, c.Callback().Data) })
	b.Handle(&tele.Btn{Unique: "vlist"}, func(c tele.Context) error { return handleViewNotesStaticList(c, c.Callback().Data) })
	
	b.Handle(&tele.Btn{Unique: "act_del"}, handleActionDelete)
	b.Handle(&tele.Btn{Unique: "act_edit"}, handleActionEdit)
}

func CheckNotesText(c tele.Context) (bool, error) {
	state, ok := userNoteState[c.Chat().ID]
	if ok {
		if state == "delete_prompt" {
			return true, handleDeleteInput(c)
		} else if state == "edit_prompt" {
			return true, handleEditInputPhase1(c)
		} else if state == "edit_typing" {
			return true, handleEditInputPhase2(c)
		}
		return true, handleNoteInput(c)
	}
	return false, nil
}