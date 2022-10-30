package app

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"TraineeGolangTestTask/repositories"
	"github.com/gin-gonic/gin"
)

func Test_getEnvOrDefault_GetOriginal(t *testing.T) {
	key := fmt.Sprintf("TEST_VAR_%d", time.Now().Unix())
	value := "some value"
	_ = os.Setenv(key, value)
	actual := getEnvOrDefault(key, "another value")
	if actual != value {
		t.Errorf("got non-original value: %s", actual)
	}

	_ = os.Unsetenv(key)
}

func Test_getEnvOrDefault_GetDefault(t *testing.T) {
	key := fmt.Sprintf("TEST_VAR_%d", time.Now().Unix())
	expectedValue := "default value"
	actual := getEnvOrDefault(key, expectedValue)
	if actual != expectedValue {
		t.Errorf("got non-default value: %s", actual)
	}
}

func Test_parseParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)
	t.Run(
		"Parsed", func(t *testing.T) {
			SubTest_parseParameters_Parsed(t, c)
		},
	)
	t.Run(
		"FailedDueToBuilderError", func(t *testing.T) {
			SubTest_parseParameters_FailedDueToBuilderError(t, c)
		},
	)
}

func SubTest_parseParameters_Parsed(t *testing.T, c *gin.Context) {
	c.Params = []gin.Param{
		{
			Key:   "transaction_id",
			Value: "1",
		},
		{
			Key:   "terminal_id",
			Value: "2",
		},
		{
			Key:   "terminal_id",
			Value: "5",
		},
		{
			Key:   "terminal_id",
			Value: "77",
		},
		{
			Key:   "status",
			Value: "accepted",
		},
		{
			Key:   "payment_type",
			Value: "cash",
		},
		{
			Key:   "date_post_from",
			Value: "2022-08-12 14:25:27",
		},
		{
			Key:   "date_post_to",
			Value: "2022-08-15 14:25:27",
		},
		{
			Key:   "payment_narrative",
			Value: "some text",
		},
	}
	builder := TransactionFilterBuilderMock{Filters: map[filterHash]string{}}
	_ = parseParameters(c, &builder)
	if !builder.hasFilterWithValue(transactionIdFilter, c.DefaultQuery("transaction_id", "")) {
		t.Errorf("filter %s is absent or has incorrect value", "transaction_id")
	}

	if !builder.hasFilterWithValue(terminalIdsFilter, strings.Join(c.QueryArray("terminal_id"), ",")) {
		t.Errorf("filter %s is absent or has incorrect value", "terminal_id")
	}

	if !builder.hasFilterWithValue(statusFilter, c.DefaultQuery("status", "")) {
		t.Errorf("filter %s is absent or has incorrect value", "status")
	}

	if !builder.hasFilterWithValue(paymentTypeFilter, c.DefaultQuery("payment_type", "")) {
		t.Errorf("filter %s is absent or has incorrect value", "payment_type")
	}

	datePostRangeExpected := c.DefaultQuery("date_post_from", "") + "-" + c.DefaultQuery("date_post_to", "")
	if !builder.hasFilterWithValue(datePostRangeFilter, datePostRangeExpected) {
		t.Errorf("filter %s-%s is absent or has incorrect value", "date_post_from", "date_post_to")
	}

	if !builder.hasFilterWithValue(paymentNarrativeFilter, c.DefaultQuery("payment_narrative", "")) {
		t.Errorf("filter %s is absent or has incorrect value", "payment_narrative")
	}
}

func SubTest_parseParameters_FailedDueToBuilderError(t *testing.T, c *gin.Context) {
	builder := TransactionFilterBuilderWithErrorMock{}
	err := parseParameters(c, &builder)
	if err == nil {
		t.Error("error is nil")
	}
}

type filterHash int

const (
	transactionIdFilter = iota
	terminalIdsFilter
	statusFilter
	paymentTypeFilter
	datePostRangeFilter
	paymentNarrativeFilter
)

type TransactionFilterBuilderMock struct {
	Filters map[filterHash]string
}

func (tf *TransactionFilterBuilderMock) AddTransactionId(value string) error {
	tf.Filters[transactionIdFilter] = value
	return nil
}

func (tf *TransactionFilterBuilderMock) AddTerminalIds(values []string) error {
	tf.Filters[terminalIdsFilter] = strings.Join(values, ",")
	return nil
}

func (tf *TransactionFilterBuilderMock) AddStatus(value string) error {
	tf.Filters[statusFilter] = value
	return nil
}

func (tf *TransactionFilterBuilderMock) AddPaymentType(value string) error {
	tf.Filters[paymentTypeFilter] = value
	return nil
}

func (tf *TransactionFilterBuilderMock) AddDatePostRange(valueFrom, valueTo string) error {
	tf.Filters[datePostRangeFilter] = valueFrom + "-" + valueTo
	return nil
}

func (tf *TransactionFilterBuilderMock) AddPaymentNarrative(value string) error {
	tf.Filters[paymentNarrativeFilter] = value
	return nil
}

func (tf *TransactionFilterBuilderMock) GetFilters() []repositories.TransactionFilter {
	return nil
}

func (tf *TransactionFilterBuilderMock) hasFilterWithValue(hash filterHash, expectedValue string) bool {
	actualValue, ok := tf.Filters[hash]
	return ok && actualValue == expectedValue
}

type TransactionFilterBuilderWithErrorMock struct {
	TransactionFilterBuilderMock
}

func (tf *TransactionFilterBuilderWithErrorMock) AddTransactionId(string) error {
	return errors.New("some error")
}
