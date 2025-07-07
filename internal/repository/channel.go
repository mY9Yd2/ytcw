package repository

import (
	model "github.com/mY9Yd2/ytcw/internal/model/api"
	"github.com/mY9Yd2/ytcw/internal/schema"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ChannelRepository interface {
	SaveChannel(channel *schema.Channel) error
	FindAll(p *model.Pagination) ([]schema.Channel, int64, error)
	SoftDeleteChannelByUploaderID(uploaderID string) error
	SoftDeleteChannelByChannelID(channelID string) error
	GetStaleChannel(d time.Duration) (*schema.Channel, error)
	UpdateChannelLastFetch(channelID uint, lastFetch time.Time) error
}

type channelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepository{db: db}
}

func (r *channelRepository) SaveChannel(channel *schema.Channel) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "uploader_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"channel",
			"channel_id",
			"uploader_id",
		}),
	}).Create(channel).Error
}

func (r *channelRepository) UpdateChannelLastFetch(channelID uint, lastFetch time.Time) error {
	return r.db.Model(&schema.Channel{}).
		Where("channel_id = ?", channelID).
		Update("last_fetch", lastFetch).Error
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

func (r *channelRepository) GetStaleChannel(d time.Duration) (*schema.Channel, error) {
	var channel schema.Channel
	cutoff := time.Now().UTC().Add(-d)

	err := r.db.Where("last_fetch IS NULL OR last_fetch < ?", cutoff).
		First(&channel).Error
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (r *channelRepository) SoftDeleteChannelByUploaderID(uploaderID string) error {
	var channel schema.Channel
	if err := r.db.Where("uploader_id = ?", uploaderID).Find(&channel).Error; err != nil {
		return err
	}
	return r.softDeleteChannel(channel)
}

func (r *channelRepository) SoftDeleteChannelByChannelID(channelID string) error {
	var channel schema.Channel
	if err := r.db.Where("channel_id = ?", channelID).Find(&channel).Error; err != nil {
		return err
	}
	return r.softDeleteChannel(channel)
}

func (r *channelRepository) softDeleteChannel(channel schema.Channel) error {
	tx := r.db.Begin()

	// Because GORM does not automatically soft-delete related videos,
	// we need to manually mark them as deleted first.
	if err := tx.Where("channel_refer = ?", channel.ID).Delete(&schema.Video{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&channel).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
