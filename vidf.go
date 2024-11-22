// go install github.com/jonfleming/vidf@latest
package main

import (
	"encoding/json"
	"flag"
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

const version = "1.3.3"

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
	var createTitleFile = flag.Bool("createTitleFile", false, "Set this flag to create a title file")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Vidf version", version)
		fmt.Println("Usage: vidf [options] <search query>")
		fmt.Println("Options:")
		fmt.Println("  -createTitleFile  Set this flag to create a title file")
		os.Exit(1)
	}

	searchQuery := strings.Join(flag.Args(), " ")

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

	// Print the title and URL of the first video
	if len(result.Items) > 0 {
		// Get the title of the video
		title := result.Items[0].Snippet.Title

		// Replace entities (&#39;) with characters
		replacer := strings.NewReplacer("&#39;", "'", "&amp;", "&", " ", " ")
		title = replacer.Replace(title)
		videoID := result.Items[0].ID.VideoID
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
		fmt.Println(videoURL)

		// Conditionally write the title to a text file
		if *createTitleFile {
			err = os.WriteFile("video_title.txt", []byte(title), 0644)
			if err != nil {
				fmt.Println("Error writing title to file:", err)
			}
		}
	} else {
		fmt.Println("No videos found.")
	}
}
