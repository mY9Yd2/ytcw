package content

import (
	"github.com/mY9Yd2/ytcw/internal/common"
)

type CategoryService interface {
	GetCategories(p *common.Pagination) ([]CategoryResponse, *common.Pagination, error)
}

type categoryService struct {
	categoryRepo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: repo,
	}
}

func (s *categoryService) GetCategories(p *common.Pagination) ([]CategoryResponse, *common.Pagination, error) {
	categories, total, err := s.categoryRepo.FindAll(p)
	if err != nil {
		return nil, nil, err
	}

	var responses []CategoryResponse
	for _, c := range categories {
		responses = append(responses, CategoryResponse{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	p.TotalRows = uint(total)
	p.TotalPages = (p.TotalRows + p.PageSize - 1) / p.PageSize

	if responses == nil {
		responses = []CategoryResponse{}
	}

	return responses, p, nil
}
