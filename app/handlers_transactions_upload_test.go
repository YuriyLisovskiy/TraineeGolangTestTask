package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"TraineeGolangTestTask/models"
	"TraineeGolangTestTask/repositories"
	"github.com/gin-gonic/gin"
)

var testData = []string{
	"TransactionId,RequestId,TerminalId,PartnerObjectId,AmountTotal,AmountOriginal,CommissionPS,CommissionClient,CommissionProvider,DateInput,DatePost,Status,PaymentType,PaymentNumber,ServiceId,Service,PayeeId,PayeeName,PayeeBankMfo,PayeeBankAccount,PaymentNarrative",
	"1,20020,3506,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 11:25:27,2022-08-12 14:25:27,accepted,cash,PS16698205,13980,Поповнення карток,14232155,pumb,254751,UA713451373919523,Перерахування коштів згідно договору про надання послуг А11/27122 від 19.11.2020 р.",
	"2,20030,3507,1111,1.00,1.00,0.00,0.00,0.00,2022-08-12 12:36:52,2022-08-12 15:36:53,declined,cash,PS16698215,13990,Поповнення карток,14332255,privat,255752,UA713461333619513,Перерахування коштів згідно договору про надання послуг А11/27123 від 19.11.2020 р.",
	"3,20040,3508,1111,3.00,3.00,0.00,0.00,-0.01,2022-08-17 9:53:43,2022-08-17 12:53:44,accepted,card,PS16698225,14000,Поповнення карток,14432355,privat,256753,UA713471293319503,Перерахування коштів згідно договору про надання послуг А11/27122 від 19.11.2020 р.",
}

var testTransactions []models.Transaction

func init() {
	for _, row := range testData[1:] {
		transaction, _ := models.NewTransactionFromCSVRow(row)
		testTransactions = append(testTransactions, *transaction)
	}
}

func TestApplication_handleTransactionsUpload_201(t *testing.T) {
	transactionRepository := newTransactionRepositoryMock()
	app := Application{
		PageSize:              5,
		TransactionRepository: transactionRepository,
	}

	requestMock, err := uploadTestRequestMock(testData)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = requestMock

	app.handleTransactionsUpload(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, actual %d", http.StatusCreated, w.Code)
	}

	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
	}

	responseDataMock := uploadResponseMock{}
	err = json.Unmarshal(data, &responseDataMock)
	if err != nil {
		t.Error(err)
	}

	expectedTrLen := len(testTransactions)
	if responseDataMock.RowsCount != expectedTrLen {
		t.Errorf("expected rows inserted %d, actual is %d", expectedTrLen, responseDataMock.RowsCount)
	}

	actualTrLen := len(transactionRepository.models)
	if actualTrLen != expectedTrLen {
		t.Errorf("expected transactions count is %d, actual is %d", expectedTrLen, actualTrLen)
	}

	for _, transaction := range testTransactions {
		if _, ok := transactionRepository.models[transaction.Id]; !ok {
			t.Errorf("transaction with id %d does not exist", transaction.Id)
		}
	}
}

func TestApplication_handleTransactionsUpload_400NotMultipartRequest(t *testing.T) {
	transactionRepository := newTransactionRepositoryMock()
	app := Application{
		PageSize:              5,
		TransactionRepository: transactionRepository,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	app.handleTransactionsUpload(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, actual %d", http.StatusCreated, w.Code)
	}
}

func TestApplication_handleTransactionsUpload_400MissingFile(t *testing.T) {
	transactionRepository := newTransactionRepositoryMock()
	app := Application{
		PageSize:              5,
		TransactionRepository: transactionRepository,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = uploadTestRequestMock([]string{})

	app.handleTransactionsUpload(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, actual %d", http.StatusCreated, w.Code)
	}
}

func TestApplication_handleTransactionsUpload_400InvalidCSVData(t *testing.T) {
	transactionRepository := newTransactionRepositoryMock()
	app := Application{
		PageSize:              5,
		TransactionRepository: transactionRepository,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = uploadTestRequestMock([]string{testData[0], "1,2,3"})

	app.handleTransactionsUpload(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, actual %d", http.StatusCreated, w.Code)
	}
}

type uploadResponseMock struct {
	RowsCount int `json:"row_count"`
}

type transactionRepositoryMock struct {
	models map[uint64]models.Transaction
}

func newTransactionRepositoryMock() *transactionRepositoryMock {
	return &transactionRepositoryMock{models: map[uint64]models.Transaction{}}
}

func (m *transactionRepositoryMock) Create(model *models.Transaction) error {
	m.models[model.Id] = *model
	return nil
}

func (m *transactionRepositoryMock) CreateBatch(dbTransaction func(repositories.TransactionRepository) error) error {
	return dbTransaction(m)
}

func (m *transactionRepositoryMock) Filter(
	filters []repositories.TransactionFilter,
	page, pageSize int,
) []models.Transaction {
	var result []models.Transaction
	for _, v := range m.models {
		result = append(result, v)
	}

	return result
}

func (m *transactionRepositoryMock) ForEach(
	filters []repositories.TransactionFilter,
	apply func(model *models.Transaction) error,
) error {
	return nil
}

func (m *transactionRepositoryMock) NewFilterBuilder() repositories.TransactionFilterBuilder {
	return nil
}

func uploadTestRequestMock(csv []string) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	if len(csv) > 0 {
		fileWriter, err := writer.CreateFormFile("file", "file.csv")
		if err != nil {
			return nil, err
		}

		for _, row := range csv {
			_, err = fileWriter.Write([]byte(fmt.Sprintf("%s\n", row)))
			if err != nil {
				return nil, err
			}
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	requestMock := httptest.NewRequest(http.MethodPost, "/", body)
	requestMock.Header.Set("Content-Type", writer.FormDataContentType())

	return requestMock, nil
}
