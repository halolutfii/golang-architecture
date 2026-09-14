package converter

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
)

func CategoryToResponse(category *entity.Category) model.CategoryResponse {
	return model.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}
