package fetcher

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os/exec"
	"time"
	"ytcw/internal/config"
	"ytcw/internal/model"
)

func fetchFromURL(url string, out chan<- model.VideoInfo, stop chan struct{}) error {
	cmd, stdout, err := startYtDLPCommand(url)
	if err != nil {
		return err
	}
	defer cmd.Wait()

	reader := bufio.NewReader(stdout)
	cutoff := time.Now().Add(-config.LoadConfig().Ytcwd.MaxVideoAge)

	for {
		select {
		case <-stop:
			_ = cmd.Process.Kill()
			return nil
		default:
		}

		line, err := readNextLine(reader)
		if err != nil {
			break
		}

		info, err := parseVideoInfo(line)
		if err != nil {
			log.Error().Err(err).Msg("Error parsing JSON")
			continue
		}

		if isOlderThan(info, cutoff) {
			close(stop)
			_ = cmd.Process.Kill()
			return nil
		}

		info.Thumbnail, err = findAvailableThumbnail(info.DisplayID)
		if err != nil {
			log.Error().Err(err).Msg("Error finding thumbnail")
		}

		out <- info
	}

	return nil
}

func FetchVideos(channel string) <-chan model.VideoInfo {
	out := make(chan model.VideoInfo)

	go func() {
		defer close(out)

		stop := make(chan struct{})

		shortsURL := fmt.Sprintf("https://www.youtube.com/channel/%s/shorts", channel)
		if err := fetchFromURL(shortsURL, out, stop); err != nil {
			log.Error().Err(err).Msg("Error fetching shorts")
		}

		stop = make(chan struct{}) // Reset

		videosURL := fmt.Sprintf("https://www.youtube.com/channel/%s/videos", channel)
		if err := fetchFromURL(videosURL, out, stop); err != nil {
			log.Error().Err(err).Msg("Error fetching videos")
		}
	}()

	return out
}

func startYtDLPCommand(url string) (*exec.Cmd, io.ReadCloser, error) {
	cmd := exec.Command("yt-dlp",
		"--no-simulate",
		"--no-download", "-j",
		url,
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	return cmd, stdout, nil
}

func readNextLine(reader *bufio.Reader) ([]byte, error) {
	return reader.ReadBytes('\n')
}

func parseVideoInfo(line []byte) (model.VideoInfo, error) {
	var info model.VideoInfo
	err := json.Unmarshal(line, &info)
	return info, err
}

func isOlderThan(info model.VideoInfo, cutoff time.Time) bool {
	return time.Unix(info.Timestamp, 0).Before(cutoff)
}

func findAvailableThumbnail(displayId string) (string, error) {
	fileNames := []string{"sddefault.webp", "hqdefault.webp", "0.webp"}
	baseURL := fmt.Sprintf("https://i.ytimg.com/vi_webp/%s/", displayId)
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	for _, fileName := range fileNames {
		thumbnailUrl := baseURL + fileName

		resp, err := client.Head(thumbnailUrl)
		if err != nil {
			continue
		}

		err = resp.Body.Close()
		if err != nil {
			return "", err
		}

		if resp.StatusCode == 200 {
			return fileName, nil
		}
	}

	return "", fmt.Errorf("no available thumbnail found")
}
