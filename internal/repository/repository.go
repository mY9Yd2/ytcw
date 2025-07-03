package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	"ytcw/internal/schema"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SaveVideo(video *schema.Video) error {
	var channel schema.Channel
	if err := r.DB.Where("uploader_id = ?", video.Channel.UploaderID).
		FirstOrCreate(&channel, &video.Channel).Error; err != nil {
		return err
	}

	video.ChannelRefer = channel.ID

	return r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "display_id"}},
		DoNothing: true,
	}).Create(video).Error
}

func (r *Repository) GetStaleChannel(d time.Duration) (*schema.Channel, error) {
	var channel schema.Channel
	cutoff := time.Now().Add(-d)

	err := r.DB.Where("last_fetch IS NULL OR last_fetch < ?", cutoff).
		First(&channel).Error
	if err != nil {
		return nil, err
	}

	return &channel, nil
}
