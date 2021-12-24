package proto

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
