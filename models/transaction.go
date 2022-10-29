package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	numberOfTransactionFields = 21

	TimeLayout = "2006-01-02 15:04:05"
)

type StatusType string

const (
	ACCEPTED StatusType = "accepted"
	DECLINED StatusType = "declined"
)

func (st *StatusType) Scan(value interface{}) error {
	*st = StatusType(value.(string))
	return nil
}

func (st *StatusType) Value() (driver.Value, error) {
	return string(*st), nil
}

type PaymentTypeType string

const (
	CASH PaymentTypeType = "cash"
	CARD PaymentTypeType = "card"
)

func (pt *PaymentTypeType) Scan(value interface{}) error {
	*pt = PaymentTypeType(value.(string))
	return nil
}

func (pt *PaymentTypeType) Value() (driver.Value, error) {
	return string(*pt), nil
}

type Transaction struct {
	Id                 uint64          `gorm:"primaryKey" json:"transaction_id"`
	RequestId          uint64          `json:"request_id"`
	TerminalId         uint64          `json:"terminal_id"`
	PartnerObjectId    uint16          `json:"partner_object_id"`
	AmountTotal        float32         `json:"amount_total"`
	AmountOriginal     float32         `json:"amount_original"`
	CommissionPS       float32         `json:"commission_ps"`
	CommissionClient   float32         `json:"commission_client"`
	CommissionProvider float32         `json:"commission_provider"`
	DateInput          time.Time       `json:"date_input"` // YYYY-MM-DD HH:MM:SS
	DatePost           time.Time       `json:"date_post"`  // YYYY-MM-DD HH:MM:SS
	Status             StatusType      `gorm:"type:rest_api_schema.status_type" json:"status"`
	PaymentType        PaymentTypeType `gorm:"type:rest_api_schema.payment_type_type" json:"payment_type"`
	PaymentNumber      string          `gorm:"size:10;check:payment_number ~ '[A-Z]{2}[0-9]{8}'" json:"payment_number"`
	ServiceId          uint64          `json:"service_id"`
	Service            string          `json:"service"`
	PayeeId            uint64          `json:"payee_id"`
	PayeeName          string          `json:"payee_name"`
	PayeeBankMfo       uint32          `json:"payee_bank_mfo"`
	PayeeBankAccount   string          `gorm:"size:17;check:payee_bank_account ~ '[A-Z]{2}[0-9]{15}'" json:"payee_bank_account"`
	PaymentNarrative   string          `json:"payment_narrative"`
}

func (t *Transaction) ToCsvRow() string {
	return fmt.Sprintf(
		"%d,%d,%d,%d,%.2f,%.2f,%.2f,%.2f,%.2f,%s,%s,%s,%s,%s,%d,%s,%d,%s,%d,%s,%s",
		t.Id,
		t.RequestId,
		t.TerminalId,
		t.PartnerObjectId,
		t.AmountTotal,
		t.AmountOriginal,
		t.CommissionPS,
		t.CommissionClient,
		t.CommissionProvider,
		t.DateInput.Format(TimeLayout),
		t.DatePost.Format(TimeLayout),
		t.Status,
		t.PaymentType,
		t.PaymentNumber,
		t.ServiceId,
		t.Service,
		t.PayeeId,
		t.PayeeName,
		t.PayeeBankMfo,
		t.PayeeBankAccount,
		t.PaymentNarrative,
	)
}

func NewTransactionFromCSVRow(row string) (*Transaction, error) {
	fields := strings.SplitN(row, ",", numberOfTransactionFields)
	fieldsLen := len(fields)
	if fieldsLen != numberOfTransactionFields {
		return nil, fmt.Errorf(
			"invalid number of transaction fields: %d required, %d got",
			numberOfTransactionFields,
			fieldsLen,
		)
	}

	var err error
	transaction := &Transaction{}
	transaction.Id, err = parseUint64(fields[0])
	if err != nil {
		return nil, err
	}

	transaction.RequestId, err = parseUint64(fields[1])
	if err != nil {
		return nil, err
	}

	transaction.TerminalId, err = parseUint64(fields[2])
	if err != nil {
		return nil, err
	}

	partnerObjectId, err := strconv.ParseUint(fields[3], 10, 16)
	if err != nil {
		return nil, err
	}

	transaction.PartnerObjectId = uint16(partnerObjectId)
	transaction.AmountTotal, err = parseFloat32(fields[4])
	if err != nil {
		return nil, err
	}

	transaction.AmountOriginal, err = parseFloat32(fields[5])
	if err != nil {
		return nil, err
	}

	transaction.CommissionPS, err = parseFloat32(fields[6])
	if err != nil {
		return nil, err
	}

	transaction.CommissionClient, err = parseFloat32(fields[7])
	if err != nil {
		return nil, err
	}

	transaction.CommissionProvider, err = parseFloat32(fields[8])
	if err != nil {
		return nil, err
	}

	transaction.DateInput, err = time.Parse(TimeLayout, fields[9])
	if err != nil {
		return nil, err
	}

	transaction.DatePost, err = time.Parse(TimeLayout, fields[10])
	if err != nil {
		return nil, err
	}

	switch status := StatusType(fields[11]); status {
	case ACCEPTED, DECLINED:
		transaction.Status = status
	default:
		return nil, fmt.Errorf("invalid status, expected %s or %s, received %v", ACCEPTED, DECLINED, status)
	}

	switch paymentType := PaymentTypeType(fields[12]); paymentType {
	case CASH, CARD:
		transaction.PaymentType = paymentType
	default:
		return nil, fmt.Errorf("invalid payment type, expected %s or %s, received %v", CASH, CARD, paymentType)
	}

	transaction.PaymentNumber = fields[13]
	transaction.ServiceId, err = parseUint64(fields[14])
	if err != nil {
		return nil, err
	}

	transaction.Service = fields[15]
	transaction.PayeeId, err = parseUint64(fields[16])
	if err != nil {
		return nil, err
	}

	transaction.PayeeName = fields[17]
	payeeBandMfo, err := strconv.ParseInt(fields[18], 10, 32)
	if err != nil {
		return nil, err
	}

	transaction.PayeeBankMfo = uint32(payeeBandMfo)
	transaction.PayeeBankAccount = fields[19]
	transaction.PaymentNarrative = fields[20]
	return transaction, nil
}

func parseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func parseFloat32(s string) (float32, error) {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}

	return float32(v), nil
}
