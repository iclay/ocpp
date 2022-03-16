package protocol

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	// "reflect"
	"testing"
	"time"
)

func RandomString(len int) string {
	var numbers = []byte{1, 2, 3, 4, 5, 7, 8, 9}
	var container string
	length := bytes.NewReader(numbers).Len()

	for i := 1; i <= len; i++ {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {

		}
		container += fmt.Sprintf("%d", numbers[random.Int64()])
	}
	return container
}

/*****************call***************/

var fnBootNotificationRequestSuccess = func() BootNotificationRequest {
	return BootNotificationRequest{
		ChargePointVendor:       "qinglianyun",
		ChargePointModel:        "lihuaye",
		ChargePointSerialNumber: RandomString(15),
		ChargeBoxSerialNumber:   RandomString(15),
		FirmwareVersion:         RandomString(15),
		Iccid:                   RandomString(15),
		Imsi:                    RandomString(15),
		MeterType:               RandomString(15),
		MeterSerialNumber:       RandomString(15),
	}
}

var fnBootNotificationRequestFailed = func() BootNotificationRequest {
	return BootNotificationRequest{
		ChargePointVendor:       "qinglianyun",
		ChargePointModel:        "lihuaye",
		ChargePointSerialNumber: RandomString(150),
		ChargeBoxSerialNumber:   RandomString(150),
		FirmwareVersion:         RandomString(150),
		Iccid:                   RandomString(150),
		Imsi:                    RandomString(150),
		MeterType:               RandomString(150),
		MeterSerialNumber:       RandomString(150),
	}
}

var fnBootNotificationResponseSuccess = func() BootNotificationResponse {
	Interval := 10
	return BootNotificationResponse{
		CurrentTime: time.Now().Format(time.RFC3339),
		Interval:    &Interval,
		Status:      "Accepted",
	}
}

var fnBootNotificationResponseFailed = func() BootNotificationResponse {
	Interval := 10
	return BootNotificationResponse{
		CurrentTime: time.Now().Format(time.RFC3339),
		Interval:    &Interval,
		Status:      "Accepted1",
	}
}

func callSuccess() *Call {
	return &Call{
		MessageTypeID: 2,
		UniqueID:      "uniqueid",
		Action:        "Authorize",
		Request:       fnBootNotificationRequestSuccess(),
	}

}

func callFailed() *Call {
	return &Call{
		MessageTypeID: 1,
		UniqueID:      "uniqueid",
		Action:        "Authorize",
		Request:       fnBootNotificationRequestSuccess(),
	}

}

func testCall(call *Call, t *testing.T) {

	err := Validate.Struct(call)
	if err != nil {
		t.Error(err)
	}
}

func TestCall(t *testing.T) {
	testCall(callSuccess(), t)
	// testCall(callFailed(), t)
}

func callResultSuccess() *CallResult {
	return &CallResult{
		MessageTypeID: 3,
		UniqueID:      "uniqueid",
		Response:      fnBootNotificationResponseSuccess(),
	}
}

func callResultFailed() *CallResult {
	return &CallResult{
		MessageTypeID: 3,
		UniqueID:      "uniqueid",
		Response:      fnBootNotificationResponseFailed(),
	}
}

func testCallResult(callResult *CallResult, t *testing.T) {

	err := Validate.Struct(callResult)
	if err != nil {
		t.Error(err)
	}
}

func TestCallResult(t *testing.T) {
	testCallResult(callResultSuccess(), t)
	// testCallResult(callResultFailed(), t)
}

func callErrorSuccess() *CallError {
	return &CallError{
		MessageTypeID:    4,
		UniqueID:         "uniqueid",
		ErrorCode:        "NotImplemented",
		ErrorDescription: "ErrorDescription",
	}

}

func callErrorFailed() *CallError {
	return &CallError{
		MessageTypeID:    5,
		UniqueID:         "uniqueid",
		ErrorCode:        "NotImplemented",
		ErrorDescription: "ErrorDescription",
	}
}

func testCallError(callError *CallError, t *testing.T) {

	err := Validate.Struct(callError)
	if err != nil {
		t.Error(err)
	}
}

func TestCallError(t *testing.T) {
	testCallError(callErrorSuccess(), t)
	// testCallError(callErrorFailed(), t)
}

/*****************BootNotification***************/
func BootNotificationRequestSuccess() *BootNotificationRequest {

	return &BootNotificationRequest{
		ChargePointVendor:       RandomString(15),
		ChargePointModel:        RandomString(15),
		ChargePointSerialNumber: RandomString(15),
		ChargeBoxSerialNumber:   RandomString(15),
		FirmwareVersion:         RandomString(15),
		Iccid:                   RandomString(15),
		Imsi:                    RandomString(15),
		MeterType:               RandomString(15),
		MeterSerialNumber:       RandomString(15),
	}

}

func BootNotificationRequestFail() *BootNotificationRequest {

	return &BootNotificationRequest{
		ChargePointVendor:       RandomString(25),
		ChargePointModel:        RandomString(25),
		ChargePointSerialNumber: RandomString(30),
		ChargeBoxSerialNumber:   RandomString(30),
		FirmwareVersion:         RandomString(55),
		Iccid:                   RandomString(25),
		Imsi:                    RandomString(25),
		MeterType:               RandomString(50),
		MeterSerialNumber:       RandomString(50),
	}
}

func testBootNotificationRequest(call *BootNotificationRequest, t *testing.T) {
	bootnotification_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &BootNotificationRequest{}
	err = json.Unmarshal(bootnotification_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestBootNotificationRequest(t *testing.T) {
	testBootNotificationRequest(BootNotificationRequestSuccess(), t)
	//testBootNotificationRequest(BootNotificationRequestFail(), t) //error

}

func BootNotificationResponseSuccess() *BootNotificationResponse {
	Interval := 10
	return &BootNotificationResponse{
		CurrentTime: time.Now().Format(time.RFC3339),
		Interval:    &Interval,
		Status:      "Accepted",
	}

}

func BootNotificationResponseFail() *BootNotificationResponse {
	Interval := -1
	return &BootNotificationResponse{
		CurrentTime: "wpqppq",
		Interval:    &Interval,
		Status:      "accepted",
	}
}

func testBootNotificationResponse(call *BootNotificationResponse, t *testing.T) {
	bootnotification_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &BootNotificationResponse{}
	err = json.Unmarshal(bootnotification_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestBootNotificationResponse(t *testing.T) {
	testBootNotificationResponse(BootNotificationResponseSuccess(), t)
	// testBootNotificationResponse(BootNotificationResponseFail(), t) //error

}

/*****************StatusNotification ***************/

func StatusNotificationRequestSuccess() *StatusNotificationRequest {
	ConnectorId := 15
	return &StatusNotificationRequest{
		ConnectorId: &ConnectorId,
		ErrorCode:   "ConnectorLockFailure",
		Info:        RandomString(40),
		Status:      "Available",
		// Timestamp:       time.Now().Format(time.RFC3339),
		Timestamp:       time.Now().Format(ISO8601),
		VendorId:        RandomString(240),
		VendorErrorCode: RandomString(40),
	}

}

func StatusNotificationRequestFail() *StatusNotificationRequest {
	ConnectorId := -5

	return &StatusNotificationRequest{
		ConnectorId:     &ConnectorId,
		ErrorCode:       "onnectorLockFailure",
		Info:            RandomString(60),
		Status:          "vailable",
		Timestamp:       "2021-12-04",
		VendorId:        RandomString(258),
		VendorErrorCode: RandomString(55),
	}

}

func testStatusNotificationRequest(call *StatusNotificationRequest, t *testing.T) {
	statusnotification_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &StatusNotificationRequest{}
	err = json.Unmarshal(statusnotification_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestStatusNotificationRequest(t *testing.T) {
	testStatusNotificationRequest(StatusNotificationRequestSuccess(), t)
	//testStatusNotificationRequest(StatusNotificationRequestFail(), t) //error
}

/*****************metervalue***************/

func MeterValueRequestSuccess() *MeterValuesRequest {
	ConnectorId := 10
	TransactionId := 1
	var meterValueReq = &MeterValuesRequest{
		ConnectorId:   &ConnectorId,
		TransactionId: &TransactionId,
	}
	var meterValue = MeterValue{
		Timestamp: time.Now().Format(time.RFC3339),
	}
	var sampledValue = SampledValue{
		Value:     RandomString(10),
		Context:   interruptionBegin,
		Format:    rawFormat,
		Measurand: energyActiveExportRegister,
		Phase:     L1,
		Location:  cable,
		Unit:      Wh,
	}
	meterValue.SampledValue = append(meterValue.SampledValue, sampledValue)
	meterValueReq.MeterValue = append(meterValueReq.MeterValue, meterValue)
	return meterValueReq
}

func MeterValueRequestFail() *MeterValuesRequest {
	ConnectorId := -1
	TransactionId := -1

	var meterValueReq = &MeterValuesRequest{
		ConnectorId:   &ConnectorId,
		TransactionId: &TransactionId,
		MeterValue:    []MeterValue{},
	}

	var meterValue = MeterValue{
		Timestamp: "2021-12-03",
	}
	var sampledValue = SampledValue{
		Value:     RandomString(10),
		Context:   ReadingContext(RandomString(10)),
		Format:    ValueFormat(RandomString(10)),
		Measurand: Measurand(RandomString(10)),
		Phase:     Phase(RandomString(10)),
		Location:  Location(RandomString(10)),
		Unit:      UnitOfMeasure(RandomString(10)),
	}
	meterValue.SampledValue = append(meterValue.SampledValue, sampledValue)
	meterValueReq.MeterValue = append(meterValueReq.MeterValue, meterValue)
	return meterValueReq
}

func MeterValueRequestFailEmptyMeterValue() *MeterValuesRequest {
	ConnectorId := 1
	TransactionId := -1

	var meterValueReq = &MeterValuesRequest{
		ConnectorId:   &ConnectorId,
		TransactionId: &TransactionId,
		MeterValue:    []MeterValue{},
	}
	return meterValueReq
}
func testMeterValueRequest(call *MeterValuesRequest, t *testing.T) {
	MeterValue_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &MeterValuesRequest{}
	err = json.Unmarshal(MeterValue_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}
func TestMeterValueRequest(t *testing.T) {
	testMeterValueRequest(MeterValueRequestSuccess(), t)
	//testMeterValueRequest(MeterValueRequestFail(), t)
	//testMeterValueRequest(MeterValueRequestFailEmptyMeterValue(), t)//测试metervalue为空的情况
}

/**************Authorize******************/

func AuthorizeRequestSuccess() *AuthorizeRequest {

	return &AuthorizeRequest{
		IdTag: IdToken(RandomString(15)),
	}

}

func AuthorizeRequestFail() *AuthorizeRequest {

	return &AuthorizeRequest{
		IdTag: IdToken(RandomString(25)),
	}

}

func testStatusAuthorizeRequest(call *AuthorizeRequest, t *testing.T) {
	authorize_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &AuthorizeRequest{}
	err = json.Unmarshal(authorize_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestAuthorizeRequest(t *testing.T) {
	testStatusAuthorizeRequest(AuthorizeRequestSuccess(), t)
	// testStatusAuthorizeRequest(AuthorizeRequestFail(), t)//error
}

func AuthorizeResponseSuccess() *AuthorizeResponse {

	return &AuthorizeResponse{
		IdTagInfo: IdTagInfo{
			ExpiryDate:  time.Now().Format(ISO8601),
			ParentIdTag: IdToken(RandomString(15)),
			Status:      authAccepted,
		},
	}

}

func AuthorizeResponseFail() *AuthorizeResponse {
	return &AuthorizeResponse{
		IdTagInfo: IdTagInfo{
			ExpiryDate:  "2021-12-13",
			ParentIdTag: IdToken(RandomString(25)),
			Status:      AuthorizationStatus(RandomString(14)),
		},
	}

}

func testStatusAuthorizeResponse(call *AuthorizeResponse, t *testing.T) {
	authorize_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &AuthorizeResponse{}
	err = json.Unmarshal(authorize_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestAuthorizeResponse(t *testing.T) {
	testStatusAuthorizeResponse(AuthorizeResponseSuccess(), t)
	//testStatusAuthorizeResponse(AuthorizeResponseFail(), t)//error
}

//// StartTransaction

func StartTransactionRequestSuccess() *StartTransactionRequest {
	MeterStart := 10
	ConnectorId := 0
	// ReservationId := 10
	return &StartTransactionRequest{
		ConnectorId: &ConnectorId,
		IdTag:       IdToken(RandomString(15)),
		MeterStart:  &MeterStart,
		// ReservationId: &ReservationId,
		Timestamp: time.Now().Format(time.RFC3339),
	}

}

func StartTransactionRequestFail() *StartTransactionRequest {
	MeterStart := 10
	ConnectorId := -1
	ReservationId := 10
	return &StartTransactionRequest{
		ConnectorId:   &ConnectorId,
		IdTag:         IdToken(RandomString(25)),
		MeterStart:    &MeterStart,
		ReservationId: &ReservationId,
		Timestamp:     "2021-12-13",
	}
}

func testStartTransactionRequest(call *StartTransactionRequest, t *testing.T) {
	StartTransaction_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &StartTransactionRequest{}
	err = json.Unmarshal(StartTransaction_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = json.Marshal(tmp)
	if err != nil {
		t.Error(err)
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestStartTransactionRequest(t *testing.T) {
	testStartTransactionRequest(StartTransactionRequestSuccess(), t)
	//testStartTransactionRequest(StartTransactionRequestFail(), t)//error
}

func StartTransactionResponseSuccess() *StartTransactionResponse {
	TransactionId := 0
	return &StartTransactionResponse{
		IdTagInfo: IdTagInfo{
			ExpiryDate:  time.Now().Format(ISO8601),
			ParentIdTag: IdToken(RandomString(15)),
			Status:      authAccepted,
		},
		TransactionId: &TransactionId,
	}

}

func StartTransactionResponseFail() *StartTransactionResponse {

	return &StartTransactionResponse{
		IdTagInfo: IdTagInfo{
			ExpiryDate:  "2021-12-13",
			ParentIdTag: IdToken(RandomString(25)),
			Status:      AuthorizationStatus(RandomString(14)),
		},
	}
}

func testStartTransactionResponse(call *StartTransactionResponse, t *testing.T) {
	StartTransaction_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &StartTransactionResponse{}
	err = json.Unmarshal(StartTransaction_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestStartTransactionResponse(t *testing.T) {
	testStartTransactionResponse(StartTransactionResponseSuccess(), t)
	// testStartTransactionResponse(StartTransactionResponseFail(), t) //error
}

/**************stopTransaction******************/

func StopTransactionRequestSuccess() *StopTransactionRequest {

	var meterValue = MeterValue{
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var sampledValue = SampledValue{
		Value:     RandomString(10),
		Context:   interruptionBegin,
		Format:    rawFormat,
		Measurand: energyActiveExportRegister,
		Phase:     L1,
		Location:  cable,
		Unit:      Wh,
	}

	meterValue.SampledValue = append(meterValue.SampledValue, sampledValue)
	MeterStop := 10
	TransactionId := 10

	return &StopTransactionRequest{
		IdTag:           IdToken(RandomString(15)),
		MeterStop:       &MeterStop,
		Timestamp:       time.Now().Format(time.RFC3339),
		TransactionId:   &TransactionId,
		Reason:          EmergencyStop,
		TransactionData: []MeterValue{meterValue},
	}
}

func StopTransactionRequestFail() *StopTransactionRequest {

	var meterValue = MeterValue{
		Timestamp: "2021-12-03",
	}
	var sampledValue = SampledValue{
		Value:     RandomString(10),
		Context:   ReadingContext(RandomString(10)),
		Format:    ValueFormat(RandomString(10)),
		Measurand: Measurand(RandomString(10)),
		Phase:     Phase(RandomString(10)),
		Location:  Location(RandomString(10)),
		Unit:      UnitOfMeasure(RandomString(10)),
	}

	meterValue.SampledValue = append(meterValue.SampledValue, sampledValue)
	MeterStop := -1
	return &StopTransactionRequest{
		IdTag:     IdToken(RandomString(21)),
		MeterStop: &MeterStop,
		//Timestamp: "2021-12-14"
		//TransactionId:   10,
		Reason:          Reason(RandomString(10)),
		TransactionData: []MeterValue{meterValue},
	}
}

func testStopTransactionRequest(call *StopTransactionRequest, t *testing.T) {
	StopTransaction_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &StopTransactionRequest{}
	err = json.Unmarshal(StopTransaction_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestStopTransactionRequest(t *testing.T) {
	testStopTransactionRequest(StopTransactionRequestSuccess(), t)
	//testStopTransactionRequest(StopTransactionRequestFail(), t)//error
}

func StopTransactionResponseSuccess() *StopTransactionResponse {

	return &StopTransactionResponse{
		IdTagInfo: IdTagInfo{
			ExpiryDate:  time.Now().Format(ISO8601),
			ParentIdTag: IdToken(RandomString(15)),
			Status:      authAccepted,
		},
	}

}

func StopTransactionResponseFail() *StopTransactionResponse {

	return &StopTransactionResponse{
		IdTagInfo: IdTagInfo{
			ExpiryDate:  "2021-12-13",
			ParentIdTag: IdToken(RandomString(25)),
			Status:      AuthorizationStatus(RandomString(14)),
		},
	}
}

func testStopTransactionResponse(call *StopTransactionResponse, t *testing.T) {
	StopTransaction_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &StopTransactionResponse{}
	err = json.Unmarshal(StopTransaction_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestStopTransactionResponse(t *testing.T) {
	testStopTransactionResponse(StopTransactionResponseSuccess(), t)
	// testStopTransactionResponse(StopTransactionResponseFail(), t) //error
}

/**************ChangeConfiguration******************/

func ChangeConfigurationRequestSuccess() *ChangeConfigurationRequest {

	return &ChangeConfigurationRequest{
		Key:   RandomString(30),
		Value: RandomString(450),
	}
}

func ChangeConfigurationRequestFail() *ChangeConfigurationRequest {

	return &ChangeConfigurationRequest{
		Key:   RandomString(100),
		Value: RandomString(650),
	}
}

func testChangeConfigurationRequest(call *ChangeConfigurationRequest, t *testing.T) {
	ChangeConfiguration_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &ChangeConfigurationRequest{}
	err = json.Unmarshal(ChangeConfiguration_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestChangeConfigurationRequest(t *testing.T) {
	testChangeConfigurationRequest(ChangeConfigurationRequestSuccess(), t)
	// testChangeConfigurationRequest(ChangeConfigurationRequestFail(), t)//error
}

func ChangeConfigurationResponseSuccess() *ChangeConfigurationResponse {

	return &ChangeConfigurationResponse{
		Status: configAccepted,
	}

}

func ChangeConfigurationResponseFail() *ChangeConfigurationResponse {

	return &ChangeConfigurationResponse{
		Status: ConfigurationStatus(RandomString(10)),
	}
}

func testChangeConfigurationResponse(call *ChangeConfigurationResponse, t *testing.T) {
	changeConfiguration_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &ChangeConfigurationResponse{}
	err = json.Unmarshal(changeConfiguration_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestChangeConfigurationResponse(t *testing.T) {
	testChangeConfigurationResponse(ChangeConfigurationResponseSuccess(), t)
	// testChangeConfigurationResponse(ChangeConfigurationResponseFail(), t) //error
}

/**************DataTransfer******************/
func DataTransferRequestSuccess() *DataTransferRequest {

	return &DataTransferRequest{
		VendorId:  RandomString(254),
		MessageId: RandomString(49),
		Data:      RandomString(100),
	}
}

func DataTransferRequestFail() *DataTransferRequest {

	return &DataTransferRequest{
		//VendorId: RandomString(256),
		MessageId: RandomString(55),
		Data:      RandomString(100),
	}
}

func testDataTransferRequest(call *DataTransferRequest, t *testing.T) {
	DataTransfer_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &DataTransferRequest{}
	err = json.Unmarshal(DataTransfer_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestDataTransferRequest(t *testing.T) {
	testDataTransferRequest(DataTransferRequestSuccess(), t)
	// testDataTransferRequest(DataTransferRequestFail(), t)//error
}

func DataTransferResponseSuccess() *DataTransferResponse {

	return &DataTransferResponse{
		Status: dataAccecpted,
		Data:   RandomString(10),
	}

}

func DataTransferResponseFail() *DataTransferResponse {

	return &DataTransferResponse{
		Status: DataTransferStatus(RandomString(10)),
		Data:   RandomString(10),
	}
}

func testDataTransferResponse(call *DataTransferResponse, t *testing.T) {
	DataTransfer, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &DataTransferResponse{}
	err = json.Unmarshal(DataTransfer, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestDataTransferResponse(t *testing.T) {
	testDataTransferResponse(DataTransferResponseSuccess(), t)
	// testDataTransferResponse(DataTransferResponseFail(), t) //error
}

/**************RemoteStartTransaction******************/

func RemoteStartTransactionRequestSuccess() *RemoteStartTransactionRequest {
	ConnectorId := 10
	ChargingProfiled := 10
	TransactionId := 10
	StackLevel := 10
	Duration := 10
	StartPeriod := 10
	var Limit float64 = 10
	NumberPhases := 10
	var MinChargingRate float64 = 114.1

	return &RemoteStartTransactionRequest{
		ConnectorId: &ConnectorId,
		IdTag:       IdToken(RandomString(18)),
		ChargingProfile: &ChargingProfile{
			ChargingProfiled:       &ChargingProfiled,
			TransactionId:          &TransactionId,
			StackLevel:             &StackLevel,
			ChargingProfilePurpose: chargePointMaxProfile,
			ChargingProfileKind:    absolute,
			RecurrencyKind:         daily,
			ValidFrom:              time.Now().Format(ISO8601),
			ValidTo:                time.Now().Format(ISO8601),
			ChargingSchedule: ChargingSchedule{
				Duration:         &Duration,
				StartSchedule:    RandomString(10),
				ChargingRateUnit: uintTypeA,
				ChargingSchedulePeriod: []ChargingSchedulePeriod{
					ChargingSchedulePeriod{
						StartPeriod:  &StartPeriod,
						Limit:        &Limit,
						NumberPhases: &NumberPhases,
					},
				},
				MinChargingRate: &MinChargingRate,
			},
		},
	}
}

func RemoteStartTransactionRequestFail() *RemoteStartTransactionRequest {
	ConnectorId := -1
	StartPeriod := -1
	var Limit float64 = -1
	NumberPhases := -1
	var MinChargingRate float64 = 114.1
	StackLevel := -1
	Duration := -1

	return &RemoteStartTransactionRequest{
		ConnectorId: &ConnectorId,
		IdTag:       IdToken(RandomString(21)),
		ChargingProfile: &ChargingProfile{
			// ChargingProfiled:10,
			//TransactionId:10,
			StackLevel:             &StackLevel,
			ChargingProfilePurpose: ChargingProfilePurposeType(RandomString(10)),
			ChargingProfileKind:    ChargingProfileKindType(RandomString(10)),
			RecurrencyKind:         RecurrencyKindType(RandomString(10)),
			//ValidFrom:RandomString(10),
			//ValidTo:RandomString(10),
			ChargingSchedule: ChargingSchedule{
				Duration: &Duration,
				//StartSchedule:RandomString(10),
				ChargingRateUnit: ChargingRateUnitType(RandomString(10)),
				ChargingSchedulePeriod: []ChargingSchedulePeriod{
					ChargingSchedulePeriod{
						StartPeriod:  &StartPeriod,
						Limit:        &Limit,
						NumberPhases: &NumberPhases,
					},
				},
				MinChargingRate: &MinChargingRate,
			},
		},
	}
}

func testRemoteStartTransactionRequest(call *RemoteStartTransactionRequest, t *testing.T) {
	RemoteStartTransaction_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &RemoteStartTransactionRequest{}
	err = json.Unmarshal(RemoteStartTransaction_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoteStartTransactionRequest(t *testing.T) {
	testRemoteStartTransactionRequest(RemoteStartTransactionRequestSuccess(), t)
	// testRemoteStartTransactionRequest(RemoteStartTransactionRequestFail(), t)//error
}

func RemoteStartTransactionrResponseSuccess() *RemoteStartTransactionResponse {

	return &RemoteStartTransactionResponse{
		Status: remoteStartAccepted,
	}

}

func RemoteStartTransactionResponseFail() *RemoteStartTransactionResponse {

	return &RemoteStartTransactionResponse{
		Status: RemoteStartStatus(RandomString(10)),
	}
}

func testRemoteStartTransactionResponse(call *RemoteStartTransactionResponse, t *testing.T) {
	RemoteStartTransaction_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &RemoteStartTransactionResponse{}
	err = json.Unmarshal(RemoteStartTransaction_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoteStartTransactionResponse(t *testing.T) {
	testRemoteStartTransactionResponse(RemoteStartTransactionrResponseSuccess(), t)
	// testRemoteStartTransactionResponse(RemoteStartTransactionResponseFail(), t) //error
}

/**************RemoteStopTransaction******************/

func RemoteStopTransactionRequestSuccess() *RemoteStopTransactionRequest {
	TransactionId := 10

	return &RemoteStopTransactionRequest{
		TransactionId: &TransactionId,
	}
}

func RemoteStopTransactionRequestFail() *RemoteStopTransactionRequest {

	return &RemoteStopTransactionRequest{}
}

func testRemoteStopTransactionRequest(call *RemoteStopTransactionRequest, t *testing.T) {
	RemoteStopTransaction_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &RemoteStopTransactionRequest{}
	err = json.Unmarshal(RemoteStopTransaction_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoteStopTransactionRequest(t *testing.T) {
	testRemoteStopTransactionRequest(RemoteStopTransactionRequestSuccess(), t)
	// testRemoteStopTransactionRequest(RemoteStopTransactionRequestFail(), t)//error
}

func RemoteStopTransactionrResponseSuccess() *RemoteStopTransactionResponse {

	return &RemoteStopTransactionResponse{
		Status: remoteStopAccepted,
	}

}

func RemoteStopTransactionResponseFail() *RemoteStopTransactionResponse {

	return &RemoteStopTransactionResponse{
		Status: RemoteStopStatus(RandomString(10)),
	}
}

func testRemoteStopTransactionResponse(call *RemoteStopTransactionResponse, t *testing.T) {
	RemoteStopTransaction_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &RemoteStopTransactionResponse{}
	err = json.Unmarshal(RemoteStopTransaction_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoteStopTransactionResponse(t *testing.T) {
	testRemoteStopTransactionResponse(RemoteStopTransactionrResponseSuccess(), t)
	// testRemoteStopTransactionResponse(RemoteStopTransactionResponseFail(), t) //error
}

/**************reset******************/

func ResetRequestSuccess() *ResetRequest {
	return &ResetRequest{
		Type: hard,
	}
}

func ResetRequestFail() *ResetRequest {

	return &ResetRequest{
		Type: ResetType(RandomString(10)),
	}
}

func testResetRequest(call *ResetRequest, t *testing.T) {
	Reset_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &ResetRequest{}
	err = json.Unmarshal(Reset_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestResetRequest(t *testing.T) {
	testResetRequest(ResetRequestSuccess(), t)
	// testResetRequest(ResetRequestFail(), t)//error
}

func ResetResponseSuccess() *ResetResponse {

	return &ResetResponse{
		Status: resetAccepted,
	}

}

func ResetResponseFail() *ResetResponse {

	return &ResetResponse{
		Status: ResetStatus(RandomString(10)),
	}
}

func testResetResponse(call *ResetResponse, t *testing.T) {
	Reset_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &ResetResponse{}
	err = json.Unmarshal(Reset_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestResetResponse(t *testing.T) {
	testResetResponse(ResetResponseSuccess(), t)
	// testResetResponse(ResetResponseFail(), t) //error
}

/**************UnlockConnector******************/

func UnlockConnectorSuccess() *UnlockConnectorRequest {
	return &UnlockConnectorRequest{
		ConnectorId: 10,
	}
}

func UnlockConnectorFail() *UnlockConnectorRequest {

	return &UnlockConnectorRequest{
		ConnectorId: 0,
	}
}

func testUnlockConnectorRequest(call *UnlockConnectorRequest, t *testing.T) {
	UnlockConnector_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &UnlockConnectorRequest{}
	err = json.Unmarshal(UnlockConnector_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestUnlockConnectorRequest(t *testing.T) {
	testUnlockConnectorRequest(UnlockConnectorSuccess(), t)
	// testUnlockConnectorRequest(UnlockConnectorFail(), t) //errors
}

func UnlockConnectorResponseSuccess() *UnlockConnectorResponse {

	return &UnlockConnectorResponse{
		Status: unlocked,
	}

}

func UnlockConnectorResponseFail() *UnlockConnectorResponse {

	return &UnlockConnectorResponse{
		Status: UnlockStatus(RandomString(10)),
	}
}

func testUnlockConnectorResponse(call *UnlockConnectorResponse, t *testing.T) {
	UnlockConnector_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &UnlockConnectorResponse{}
	err = json.Unmarshal(UnlockConnector_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestUnlockConnectorResponse(t *testing.T) {
	testUnlockConnectorResponse(UnlockConnectorResponseSuccess(), t)
	// testUnlockConnectorResponse(UnlockConnectorResponseFail(), t) //error
}

/**************************SetChargingProfile*******************************/
func SetChargingProfileSuccess() *SetChargingProfileRequest {
	ConnectorId := 10
	ChargingProfiled := 10
	TransactionId := 10
	StackLevel := 10
	StartPeriod := 10
	var Limit float64 = 10
	NumberPhases := 10
	var MinChargingRate float64 = 114.1
	Duration := 10

	return &SetChargingProfileRequest{
		ConnectorId: &ConnectorId,
		ChargingProfile: ChargingProfile{
			ChargingProfiled:       &ChargingProfiled,
			TransactionId:          &TransactionId,
			StackLevel:             &StackLevel,
			ChargingProfilePurpose: chargePointMaxProfile,
			ChargingProfileKind:    absolute,
			RecurrencyKind:         daily,
			ValidFrom:              time.Now().Format(ISO8601),
			ValidTo:                time.Now().Format(ISO8601),
			ChargingSchedule: ChargingSchedule{
				Duration:         &Duration,
				StartSchedule:    RandomString(10),
				ChargingRateUnit: uintTypeA,
				ChargingSchedulePeriod: []ChargingSchedulePeriod{
					ChargingSchedulePeriod{
						StartPeriod:  &StartPeriod,
						Limit:        &Limit,
						NumberPhases: &NumberPhases,
					},
				},
				MinChargingRate: &MinChargingRate,
			},
		},
	}
}

func SetChargingProfileFail() *SetChargingProfileRequest {
	ConnectorId := -1
	StackLevel := -1
	Duration := -1
	StartPeriod := -1
	var Limit float64 = -1
	NumberPhases := -1
	var MinChargingRate float64 = 114.1

	return &SetChargingProfileRequest{
		ConnectorId: &ConnectorId,
		ChargingProfile: ChargingProfile{
			// ChargingProfiled:10,
			//TransactionId:10,
			StackLevel:             &StackLevel,
			ChargingProfilePurpose: ChargingProfilePurposeType(RandomString(10)),
			ChargingProfileKind:    ChargingProfileKindType(RandomString(10)),
			RecurrencyKind:         RecurrencyKindType(RandomString(10)),
			//ValidFrom:RandomString(10),
			//ValidTo:RandomString(10),
			ChargingSchedule: ChargingSchedule{
				Duration: &Duration,
				//StartSchedule:RandomString(10),
				ChargingRateUnit: ChargingRateUnitType(RandomString(10)),
				ChargingSchedulePeriod: []ChargingSchedulePeriod{
					ChargingSchedulePeriod{
						StartPeriod:  &StartPeriod,
						Limit:        &Limit,
						NumberPhases: &NumberPhases,
					},
				},
				MinChargingRate: &MinChargingRate,
			},
		},
	}
}

func testSetChargingProfileRequest(call *SetChargingProfileRequest, t *testing.T) {
	SetChargingProfile_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	//t.Log(string(SetChargingProfile_reqbyte))
	var tmp = &SetChargingProfileRequest{}
	err = json.Unmarshal(SetChargingProfile_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	//t.Logf("%+v", tmp)
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestSetChargingProfileRequest(t *testing.T) {
	testSetChargingProfileRequest(SetChargingProfileSuccess(), t)
	// testSetChargingProfileRequest(SetChargingProfileFail(), t)//error
}

func SetChargingProfileResponseSuccess() *SetChargingProfileResponse {

	return &SetChargingProfileResponse{
		Status: chargingProfileAccepted,
	}

}

func SetChargingProfileResponseFail() *SetChargingProfileResponse {

	return &SetChargingProfileResponse{
		Status: ChargingProfileStatus(RandomString(10)),
	}
}

func testSetChargingProfileResponse(call *SetChargingProfileResponse, t *testing.T) {
	SetChargingProfile_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &SetChargingProfileResponse{}
	err = json.Unmarshal(SetChargingProfile_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestSetChargingProfileResponse(t *testing.T) {
	testSetChargingProfileResponse(SetChargingProfileResponseSuccess(), t)
	// testSetChargingProfileResponse(SetChargingProfileResponseFail(), t) //error
}

/**************************SendLocalList*******************************/

func SetSendLocalListSuccess() *SendLocalListRequest {
	ListVersion := 0

	return &SendLocalListRequest{
		ListVersion: &ListVersion,
		LocalAuthorizationList: []AuthorizationData{
			AuthorizationData{
				IdTag: RandomString(10),
				IdTagInfo: IdTagInfo{
					ExpiryDate:  time.Now().Format(ISO8601),
					ParentIdTag: IdToken(RandomString(10)),
					Status:      authExpired,
				},
			},
		},
		UpdateType: UpdateTypeDifferential,
	}
}

func SetSendLocalListFail() *SendLocalListRequest {
	ListVersion := -1
	return &SendLocalListRequest{
		ListVersion: &ListVersion,
		LocalAuthorizationList: []AuthorizationData{
			AuthorizationData{
				IdTag: RandomString(10),
				IdTagInfo: IdTagInfo{
					ExpiryDate:  RandomString(10),
					ParentIdTag: IdToken(RandomString(100)),
					Status:      AuthorizationStatus(RandomString(100)),
				},
			},
		},
		UpdateType: UpdateType(RandomString(100)),
	}

}

func testSendLocalListRequest(call *SendLocalListRequest, t *testing.T) {
	SendLocalList_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &SendLocalListRequest{}
	err = json.Unmarshal(SendLocalList_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestSendLocalListRequest(t *testing.T) {
	testSendLocalListRequest(SetSendLocalListSuccess(), t)
	// testSendLocalListRequest(SetSendLocalListFail(), t) //error
}

func SetSendLocalListResponseSuccess() *SendLocalListResponse {

	return &SendLocalListResponse{
		Status: UpdateStatusAccepted,
	}

}

func SetSendLocalListResponseFail() *SendLocalListResponse {

	return &SendLocalListResponse{
		Status: UpdateStatus(RandomString(10)),
	}
}

func testSendLocalListResponse(call *SendLocalListResponse, t *testing.T) {
	SendLocalList_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &SendLocalListResponse{}
	err = json.Unmarshal(SendLocalList_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestSendLocalListResponse(t *testing.T) {
	testSendLocalListResponse(SetSendLocalListResponseSuccess(), t)
	// testSendLocalListResponse(SetSendLocalListResponseFail(), t) //error
}

/**************************SendLocalList*******************************/

func SetGetLocalListVersionResponseSuccess() *GetLocalListVersionResponse {
	ListVersion := 10

	return &GetLocalListVersionResponse{
		ListVersion: &ListVersion,
	}

}

func SetGetLocalListVersionResponseFail() *GetLocalListVersionResponse {
	ListVersion := -1
	return &GetLocalListVersionResponse{
		ListVersion: &ListVersion,
	}
}

func testGetLocalListVersionResponse(call *GetLocalListVersionResponse, t *testing.T) {
	GetLocalListVersion_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &GetLocalListVersionResponse{}
	err = json.Unmarshal(GetLocalListVersion_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestGetLocalListVersionResponse(t *testing.T) {
	testGetLocalListVersionResponse(SetGetLocalListVersionResponseSuccess(), t)
	// testGetLocalListVersionResponse(SetGetLocalListVersionResponseFail(), t) //error
}

/**************************SendLocalList*******************************/

func SetGetConfigurationSuccess() *GetConfigurationRequest {
	return &GetConfigurationRequest{
		Key: []string{RandomString(50)},
	}

}

func SetGetConfigurationFail() *GetConfigurationRequest {
	return &GetConfigurationRequest{
		Key: []string{RandomString(100)},
	}
}

func testGetConfigurationRequest(call *GetConfigurationRequest, t *testing.T) {
	GetConfiguration_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &GetConfigurationRequest{}
	err = json.Unmarshal(GetConfiguration_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestGetConfigurationRequest(t *testing.T) {
	testGetConfigurationRequest(SetGetConfigurationSuccess(), t)
	// testGetConfigurationRequest(SetGetConfigurationFail(), t) //error
}

func SetGetConfigurationResponseSuccess() *GetConfigurationResponse {
	return &GetConfigurationResponse{
		ConfigurationKey: []ConfigurationKey{
			ConfigurationKey{
				Key:      RandomString(30),
				Readonly: true,
				Value:    RandomString(30),
			},
		},
		UnknownKey: []string{RandomString(30)},
	}
}

func SetGetConfigurationResponseFail() *GetConfigurationResponse {
	return &GetConfigurationResponse{
		ConfigurationKey: []ConfigurationKey{
			ConfigurationKey{
				Key:      RandomString(70),
				Readonly: true,
				Value:    RandomString(550),
			},
		},
		UnknownKey: []string{RandomString(70)},
	}
}

func testGetConfigurationResponse(call *GetConfigurationResponse, t *testing.T) {
	GetConfiguration_resbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &GetConfigurationResponse{}
	err = json.Unmarshal(GetConfiguration_resbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func TestGetConfigurationResponse(t *testing.T) {
	testGetConfigurationResponse(SetGetConfigurationResponseSuccess(), t)
	// testGetConfigurationResponse(SetGetConfigurationResponseFail(), t) //error
}

// func TestReserveNowRequest(t *testing.T) {
// 	ConnectorId := -1
// 	// ReservationId := 0

// 	a := ReserveNowRequest{
// 		ConnectorId: &ConnectorId,
// 		ExpiryDate:  time.Now().Format(time.RFC3339),
// 		IdTag:       "idtag",
// 		ParentIdTag: "ParentIdTag",
// 		// ReservationId: &ReservationId,
// 	}
// 	c, _ := json.Marshal(a)
// 	fmt.Println("###########")
// 	fmt.Println(string(c))
// 	var tmp = &ReserveNowRequest{}
// 	err := json.Unmarshal([]byte("{\"connectorId\":0,\"expiryDate\":\"2022-03-15T15:19:25+08:00\",\"idTag\":\"idtag\",\"parentIdTag\":\"ParentIdTag\"}"), tmp)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	fmt.Printf("%+v\n", tmp)
// 	err = Validate.Struct(tmp)
// 	if err != nil {
// 		t.Error(err,"111")
// 	}
// 	callResult := &CallResult{
// 		MessageTypeID: CALL_RESULT,
// 		UniqueID:      "1234567",
// 		Response:      tmp,
// 	}
// 	fmt.Printf("%+v\n", callResult)
// }

// func TestChangeAvailabilityRequest(t *testing.T) {
// 	ConnectorId := 0
// 	a := ChangeAvailabilityRequest{
// 		ConnectorId: &ConnectorId,
// 		Type:        AvailabilityTypeOperative,
// 	}
// 	c, _ := json.Marshal(a)
// 	fmt.Println("###########")
// 	fmt.Println(string(c))
// 	var tmp = &ChangeAvailabilityRequest{}
// 	err := json.Unmarshal(c, tmp)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	fmt.Printf("%+v\n", tmp)
// 	err = Validate.Struct(tmp)
// 	if err != nil {
// 		t.Error(err, "111")
// 	}
// 	callResult := &CallResult{
// 		MessageTypeID: CALL_RESULT,
// 		UniqueID:      "1234567",
// 		Response:      tmp,
// 	}
// 	fmt.Printf("%+v\n", callResult)
// }

func TestGetCompositeScheduleRequest(t *testing.T) {
	ConnectorId := 1
	Duration := 1
	a := GetCompositeScheduleRequest{
		ConnectorId:      &ConnectorId,
		Duration:         &Duration,
		ChargingRateUnit: uintTypeA,
	}
	c, _ := json.Marshal(a)
	fmt.Println("###########")
	fmt.Println(string(c))
	var tmp = &GetCompositeScheduleRequest{}
	err := json.Unmarshal(c, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", tmp)
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err, "111")
	}
	callResult := &CallResult{
		MessageTypeID: CALL_RESULT,
		UniqueID:      "1234567",
		Response:      tmp,
	}
	fmt.Printf("%+v\n", callResult)
}

// func TestGetDiagnosticsRequest(t *testing.T) {
// 	Retries := 1
// 	RetryInterval := 1
// 	a := GetDiagnosticsRequest{
// 		Location:      "{a}",
// 		Retries:       &Retries,
// 		RetryInterval: &RetryInterval,
// 		StartTime:     "bb",
// 		StopTime:      "cc",
// 	}
// 	c, _ := json.Marshal(a)
// 	fmt.Println("###########")
// 	fmt.Println(string(c))
// 	var tmp = &GetDiagnosticsRequest{}
// 	err := json.Unmarshal(c, tmp)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	fmt.Printf("%+v\n", tmp)
// 	err = Validate.Struct(tmp)
// 	if err != nil {
// 		t.Error(err, "111")
// 	}
// 	callResult := &CallResult{
// 		MessageTypeID: CALL_RESULT,
// 		UniqueID:      "1234567",
// 		Response:      tmp,
// 	}
// 	fmt.Printf("%+v\n", callResult)
// }
