package content

import "github.com/mY9Yd2/ytcw/internal/common"

type VideoService interface {
	GetVideos(p *common.Pagination) ([]VideoResponse, *common.Pagination, error)
}

type videoService struct {
	videoRepo VideoRepository
}

func NewVideoService(repo VideoRepository) VideoService {
	return &videoService{
		videoRepo: repo,
	}
}

func (s *videoService) GetVideos(p *common.Pagination) ([]VideoResponse, *common.Pagination, error) {
	videos, total, err := s.videoRepo.FindAll(p)
	if err != nil {
		return nil, nil, err
	}

	var responses []VideoResponse
	for _, video := range videos {
		var category *CategoryResponse
		if video.Channel.Category != nil {
			category = &CategoryResponse{
				ID:   video.Channel.Category.ID,
				Name: video.Channel.Category.Name,
			}
		}

		channel := ChannelSummary{
			ID:       video.Channel.ID,
			Channel:  video.Channel.Channel,
			Category: category,
		}

		responses = append(responses, VideoResponse{
			ID:        video.ID,
			Timestamp: video.Timestamp,
			FullTitle: video.FullTitle,
			DisplayID: video.DisplayID,
			Duration:  video.Duration,
			Thumbnail: video.Thumbnail,
			Language:  video.Language,
			Channel:   channel,
			VideoType: video.VideoType,
		})
	}

	p.TotalRows = uint(total)
	p.TotalPages = (p.TotalRows + p.PageSize - 1) / p.PageSize

	if responses == nil {
		responses = []VideoResponse{}
	}

	return responses, p, nil
}
