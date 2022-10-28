package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	numberOfTransactionFields = 21

	TimeLayout = "2006-01-02 15:04:05"
)

type Transaction struct {
	Id                 uint64 `gorm:"primaryKey"`
	RequestId          uint64
	TerminalId         uint64
	PartnerObjectId    uint16
	AmountTotal        float32
	AmountOriginal     float32
	CommissionPS       float32
	CommissionClient   float32
	CommissionProvider float32
	DateInput          time.Time // YYYY-MM-DD HH:MM:SS
	DatePost           time.Time // YYYY-MM-DD HH:MM:SS
	Status             string    `gorm:"size:8;check:status IN ('accepted', 'declined')"`
	PaymentType        string    `gorm:"size:4;check:payment_type IN ('cash', 'card')"`
	PaymentNumber      string    `gorm:"size:10;check:payment_number ~ 'PS[0-9]{8}'"`
	ServiceId          uint64
	Service            string
	PayeeId            uint64
	PayeeName          string
	PayeeBankMfo       uint32
	PayeeBankAccount   string `gorm:"size:17;check:payee_bank_account ~ 'UA[0-9]{15}'"`
	PaymentNarrative   string
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

	transaction.Status = fields[11]
	transaction.PaymentType = fields[12]
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
