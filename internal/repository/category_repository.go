package repository

import (
	"golang-clean-architecture/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	Repository[entity.Category]
	Log *logrus.Logger
}

func NewCategoryRepository(log *logrus.Logger) *CategoryRepository {
	return &CategoryRepository{
		Log: log,
	}
}

func (r *CategoryRepository) FindAllParents(tx *gorm.DB) ([]entity.Category, error) {
	var categories []entity.Category
	if err := tx.Where("parent_id is null ").Preload("Children").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// func (r *CategoryRepository) FindAllChildren(tx *gorm.DB, parentId string) ([]entity.Category, error) {
// 	var categories []entity.Category
// 	if err := tx.Where("parent_id = ?", parentId).Find(&categories).Error; err != nil {
// 		return nil, err
// 	}
// 	return categories, nil
// }
