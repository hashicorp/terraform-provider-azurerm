package loader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// PRFile represents a file in a GitHub PR
type PRFile struct {
	Filename  string `json:"filename"`
	Status    string `json:"status"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Changes   int    `json:"changes"`
	Patch     string `json:"patch"`
}

// PRInfo represents GitHub PR information
type PRInfo struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Base   struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"base"`
	Head struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
}

// GitHubLoader loads changes from GitHub API
type GitHubLoader struct {
	prNumber int
}

// Load loads changes from GitHub API and returns a ChangeSet
func (l *GitHubLoader) Load() (*ChangeSet, error) {
	cs := NewChangeSet()

	token := os.Getenv("GITHUB_TOKEN")
	owner, name := getRepoInfo()

	prNum := l.prNumber

	log.Printf("Fetching PR #%d changes from GitHub API (%s/%s)...", prNum, owner, name)

	files, err := fetchPRFiles(token, owner, name, prNum)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PR files: %w", err)
	}

	for _, file := range files {
		if !isServiceFile(file.Filename) {
			continue
		}

		// Normalize the file path for consistent tracking
		normalizedPath := normalizeFilePath(file.Filename)

		// Check if this is a new file
		isNewFile := file.Status == "added"

		if file.Patch != "" {
			if err := cs.parsePatch(normalizedPath, file.Patch); err != nil {
				log.Printf("Warning: failed to parse patch for %s: %v", file.Filename, err)
			}
		}

		cs.changedFiles[normalizedPath] = true
		if isNewFile {
			cs.newFiles[normalizedPath] = true
		}
	}

	log.Printf("✓ Found %d changed files from GitHub API", len(cs.changedFiles))
	return cs, nil
}

// fetchPRFiles fetches the list of changed files from GitHub API
func fetchPRFiles(token, owner, name string, prNum int) ([]PRFile, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/files", owner, name, prNum)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Warning: failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("GitHub API returned status %d, failed to read body: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var files []PRFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	return files, nil
}

// getRepoInfo gets the repository owner and name
func getRepoInfo() (owner, name string) {
	return "hashicorp", "terraform-provider-azurerm"
}

// fetchPRInfo fetches PR information from GitHub API
func fetchPRInfo(token, owner, name string, prNum int) (*PRInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, name, prNum)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Warning: failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("GitHub API returned status %d, failed to read body: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var prInfo PRInfo
	if err := json.NewDecoder(resp.Body).Decode(&prInfo); err != nil {
		return nil, err
	}

	return &prInfo, nil
}
