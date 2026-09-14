package usecase

import (
	"context"
	"encoding/json"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	CategoryRepository *repository.CategoryRepository
	Redis              *redis.Client
}

func NewCategoryUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	CategoryRepository *repository.CategoryRepository, redis *redis.Client) *CategoryUseCase {
	return &CategoryUseCase{
		DB:                 db,
		Log:                logger,
		Validate:           validate,
		CategoryRepository: CategoryRepository,
		Redis:              redis,
	}
}

const categoriesCacheKey = "categories"

func (c *CategoryUseCase) List(ctx context.Context) (string, error) {
	var responses []model.CategoryResponse
	// check in redis
	value, err := c.Redis.Get(ctx, categoriesCacheKey).Result()
	// if exixst, return the data
	if err == nil {
		return value, nil
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// if not found in redis, query to database
	// load parent categories
	parents, err := c.CategoryRepository.FindAllParents(tx)
	if err != nil {
		c.Log.WithError(err).Error("failed to load parent categories")
		return "", fiber.ErrInternalServerError
	}

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

	// Best-effort cache write. If Redis is down, log and continue.
	jsonValue, err := json.Marshal(fiber.Map{
		"data": responses,
	})
	if err != nil {
		c.Log.WithError(err).Warn("failed to marshal categories for redis")
		return "", fiber.ErrInternalServerError
	}

	if err := c.Redis.Set(ctx, categoriesCacheKey, jsonValue, time.Hour*1).Err(); err != nil {
		c.Log.WithError(err).Warn("failed to save categories to redis")
	}

	return string(jsonValue), nil
}
