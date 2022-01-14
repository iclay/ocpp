package protocol

type StartTransactionRequest struct {
	ConnectorId   int     `json:"type" validate:"required,gte=0"`
	IdTag         IdToken `json:"idTag" validate:"required,max=20"`
	MeterStart    int     `json:"meterStart" validate:"required,gte=0"`
	ReservationId int     `json:"reservationId,omitempty" validate:"omitempty"`
	Timestamp     string  `json:"timestamp" validate:"required,dateTime"`
}

func (StartTransactionRequest) Action() string {
	return StartTransactionName
}

type StartTransactionResponse struct {
	IdTagInfo     IdTagInfo `json:"idTagInfo" validate:"required"`
	TransactionId int       `json:"transactionId" validate:"required"`
}

func (StartTransactionResponse) Action() string {
	return StartTransactionName
}
