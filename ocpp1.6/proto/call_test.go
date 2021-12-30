package proto

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
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
func callSuccess() *Call {
	return &Call{
		MessageTypeID: 2,
		UniqueID:      "uniqueid",
		Action:        "Authorize",
		// Request:       "request",
	}

}

func callFaildDueMessageTypeID() *Call {
	return &Call{
		MessageTypeID: 5,
		UniqueID:      "uniqueid",
		Action:        "Authorize",
		Request:       "request",
	}
}

func callFaildDueUniqueID() *Call {
	return &Call{
		MessageTypeID: 5,
		UniqueID:      RandomString(38), //length>36
		Action:        "Authorize",
		Request:       "request",
	}
}

func callFaildDueAction() *Call {
	return &Call{
		MessageTypeID: 5,
		UniqueID:      "uniqueid",
		Action:        RandomString(38), //length>36
		Request:       "request",
	}
}

func testCall(call *Call, t *testing.T) {
	call_byte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var callTmp = &Call{}
	call_byte[0] = '{'
	call_byte[len(call_byte)-1] = '}'
	err = json.Unmarshal(call_byte, callTmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(callTmp)
	if err != nil {
		t.Error(err)
	}
}

func Test_Call(t *testing.T) {
	// testCall(callSuccess(), t)
	//testCall(callFaildDueMessageTypeID(), t)
	//testCall(callFaildDueUniqueID(), t)//error
	//testCall(callFaildDueAction(), t)//error
}

func callResultSuccess() *CallResult {
	return &CallResult{
		MessageTypeID: 3,
		UniqueID:      "uniqueid",
		Response:      "response",
	}

}

func callResultFaildDueMessageTypeID() *CallResult {
	return &CallResult{
		MessageTypeID: 5,
		UniqueID:      "uniqueid",
		Response:      "response",
	}
}

func callResultFaildDueUniqueID() *CallResult {
	return &CallResult{
		MessageTypeID: 5,
		UniqueID:      RandomString(38), //length>36
		Response:      "response",
	}
}

func testCallResult(call *CallResult, t *testing.T) {
	call_byte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var callTmp = &CallResult{}
	err = json.Unmarshal(call_byte, callTmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(callTmp)
	if err != nil {
		t.Error(err)
	}
}

func Test_CallResult(t *testing.T) {
	testCallResult(callResultSuccess(), t)
	//testCallResult(callResultFaildDueMessageTypeID(), t) //error
	//testCallResult(callResultFaildDueUniqueID(), t) //error
}

func callErrorSuccess() *CallError {
	return &CallError{
		MessageTypeID:    4,
		UniqueID:         "uniqueid",
		ErrorCode:        "NotImplemented",
		ErrorDescription: "ErrorDescription",
	}

}

func callErrorFaildDueMessageTypeID() *CallError {
	return &CallError{
		MessageTypeID:    3,
		UniqueID:         "uniqueid",
		ErrorCode:        "NotImplemented",
		ErrorDescription: "ErrorDescription",
	}
}

func callErrorFaildDueUniqueID() *CallError {
	return &CallError{
		MessageTypeID:    4,
		UniqueID:         RandomString(38), //length>36
		ErrorCode:        "NotImplemented",
		ErrorDescription: "ErrorDescription",
	}
}

func callErrorFaildDueErrcode() *CallError {
	return &CallError{
		MessageTypeID:    4,
		UniqueID:         "uniqueid", //length>36
		ErrorCode:        "notImplemented",
		ErrorDescription: "ErrorDescription",
	}
}

func testCallError(call *CallError, t *testing.T) {
	call_byte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var callTmp = &CallError{}
	err = json.Unmarshal(call_byte, callTmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(callTmp)
	if err != nil {
		t.Error(err)
	}
}

func Test_CallErrorResult(t *testing.T) {
	// testCallError(callErrorSuccess(), t)
	//testCallError(callErrorFaildDueMessageTypeID(), t) //error
	//testCallError(callErrorFaildDueUniqueID(), t) //error
	//testCallError(callErrorFaildDueErrcode(), t) //error
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

func Test_BootNotificationRequest(t *testing.T) {
	testBootNotificationRequest(BootNotificationRequestSuccess(), t)
	//testBootNotificationRequest(BootNotificationRequestFail(), t) //error

}

func BootNotificationResponseSuccess() *BootNotificationResponse {

	return &BootNotificationResponse{
		CurrentTime: time.Now().Format(time.RFC3339),
		Interval:    10,
		Status:      "Accepted",
	}

}

func BootNotificationResponseFail() *BootNotificationResponse {

	return &BootNotificationResponse{
		CurrentTime: "wpqppq",
		Interval:    -1,
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

func Test_BootNotificationResponse(t *testing.T) {
	testBootNotificationResponse(BootNotificationResponseSuccess(), t)
	// testBootNotificationResponse(BootNotificationResponseFail(), t) //error

}

/*****************StatusNotification ***************/

func StatusNotificationRequestSuccess() *StatusNotificationRequest {

	return &StatusNotificationRequest{
		ConnectorId: 15,
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

	return &StatusNotificationRequest{
		ConnectorId:     -5,
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

func Test_StatusNotificationRequest(t *testing.T) {
	testStatusNotificationRequest(StatusNotificationRequestSuccess(), t)
	//testStatusNotificationRequest(StatusNotificationRequestFail(), t) //error
}

/*****************metervalue***************/

func MeterValueRequestSuccess() *MeterValuesRequest {
	var meterValueReq = &MeterValuesRequest{
		ConnectorId:   10,
		TransactionId: 1,
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

	var meterValueReq = &MeterValuesRequest{
		ConnectorId:   -1,
		TransactionId: -1,
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

	var meterValueReq = &MeterValuesRequest{
		ConnectorId:   1,
		TransactionId: -1,
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
func Test_MeterValueRequest(t *testing.T) {
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

func Test_AuthorizeRequest(t *testing.T) {
	testStatusAuthorizeRequest(AuthorizeRequestSuccess(), t)
	// testStatusAuthorizeRequest(AuthorizeRequestFail(), t)//error
}

func AuthorizeResponseSuccess() *AuthorizeResponse {

	return &AuthorizeResponse{
		IdTagInfo: IdTagInfo{
			ExpiryDate:  "2021-12-13",
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

func Test_AuthorizeResponse(t *testing.T) {
	testStatusAuthorizeResponse(AuthorizeResponseSuccess(), t)
	//testStatusAuthorizeResponse(AuthorizeResponseFail(), t)//error
}

//// StartTransaction

func StartTransactionRequestSuccess() *StartTransactionRequest {

	return &StartTransactionRequest{
		ConnectorId:   10,
		IdTag:         IdToken(RandomString(15)),
		MeterStart:    10,
		ReservationId: 10,
		Timestamp:     time.Now().Format(time.RFC3339),
	}

}

func StartTransactionRequestFail() *StartTransactionRequest {

	return &StartTransactionRequest{
		ConnectorId:   -1,
		IdTag:         IdToken(RandomString(25)),
		MeterStart:    -1,
		ReservationId: 10,
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
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func Test_StartTransactionRequest(t *testing.T) {
	testStartTransactionRequest(StartTransactionRequestSuccess(), t)
	//testStartTransactionRequest(StartTransactionRequestFail(), t)//error
}

func StartTransactionResponseSuccess() *StartTransactionResponse {

	return &StartTransactionResponse{
		IdTagInfo: IdTagInfo{
			ExpiryDate:  "2021-12-13",
			ParentIdTag: IdToken(RandomString(15)),
			Status:      authAccepted,
		},
		TransactionId: 1,
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

func Test_StartTransactionResponse(t *testing.T) {
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

	return &StopTransactionRequest{
		IdTag:           IdToken(RandomString(15)),
		MeterStop:       10,
		Timestamp:       time.Now().Format(time.RFC3339),
		TransactionId:   10,
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
	return &StopTransactionRequest{
		IdTag:     IdToken(RandomString(21)),
		MeterStop: -1,
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

func Test_StopTransactionRequest(t *testing.T) {
	testStopTransactionRequest(StopTransactionRequestSuccess(), t)
	//testStopTransactionRequest(StopTransactionRequestFail(), t)//error
}

func StopTransactionResponseSuccess() *StopTransactionResponse {

	return &StopTransactionResponse{
		IdTagInfo: IdTagInfo{
			ExpiryDate:  "2021-12-13",
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

func Test_StopTransactionResponse(t *testing.T) {
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

func Test_ChangeConfigurationRequest(t *testing.T) {
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

func Test_changeConfigurationResponse(t *testing.T) {
	// testChangeConfigurationResponse(ChangeConfigurationResponseSuccess(), t)
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

func Test_DataTransferRequest(t *testing.T) {
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

func Test_DataTransferResponse(t *testing.T) {
	testDataTransferResponse(DataTransferResponseSuccess(), t)
	// testDataTransferResponse(DataTransferResponseFail(), t) //error
}

/**************RemoteStartTransaction******************/

func RemoteStartTransactionRequestSuccess() *RemoteStartTransactionRequest {
	return &RemoteStartTransactionRequest{
		ConnectorId: 10,
		IdTag:       IdToken(RandomString(18)),
		ChargingProfile: ChargingProfile{
			ChargingProfiled:       10,
			TransactionId:          10,
			StackLevel:             10,
			ChargingProfilePurpose: chargePointMaxProfile,
			ChargingProfileKind:    absolute,
			RecurrencyKind:         daily,
			ValidFrom:              time.Now().Format(ISO8601),
			ValidTo:                time.Now().Format(ISO8601),
			ChargingSchedule: ChargingSchedule{
				Duration:         10,
				StartSchedule:    RandomString(10),
				ChargingRateUnit: uintTypeA,
				ChargingSchedulePeriod: []ChargingSchedulePeriod{
					ChargingSchedulePeriod{
						StartPeriod:  10,
						Limit:        10,
						NumberPhases: 10,
					},
				},
				MinChargingRate: 114.1,
			},
		},
	}
}

func RemoteStartTransactionRequestFail() *RemoteStartTransactionRequest {

	return &RemoteStartTransactionRequest{
		ConnectorId: -1,
		IdTag:       IdToken(RandomString(21)),
		ChargingProfile: ChargingProfile{
			// ChargingProfiled:10,
			//TransactionId:10,
			StackLevel:             -1,
			ChargingProfilePurpose: ChargingProfilePurposeType(RandomString(10)),
			ChargingProfileKind:    ChargingProfileKindType(RandomString(10)),
			RecurrencyKind:         RecurrencyKindType(RandomString(10)),
			//ValidFrom:RandomString(10),
			//ValidTo:RandomString(10),
			ChargingSchedule: ChargingSchedule{
				Duration: -1,
				//StartSchedule:RandomString(10),
				ChargingRateUnit: ChargingRateUnitType(RandomString(10)),
				ChargingSchedulePeriod: []ChargingSchedulePeriod{
					ChargingSchedulePeriod{
						StartPeriod:  -1,
						Limit:        -1,
						NumberPhases: -1,
					},
				},
				MinChargingRate: 114.1,
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

func Test_RemoteStartTransactionRequest(t *testing.T) {
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

func Test_RemoteStartTransactionResponse(t *testing.T) {
	testRemoteStartTransactionResponse(RemoteStartTransactionrResponseSuccess(), t)
	// testRemoteStartTransactionResponse(RemoteStartTransactionResponseFail(), t) //error
}

/**************RemoteStopTransaction******************/

func RemoteStopTransactionRequestSuccess() *RemoteStopTransactionRequest {
	return &RemoteStopTransactionRequest{
		TransactionId: 10,
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

func Test_RemoteStopTransactionRequest(t *testing.T) {
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

func Test_RemoteStopTransactionResponse(t *testing.T) {
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

func Test_ResetRequest(t *testing.T) {
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

func Test_ResetResponse(t *testing.T) {
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

func Test_UnlockConnectorRequest(t *testing.T) {
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

func Test_UnlockConnectorResponse(t *testing.T) {
	testUnlockConnectorResponse(UnlockConnectorResponseSuccess(), t)
	// testUnlockConnectorResponse(UnlockConnectorResponseFail(), t) //error
}

/**************************SetChargingProfile*******************************/
func SetChargingProfileSuccess() *SetChargingProfileRequest {
	return &SetChargingProfileRequest{
		ConnectorId: 10,
		ChargingProfile: ChargingProfile{
			ChargingProfiled:       10,
			TransactionId:          10,
			StackLevel:             10,
			ChargingProfilePurpose: chargePointMaxProfile,
			ChargingProfileKind:    absolute,
			RecurrencyKind:         daily,
			ValidFrom:              time.Now().Format(ISO8601),
			ValidTo:                time.Now().Format(ISO8601),
			ChargingSchedule: ChargingSchedule{
				Duration:         10,
				StartSchedule:    RandomString(10),
				ChargingRateUnit: uintTypeA,
				ChargingSchedulePeriod: []ChargingSchedulePeriod{
					ChargingSchedulePeriod{
						StartPeriod:  10,
						Limit:        10,
						NumberPhases: 10,
					},
				},
				MinChargingRate: 114.1,
			},
		},
	}
}

func SetChargingProfileFail() *SetChargingProfileRequest {

	return &SetChargingProfileRequest{
		ConnectorId: -1,
		ChargingProfile: ChargingProfile{
			// ChargingProfiled:10,
			//TransactionId:10,
			StackLevel:             -1,
			ChargingProfilePurpose: ChargingProfilePurposeType(RandomString(10)),
			ChargingProfileKind:    ChargingProfileKindType(RandomString(10)),
			RecurrencyKind:         RecurrencyKindType(RandomString(10)),
			//ValidFrom:RandomString(10),
			//ValidTo:RandomString(10),
			ChargingSchedule: ChargingSchedule{
				Duration: -1,
				//StartSchedule:RandomString(10),
				ChargingRateUnit: ChargingRateUnitType(RandomString(10)),
				ChargingSchedulePeriod: []ChargingSchedulePeriod{
					ChargingSchedulePeriod{
						StartPeriod:  -1,
						Limit:        -1,
						NumberPhases: -1,
					},
				},
				MinChargingRate: 114.1,
			},
		},
	}
}

func testSetChargingProfileRequest(call *SetChargingProfileRequest, t *testing.T) {
	SetChargingProfile_reqbyte, err := json.Marshal(call)
	if err != nil {
		t.Error(err)
	}
	var tmp = &SetChargingProfileRequest{}
	err = json.Unmarshal(SetChargingProfile_reqbyte, tmp)
	if err != nil {
		t.Error(err)
		return
	}
	err = Validate.Struct(tmp)
	if err != nil {
		t.Error(err)
	}
}

func Test_SetChargingProfileRequest(t *testing.T) {
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

func Test_SetChargingProfileResponse(t *testing.T) {
	testSetChargingProfileResponse(SetChargingProfileResponseSuccess(), t)
	// testSetChargingProfileResponse(SetChargingProfileResponseFail(), t) //error
}

func testOcppSuccess(t *testing.T) {
	if ocpptrait, ok := OCPP16M.GetTraitAction("BootNotification"); !ok {
		return
	} else {
		reqType := ocpptrait.RequestType() //根据type构造实例， 调用对应的方法
		reqByte, err := json.Marshal(BootNotificationRequestSuccess())
		if err != nil {
			t.Error(err)
			return
		}
		request := reflect.New(reqType).Interface()
		err = json.Unmarshal(reqByte, request)
		if err != nil {
			t.Error(err)
			return
		}
		err = Validate.Struct(request)
		if err != nil {
			t.Error(err)
			return
		}

	}
}

func testOcppFail(t *testing.T) {
	if ocpptrait, ok := OCPP16M.GetTraitAction("BootNotification"); !ok {
		t.Log(ocpptrait)
		t.Log("not support")
		return
	} else {
		reqType := ocpptrait.RequestType() //根据type构造实例， 调用对应的方法
		reqByte, err := json.Marshal(BootNotificationRequestFail())
		if err != nil {
			t.Error(err)
			return
		}
		request := reflect.New(reqType).Interface()
		err = json.Unmarshal(reqByte, request)
		if err != nil {
			t.Error(err)
			return
		}
		err = Validate.Struct(request)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func Test_OccpTrait(t *testing.T) {
	testOcppSuccess(t)
	// testOcppFail(t)
}

func parseMessage(wsmsg []byte) ([]interface{}, error) {
	var fields []interface{}
	err := json.Unmarshal(wsmsg, &fields)
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func Test_OccpCall(t *testing.T) {

	call := &Call{
		MessageTypeID: 2,
		UniqueID:      "uniqueid",
		Action:        "Authorize",
		// Request:       json.RawMessage(`{"chargePointVendor": "VendorX", "chargePointModel": "SingleSocketCharger"}`),
		//Request: json.RawMessage(`{"chargePointVendsddddddddddddsdsdsdsdsdsdor": "VendorX", "chargePointModel": "SingleSocketCharger"}`),
	}
	call_byte, err := call.MarshalJSON()
	if err != nil {
		t.Error(err)
		return
	}

	fields, err := parseMessage(call_byte)
	if err != nil {
		t.Error(err)
		return
	}
	if ocpptrait, ok := OCPP16M.GetTraitAction("BootNotification"); !ok {
		return
	} else {
		reqType := ocpptrait.RequestType() //根据type构造实例， 调用对应的方法
		reqByte, err := json.Marshal(fields[3])
		if err != nil {
			t.Error(err)
			return
		}
		request := reflect.New(reqType).Interface()
		err = json.Unmarshal(reqByte, &request)
		if err != nil {
			t.Error(err, string(reqByte[0:5]))
			return
		}
		err = Validate.Struct(request)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
