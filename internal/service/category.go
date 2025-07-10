package service

import (
	model "github.com/mY9Yd2/ytcw/internal/model/api"
	"github.com/mY9Yd2/ytcw/internal/repository"
)

type CategoryService interface {
	GetCategories(p *model.Pagination) ([]model.CategoryResponse, *model.Pagination, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: repo,
	}
}

func (s *categoryService) GetCategories(p *model.Pagination) ([]model.CategoryResponse, *model.Pagination, error) {
	categories, total, err := s.categoryRepo.FindAll(p)
	if err != nil {
		return nil, nil, err
	}

	var responses []model.CategoryResponse
	for _, c := range categories {
		responses = append(responses, model.CategoryResponse{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	p.TotalRows = uint(total)
	p.TotalPages = (p.TotalRows + p.PageSize - 1) / p.PageSize

	if responses == nil {
		responses = []model.CategoryResponse{}
	}

	return responses, p, nil
}
