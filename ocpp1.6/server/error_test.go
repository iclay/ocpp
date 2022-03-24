package server

// import (
// 	"bytes"
// 	"crypto/rand"
// 	"fmt"
// 	"math/big"
// 	"ocpp16/protocol"
// 	"testing"
// 	// "time"
// 	"errors"
// )

// func RandomString(len int) string {
// 	var numbers = []byte{1, 2, 3, 4, 5, 7, 8, 9}
// 	var container string
// 	length := bytes.NewReader(numbers).Len()

// 	for i := 1; i <= len; i++ {
// 		random, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
// 		if err != nil {

// 		}
// 		container += fmt.Sprintf("%d", numbers[random.Int64()])
// 	}
// 	return container
// }

// func callFaildDueUniqueID() *proto.Call {
// 	return &proto.Call{
// 		MessageTypeID: 2,
// 		UniqueID:      RandomString(36), //length>36
// 		Action:        "Authorize",
// 		// Request:       proto.Request{},
// 	}
// }
// func Test_call(t* testing.T) {
// 	validate := proto.Validate
//     err := validate.Struct(callFaildDueUniqueID())
// 	if err != nil {
// 		err = errors.New("new error")
// 		t.Errorf("%+v",checkValidatorError(err,"auth"))
// 	}
// }

// func BootNotificationRequestFail() *proto.BootNotificationRequest {

// 	return &proto.BootNotificationRequest{
// 		ChargePointVendor:       RandomString(15),
// 		ChargePointModel:        RandomString(25),
// 		ChargePointSerialNumber: RandomString(30),
// 		ChargeBoxSerialNumber:   RandomString(30),
// 		FirmwareVersion:         RandomString(55),
// 		Iccid:                   RandomString(25),
// 		Imsi:                    RandomString(25),
// 		MeterType:               RandomString(50),
// 		MeterSerialNumber:       RandomString(50),
// 	}
// }

// func Test_BootNotificationRequestFail(t *testing.T) {
// 	validate := proto.Validate
// 	err := validate.Struct(BootNotificationRequestFail())
// 	if err != nil {
// 		t.Errorf("%+v", checkValidatorError(err, "auth"))
// 	}
// }

// func StatusNotificationRequestSuccess() *proto.StatusNotificationRequest {

// 	return &proto.StatusNotificationRequest{
// 		ConnectorId: 15,
// 		ErrorCode:   "ConnectorLockFailure",
// 		Info:        RandomString(40),
// 		Status:      "Available",
// 		// Timestamp:       time.Now().Format(time.RFC3339),
// 		// Timestamp:       time.Now().Format(time.RFC3339),
// 		Timestamp:       "2010100",
// 		VendorId:        RandomString(240),
// 		VendorErrorCode: RandomString(40),
// 	}

// }

// func Test_StatusNotificationRequestSuccess(t* testing.T) {
// 	validate := proto.Validate
//     err := validate.Struct(StatusNotificationRequestSuccess())
// 	if err != nil {
// 		t.Errorf("%+v",checkValidatorError(err,"auth"))
// 	}
// }

// func MeterValueRequestFailEmptyMeterValue() *proto.MeterValuesRequest {

// 	var meterValueReq = &proto.MeterValuesRequest{
// 		ConnectorId:   1,
// 		TransactionId: -1,
// 		MeterValue:    []proto.MeterValue{},
// 	}
// 	return meterValueReq
// }

// func Test_MeterValueRequestFail(t* testing.T) {
// 	validate := proto.Validate
//     err := validate.Struct(MeterValueRequestFailEmptyMeterValue())
// 	if err != nil {
// 		t.Errorf("%+v",checkValidatorError(err,"auth"))
// 	}
// }
