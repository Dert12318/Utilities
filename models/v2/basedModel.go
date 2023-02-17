package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	BaseModel struct {
		CreatedAt time.Time `json:"createdAt"`
		CreatedBy string    `json:"createdBy"`
		UpdatedAt time.Time `json:"updatedAt"`
		UpdatedBy string    `json:"updatedBy"`
	}

	BaseModelSoftDelete struct {
		BaseModel
		DeletedAt gorm.DeletedAt `json:"deletedAt"`
		DeletedBy *string        `json:"deletedBy"`
	}
)
