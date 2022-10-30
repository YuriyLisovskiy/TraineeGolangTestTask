package app

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"TraineeGolangTestTask/models"
	"github.com/gin-gonic/gin"
)

func TestApplication_handleTransactionsAsCsv_200(t *testing.T) {
	transactionRepository := newTransactionRepositoryMock(
		[]models.Transaction{
			testTransactions[0],
			testTransactions[1],
		},
	)
	app := Application{
		PageSize:              2,
		TransactionRepository: transactionRepository,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	app.handleTransactionsAsCsv(c)

	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
	}

	actualCsv := string(data)
	expectedCsv := fmt.Sprintf("%s\n%s\n%s\n", testData[0], testData[1], testData[2])
	if actualCsv != expectedCsv {
		t.Errorf("expected csv:\n%s\nactual csv:\n%s", expectedCsv, actualCsv)
	}
}
