package model

type CategoryResponse struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Children []CategoryResponse `json:"children"`
}
