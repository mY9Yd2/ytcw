package service

import (
	"github.com/google/uuid"
	model "github.com/mY9Yd2/ytcw/internal/model/api"
	"github.com/mY9Yd2/ytcw/internal/repository"
)

type ChannelService interface {
	GetChannels(p *model.Pagination, categoryName string) ([]model.ChannelResponse, *model.Pagination, error)
}

type channelService struct {
	channelRepo repository.ChannelRepository
}

func NewChannelService(repo repository.ChannelRepository) ChannelService {
	return &channelService{
		channelRepo: repo,
	}
}

func (s *channelService) GetChannels(p *model.Pagination, categoryName string) ([]model.ChannelResponse, *model.Pagination, error) {
	channels, total, err := s.channelRepo.FindAll(p, categoryName)
	if err != nil {
		return nil, nil, err
	}

	var responses []model.ChannelResponse
	for _, ch := range channels {
		var category *model.CategoryResponse
		if ch.Category != nil && ch.Category.ID != uuid.Nil {
			category = &model.CategoryResponse{
				ID:   ch.Category.ID,
				Name: ch.Category.Name,
			}
		}

		responses = append(responses, model.ChannelResponse{
			ID:            ch.ID,
			UploaderID:    ch.UploaderID,
			ChannelID:     ch.ChannelID,
			Channel:       ch.Channel,
			LastFetch:     ch.LastFetch,
			DisabledAt:    ch.DisabledAt,
			DisabledUntil: ch.DisabledUntil,
			Category:      category,
		})
	}

	p.TotalRows = uint(total)
	p.TotalPages = (p.TotalRows + p.PageSize - 1) / p.PageSize

	if responses == nil {
		responses = []model.ChannelResponse{}
	}

	return responses, p, nil
}
