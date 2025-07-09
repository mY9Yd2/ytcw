package repository

import (
	model "github.com/mY9Yd2/ytcw/internal/model/api"
	"github.com/mY9Yd2/ytcw/internal/schema"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VideoRepository interface {
	SaveVideo(video *schema.Video) error
	FindAll(p *model.Pagination) ([]schema.Video, int64, error)
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

func (r *videoRepository) FindAll(p *model.Pagination) ([]schema.Video, int64, error) {
	var videos []schema.Video
	var total int64

	db := r.db.Model(&schema.Video{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Preload("Channel").
		Preload("Channel.Category").
		Limit(int(p.Limit())).
		Offset(int(p.Offset())).
		Order("created_at DESC").
		Find(&videos).Error

	return videos, total, err
}
