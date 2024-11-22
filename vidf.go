package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

const (
	youtubeAPIURL = "https://www.googleapis.com/youtube/v3/search"
)

type YouTubeResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a search query as a command-line argument.")
		os.Exit(1)
	}

	searchQuery := strings.Join(os.Args[1:], " ")

	// Load API key from .env file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	envPath := filepath.Join(homeDir, ".config", "fabric", ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		fmt.Println("YOUTUBE_API_KEY not found in .env file")
		os.Exit(1)
	}

	// Construct the API request URL
	params := url.Values{}
	params.Add("part", "snippet")
	params.Add("q", searchQuery)
	params.Add("type", "video")
	params.Add("key", apiKey)

	fullURL := youtubeAPIURL + "?" + params.Encode()

	// Make the API request
	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		os.Exit(1)
	}

	// Parse the JSON response
	var result YouTubeResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	// Print the URL of the first video
	if len(result.Items) > 0 {
		// Get the title of the video
		title := result.Items[0].Snippet.Title

		// Replace entities (&#39;) with characters
		replacer := strings.NewReplacer("&#39;", "'", "&quot;", "\"", "&amp;", "&")
		title = replacer.Replace(title)
		videoID := result.Items[0].ID.VideoID
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s&title=%s", videoID, title)
		fmt.Println(videoURL)
	} else {
		fmt.Println("No videos found.")
	}
}
