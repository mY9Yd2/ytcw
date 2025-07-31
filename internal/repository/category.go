package repository

import (
	"errors"
	"github.com/google/uuid"
	model "github.com/mY9Yd2/ytcw/internal/model/api"
	"github.com/mY9Yd2/ytcw/internal/schema"
	"gorm.io/gorm"
)

type CategoryNotFoundError struct {
	CategoryName string
}

type CategoryNotEmptyError struct {
	CategoryName string
}

type CategoryRepository interface {
	SaveCategory(categoryName string) (uuid.UUID, error)
	DeleteCategory(categoryName string) error
	IsCategoryEmpty(categoryName string) (bool, error)
	FindByName(categoryName string) (*schema.Category, error)
	FindAll(p *model.Pagination) ([]schema.Category, int64, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) SaveCategory(categoryName string) (uuid.UUID, error) {
	c := &schema.Category{
		Name: categoryName,
	}

	if err := r.db.Where("name ILIKE ?", categoryName).
		FirstOrCreate(&c).Error; err != nil {
		return uuid.Nil, err
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

func (r *categoryRepository) IsCategoryEmpty(categoryName string) (bool, error) {
	var count int64

	if err := r.db.Model(&schema.Channel{}).
		Joins("JOIN categories ON channels.category_refer = categories.id").
		Where("categories.name ILIKE ?", categoryName).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *categoryRepository) DeleteCategory(categoryName string) error {
	category, err := r.FindByName(categoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &CategoryNotFoundError{CategoryName: categoryName}
		}
		return err
	}

	isEmpty, err := r.IsCategoryEmpty(categoryName)
	if err != nil {
		return err
	}

	if !isEmpty {
		return &CategoryNotEmptyError{CategoryName: categoryName}
	}

	return r.db.Delete(category).Error
}

func (r *categoryRepository) FindByName(categoryName string) (*schema.Category, error) {
	var category schema.Category
	if err := r.db.Where("name ILIKE ?", categoryName).
		First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (e *CategoryNotFoundError) Error() string {
	return "category not found"
}

func (e *CategoryNotEmptyError) Error() string {
	return "cannot delete category: category contains channels"
}
