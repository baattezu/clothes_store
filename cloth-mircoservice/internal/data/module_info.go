package data

import (
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type ClotheInfo struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ClothName string    `json:"cloth_name"` // Corrected field name
	ClothCost int32     `json:"cloth_cost"` // Corrected field name
	ClothSize string    `json:"cloth_size"` // Corrected field name
	Version   int       `json:"version"`    // Corrected type to int
}
