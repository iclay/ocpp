package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {

	Validate.RegisterValidation("remoteStopStatus", func(f validator.FieldLevel) bool {
		status := RemoteStopStatus(f.Field().String())
		switch status {
		case remoteStopAccepted, remoteStopRejected:
			return true
		default:
			return false
		}
	})
}

type RemoteStopTransactionRequest struct {
	TransactionId *int `json:"transactionId" validate:"required"`
}

func (RemoteStopTransactionRequest) Action() string {
	return RemoteStopTransactionName
}

type RemoteStopStatus string

const (
	remoteStopAccepted RemoteStopStatus = "Accepted"
	remoteStopRejected RemoteStopStatus = "Rejected"
)

type RemoteStopTransactionResponse struct {
	Status RemoteStopStatus `json:"status" validate:"required,remoteStopStatus"`
}

func (RemoteStopTransactionResponse) Action() string {
	return RemoteStopTransactionName
}
