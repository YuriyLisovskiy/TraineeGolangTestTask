package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"TraineeGolangTestTask/models"
	"github.com/gin-gonic/gin"
)

func TestApplication_handleTransactionsAsJson_200GetAll(t *testing.T) {
	transactionRepository := newTransactionRepositoryMock(
		[]models.Transaction{
			testTransactions[0],
			testTransactions[1],
		},
	)
	app := Application{
		PageSize:              5,
		TransactionRepository: transactionRepository,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	app.handleTransactionsAsJson(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, actual %d", http.StatusOK, w.Code)
	}

	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
	}

	responseBody := transactionsAsJsonResponseMock{}
	err = json.Unmarshal(data, &responseBody)
	if err != nil {
		t.Error(err)
	}

	expectedCount := 2
	if responseBody.Count != expectedCount {
		t.Errorf("expected count %d, actual count %d", expectedCount, responseBody.Count)
	}

	resultsLen := len(responseBody.Results)
	if responseBody.Count != resultsLen {
		t.Errorf("count is not equals to results len: %d != %d", responseBody.Count, resultsLen)
	}

	if responseBody.PreviousPage != nil {
		t.Errorf("expected nil previous page, actual %d", *responseBody.PreviousPage)
	}

	if responseBody.NextPage != nil {
		t.Errorf("expected nil next page, actual %d", *responseBody.NextPage)
	}
}

func TestApplication_handleTransactionsAsJson_400InvalidPage(t *testing.T) {
	runTestApplication_handleTransactionsAsJson_400(t, "hello")
}

func TestApplication_handleTransactionsAsJson_400NonPositivePage(t *testing.T) {
	runTestApplication_handleTransactionsAsJson_400(t, "0")
}

func runTestApplication_handleTransactionsAsJson_400(t *testing.T, pageParamToTest string) {
	transactionRepository := newTransactionRepositoryMock(
		[]models.Transaction{
			testTransactions[0],
			testTransactions[1],
		},
	)
	app := Application{
		PageSize:              5,
		TransactionRepository: transactionRepository,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?page=%s", pageParamToTest), nil)

	app.handleTransactionsAsJson(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, actual %d", http.StatusBadRequest, w.Code)
	}

	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
	}

	responseBody := MessageResponseMock{}
	err = json.Unmarshal(data, &responseBody)
	if err != nil {
		t.Error(err)
	}
}

type transactionsAsJsonResponseMock struct {
	Count        int
	NextPage     *int
	PreviousPage *int
	Results      []models.Transaction
}
