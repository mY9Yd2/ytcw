package service

import (
	model "ytcw/internal/model/api"
	"ytcw/internal/repository"
)

type ChannelService interface {
	GetChannels() ([]model.ChannelResponse, error)
}

type channelService struct {
	channelRepo repository.ChannelRepository
}

func NewChannelService(repo repository.ChannelRepository) ChannelService {
	return &channelService{
		channelRepo: repo,
	}
}

func (r *channelService) GetChannels() ([]model.ChannelResponse, error) {
	channels, err := r.channelRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []model.ChannelResponse
	for _, ch := range channels {
		var category *model.CategoryResponse
		if ch.Category != nil && ch.Category.ID != 0 {
			category = &model.CategoryResponse{
				ID:   ch.Category.ID,
				Name: ch.Category.Name,
			}
		}

		responses = append(responses, model.ChannelResponse{
			ID:         ch.ID,
			UploaderID: ch.UploaderID,
			ChannelID:  ch.ChannelID,
			Channel:    ch.Channel,
			LastFetch:  ch.LastFetch,
			DisabledAt: ch.DisabledAt,
			Category:   category,
		})
	}

	if responses == nil {
		responses = []model.ChannelResponse{}
	}

	return responses, nil
}
