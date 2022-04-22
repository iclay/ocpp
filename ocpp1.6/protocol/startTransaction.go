package protocol

type StartTransactionRequest struct {
	ConnectorId   *int    `json:"connectorId" validate:"required,gte=0"`
	IdTag         IdToken `json:"idTag" validate:"required,max=20"`
	MeterStart    *int    `json:"meterStart" validate:"required,gte=0"`
	ReservationId *int    `json:"reservationId,omitempty" validate:"omitempty"`
	Timestamp     string  `json:"timestamp" validate:"required,dateTime"`
}

func (StartTransactionRequest) Action() string {
	return StartTransactionName
}
func (r *StartTransactionRequest) Reset() {
	r.ConnectorId = nil
	r.IdTag = ""
	r.MeterStart = nil
	r.ReservationId = nil
	r.Timestamp = ""
}

type StartTransactionResponse struct {
	IdTagInfo     IdTagInfo `json:"idTagInfo" validate:"required"`
	TransactionId *int      `json:"transactionId" validate:"required"`
}

func (StartTransactionResponse) Action() string {
	return StartTransactionName
}

func (r *StartTransactionResponse) Reset() {
	r.IdTagInfo = IdTagInfo{}
	r.TransactionId = nil
}
