package repository

import (
	"github.com/google/uuid"
	"github.com/mY9Yd2/ytcw/internal/schema"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryRepository interface {
	SaveCategory(category string) (uuid.UUID, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) SaveCategory(category string) (uuid.UUID, error) {
	c := &schema.Category{
		Name: category,
	}

	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).FirstOrCreate(&c).Error; err != nil {
		return c.ID, err
	}

	return c.ID, nil
}
