package usecase

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	CategoryRepository *repository.CategoryRepository
}

func NewCategoryUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	CategoryRepository *repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		DB:                 db,
		Log:                logger,
		Validate:           validate,
		CategoryRepository: CategoryRepository,
	}
}

func (c *CategoryUseCase) List(ctx context.Context) ([]model.CategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// load parent categories
	parents, err := c.CategoryRepository.FindAllParents(tx)
	if err != nil {
		c.Log.WithError(err).Error("failed to load parent categories")
		return nil, fiber.ErrInternalServerError
	}

	var responses []model.CategoryResponse
	// iterate each parent
	for _, parent := range parents {
		response := converter.CategoryToResponse(&parent)

		var childResponses []model.CategoryResponse
		for _, child := range parent.Children {
			responseChild := converter.CategoryToResponse(&child)
			childResponses = append(childResponses, responseChild)
		}

		response.Children = childResponses

		responses = append(responses, response)
	}

	return responses, nil
}
