package bot

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

	tele "gopkg.in/telebot.v3"
)

var userUploadState = make(map[int64]bool)

type ProjectEndpoint struct {
	ID    string
	Label string
	Path  string
	Icon  string
}

var projectEndpoints = []ProjectEndpoint{
	{"airdrop", "Airdrop", "/allairdrop", "🪂"},
	{"cryptocommunity", "Crypto Community", "/cryptocommunity", "🪙"},
	{"aitools", "AI Tools", "/aitools", "🤖"},
	{"web3tools", "Web3 Tools", "/web3tools", "🌐"},
	{"githubrepo", "Github Repo", "/githubrepo", "🐙"},
}

func checkAuth(c tele.Context) bool {
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	if chatIDStr == "" {
		return false
	}
	expectedID, _ := strconv.ParseInt(chatIDStr, 10, 64)
	return c.Chat().ID == expectedID
}

var apiBaseURL string

func getBaseURL() (string, error) {
	if apiBaseURL != "" {
		return apiBaseURL, nil
	}
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		return "", fmt.Errorf("API_BASE_URL is not set in .env")
	}
	apiBaseURL = baseURL
	return apiBaseURL, nil
}

func handleSpeedTest(c tele.Context) error {
	c.Respond()
	if !checkAuth(c) {
		return c.Send("❌ Unauthorized access.")
	}

	c.Send("⏳ Testing API speed...")
	baseURL, err := getBaseURL()
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Configuration Error: %v", err))
	}
	endpoints := []string{"/allairdrop", "/profilelink", "/postslink", "/cryptocommunity", "/price", "/portfolio", "/aitools", "/web3tools", "/githubrepo"}

	results := "⚡ API Speed Test Results\n\n"
	allNormal := true

	for _, ep := range endpoints {
		start := time.Now()
		resp, err := http.Get(baseURL + ep)
		if err != nil {
			allNormal = false
			results += fmt.Sprintf("🔗 %s : Error (%v)\n", ep, err)
			continue
		}
		resp.Body.Close()
		duration := time.Since(start).Milliseconds()
		
		if resp.StatusCode == 200 {
			results += fmt.Sprintf("🔗 %s : %d ms\n", ep, duration)
		} else {
			allNormal = false
			results += fmt.Sprintf("🔗 %s : %d ms (Status %d)\n", ep, duration, resp.StatusCode)
		}
	}

	results += "\nStatus: "
	if allNormal {
		results += "All endpoints responded normally."
	} else {
		results += "Some endpoints experienced issues."
	}

	return c.Send(results)
}

func handleCheckMissingImages(c tele.Context) error {
	c.Respond()
	if !checkAuth(c) {
		return c.Send("❌ Unauthorized access.")
	}

	menu := &tele.ReplyMarkup{}
	var rows []tele.Row

	for _, ep := range projectEndpoints {
		if ep.ID == "githubrepo" { continue }
		btn := menu.Data(ep.Icon+" "+ep.Label, "exe_img_chk", ep.ID)
		if len(rows) > 0 && len(rows[len(rows)-1]) < 2 {
			rows[len(rows)-1] = append(rows[len(rows)-1], btn)
		} else {
			rows = append(rows, menu.Row(btn))
		}
	}

	btnCancel := menu.Data("❌ Cancel", "btn_cancel_chk")
	rows = append(rows, menu.Row(btnCancel))
	
	menu.Inline(rows...)
	return c.Edit("🔍 *Check Missing Images*\nPilih project mana yang mau di-scan:", menu, tele.ModeMarkdown)
}

func handleExecuteImageCheck(c tele.Context) error {
	c.Respond()
	if !checkAuth(c) {
		return c.Send("❌ Unauthorized access.")
	}

	projectID := c.Callback().Data
	var selectedEp *ProjectEndpoint
	for _, ep := range projectEndpoints {
		if ep.ID == projectID {
			selectedEp = &ep
			break
		}
	}

	if selectedEp == nil {
		return c.Edit("❌ Invalid project selected.")
	}

	c.Send(fmt.Sprintf("🔍 Checking missing images for *%s*...", selectedEp.Label), tele.ModeMarkdown)

	baseURL, err := getBaseURL()
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Configuration Error: %v", err))
	}
	urlStr := baseURL + selectedEp.Path

	go func() {
		totalMissing := 0
		var detailsBlocks []string

		resp, err := http.Get(urlStr)
		if err != nil {
			log.Printf("Error fetching %s: %v\n", urlStr, err)
			c.Send(fmt.Sprintf("❌ Error fetching %s: %v", selectedEp.Label, err))
			return
		}

		var data struct {
			Data []struct {
				Name          string `json:"name"`
				Image         string `json:"image"`
				Logo          string `json:"logo"`
				ImageURL      string `json:"image_url"`
				ImgURL        string `json:"img_url"`
				ImageURICamel string `json:"imageUrl"`
				ImgURICamel   string `json:"imgURL"`
			} `json:"data"`
		}

		err = json.NewDecoder(resp.Body).Decode(&data)
		resp.Body.Close()
		if err != nil {
			log.Printf("Error decoding JSON from %s: %v\n", urlStr, err)
			c.Send(fmt.Sprintf("❌ Error decoding JSON from %s", selectedEp.Label))
			return
		}

		var blockDetails string
		for _, item := range data.Data {
			imgURL := item.Image
			if imgURL == "" { imgURL = item.Logo }
			if imgURL == "" { imgURL = item.ImageURL }
			if imgURL == "" { imgURL = item.ImgURL }
			if imgURL == "" { imgURL = item.ImageURICamel }
			if imgURL == "" { imgURL = item.ImgURICamel }
			if imgURL == "" { continue }

			req, err := http.NewRequest("GET", imgURL, nil)
			if err != nil { continue }
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
			
			client := &http.Client{Timeout: 10 * time.Second}
			imgResp, err := client.Do(req)

			if err != nil || imgResp.StatusCode != 200 {
				totalMissing++
				blockDetails += fmt.Sprintf("- Name: \"%s\"\n", item.Name)
			}
			if imgResp != nil {
				imgResp.Body.Close()
			}
		}

		if blockDetails != "" {
			detailsBlocks = append(detailsBlocks, fmt.Sprintf("[%s]\n%s", selectedEp.Label, blockDetails))
		}

		msg := fmt.Sprintf("🔍 Image Check Complete for *%s*!\n\nTotal Broken Images: %d\n", selectedEp.Label, totalMissing)
		if len(detailsBlocks) > 0 {
			msg += "\nDetails:\n" + strings.Join(detailsBlocks, "\n")
		} else {
			msg += "\nDetails: All images are safe!"
		}

		if len(msg) > 4000 { msg = msg[:4000] + "\n... (truncated)" }
		c.Send(msg)
	}()

	return nil
}

func handleCheckInvalidLink(c tele.Context) error {
	c.Respond()
	if !checkAuth(c) {
		return c.Send("❌ Unauthorized access.")
	}

	menu := &tele.ReplyMarkup{}
	var rows []tele.Row

	for _, ep := range projectEndpoints {
		btn := menu.Data(ep.Icon+" "+ep.Label, "exe_lnk_chk", ep.ID)
		if len(rows) > 0 && len(rows[len(rows)-1]) < 2 {
			rows[len(rows)-1] = append(rows[len(rows)-1], btn)
		} else {
			rows = append(rows, menu.Row(btn))
		}
	}

	btnCancel := menu.Data("❌ Cancel", "btn_cancel_chk")
	rows = append(rows, menu.Row(btnCancel))
	
	menu.Inline(rows...)
	return c.Edit("🔗 *Check Invalid Links*\nPilih project mana yang mau di-scan:", menu, tele.ModeMarkdown)
}

func handleExecuteLinkCheck(c tele.Context) error {
	c.Respond()
	if !checkAuth(c) {
		return c.Send("❌ Unauthorized access.")
	}

	projectID := c.Callback().Data
	var selectedEp *ProjectEndpoint
	for _, ep := range projectEndpoints {
		if ep.ID == projectID {
			selectedEp = &ep
			break
		}
	}

	if selectedEp == nil {
		return c.Edit("❌ Invalid project selected.")
	}

	c.Send(fmt.Sprintf("🔗 Checking invalid links for *%s*...", selectedEp.Label), tele.ModeMarkdown)

	baseURL, err := getBaseURL()
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Configuration Error: %v", err))
	}
	urlStr := baseURL + selectedEp.Path

	go func() {
		totalInvalid := 0
		var detailsBlocks []string

		resp, err := http.Get(urlStr)
		if err != nil {
			log.Printf("Error fetching %s: %v\n", urlStr, err)
			c.Send(fmt.Sprintf("❌ Error fetching %s: %v", selectedEp.Label, err))
			return
		}

		var data struct {
			Data []struct {
				Name         string `json:"name"`
				Link         string `json:"link"`
				LinkURL      string `json:"link_url"`
				Website      string `json:"website"`
				LinkURLCamel string `json:"linkURL"`
				RepoURL      string `json:"repo_url"`
				Twitter      string `json:"twitter"`
				Discord      string `json:"discord"`
				Telegram     string `json:"telegram"`
				LinkTwitter  string `json:"link_twitter"`
				LinkDiscord  string `json:"link_discord"`
				LinkTelegram string `json:"link_telegram"`
				LinkClaim    string `json:"link_claim"`
				VideoURL     string `json:"video_url"`
				Instagram    string `json:"instagram"`
				Youtube      string `json:"youtube"`
			} `json:"data"`
		}

		err = json.NewDecoder(resp.Body).Decode(&data)
		resp.Body.Close()
		if err != nil {
			log.Printf("Error decoding JSON from %s: %v\n", urlStr, err)
			c.Send(fmt.Sprintf("❌ Error decoding JSON from %s", selectedEp.Label))
			return
		}

		var blockDetails string
		for _, item := range data.Data {
			primaryLink := item.Link
			if primaryLink == "" { primaryLink = item.LinkURL }
			if primaryLink == "" { primaryLink = item.Website }
			if primaryLink == "" { primaryLink = item.LinkURLCamel }

			linksToCheck := []struct {
				Name string
				URL  string
			}{
				{"Primary", primaryLink},
				{"Repo", item.RepoURL},
				{"Twitter", item.Twitter},
				{"Discord", item.Discord},
				{"Telegram", item.Telegram},
				{"Twitter", item.LinkTwitter},
				{"Discord", item.LinkDiscord},
				{"Telegram", item.LinkTelegram},
				{"Claim", item.LinkClaim},
				{"Video", item.VideoURL},
				{"Instagram", item.Instagram},
				{"Youtube", item.Youtube},
			}

			for _, l := range linksToCheck {
				link := l.URL
				if link == "" { continue }

				parsedURL, err := url.ParseRequestURI(link)
				if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
					totalInvalid++
					blockDetails += fmt.Sprintf("- Name: \"%s\" (Invalid %s format: %s)\n", item.Name, l.Name, link)
					continue
				}

				req, err := http.NewRequest("GET", link, nil)
				if err != nil { continue }
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
				
				client := &http.Client{Timeout: 10 * time.Second}
				linkResp, err := client.Do(req)

				if err != nil || linkResp.StatusCode >= 400 {
					totalInvalid++
					blockDetails += fmt.Sprintf("- Name: \"%s\" (%s Link: %s)\n", item.Name, l.Name, link)
				}
				if linkResp != nil {
					linkResp.Body.Close()
				}
			}
		}

		if blockDetails != "" {
			detailsBlocks = append(detailsBlocks, fmt.Sprintf("[%s]\n%s", selectedEp.Label, blockDetails))
		}

		msg := fmt.Sprintf("🔗 Invalid Link Check Complete for *%s*!\n\nTotal Invalid Links: %d\n", selectedEp.Label, totalInvalid)
		if len(detailsBlocks) > 0 {
			msg += "\nDetails:\n" + strings.Join(detailsBlocks, "\n")
		} else {
			msg += "\nDetails: All links are valid!"
		}

		if len(msg) > 4000 { msg = msg[:4000] + "\n... (truncated)" }
		c.Send(msg)
	}()

	return nil
}

func handleCDNInit(c tele.Context) error {
	c.Respond()
	if !checkAuth(c) {
		return c.Send("❌ Unauthorized access.")
	}

	userUploadState[c.Chat().ID] = true
	return c.Send("🖼️ GitHub CDN Upload\n\nPlease send me the photo you want to upload. (It will be uploaded to your configured GitHub repo).")
}

func handlePhotoUpload(c tele.Context) error {
	if !checkAuth(c) {
		return nil // Ignore silently if unauthorized
	}

	if !userUploadState[c.Chat().ID] {
		return nil // Not in upload state
	}

	// Reset state
	userUploadState[c.Chat().ID] = false
	c.Send("⏳ Processing and uploading photo to GitHub...")

	photo := c.Message().Photo
	if photo == nil {
		return c.Send("❌ No photo found in the message.")
	}

	file, err := TelegramBot.FileByID(photo.FileID)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Failed to get photo: %v", err))
	}

	rc, err := TelegramBot.File(&file)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Failed to download photo: %v", err))
	}
	defer rc.Close()

	buf, err := io.ReadAll(rc)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Failed to read photo data: %v", err))
	}

	// GitHub upload variables
	token := os.Getenv("GITHUB_TOKEN")
	repoOwner := os.Getenv("GITHUB_USERNAME")
	repoName := os.Getenv("GITHUB_REPO")
	uploadDir := os.Getenv("GITHUB_UPLOAD_DIR")

	if token == "" || repoOwner == "" || repoName == "" {
		return c.Send("❌ GitHub CDN configuration is incomplete in .env.")
	}

	if uploadDir == "" {
		uploadDir = "images"
	}

	now := time.Now()
	folderPath := fmt.Sprintf("%s/%d", uploadDir, now.Year())
	filename := fmt.Sprintf("%d_upload.jpg", now.Unix())
	path := fmt.Sprintf("%s/%s", folderPath, filename)
	
	urlStr := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", repoOwner, repoName, path)
	
	payload := map[string]interface{}{
		"message": "Upload via Nww Telegram Bot",
		"content": base64.StdEncoding.EncodeToString(buf),
	}
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("PUT", urlStr, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Request creation failed: %v", err))
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Upload failed: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		
		var returnedPath = path
		var sha = ""
		if content, ok := result["content"].(map[string]interface{}); ok {
			if pathVal, exists := content["path"].(string); exists {
				returnedPath = pathVal
			}
			if shaVal, exists := content["sha"].(string); exists {
				sha = shaVal
			}
		}

		parts := strings.Split(returnedPath, "/")
		for i, p := range parts {
			parts[i] = url.PathEscape(p)
		}
		escapedPath := strings.Join(parts, "/")

		finalURL := fmt.Sprintf(
			"https://%s.github.io/%s/%s",
			repoOwner,
			repoName,
			escapedPath,
		)

		img := models.Image{
			ID:       primitive.NewObjectID(),
			Filename: filename,
			URL:      finalURL,
			Size:     int64(len(buf)),
			Sha:      sha,
			Path:     returnedPath,
		}

		_, err := config.Database.Collection("images").InsertOne(context.Background(), img)
		if err != nil {
			return c.Send(fmt.Sprintf("❌ Upload successful to GitHub, but failed to save to database: %v", err))
		}

		return c.Send(fmt.Sprintf("✅ Upload Successful!\n\nURL: %s", finalURL))
	}

	body, _ := io.ReadAll(resp.Body)
	return c.Send(fmt.Sprintf("❌ GitHub API Error: %d\n%s", resp.StatusCode, string(body)))
}