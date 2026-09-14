package entity

type Category struct {
	ID       string     `gorm:"column:id;primaryKey"`
	Name     string     `gorm:"column:name"`
	ParentId string     `gorm:"column:parent_id"`
	Parent   *Category  `gorm:"foreignKey:parent_id;references:id"`
	Children []Category `gorm:"foreignKey:parent_id;references:id"`
}

func (a *Category) TableName() string {
	return "categories"
}
