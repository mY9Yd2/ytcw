package fetcher

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mY9Yd2/ytcw/internal/config"
	model "github.com/mY9Yd2/ytcw/internal/model/fetcher"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type fetchOptions struct {
	checkCutoff  bool
	addThumbnail bool
}

type ytdlp struct {
	baseArgs     []string
	url          string
	fetchOptions fetchOptions
}

func newYtDlp() *ytdlp {
	return &ytdlp{
		baseArgs: []string{
			"--no-simulate",
			"--no-download",
			"-j",
		},
	}
}

func (y *ytdlp) WithURL(url string) *ytdlp {
	y.url = url
	return y
}

func (y *ytdlp) AddArg(arg string) *ytdlp {
	y.baseArgs = append(y.baseArgs, arg)
	return y
}

func (y *ytdlp) Cmd() *exec.Cmd {
	args := append(y.baseArgs, y.url)
	return exec.Command("yt-dlp", args...)
}

func (y *ytdlp) channelBaseURL(channel string) string {
	if strings.HasPrefix(channel, "@") {
		return fmt.Sprintf("https://www.youtube.com/%s", channel)
	}
	return fmt.Sprintf("https://www.youtube.com/channel/%s", channel)
}

func (y *ytdlp) SetFetchOpts(options fetchOptions) *ytdlp {
	y.fetchOptions = options
	return y
}

func (y *ytdlp) SetChannelURL(channel string) *ytdlp {
	y.WithURL(y.channelBaseURL(channel))
	return y
}

func (y *ytdlp) SetShortsURL(channel string) *ytdlp {
	y.WithURL(y.channelBaseURL(channel) + "/shorts")
	return y
}

func (y *ytdlp) SetVideosURL(channel string) *ytdlp {
	y.WithURL(y.channelBaseURL(channel) + "/videos")
	return y
}

func (y *ytdlp) fetch(logger zerolog.Logger, out chan<- model.VideoInfo, stop chan struct{}) error {
	cmd, stdout, err := y.startYtDLPCommand()
	if err != nil {
		return err
	}

	defer func() {
		_ = cmd.Wait()
	}()

	reader := bufio.NewReader(stdout)
	cutoff := time.Now().UTC().Add(-config.GetConfig().Fetcher.MaxVideoAge)

	return y.streamVideos(logger, cmd, reader, cutoff, out, stop)
}

func (y *ytdlp) streamVideos(
	logger zerolog.Logger,
	cmd *exec.Cmd,
	reader *bufio.Reader,
	cutoff time.Time,
	out chan<- model.VideoInfo,
	stop chan struct{},
) error {
	for {
		select {
		case <-stop:
			_ = cmd.Process.Kill()
			return nil
		default:
		}

		line, err := y.readNextLine(reader)
		if err != nil {
			break
		}

		info, err := y.parseVideoInfo(line)
		if err != nil {
			logger.Error().Err(err).Msg("Error parsing JSON")
			continue
		}

		if y.fetchOptions.checkCutoff && y.isOlderThan(info.Timestamp, cutoff) {
			close(stop)
			_ = cmd.Process.Kill()
			return nil
		}

		if y.fetchOptions.addThumbnail {
			info.Thumbnail, err = y.findAvailableThumbnail(info.DisplayID)
			if err != nil {
				logger.Error().Err(err).Msg("Error finding thumbnail")
			}
		}

		out <- info
	}
	return nil
}

func (y *ytdlp) startYtDLPCommand() (*exec.Cmd, io.ReadCloser, error) {
	cmd := y.Cmd()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	return cmd, stdout, nil
}

func (y *ytdlp) readNextLine(reader *bufio.Reader) ([]byte, error) {
	return reader.ReadBytes('\n')
}

func (y *ytdlp) parseVideoInfo(line []byte) (model.VideoInfo, error) {
	var info model.VideoInfo
	err := json.Unmarshal(line, &info)
	return info, err
}

func (y *ytdlp) isOlderThan(timestamp int64, cutoff time.Time) bool {
	return time.Unix(timestamp, 0).Before(cutoff)
}

func (y *ytdlp) findAvailableThumbnail(displayId string) (string, error) {
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
