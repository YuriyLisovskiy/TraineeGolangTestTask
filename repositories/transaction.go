package repositories

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"TraineeGolangTestTask/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(model *models.Transaction) error
	CreateBatch(dbTransaction func(TransactionRepository) error) error
	Filter(filters []TransactionFilter, page, pageSize int) []models.Transaction
	ForEach(filters []TransactionFilter, apply func(model *models.Transaction) error) error

	NewFilterBuilder() TransactionFilterBuilder
}

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{db: db}
}

func (tr *TransactionRepositoryImpl) Create(model *models.Transaction) error {
	tx := tr.db.Create(model)
	return tx.Error
}

func (tr *TransactionRepositoryImpl) CreateBatch(dbTransaction func(TransactionRepository) error) error {
	return tr.db.Transaction(
		func(tx *gorm.DB) error {
			return dbTransaction(tr)
		},
	)
}

// Filter returns paginated result that is a list of transactions with applied filters.
// If page or pageSize is less than or equals to zero, pagination is ignored.
func (tr *TransactionRepositoryImpl) Filter(filters []TransactionFilter, page, pageSize int) []models.Transaction {
	tx := tr.db.Model(&models.Transaction{})
	applyFilters(tx, filters)
	var transactions []models.Transaction
	if page > 0 && pageSize > 0 {
		tx.Limit(pageSize).Offset((page - 1) * pageSize)
	}

	tx.Find(&transactions)
	return transactions
}

func (tr *TransactionRepositoryImpl) ForEach(
	filters []TransactionFilter,
	apply func(model *models.Transaction) error,
) error {
	tx := tr.db.Model(&models.Transaction{})
	applyFilters(tx, filters)
	rows, err := tx.Rows()
	if err != nil {
		return err
	}

	for rows.Next() {
		transaction := models.Transaction{}
		err = tr.db.ScanRows(rows, &transaction)
		if err != nil {
			return err
		}

		err = apply(&transaction)
		if err != nil {
			return err
		}
	}

	return nil
}

func applyFilters(tx *gorm.DB, filters []TransactionFilter) {
	for _, filter := range filters {
		filter(tx)
	}
}

func (tr *TransactionRepositoryImpl) NewFilterBuilder() TransactionFilterBuilder {
	return &TransactionFilterBuilderImpl{}
}

type TransactionFilter func(tx *gorm.DB)

type TransactionFilterBuilder interface {
	AddTransactionId(value string) error
	AddTerminalIds(values []string) error
	AddStatus(value string) error
	AddPaymentType(value string) error
	AddDatePostRange(valueFrom, valueTo string) error
	AddPaymentNarrative(value string) error
	GetFilters() []TransactionFilter
}

type TransactionFilterBuilderImpl struct {
	filters []TransactionFilter
}

func (tf *TransactionFilterBuilderImpl) AddTransactionId(value string) error {
	if value != "" {
		transactionId, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}

		tf.filters = append(
			tf.filters, func(tx *gorm.DB) {
				tx.Where("id = ?", transactionId)
			},
		)
	}

	return nil
}

func (tf *TransactionFilterBuilderImpl) AddTerminalIds(values []string) error {
	if len(values) > 0 {
		var ids []uint64
		for _, stringId := range values {
			id, err := strconv.ParseUint(stringId, 10, 64)
			if err != nil {
				return err
			}

			ids = append(ids, id)
		}

		tf.filters = append(
			tf.filters, func(tx *gorm.DB) {
				tx.Where("(terminal_id) IN ?", ids)
			},
		)
	}

	return nil
}

func (tf *TransactionFilterBuilderImpl) AddStatus(value string) error {
	if value != "" {
		switch status := models.StatusType(value); status {
		case models.ACCEPTED, models.DECLINED:
			tf.filters = append(
				tf.filters, func(tx *gorm.DB) {
					tx.Where("status = ?", status)
				},
			)
		default:
			return fmt.Errorf(
				"value of \"status\" parameter should be either \"%v\" or \"%v\"",
				models.ACCEPTED,
				models.DECLINED,
			)
		}
	}

	return nil
}

func (tf *TransactionFilterBuilderImpl) AddPaymentType(value string) error {
	if value != "" {
		switch paymentType := models.PaymentTypeType(value); paymentType {
		case models.CASH, models.CARD:
			tf.filters = append(
				tf.filters, func(tx *gorm.DB) {
					tx.Where("payment_type = ?", paymentType)
				},
			)
		default:
			return fmt.Errorf(
				"value of \"payment_type\" parameter should be either \"%v\" or \"%v\"",
				models.CASH,
				models.CARD,
			)
		}
	}

	return nil
}

func (tf *TransactionFilterBuilderImpl) AddDatePostRange(valueFrom, valueTo string) error {
	if valueFrom != "" {
		if valueTo == "" {
			return errors.New("parameter \"date_post_to\" is required when using \"date_post_from\"")
		}

		from, err := time.Parse(models.TimeLayout, valueFrom)
		if err != nil {
			return err
		}

		to, err := time.Parse(models.TimeLayout, valueTo)
		if err != nil {
			return err
		}

		tf.filters = append(
			tf.filters, func(tx *gorm.DB) {
				tx.Where("date_post BETWEEN ? AND ?", from, to)
			},
		)
	}

	return nil
}

func (tf *TransactionFilterBuilderImpl) AddPaymentNarrative(value string) error {
	if value != "" {
		tf.filters = append(
			tf.filters, func(tx *gorm.DB) {
				tx.Where("payment_narrative LIKE ?", fmt.Sprintf("%%%s%%", value))
			},
		)
	}

	return nil
}

func (tf *TransactionFilterBuilderImpl) GetFilters() []TransactionFilter {
	return tf.filters
}
