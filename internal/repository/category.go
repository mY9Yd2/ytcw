package repository

import (
	"github.com/google/uuid"
	model "github.com/mY9Yd2/ytcw/internal/model/api"
	"github.com/mY9Yd2/ytcw/internal/schema"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	SaveCategory(category string) (uuid.UUID, error)
	FindAll(p *model.Pagination) ([]schema.Category, int64, error)
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

	if err := r.db.Where(schema.Category{Name: category}).
		FirstOrCreate(&c).Error; err != nil {
		return c.ID, err
	}

	return c.ID, nil
}

func (r *categoryRepository) FindAll(p *model.Pagination) ([]schema.Category, int64, error) {
	var categories []schema.Category
	var total int64

	db := r.db.Model(&schema.Category{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Limit(int(p.Limit())).
		Offset(int(p.Offset())).
		Order("name ASC").
		Find(&categories).Error

	return categories, total, err
}
