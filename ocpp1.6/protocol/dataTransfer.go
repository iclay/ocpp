package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {
	Validate.RegisterValidation("dataTransferStatus", func(f validator.FieldLevel) bool {
		status := DataTransferStatus(f.Field().String())
		switch status {
		case dataAccecpted, dataRejected, dataUnknownMessageId, dataUnknownVendorId:
			return true
		default:
			return false
		}
	})
}

type DataTransferRequest struct {
	VendorId  string `json:"vendorId"  validate:"required,max=255"`
	MessageId string `json:"messageId,omitempty" validate:"omitempty,max=50"`
	Data      string `json:"data,omitempty" validate:"omitempty"`
}

func (DataTransferRequest) Action() string {
	return DataTransferName
}

type DataTransferStatus string

const (
	dataAccecpted        DataTransferStatus = "Accepted"
	dataRejected         DataTransferStatus = "Rejected"
	dataUnknownMessageId DataTransferStatus = "UnknownMessageId"
	dataUnknownVendorId  DataTransferStatus = "UnknownVendorId"
)

type DataTransferResponse struct {
	Status DataTransferStatus `json:"status" validate:"required,dataTransferStatus"`
	Data   string             `json:"data,omitempty" validate:"omitempty"`
}

func (DataTransferResponse) Action() string {
	return DataTransferName
}
