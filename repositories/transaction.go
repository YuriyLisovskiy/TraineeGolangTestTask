package repositories

import (
	"fmt"
	"time"

	"TraineeGolangTestTask/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(model *models.Transaction) error
	CreateBatch(dbTransaction func(TransactionRepository) error) error
	Filter(filters []TransactionFilter, page, pageSize int) []models.Transaction
	ForEach(filters []TransactionFilter, apply func(model *models.Transaction) error) error
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
	for _, filter := range filters {
		filter(tx)
	}

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
	for _, filter := range filters {
		filter(tx)
	}

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

type TransactionFilter func(tx *gorm.DB)

func FilterByTransactionId(id uint64) TransactionFilter {
	return func(tx *gorm.DB) {
		tx.Where("id = ?", id)
	}
}

func FilterByTerminalId(ids []uint64) TransactionFilter {
	return func(tx *gorm.DB) {
		tx.Where("(terminal_id) IN ?", ids)
	}
}

func FilterByStatus(status models.StatusType) TransactionFilter {
	return func(tx *gorm.DB) {
		tx.Where("status = ?", status)
	}
}

func FilterByPaymentType(paymentType models.PaymentTypeType) TransactionFilter {
	return func(tx *gorm.DB) {
		tx.Where("payment_type = ?", paymentType)
	}
}

func FilterByDatePostTimeRange(from, to time.Time) TransactionFilter {
	return func(tx *gorm.DB) {
		tx.Where("date_post BETWEEN ? AND ?", from, to)
	}
}

func ContainsTextInPaymentNarrative(text string) TransactionFilter {
	return func(tx *gorm.DB) {
		tx.Where("payment_narrative LIKE ?", fmt.Sprintf("%%%s%%%", text))
	}
}
