package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	BaseModel struct {
		CreatedAt time.Time `json:"created_at"`
		CreatedBy string    `json:"created_by"`
		UpdatedAt time.Time `json:"updated_at"`
		UpdatedBy string    `json:"updated_by"`
	}

	BaseModelSoftDelete struct {
		BaseModel
		DeletedAt gorm.DeletedAt `json:"deleted_at"`
		DeletedBy *string        `json:"deleted_by"`
	}
)
