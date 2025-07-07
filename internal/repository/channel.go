package repository

import (
	"gorm.io/gorm"
	model "ytcw/internal/model/api"
	"ytcw/internal/schema"
)

type ChannelRepository interface {
	FindAll(p *model.Pagination) ([]schema.Channel, int64, error)
}

type channelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepository{db: db}
}

func (r *channelRepository) FindAll(p *model.Pagination) ([]schema.Channel, int64, error) {
	var channels []schema.Channel
	var total int64

	db := r.db.Model(&schema.Channel{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Preload("Category").
		Limit(int(p.Limit())).
		Offset(int(p.Offset())).
		Find(&channels).Error

	return channels, total, err
}
