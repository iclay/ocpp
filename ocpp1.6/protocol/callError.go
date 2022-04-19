package protocol

import (
	"encoding/json"
)

type ErrorDetails struct{}

type CallError struct {
	MessageTypeID    MessageType  `json:"messageTypeId" validate:"required,eq=4"`
	UniqueID         string       `json:"uniqueId" validate:"required,max=36"`
	ErrorCode        ErrCodeType  `json:"errorCode" validate:"errorCode"`
	ErrorDescription string       `json:"errorDescription" validate:"required"`
	ErrorDetails     ErrorDetails `json:"errorDetails" validate:"omitempty"`
}

func (ce CallError) Action() string {
	return "CallError"
}
func (ce *CallError) MessageType() MessageType {
	return ce.MessageTypeID
}

func (ce *CallError) UID() string {
	return ce.UniqueID
}

func (ce *CallError) String() string {
	callErrBytes, _ := json.Marshal(ce)
	return string(callErrBytes)
}

func (ce *CallError) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 5)
	fields[0], fields[1], fields[2], fields[3], fields[4] = int(ce.MessageTypeID), ce.UniqueID, ce.ErrorCode, ce.ErrorDescription, ce.ErrorDetails
	return json.Marshal(fields)
}

func (ce *CallError) UnmarshalErrorDetails(m interface{}) error {
	errorDetails, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(errorDetails, &ce.ErrorDetails)
}
