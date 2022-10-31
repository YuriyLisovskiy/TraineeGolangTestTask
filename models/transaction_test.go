package models

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"
)

func TestTransaction_ToCsvRow(t *testing.T) {
	transaction := Transaction{
		Id:                 1,
		RequestId:          20020,
		TerminalId:         3506,
		PartnerObjectId:    1111,
		AmountTotal:        1.00,
		AmountOriginal:     1.00,
		CommissionPS:       0.00,
		CommissionClient:   0.00,
		CommissionProvider: 0.00,
		DateInput:          time.Date(2022, time.August, 12, 11, 25, 27, 0, time.UTC),
		DatePost:           time.Date(2022, time.August, 12, 14, 25, 27, 0, time.UTC),
		Status:             "accepted",
		PaymentType:        "cash",
		PaymentNumber:      "PS16698205",
		ServiceId:          13980,
		Service:            "Поповнення карток",
		PayeeId:            14232155,
		PayeeName:          "pumb",
		PayeeBankMfo:       254751,
		PayeeBankAccount:   "UA713451373919523",
		PaymentNarrative:   "Перерахування коштів згідно договору про надання послуг А11/27122 від 19.11.2020 р.",
	}
	expected := "1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів згідно договору про надання послуг А11/27122 від 19.11.2020 р."
	actual := transaction.ToCsvRow()
	if expected != actual {
		t.Errorf("\n%v\n!=\n%v", expected, actual)
	}
}

func TestNewTransactionFromCSVRow_Success(t *testing.T) {
	csvData := "1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів згідно договору про надання послуг А11/27122 від 19.11.2020 р."
	transaction, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err != nil {
		t.Error(err)
	}

	expectedId := uint64(1)
	if transaction.Id != expectedId {
		t.Errorf("%v != %v", transaction.Id, expectedId)
	}

	expectedRequestId := uint64(20020)
	if transaction.RequestId != expectedRequestId {
		t.Errorf("%v != %v", transaction.RequestId, expectedRequestId)
	}

	expectedTerminalId := uint64(3506)
	if transaction.TerminalId != expectedTerminalId {
		t.Errorf("%v != %v", transaction.TerminalId, expectedTerminalId)
	}

	expectedPartnerObjectId := uint16(1111)
	if transaction.PartnerObjectId != expectedPartnerObjectId {
		t.Errorf("%v != %v", transaction.PartnerObjectId, expectedPartnerObjectId)
	}

	expectedAmountTotal := float32(1.00)
	if transaction.AmountTotal != expectedAmountTotal {
		t.Errorf("%v != %v", transaction.AmountTotal, expectedAmountTotal)
	}

	expectedAmountOriginal := float32(1.00)
	if transaction.AmountOriginal != expectedAmountOriginal {
		t.Errorf("%v != %v", transaction.AmountOriginal, expectedAmountOriginal)
	}

	expectedCommissionPS := float32(0.00)
	if transaction.CommissionPS != expectedCommissionPS {
		t.Errorf("%v != %v", transaction.CommissionPS, expectedCommissionPS)
	}

	expectedCommissionClient := float32(0.00)
	if transaction.CommissionClient != expectedCommissionClient {
		t.Errorf("%v != %v", transaction.CommissionClient, expectedCommissionClient)
	}

	expectedCommissionProvider := float32(0.00)
	if transaction.CommissionProvider != expectedCommissionProvider {
		t.Errorf("%v != %v", transaction.CommissionProvider, expectedCommissionProvider)
	}

	expectedDateInput, _ := time.Parse(TimeLayout, "2022-08-12 11:25:27")
	if transaction.DateInput != expectedDateInput {
		t.Errorf("%v != %v", transaction.DateInput, expectedDateInput)
	}

	expectedDatePost, _ := time.Parse(TimeLayout, "2022-08-12 14:25:27")
	if transaction.DatePost != expectedDatePost {
		t.Errorf("%v != %v", transaction.DatePost, expectedDatePost)
	}

	expectedStatus := ACCEPTED
	if transaction.Status != expectedStatus {
		t.Errorf("%v != %v", transaction.Status, expectedStatus)
	}

	expectedPaymentType := CASH
	if transaction.PaymentType != expectedPaymentType {
		t.Errorf("%v != %v", transaction.PaymentType, expectedPaymentType)
	}

	expectedPaymentNumber := "PS16698205"
	if transaction.PaymentNumber != expectedPaymentNumber {
		t.Errorf("%v != %v", transaction.PaymentNumber, expectedPaymentNumber)
	}

	expectedServiceId := uint64(13980)
	if transaction.ServiceId != expectedServiceId {
		t.Errorf("%v != %v", transaction.ServiceId, expectedServiceId)
	}

	expectedService := "Поповнення карток"
	if transaction.Service != expectedService {
		t.Errorf("%v != %v", transaction.Service, expectedService)
	}

	expectedPayeeId := uint64(14232155)
	if transaction.PayeeId != expectedPayeeId {
		t.Errorf("%v != %v", transaction.PayeeId, expectedPayeeId)
	}

	expectedPayeeName := "pumb"
	if transaction.PayeeName != expectedPayeeName {
		t.Errorf("%v != %v", transaction.PayeeName, expectedPayeeName)
	}

	expectedPayeeBankMfo := uint32(254751)
	if transaction.PayeeBankMfo != expectedPayeeBankMfo {
		t.Errorf("%v != %v", transaction.PayeeBankMfo, expectedPayeeBankMfo)
	}

	expectedPayeeBankAccount := "UA713451373919523"
	if transaction.PayeeBankAccount != expectedPayeeBankAccount {
		t.Errorf("%v != %v", transaction.PayeeBankAccount, expectedPayeeBankAccount)
	}

	expectedPaymentNarrative := "Перерахування коштів згідно договору про надання послуг А11/27122 від 19.11.2020 р."
	if transaction.PaymentNarrative != expectedPaymentNarrative {
		t.Errorf("%v != %v", transaction.PaymentNarrative, expectedPaymentNarrative)
	}
}

func TestNewTransactionFromCSVRow_MissingFieldInCSVRow(t *testing.T) {
	csvData := "1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523"
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidId(t *testing.T) {
	csvData := "-1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів"
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidRequestId(t *testing.T) {
	csvData := "1,-1,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів"
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidTerminalId(t *testing.T) {
	csvData := "1,20020,-1,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів"
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidPartnerObjectId(t *testing.T) {
	csvData := "1,20020,3506,65536,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів"
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidAmountTotal(t *testing.T) {
	csvData := fmt.Sprintf(
		"1,20020,3506,1111,%f,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів",
		math.MaxFloat32+math.MaxFloat32*0.0001,
	)
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidAmountOriginal(t *testing.T) {
	csvData := fmt.Sprintf(
		"1,20020,3506,1111,1.00,%f,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів",
		math.MaxFloat32+math.MaxFloat32*0.0001,
	)
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidCommissionPS(t *testing.T) {
	csvData := fmt.Sprintf(
		"1,20020,3506,1111,1.00,1.00,%f,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів",
		math.MaxFloat32+math.MaxFloat32*0.0001,
	)
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidCommissionClient(t *testing.T) {
	csvData := fmt.Sprintf(
		"1,20020,3506,1111,1.00,1.00,0.00,%f,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів",
		math.MaxFloat32+math.MaxFloat32*0.0001,
	)
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidCommissionProvider(t *testing.T) {
	csvData := fmt.Sprintf(
		"1,20020,3506,1111,1.00,1.00,0.00,0.00,%f,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів",
		math.MaxFloat32+math.MaxFloat32*0.0001,
	)
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidDateInput(t *testing.T) {
	csvData := "1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів"
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidDatePost(t *testing.T) {
	csvData := "1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-25-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів"
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidStatus(t *testing.T) {
	csvData := "1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,rejected,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів"
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidPaymentType(t *testing.T) {
	csvData := "1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,credit,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів"
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidServiceId(t *testing.T) {
	csvData := "1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,-1,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів"
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidPayeeId(t *testing.T) {
	csvData := "1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,-1,pumb,254751,UA713451373919523,Перерахування коштів"
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func TestNewTransactionFromCSVRow_InvalidPayeeBankMfo(t *testing.T) {
	csvData := fmt.Sprintf(
		"1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,%f,UA713451373919523,Перерахування коштів",
		math.MaxFloat32+math.MaxFloat32*0.0001,
	)
	_, err := NewTransactionFromCSVRow(strings.Split(csvData, ","))
	if err == nil {
		t.Error("error is nil")
	}
}

func Test_parseFloat32_Parsed(t *testing.T) {
	expected := float32(2345678.987654)
	actual, err := parseFloat32(fmt.Sprintf("%f", expected))
	if err != nil {
		t.Error(err)
	}

	if expected != actual {
		t.Errorf("%f != %f", expected, actual)
	}
}

func Test_parseFloat32_OverflowError(t *testing.T) {
	_, err := parseFloat32(fmt.Sprintf("%f", math.MaxFloat32+math.MaxFloat32*0.0001))
	if err == nil {
		t.Error("error is nil")
	}
}

func Test_parseUint64_PositiveIntParsed(t *testing.T) {
	expected := uint64(math.MaxUint64)
	actual, err := parseUint64(fmt.Sprintf("%d", expected))
	if err != nil {
		t.Error(err)
	}

	if expected != actual {
		t.Errorf("%d != %d", expected, actual)
	}
}

func Test_parseUint64_ZeroParsed(t *testing.T) {
	expected := uint64(0)
	actual, err := parseUint64(fmt.Sprintf("%d", expected))
	if err != nil {
		t.Error(err)
	}

	if expected != actual {
		t.Errorf("%d != %d", expected, actual)
	}
}

func Test_parseUint64_NegativeIntError(t *testing.T) {
	_, err := parseUint64(fmt.Sprintf("%d", -1))
	if err == nil {
		t.Error("error is nil")
	}
}
