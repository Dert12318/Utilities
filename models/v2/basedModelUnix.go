package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/soft_delete"
)

type (
	BaseModelUnix struct {
		CreatedAt int64  `json:"createdAt"`
		CreatedBy string `json:"createdBy"`
		UpdatedAt int64  `json:"updatedAt"`
		UpdatedBy string `json:"updatedBy"`
	}

	BaseModelSoftDeleteUnix struct {
		BaseModelUnix
		DeletedAt soft_delete.DeletedAt `json:"deletedAt"`
		DeletedBy *string               `json:"deletedBy"`
	}

	BaseModelUnixMilli struct {
		CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime:milli"`
		CreatedBy string `json:"createdBy"`
		UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime:milli"`
		UpdatedBy string `json:"updatedBy"`
	}

	BaseModelSoftDeleteUnixMilli struct {
		BaseModelUnixMilli
		DeletedAt soft_delete.DeletedAt `json:"deletedAt" gorm:"softDelete:milli"`
		DeletedBy *string               `json:"deletedBy"`
	}

	UnixMilliSerializer struct{}

	// BaseModelUnixMilliSerializer unixMilli to timestamp
	BaseModelUnixMilliSerializer struct {
		CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime:milli;serializer:unixMilliTime;type:time"`
		CreatedBy string `json:"createdBy"`
		UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime:false;serializer:unixMilliTime;type:time"`
		UpdatedBy string `json:"updatedBy"`
	}

	BaseModelSoftDeleteUnixMilliSerializer struct {
		BaseModelUnixMilliSerializer
		DeletedAt gorm.DeletedAt `json:"deletedAt"`
		DeletedBy *string        `json:"deletedBy"`
	}
)

// Scan implements serializer interface
func (UnixMilliSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	t := sql.NullTime{}
	if err = t.Scan(dbValue); err == nil && t.Valid {
		err = field.Set(ctx, dst, t.Time.UnixMilli())
	}

	return
}

// Value implements serializer interface
func (UnixMilliSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (result interface{}, err error) {
	rv := reflect.ValueOf(fieldValue)
	switch v := fieldValue.(type) {
	case int64, int, uint, uint64, int32, uint32, int16, uint16:
		result = time.UnixMilli(reflect.Indirect(rv).Int())
	case *int64, *int, *uint, *uint64, *int32, *uint32, *int16, *uint16:
		if rv.IsZero() {
			return nil, nil
		}
		result = time.UnixMilli(reflect.Indirect(rv).Int())
	default:
		err = fmt.Errorf("invalid field type %#v for UnixMilliSerializer, only int, uint supported", v)
	}
	return
}
