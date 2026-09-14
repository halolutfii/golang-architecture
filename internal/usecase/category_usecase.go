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

func (c *CategoryUseCase) List(ctx context.Context) ([]model.CategoryResponse, error) {
	// Try to serve from cache first. A cache miss or a Redis outage must NOT
	// fail the request: we simply fall through to the database.
	if cached, ok := c.getCategoriesFromCache(ctx); ok {
		return cached, nil
	}

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

	// Best-effort cache write. If Redis is down, log and continue.
	c.saveCategoriesToCache(ctx, responses)

	return responses, nil
}

// getCategoriesFromCache returns cached categories and true on a cache hit.
// On a miss or any Redis error, it returns false so the caller falls back to DB.
func (c *CategoryUseCase) getCategoriesFromCache(ctx context.Context) ([]model.CategoryResponse, bool) {
	value, err := c.Redis.Get(ctx, categoriesCacheKey).Result()
	if err != nil {
		// redis.Nil = normal cache miss; anything else = Redis problem.
		if err != redis.Nil {
			c.Log.WithError(err).Warn("failed to read categories from redis, falling back to database")
		}
		return nil, false
	}

	var responses []model.CategoryResponse
	if err := json.Unmarshal([]byte(value), &responses); err != nil {
		c.Log.WithError(err).Warn("failed to unmarshal cached categories, falling back to database")
		return nil, false
	}

	return responses, true
}

// saveCategoriesToCache stores categories in Redis. Failures are logged but
// never returned, so a Redis outage cannot break the endpoint.
func (c *CategoryUseCase) saveCategoriesToCache(ctx context.Context, responses []model.CategoryResponse) {
	jsonValue, err := json.Marshal(responses)
	if err != nil {
		c.Log.WithError(err).Warn("failed to marshal categories for redis")
		return
	}

	if err := c.Redis.Set(ctx, categoriesCacheKey, jsonValue, time.Hour*1).Err(); err != nil {
		c.Log.WithError(err).Warn("failed to save categories to redis")
	}
}
