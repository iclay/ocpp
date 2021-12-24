package websocket

import (
	"fmt"
	"ocpp16/proto"

	validator "github.com/go-playground/validator/v10"
)

type Error struct {
	ErrorCode        proto.ErrCodeType
	ErrorDescription string
	ErrorDetails     interface{}
}

func (e *Error) Error() string {
	return fmt.Sprintf("errcode(%v), ErrorDescription(%v),ErrorDetails(%v)", e.ErrorCode, e.ErrorDescription, e.ErrorDetails)
}

func occurenceConstraintViolation(fieldError validator.FieldError, action string) *Error {
	return &Error{
		ErrorCode:        proto.OccurenceConstraintViolation,
		ErrorDescription: fmt.Sprintf("action:%v, field %s must required but not found,but it has been omitted", action, fieldError.Namespace()),
		ErrorDetails:     nil,
	}
}

func genericError(fieldErrors validator.ValidationErrors, action string) *Error {
	return &Error{
		ErrorCode:        proto.GenericError,
		ErrorDescription: fmt.Sprintf("action:%v,error:%v", action, fieldErrors.Error()),
		ErrorDetails:     nil,
	}
}

func propertyConstraintViolation(fieldError validator.FieldError, condition string, action string) *Error {
	return &Error{
		ErrorCode:        proto.PropertyConstraintViolation,
		ErrorDescription: fmt.Sprintf("action:%v, field %s must %v %v, but the value passed is %v", action, fieldError.Namespace(), condition, fieldError.Param(), fieldError.Value()),
		ErrorDetails:     nil,
	}
}

//when multiple restriction rule errors occur in the same field, only the first error will be returned
func checkValidatorError(err error, action string) *Error {
	if validatorErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validatorError := range validatorErrors {
			switch validatorError.ActualTag() {
			case "lt", "gt", "lte", "gte", "min", "max", "eq", "ne":
				return propertyConstraintViolation(validatorError, validatorError.ActualTag(), action)
			case "required":
				return occurenceConstraintViolation(validatorError, action)
			default:
				return genericError(validatorErrors, action)
			}
		}
	}
	return &Error{
		ErrorCode:        proto.CallInternalError,
		ErrorDescription: err.Error(),
		ErrorDetails:     nil,
	}
}
