package fetcher

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"ytcw/internal/config"
	"ytcw/internal/model"
)

type FetchOptions struct {
	CheckCutoff  bool
	AddThumbnail bool
}

type Ytdlp struct {
	BaseArgs []string
	URL      string
}

func NewYtDlp() *Ytdlp {
	return &Ytdlp{
		BaseArgs: []string{
			"--no-simulate",
			"--no-download",
			"-j",
		},
	}
}

func (y *Ytdlp) WithURL(url string) *Ytdlp {
	y.URL = url
	return y
}

func (y *Ytdlp) AddArg(arg string) *Ytdlp {
	y.BaseArgs = append(y.BaseArgs, arg)
	return y
}

func (y *Ytdlp) Cmd() *exec.Cmd {
	args := append(y.BaseArgs, y.URL)
	return exec.Command("yt-dlp", args...)
}

func (y *Ytdlp) channelBaseURL(channel string) string {
	if strings.HasPrefix(channel, "@") {
		return fmt.Sprintf("https://www.youtube.com/%s", channel)
	}
	return fmt.Sprintf("https://www.youtube.com/channel/%s", channel)
}

func (y *Ytdlp) Channel(channel string) *Ytdlp {
	y.WithURL(y.channelBaseURL(channel))
	return y
}

func (y *Ytdlp) Shorts(channel string) *Ytdlp {
	y.WithURL(y.channelBaseURL(channel) + "/shorts")
	return y
}

func (y *Ytdlp) Videos(channel string) *Ytdlp {
	y.WithURL(y.channelBaseURL(channel) + "/videos")
	return y
}

func fetch(ytdlp Ytdlp, out chan<- model.VideoInfo, stop chan struct{}, opts FetchOptions) error {
	cmd, stdout, err := startYtDLPCommand(ytdlp)
	if err != nil {
		return err
	}
	defer cmd.Wait()

	reader := bufio.NewReader(stdout)
	cutoff := time.Now().UTC().Add(-config.LoadConfig().Ytcwd.MaxVideoAge)

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

		if opts.CheckCutoff && isOlderThan(info.Timestamp, cutoff) {
			close(stop)
			_ = cmd.Process.Kill()
			return nil
		}

		if opts.AddThumbnail {
			info.Thumbnail, err = findAvailableThumbnail(info.DisplayID)
			if err != nil {
				log.Error().Err(err).Msg("Error finding thumbnail")
			}
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
		fetchOpts := FetchOptions{
			CheckCutoff:  true,
			AddThumbnail: true,
		}

		shorts := NewYtDlp().Shorts(channel)
		if err := fetch(*shorts, out, stop, fetchOpts); err != nil {
			log.Error().Err(err).Msg("Error fetching shorts")
		}

		stop = make(chan struct{}) // Reset

		videos := NewYtDlp().Videos(channel)
		if err := fetch(*videos, out, stop, fetchOpts); err != nil {
			log.Error().Err(err).Msg("Error fetching videos")
		}
	}()

	return out
}

func startYtDLPCommand(ytdlp Ytdlp) (*exec.Cmd, io.ReadCloser, error) {
	cmd := ytdlp.Cmd()

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

func isOlderThan(timestamp int64, cutoff time.Time) bool {
	return time.Unix(timestamp, 0).Before(cutoff)
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

func GetChannelInfo(channel string) model.ChannelInfo {
	out := make(chan model.VideoInfo)
	stop := make(chan struct{})

	go func() {
		defer close(out)

		ytdlp := NewYtDlp().
			AddArg("-I:10").
			Channel(channel)

		if err := fetch(*ytdlp, out, stop, FetchOptions{
			CheckCutoff:  false,
			AddThumbnail: false,
		}); err != nil {
			log.Error().Err(err).Msg("Error fetching channel info")
		}
	}()

	info, ok := <-out
	if !ok {
		log.Fatal().Msg("No video info found")
	}
	close(stop)

	return model.ChannelInfo{
		UploaderID: info.UploaderID,
		ChannelID:  info.ChannelID,
		Channel:    info.Channel,
	}
}
