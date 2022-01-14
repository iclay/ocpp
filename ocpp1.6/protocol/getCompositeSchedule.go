package protocol

import validator "github.com/go-playground/validator/v10"

type GetCompositeScheduleStatus string

const (
	GetCompositeScheduleStatusAccepted GetCompositeScheduleStatus = "Accepted"
	GetCompositeScheduleStatusRejected GetCompositeScheduleStatus = "Rejected"
)

func init() {
	Validate.RegisterValidation("compositeScheduleStatus", func(f validator.FieldLevel) bool {
		status := GetCompositeScheduleStatus(f.Field().String())
		switch status {
		case GetCompositeScheduleStatusAccepted, GetCompositeScheduleStatusRejected:
			return true
		default:
			return false
		}
	})
}

type GetCompositeScheduleRequest struct {
	ConnectorId      int                  `json:"connectorId" validate:"gte=0"`
	Duration         int                  `json:"duration" validate:"gte=0"`
	ChargingRateUnit ChargingRateUnitType `json:"chargingRateUnit,omitempty" validate:"omitempty,chargingRateUnit"`
}

func (GetCompositeScheduleRequest) Action() string {
	return GetCompositeScheduleName
}

type GetCompositeScheduleResponse struct {
	Status           GetCompositeScheduleStatus `json:"status" validate:"required,compositeScheduleStatus"`
	ConnectorId      int                        `json:"connectorId,omitempty" validate:"omitempty,gt=0"`
	ScheduleStart    string                     `json:"scheduleStart,omitempty" validate:"omitempty,dateTime"`
	ChargingSchedule ChargingSchedule           `json:"chargingSchedule,omitempty" validate:"omitempty"`
}

func (GetCompositeScheduleResponse) Action() string {
	return GetCompositeScheduleName
}
