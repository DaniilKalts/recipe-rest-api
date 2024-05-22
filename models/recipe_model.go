package models

import "time"

type Ingredient struct {
	Ingredient string `json:"ingredient"`
	Quantity   string `json:"quantity"`
}

type Recipe struct {
	Id           int          `json:"id"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []string     `json:"instructions"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}