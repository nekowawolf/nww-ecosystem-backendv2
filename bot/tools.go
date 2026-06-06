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

func checkAuth(c tele.Context) bool {
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	if chatIDStr == "" {
		return false
	}
	expectedID, _ := strconv.ParseInt(chatIDStr, 10, 64)
	return c.Chat().ID == expectedID
}

func getBaseURL() (string, error) {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		return "", fmt.Errorf("API_BASE_URL is not set in .env")
	}
	return baseURL, nil
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
	endpoints := []string{"/allairdrop", "/profilelink", "/postslink", "/cryptocommunity", "/price", "/portfolio"}

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

	c.Send("🔍 Checking for missing images. This might take a while...")

	baseURL, err := getBaseURL()
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Configuration Error: %v", err))
	}

	endpoints := map[string]string{
		"Airdrop":          baseURL + "/allairdrop",
		"Crypto Community": baseURL + "/cryptocommunity",
	}

	totalMissing := 0
	details := ""

	for label, url := range endpoints {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error fetching %s: %v\n", url, err)
			continue
		}
		defer resp.Body.Close()

		var data struct {
			Data []struct {
				Name     string `json:"name"`
				Image    string `json:"image"`
				Logo     string `json:"logo"`
				ImageURL string `json:"image_url"`
				ImgURL   string `json:"img_url"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Printf("Error decoding JSON from %s: %v\n", url, err)
			continue
		}

		for _, item := range data.Data {
			imgURL := item.Image
			if imgURL == "" {
				imgURL = item.Logo
			}
			if imgURL == "" {
				imgURL = item.ImageURL
			}
			if imgURL == "" {
				imgURL = item.ImgURL
			}
			if imgURL == "" {
				continue
			}

			// Ping the image URL
			req, err := http.NewRequest("GET", imgURL, nil)
			if err != nil {
				continue
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
			
			client := &http.Client{Timeout: 10 * time.Second}
			imgResp, err := client.Do(req)

			if err != nil || imgResp.StatusCode != 200 {
				totalMissing++
				details += fmt.Sprintf("- [%s] Name: \"%s\"\n", label, item.Name)
			}
			if imgResp != nil {
				imgResp.Body.Close()
			}
		}
	}

	msg := "🔍 Image Check Complete!\n\n"
	msg += fmt.Sprintf("Total Broken Images: %d\n", totalMissing)

	if totalMissing > 0 {
		msg += "\nDetails:\n" + details
	} else {
		msg += "\nDetails: All images are safe!"
	}

	return c.Send(msg)
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