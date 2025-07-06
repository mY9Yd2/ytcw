package repository

import (
	"gorm.io/gorm"
	"ytcw/internal/schema"
)

type ChannelRepository interface {
	FindAll() ([]schema.Channel, error)
}

type channelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepository{db: db}
}

func (r *channelRepository) FindAll() ([]schema.Channel, error) {
	var channels []schema.Channel
	return channels, r.db.Preload("Category").Find(&channels).Error
}
