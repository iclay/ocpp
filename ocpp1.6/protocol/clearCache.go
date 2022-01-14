package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {
	Validate.RegisterValidation("cacheStatus", func(f validator.FieldLevel) bool {
		status := ClearCacheStatus(f.Field().String())
		switch status {
		case ClearCacheStatusAccepted, ClearCacheStatusRejected:
			return true
		default:
			return false
		}
	})
}

type ClearCacheStatus string

const (
	ClearCacheStatusAccepted ClearCacheStatus = "Accepted"
	ClearCacheStatusRejected ClearCacheStatus = "Rejected"
)

type ClearCacheRequest struct{}

func (ClearCacheRequest) Action() string {
	return ClearCacheName
}

type ClearCacheResponse struct {
	Status ClearCacheStatus `json:"status" validate:"required,cacheStatus"`
}

func (ClearCacheResponse) Action() string {
	return ClearCacheName
}
