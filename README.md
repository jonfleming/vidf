# Vidf

Vidf is a command-line tool written in Go that allows you to search for YouTube videos using the YouTube Data API. It retrieves the first video result for a given search query and prints the video's URL.

## Features

- Search for YouTube videos using a command-line query.
- Retrieves and displays the URL of the first video result.
- Utilizes the YouTube Data API for fetching video data.

## Installation

To install Vidf, you need to have Go installed on your system. You can install Vidf using the following command:

```bash
go install github.com/jonfleming/vidf@latest
```

### Dependencies

Vidf works well with the following tools:

- `fabric`: You can find `fabric` at [https://github.com/danielmiessler/fabric](https://github.com/danielmiessler/fabric). 
   ```bash
   go install github.com/danielmiessler/yt@latest
   ```

- `yt`: Install **YouTub Transcripts** using the following command:

  ```bash
  go install github.com/danielmiessler/yt@latest
  ```

## Usage

1. Ensure you have a YouTube Data API key. Save it in a `.env` file located at `~/.config/fabric/.env` with the following content:

   ```
   YOUTUBE_API_KEY=your_api_key_here
   ```

2. Run the Vidf command with a search query:

   ```bash
   vidf "your search query"
   ```

   Replace `"your search query"` with the actual search term you want to use.

3. The program will output the URL of the first video result.

### Example: Downloading and Summarizing a YouTube Transcript

You can use the following command to download a YouTube transcript and summarize it using `fabric`:

```bash
url=$(vidf "Tesla Optimus Robot Explained")
yt $url | fabric -p summarize
```

## Version

Current version: 1.3.1

## License

This project is licensed under the MIT License.
