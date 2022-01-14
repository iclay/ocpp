package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {
	Validate.RegisterValidation("errorCode", func(f validator.FieldLevel) bool {
		errcode := ErrCodeType(f.Field().String())
		switch errcode {
		case NotImplemented, NotSupported, CallInternalError, ProtocolError,
			SecurityError, FormationViolation, PropertyConstraintViolation, OccurenceConstraintViolation,
			TypeConstraintViolation, GenericError:
			return true
		}
		return false
	})
}

type ErrCodeType string

const (
	NotImplemented               ErrCodeType = "NotImplemented"
	NotSupported                 ErrCodeType = "NotSupported"
	CallInternalError            ErrCodeType = "InternalError"
	ProtocolError                ErrCodeType = "ProtocolError"
	SecurityError                ErrCodeType = "SecurityError"
	FormationViolation           ErrCodeType = "FormationViolation"
	PropertyConstraintViolation  ErrCodeType = "PropertyConstraintViolation"
	OccurenceConstraintViolation ErrCodeType = "OccurenceConstraintViolation"
	TypeConstraintViolation      ErrCodeType = "TypeConstraintViolation"
	GenericError                 ErrCodeType = "GenericError"
)
