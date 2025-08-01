package repository

import (
	"github.com/google/uuid"
	model "github.com/mY9Yd2/ytcw/internal/model/api"
	"github.com/mY9Yd2/ytcw/internal/schema"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ChannelRepository interface {
	SaveChannel(channel *schema.Channel) error
	FindAll(p *model.Pagination, categoryName string) ([]schema.Channel, int64, error)
	GetChannelByUploaderID(uploaderID string) (*schema.Channel, error)
	GetChannelByChannelID(handle string) (*schema.Channel, error)
	SoftDeleteChannelByUploaderID(uploaderID string) error
	SoftDeleteChannelByChannelID(channelID string) error
	GetStaleChannel(d time.Duration) (*schema.Channel, error)
	UpdateChannelLastFetch(channelID uuid.UUID, lastFetch time.Time) error
	DisableChannelByUploaderID(uploaderID string, at, until time.Time) error
	DisableChannelByChannelID(channelID string, at, until time.Time) error
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
			"category_refer",
		}),
	}).Create(channel).Error
}

func (r *channelRepository) UpdateChannelLastFetch(channelID uuid.UUID, lastFetch time.Time) error {
	return r.db.Model(&schema.Channel{}).
		Where("id = ?", channelID).
		Update("last_fetch", lastFetch).Error
}

func (r *channelRepository) FindAll(p *model.Pagination, categoryName string) ([]schema.Channel, int64, error) {
	var channels []schema.Channel
	var total int64

	db := r.db.Model(&schema.Channel{})

	if categoryName != "" {
		db = db.Joins("JOIN categories ON categories.id = channels.category_refer").
			Where("categories.name ILIKE ?", categoryName)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Preload("Category").
		Limit(int(p.Limit())).
		Offset(int(p.Offset())).
		Order("channel ASC").
		Find(&channels).Error

	return channels, total, err
}

func (r *channelRepository) GetChannelByUploaderID(uploaderID string) (*schema.Channel, error) {
	var channel schema.Channel
	if err := r.db.Where("uploader_id ILIKE ?", uploaderID).First(&channel).Error; err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *channelRepository) GetChannelByChannelID(handle string) (*schema.Channel, error) {
	var channel schema.Channel
	if err := r.db.Where("channel_id = ?", handle).First(&channel).Error; err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *channelRepository) GetStaleChannel(d time.Duration) (*schema.Channel, error) {
	var channel schema.Channel
	now := time.Now().UTC()
	cutoff := now.Add(-d)

	err := r.db.Where("last_fetch IS NULL OR last_fetch < ?", cutoff).
		Where("disabled_at IS NULL OR disabled_at > ? OR disabled_until <= ?", now, now).
		Order(clause.OrderBy{Expression: clause.Expr{SQL: "RANDOM()"}}).
		Take(&channel).Error
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (r *channelRepository) SoftDeleteChannelByUploaderID(uploaderID string) error {
	var channel schema.Channel
	if err := r.db.Where("uploader_id ILIKE ?", uploaderID).Find(&channel).Error; err != nil {
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

func (r *channelRepository) DisableChannelByUploaderID(uploaderID string, at, until time.Time) error {
	if err := r.db.Model(&schema.Channel{}).
		Where("uploader_id ILIKE ?", uploaderID).
		Updates(map[string]interface{}{
			"disabled_at":    at,
			"disabled_until": until,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (r *channelRepository) DisableChannelByChannelID(channelID string, at, until time.Time) error {
	if err := r.db.Model(&schema.Channel{}).
		Where("channel_id = ?", channelID).
		Updates(map[string]interface{}{
			"disabled_at":    at,
			"disabled_until": until,
		}).Error; err != nil {
		return err
	}
	return nil
}
