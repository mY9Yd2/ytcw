package content

import (
	"github.com/google/uuid"
	"github.com/mY9Yd2/ytcw/internal/common"
)

type ChannelService interface {
	GetChannels(p *common.Pagination, categoryName string) ([]ChannelResponse, *common.Pagination, error)
}

type channelService struct {
	channelRepo ChannelRepository
}

func NewChannelService(repo ChannelRepository) ChannelService {
	return &channelService{
		channelRepo: repo,
	}
}

func (s *channelService) GetChannels(p *common.Pagination, categoryName string) ([]ChannelResponse, *common.Pagination, error) {
	channels, total, err := s.channelRepo.FindAll(p, categoryName)
	if err != nil {
		return nil, nil, err
	}

	var responses []ChannelResponse
	for _, ch := range channels {
		var category *CategoryResponse
		if ch.Category != nil && ch.Category.ID != uuid.Nil {
			category = &CategoryResponse{
				ID:   ch.Category.ID,
				Name: ch.Category.Name,
			}
		}

		responses = append(responses, ChannelResponse{
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
		responses = []ChannelResponse{}
	}

	return responses, p, nil
}
