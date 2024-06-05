package data

import (
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type ClotheInfo struct {
	ID             int       `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ModuleName     string    `json:"module_name"`
	ModuleDuration int32     `json:"module_duration"`
	ExamType       string    `json:"exam_type"`
	Version        string    `json:"version"`
}
