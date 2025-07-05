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

func (r *Repository) SaveCategory(category string) (uint, error) {
	c := &schema.Category{
		Name: category,
	}

	if err := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).FirstOrCreate(&c).Error; err != nil {
		return c.ID, err
	}

	return c.ID, nil
}

func (r *Repository) SaveChannel(channel *schema.Channel) error {
	return r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "uploader_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"channel",
			"channel_id",
			"uploader_id",
		}),
	}).Create(channel).Error
}

func (r *Repository) SaveVideo(video *schema.Video) error {
	var channel schema.Channel
	if err := r.DB.Where("uploader_id = ?", video.Channel.UploaderID).
		First(&channel, &video.Channel).Error; err != nil {
		return err
	}

	video.ChannelRefer = channel.ID

	return r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "display_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"full_title",
			"duration",
			"thumbnail",
			"timestamp",
		}),
	}).Create(video).Error
}

func (r *Repository) GetStaleChannel(d time.Duration) (*schema.Channel, error) {
	var channel schema.Channel
	cutoff := time.Now().UTC().Add(-d)

	err := r.DB.Where("last_fetch IS NULL OR last_fetch < ?", cutoff).
		First(&channel).Error
	if err != nil {
		return nil, err
	}

	return &channel, nil
}
