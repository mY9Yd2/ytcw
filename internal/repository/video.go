package repository

import (
	"github.com/mY9Yd2/ytcw/internal/schema"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VideoRepository interface {
	SaveVideo(video *schema.Video) error
}

type videoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &videoRepository{db: db}
}

func (r *videoRepository) SaveVideo(video *schema.Video) error {
	var channel schema.Channel
	if err := r.db.Where("uploader_id = ?", video.Channel.UploaderID).
		First(&channel, &video.Channel).Error; err != nil {
		return err
	}

	video.ChannelRefer = channel.ID

	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "display_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"full_title",
			"duration",
			"thumbnail",
			"timestamp",
		}),
	}).Create(video).Error
}
