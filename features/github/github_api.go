package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
)

var useBackupTokenUntil int64 = 0

type githubRepoResponse struct {
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	Language        string    `json:"language"`
	PushedAt        time.Time `json:"pushed_at"`
	Owner           struct {
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
}

func doGithubRequest(url string) ([]byte, error) {
	tokens := []string{os.Getenv("GITHUB_TOKEN"), os.Getenv("GITHUB_TOKEN2")}
	var validTokens []string
	for _, t := range tokens {
		if t != "" {
			validTokens = append(validTokens, t)
		}
	}

	currentTokenIndex := 0
	if len(validTokens) > 1 && time.Now().UnixMilli() < useBackupTokenUntil {
		currentTokenIndex = 1
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if len(validTokens) > 0 {
		req.Header.Set("Authorization", "Bearer "+validTokens[currentTokenIndex])
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 || resp.StatusCode == 401 {
		if len(validTokens) > 1 {
			if currentTokenIndex == 0 {
				useBackupTokenUntil = time.Now().UnixMilli() + 5*60*60*1000 // Switch for 5 hours
				fmt.Println("Github API rate limit hit on Token 1. Switching to Token 2 for 5 hours.")
			} else {
				useBackupTokenUntil = 0
				fmt.Println("Github API rate limit hit on Token 2. Reverting to Token 1.")
			}
			
			req.Header.Set("Authorization", "Bearer "+validTokens[(currentTokenIndex+1)%2])
			resp2, err2 := client.Do(req)
			if err2 != nil {
				return nil, err2
			}
			defer resp2.Body.Close()
			resp = resp2
		}
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github api returned status %d for url %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func fixMarkdownLinks(content, repoPath, owner, repoName, defaultBranch string) string {
	dir := path.Dir(repoPath)

	imgSrcRegex := regexp.MustCompile(`(?i)(src=["'])([^"']+)["']`)
	content = imgSrcRegex.ReplaceAllStringFunc(content, func(m string) string {
		matches := imgSrcRegex.FindStringSubmatch(m)
		if len(matches) == 3 {
			urlStr := matches[2]
			if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") && !strings.HasPrefix(urlStr, "mailto:") && !strings.HasPrefix(urlStr, "data:") {
				resolvedPath := path.Join(dir, urlStr)
				newUrl := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repoName, defaultBranch, resolvedPath)
				return matches[1] + newUrl + "\""
			}
		}
		return m
	})

	mdImgRegex := regexp.MustCompile(`(!\[[^\]]*\])\(([^)]+)\)`)
	content = mdImgRegex.ReplaceAllStringFunc(content, func(m string) string {
		matches := mdImgRegex.FindStringSubmatch(m)
		if len(matches) == 3 {
			urlStr := matches[2]
			if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") && !strings.HasPrefix(urlStr, "data:") {
				resolvedPath := path.Join(dir, urlStr)
				newUrl := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repoName, defaultBranch, resolvedPath)
				return matches[1] + "(" + newUrl + ")"
			}
		}
		return m
	})

	aHrefRegex := regexp.MustCompile(`(?i)(href=["'])([^"']+)["']`)
	content = aHrefRegex.ReplaceAllStringFunc(content, func(m string) string {
		matches := aHrefRegex.FindStringSubmatch(m)
		if len(matches) == 3 {
			urlStr := matches[2]
			if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") && !strings.HasPrefix(urlStr, "mailto:") && !strings.HasPrefix(urlStr, "#") {
				resolvedPath := path.Join(dir, urlStr)
				newUrl := fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s", owner, repoName, defaultBranch, resolvedPath)
				return matches[1] + newUrl + "\""
			}
		}
		return m
	})

	mdLinkRegex := regexp.MustCompile(`(\[[^\]]*\])\(([^)]+)\)`)
	content = mdLinkRegex.ReplaceAllStringFunc(content, func(m string) string {
		matches := mdLinkRegex.FindStringSubmatch(m)
		if len(matches) == 3 {
			urlStr := matches[2]
			if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") && !strings.HasPrefix(urlStr, "#") && !strings.HasPrefix(urlStr, "mailto:") {
				resolvedPath := path.Join(dir, urlStr)
				newUrl := fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s", owner, repoName, defaultBranch, resolvedPath)
				return matches[1] + "(" + newUrl + ")"
			}
		}
		return m
	})

	return content
}

func FetchGithubRepoStats(owner, repoName string) (*GithubStats, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repoName)
	body, err := doGithubRequest(url)
	if err != nil {
		return nil, err
	}

	var data githubRepoResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	stats := &GithubStats{
		Stars:      data.StargazersCount,
		Forks:      data.ForksCount,
		Language:   data.Language,
		ImageURL:   data.Owner.AvatarURL,
		LastUpdate: data.PushedAt,
	}

	return stats, nil
}

func FetchGithubRepoDetails(owner, repoName string) (map[string]interface{}, []MdFile, error) {
	repoUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repoName)
	repoBody, err := doGithubRequest(repoUrl)
	if err != nil {
		return nil, nil, err
	}
	
	var repoData map[string]interface{}
	if err := json.Unmarshal(repoBody, &repoData); err != nil {
		return nil, nil, err
	}

	defaultBranch := "main"
	if branch, ok := repoData["default_branch"].(string); ok {
		defaultBranch = branch
	}

	treeUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/%s?recursive=1", owner, repoName, defaultBranch)
	treeBody, err := doGithubRequest(treeUrl)
	if err != nil {
		return nil, nil, err
	}

	type TreeItem struct {
		Path string `json:"path"`
		Type string `json:"type"`
	}

	var treeResponse struct {
		Tree []TreeItem `json:"tree"`
	}

	if err := json.Unmarshal(treeBody, &treeResponse); err != nil {
		return nil, nil, err
	}

	type ContentItem struct {
		Name        string
		Type        string
		DownloadURL string
		Path        string
	}

	var filesToDownload []ContentItem
	for _, item := range treeResponse.Tree {
		if item.Type == "blob" {
			parts := strings.Split(item.Path, "/")
			name := parts[len(parts)-1]
			lowerName := strings.ToLower(name)
			
			if strings.HasSuffix(lowerName, ".md") || strings.HasSuffix(lowerName, ".mdx") || lowerName == "license" || lowerName == "code_of_conduct" {
				downloadURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repoName, defaultBranch, item.Path)
				filesToDownload = append(filesToDownload, ContentItem{
					Name:        name,
					Type:        "file",
					DownloadURL: downloadURL,
					Path:        item.Path,
				})
			}
		}
	}

	var downloadedMdFiles []MdFile
	var mu sync.Mutex

	var downloadWg sync.WaitGroup
	for _, file := range filesToDownload {
		downloadWg.Add(1)
		go func(f ContentItem) {
			defer downloadWg.Done()
			
			req, err := http.NewRequest("GET", f.DownloadURL, nil)
			if err != nil {
				return
			}
			
			tokens := []string{os.Getenv("GITHUB_TOKEN"), os.Getenv("GITHUB_TOKEN2")}
			var validTokens []string
			for _, t := range tokens {
				if t != "" {
					validTokens = append(validTokens, t)
				}
			}
			
			currentTokenIndex := 0
			if len(validTokens) > 0 {
				if len(validTokens) > 1 && time.Now().UnixMilli() < useBackupTokenUntil {
					currentTokenIndex = 1
				}
				req.Header.Set("Authorization", "Bearer "+validTokens[currentTokenIndex])
			}
			
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				return
			}
			
			if (resp.StatusCode == 403 || resp.StatusCode == 401) && len(validTokens) > 1 {
				resp.Body.Close()
				
				if currentTokenIndex == 0 {
					useBackupTokenUntil = time.Now().UnixMilli() + 5*60*60*1000
					fmt.Println("Raw download rate limit hit on Token 1. Switching to Token 2 for 5 hours.")
				} else {
					useBackupTokenUntil = 0
					fmt.Println("Raw download rate limit hit on Token 2. Reverting to Token 1.")
				}
				
				req.Header.Set("Authorization", "Bearer "+validTokens[(currentTokenIndex+1)%2])
				resp2, err2 := client.Do(req)
				if err2 != nil {
					return
				}
				resp = resp2
			}
			defer resp.Body.Close()
			
			if resp.StatusCode == 200 {
				contentBytes, err := io.ReadAll(resp.Body)
				if err == nil {
					finalContent := fixMarkdownLinks(string(contentBytes), f.Path, owner, repoName, defaultBranch)
					mu.Lock()
					downloadedMdFiles = append(downloadedMdFiles, MdFile{
						Name:    f.Name,
						Content: finalContent,
					})
					mu.Unlock()
				}
			}
		}(file)
	}

	downloadWg.Wait()

	uniqueFilesMap := make(map[string]MdFile)
	for _, f := range downloadedMdFiles {
		lowerName := strings.ToLower(f.Name)
		if _, exists := uniqueFilesMap[lowerName]; !exists {
			uniqueFilesMap[lowerName] = f
		}
	}

	var finalMdFiles []MdFile
	for _, f := range uniqueFilesMap {
		finalMdFiles = append(finalMdFiles, f)
	}

	getPriority := func(name string) int {
		lower := strings.ToLower(name)
		if strings.HasPrefix(lower, "readme") {
			return 1
		}
		if strings.HasPrefix(lower, "code_of_conduct") {
			return 2
		}
		if strings.HasPrefix(lower, "contributing") {
			return 3
		}
		if strings.HasPrefix(lower, "license") {
			return 4
		}
		return 5
	}

	for i := 0; i < len(finalMdFiles)-1; i++ {
		for j := i + 1; j < len(finalMdFiles); j++ {
			pI := getPriority(finalMdFiles[i].Name)
			pJ := getPriority(finalMdFiles[j].Name)
			
			swap := false
			if pI != pJ {
				swap = pI > pJ
			} else {
				swap = strings.Compare(finalMdFiles[i].Name, finalMdFiles[j].Name) > 0
			}
			
			if swap {
				finalMdFiles[i], finalMdFiles[j] = finalMdFiles[j], finalMdFiles[i]
			}
		}
	}

	return repoData, finalMdFiles, nil
}

func FetchGithubCommits(owner, repoName, perPage string) ([]interface{}, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?per_page=%s", owner, repoName, perPage)
	body, err := doGithubRequest(url)
	if err != nil {
		return nil, err
	}

	var data []interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}