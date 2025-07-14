package service

import (
	model "github.com/mY9Yd2/ytcw/internal/model/api"
	"github.com/mY9Yd2/ytcw/internal/repository"
)

type VideoService interface {
	GetVideos(p *model.Pagination) ([]model.VideoResponse, *model.Pagination, error)
}

type videoService struct {
	videoRepo repository.VideoRepository
}

func NewVideoService(repo repository.VideoRepository) VideoService {
	return &videoService{
		videoRepo: repo,
	}
}

func (s *videoService) GetVideos(p *model.Pagination) ([]model.VideoResponse, *model.Pagination, error) {
	videos, total, err := s.videoRepo.FindAll(p)
	if err != nil {
		return nil, nil, err
	}

	var responses []model.VideoResponse
	for _, video := range videos {
		var category *model.CategoryResponse
		if video.Channel.Category != nil {
			category = &model.CategoryResponse{
				ID:   video.Channel.Category.ID,
				Name: video.Channel.Category.Name,
			}
		}

		channel := model.ChannelSummary{
			ID:       video.Channel.ID,
			Channel:  video.Channel.Channel,
			Category: category,
		}

		responses = append(responses, model.VideoResponse{
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
		responses = []model.VideoResponse{}
	}

	return responses, p, nil
}
