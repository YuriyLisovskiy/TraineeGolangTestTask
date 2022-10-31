package repositories

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"TraineeGolangTestTask/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestTransactionRepositoryImpl_Filter(t *testing.T) {
	db, err := openTestDb()
	if err != nil {
		t.Error(err)
	}

	db.Exec("CREATE SCHEMA rest_api")
	_ = models.MigrateAll(db)

	for _, tr := range testTransactions {
		db.Create(&tr)
	}

	defer func() {
		for _, tr := range testTransactions {
			db.Delete(&tr, "id = ?", tr.Id)
		}
	}()

	repo := &TransactionRepositoryImpl{db: db}

	t.Run(
		"ByTransactionId", func(t *testing.T) {
			SubTestTransactionRepositoryImpl_FilterByTransactionId(t, repo)
		},
	)

	t.Run(
		"ByTerminalId", func(t *testing.T) {
			SubTestTransactionRepositoryImpl_FilterByTerminalId(t, repo)
		},
	)

	t.Run(
		"ByStatus", func(t *testing.T) {
			SubTestTransactionRepositoryImpl_FilterByStatus(t, repo)
		},
	)

	t.Run(
		"ByPaymentType", func(t *testing.T) {
			SubTestTransactionRepositoryImpl_FilterByPaymentType(t, repo)
		},
	)

	t.Run(
		"ByAddDatePostRange", func(t *testing.T) {
			SubTestTransactionRepositoryImpl_FilterByAddDatePostRange(t, repo)
		},
	)

	t.Run(
		"ByPaymentNarrative", func(t *testing.T) {
			SubTestTransactionRepositoryImpl_FilterByPaymentNarrative(t, repo)
		},
	)

	t.Run(
		"ByStatusAndDatePostRange", func(t *testing.T) {
			SubTestTransactionRepositoryImpl_FilterByStatusAndDatePostRange(t, repo)
		},
	)

	t.Run(
		"IgnorePagination_IncorrectPage", func(t *testing.T) {
			SubTestTransactionRepositoryImpl_FilterIgnorePagination_IncorrectPage(t, repo)
		},
	)

	t.Run(
		"IgnorePagination_IncorrectPageSize", func(t *testing.T) {
			SubTestTransactionRepositoryImpl_FilterIgnorePagination_IncorrectPageSize(t, repo)
		},
	)
}

func SubTestTransactionRepositoryImpl_FilterByTransactionId(t *testing.T, repo *TransactionRepositoryImpl) {
	builder := repo.NewFilterBuilder()
	_ = builder.AddTransactionId(fmt.Sprintf("%d", testTransactions[0].Id))
	transactions := repo.Filter(builder.GetFilters(), 1, 10)
	checkFilterSingleResult(t, transactions, 0)
}

func SubTestTransactionRepositoryImpl_FilterByTerminalId(t *testing.T, repo *TransactionRepositoryImpl) {
	builder := repo.NewFilterBuilder()
	_ = builder.AddTerminalIds(
		[]string{
			fmt.Sprintf("%d", testTransactions[0].TerminalId),
			fmt.Sprintf("%d", testTransactions[2].TerminalId),
		},
	)
	transactions := repo.Filter(builder.GetFilters(), 1, 10)
	expectedLen := 2
	actualLen := len(transactions)
	if actualLen != expectedLen {
		t.Errorf("expected len %d, actual len %d", expectedLen, actualLen)
	}

	if transactions[0].Id != testTransactions[0].Id {
		t.Errorf("expected 1st transaction id %d, actual %d", testTransactions[0].Id, transactions[0].Id)
	}

	if transactions[1].Id != testTransactions[2].Id {
		t.Errorf("expected 2nd transaction id %d, actual %d", testTransactions[2].Id, transactions[1].Id)
	}
}

func SubTestTransactionRepositoryImpl_FilterByStatus(t *testing.T, repo *TransactionRepositoryImpl) {
	builder := repo.NewFilterBuilder()
	_ = builder.AddStatus(string(testTransactions[1].Status))
	transactions := repo.Filter(builder.GetFilters(), 1, 10)
	checkFilterSingleResult(t, transactions, 1)
}

func SubTestTransactionRepositoryImpl_FilterByPaymentType(t *testing.T, repo *TransactionRepositoryImpl) {
	builder := repo.NewFilterBuilder()
	_ = builder.AddPaymentType(string(testTransactions[2].PaymentType))
	transactions := repo.Filter(builder.GetFilters(), 1, 10)
	checkFilterSingleResult(t, transactions, 2)
}

func SubTestTransactionRepositoryImpl_FilterByAddDatePostRange(t *testing.T, repo *TransactionRepositoryImpl) {
	builder := repo.NewFilterBuilder()
	_ = builder.AddDatePostRange("2022-08-12 14:25:27", "2022-08-15 13:02:10")
	transactions := repo.Filter(builder.GetFilters(), 1, 10)
	expectedLen := 2
	actualLen := len(transactions)
	if actualLen != expectedLen {
		t.Errorf("expected len %d, actual len %d", expectedLen, actualLen)
	}

	if transactions[0].Id != testTransactions[0].Id {
		t.Errorf("expected 1st transaction id %d, actual %d", testTransactions[0].Id, transactions[0].Id)
	}

	if transactions[1].Id != testTransactions[1].Id {
		t.Errorf("expected 2nd transaction id %d, actual %d", testTransactions[1].Id, transactions[1].Id)
	}
}

func SubTestTransactionRepositoryImpl_FilterByPaymentNarrative(t *testing.T, repo *TransactionRepositoryImpl) {
	builder := repo.NewFilterBuilder()
	_ = builder.AddPaymentNarrative("А11/27123 від 19.11.2020 р.")
	transactions := repo.Filter(builder.GetFilters(), 1, 10)
	checkFilterSingleResult(t, transactions, 1)
}

func SubTestTransactionRepositoryImpl_FilterByStatusAndDatePostRange(
	t *testing.T,
	repo *TransactionRepositoryImpl,
) {
	builder := repo.NewFilterBuilder()
	_ = builder.AddStatus(string(models.DECLINED))
	_ = builder.AddDatePostRange("2022-08-12 14:25:27", "2022-08-15 13:02:10")
	transactions := repo.Filter(builder.GetFilters(), 1, 10)
	checkFilterSingleResult(t, transactions, 1)
}

func SubTestTransactionRepositoryImpl_FilterIgnorePagination_IncorrectPage(
	t *testing.T,
	repo *TransactionRepositoryImpl,
) {
	transactions := repo.Filter([]TransactionFilter{}, -1, 2)
	expectedCount := 3
	actualCount := len(transactions)
	if actualCount != expectedCount {
		t.Errorf("expected count %d, actual %d", expectedCount, expectedCount)
	}
}

func SubTestTransactionRepositoryImpl_FilterIgnorePagination_IncorrectPageSize(
	t *testing.T,
	repo *TransactionRepositoryImpl,
) {
	transactions := repo.Filter([]TransactionFilter{}, 1, 0)
	expectedCount := 3
	actualCount := len(transactions)
	if actualCount != expectedCount {
		t.Errorf("expected count %d, actual %d", expectedCount, expectedCount)
	}
}

func checkFilterSingleResult(t *testing.T, transactions []models.Transaction, idToCheck int) {
	expectedLen := 1
	actualLen := len(transactions)
	if actualLen != expectedLen {
		t.Errorf("expected len %d, actual len %d", expectedLen, actualLen)
	}

	if transactions[0].Id != testTransactions[idToCheck].Id {
		t.Errorf("expected transaction id %d, actual %d", testTransactions[idToCheck].Id, transactions[0].Id)
	}
}

func TestNewTransactionRepository_IsNotNil(t *testing.T) {
	if NewTransactionRepository(nil) == nil {
		t.Error("transaction repository is nil")
	}
}

func TestTransactionFilterBuilderImpl_AddTransactionId_Added(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddTransactionId("10")
	if err != nil {
		t.Error(err)
	}

	if len(builder.filters) != 1 {
		t.Error("filter was not added")
	}
}

func TestTransactionFilterBuilderImpl_AddTransactionId_InvalidValue(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddTransactionId("hello")
	if err == nil {
		t.Error("error is nil")
	}
}

func TestTransactionFilterBuilderImpl_AddTransactionId_NegativeId(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddTransactionId("-1")
	if err == nil {
		t.Error("error is nil")
	}
}

func TestTransactionFilterBuilderImpl_AddTransactionId_Empty(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddTransactionId("")
	if err != nil {
		t.Error(err)
	}

	if len(builder.filters) != 0 {
		t.Error("filter was added")
	}
}

func TestTransactionFilterBuilderImpl_AddTerminalIds_Added(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddTerminalIds([]string{"10", "11"})
	if err != nil {
		t.Error(err)
	}

	if len(builder.filters) != 1 {
		t.Error("filter was not added")
	}
}

func TestTransactionFilterBuilderImpl_AddTerminalIds_InvalidValue(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddTerminalIds([]string{"1", "hello"})
	if err == nil {
		t.Error("error is nil")
	}
}

func TestTransactionFilterBuilderImpl_AddTerminalIds_NegativeId(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddTerminalIds([]string{"-1"})
	if err == nil {
		t.Error("error is nil")
	}
}

func TestTransactionFilterBuilderImpl_AddTerminalIds_EmptyList(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddTerminalIds([]string{})
	if err != nil {
		t.Error(err)
	}

	if len(builder.filters) != 0 {
		t.Error("filter was added")
	}
}

func TestTransactionFilterBuilderImpl_AddTerminalIds_EmptyValueInList(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddTerminalIds([]string{""})
	if err == nil {
		t.Error("error is nil")
	}
}

func TestTransactionFilterBuilderImpl_AddStatus_Accepted(t *testing.T) {
	testAddStatusSuccess(t, models.ACCEPTED)
}

func TestTransactionFilterBuilderImpl_AddStatus_Declined(t *testing.T) {
	testAddStatusSuccess(t, models.DECLINED)
}

func testAddStatusSuccess(t *testing.T, status models.StatusType) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddStatus(string(status))
	if err != nil {
		t.Error(err)
	}

	if len(builder.filters) != 1 {
		t.Error("filter was not added")
	}
}

func TestTransactionFilterBuilderImpl_AddStatus_Invalid(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddStatus("rejected")
	if err == nil {
		t.Error("error is nil")
	}
}

func TestTransactionFilterBuilderImpl_AddStatus_Empty(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddStatus("")
	if err != nil {
		t.Error(err)
	}

	if len(builder.filters) != 0 {
		t.Error("filter was added")
	}
}

func TestTransactionFilterBuilderImpl_AddPaymentType_Cash(t *testing.T) {
	testAddPaymentTypeSuccess(t, models.CASH)
}

func TestTransactionFilterBuilderImpl_AddPaymentType_Card(t *testing.T) {
	testAddPaymentTypeSuccess(t, models.CARD)
}

func testAddPaymentTypeSuccess(t *testing.T, paymentType models.PaymentTypeType) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddPaymentType(string(paymentType))
	if err != nil {
		t.Error(err)
	}

	if len(builder.filters) != 1 {
		t.Error("filter was not added")
	}
}

func TestTransactionFilterBuilderImpl_AddPaymentType_Invalid(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddPaymentType("money")
	if err == nil {
		t.Error("error is nil")
	}
}

func TestTransactionFilterBuilderImpl_AddPaymentType_Empty(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddPaymentType("")
	if err != nil {
		t.Error(err)
	}

	if len(builder.filters) != 0 {
		t.Error("filter was added")
	}
}

func TestTransactionFilterBuilderImpl_AddDatePostRange_Added(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddDatePostRange("2022-08-12 14:25:27", "2022-08-18 16:25:27")
	if err != nil {
		t.Error(err)
	}

	if len(builder.filters) != 1 {
		t.Error("filter was not added")
	}
}

func TestTransactionFilterBuilderImpl_AddDatePostRange_ValueToIsEmpty(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddDatePostRange("2022-08-12 14:25:27", "")
	if err == nil {
		t.Error("error is nil")
	}
}

func TestTransactionFilterBuilderImpl_AddDatePostRange_InvalidValueFrom(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddDatePostRange("2022-08 14:25:27", "2022-08-17 14:25:27")
	if err == nil {
		t.Error("error is nil")
	}
}

func TestTransactionFilterBuilderImpl_AddDatePostRange_InvalidValueTo(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddDatePostRange("2022-08-12 14:25:27", "2022-08-17 25:27")
	if err == nil {
		t.Error("error is nil")
	}
}

func TestTransactionFilterBuilderImpl_AddDatePostRange_Empty(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddDatePostRange("", "")
	if err != nil {
		t.Error(err)
	}

	if len(builder.filters) != 0 {
		t.Error("filter was added")
	}
}

func TestTransactionFilterBuilderImpl_AddPaymentNarrative_Added(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddPaymentNarrative("some text")
	if err != nil {
		t.Error(err)
	}

	if len(builder.filters) != 1 {
		t.Error("filter was not added")
	}
}

func TestTransactionFilterBuilderImpl_AddPaymentNarrative_Empty(t *testing.T) {
	builder := TransactionFilterBuilderImpl{}
	err := builder.AddPaymentNarrative("")
	if err != nil {
		t.Error(err)
	}

	if len(builder.filters) != 0 {
		t.Error("filter was added")
	}
}

func openTestDb() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB_NAME"),
		os.Getenv("POSTGRES_PORT"),
	)
	return gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: fmt.Sprintf("%s.", os.Getenv("POSTGRES_USER_SCHEMA")),
			},
		},
	)
}

var testData = []string{
	"TransactionId,RequestId,TerminalId,PartnerObjectId,AmountTotal,AmountOriginal,CommissionPS,CommissionClient,CommissionProvider,DateInput,DatePost,Status,PaymentType,PaymentNumber,ServiceId,Service,PayeeId,PayeeName,PayeeBankMfo,PayeeBankAccount,PaymentNarrative",
	"1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів згідно договору про надання послуг А11/27122 від 19.11.2020 р.",
	"2,20030,3507,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 12:36:52,2022-08-12 15:36:53,declined,cash,PS16698215,13990,Поповнення карток,14332255,privat,255752,UA713461333619513,Перерахування коштів згідно договору про надання послуг А11/27123 від 19.11.2020 р.",
	"3,20040,3508,1111,3.00,3.00,0.00,0.00,-0.01,2022-08-17 9:53:43,2022-08-17 12:53:44,accepted,card,PS16698225,14000,Поповнення карток,14432355,privat,256753,UA713471293319503,Перерахування коштів згідно договору про надання послуг А11/27122 від 19.11.2020 р.",
}

var testTransactions []models.Transaction

func init() {
	for _, row := range testData[1:] {
		transaction, _ := models.NewTransactionFromCSVRow(strings.Split(row, ","))
		testTransactions = append(testTransactions, transaction)
	}
}
