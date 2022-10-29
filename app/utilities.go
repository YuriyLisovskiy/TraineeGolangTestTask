package app

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"TraineeGolangTestTask/models"
	"TraineeGolangTestTask/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ConnectToPostgreSQLWithEnv() (*gorm.DB, error) {
	port, err := strconv.Atoi(getenvOrDefault("POSTGRES_PORT", "5432"))
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		getenvOrDefault("POSTGRES_HOST", "localhost"),
		getenvOrDefault("POSTGRES_USER", "postgres"),
		getenvOrDefault("POSTGRES_PASSWORD", ""),
		getenvOrDefault("POSTGRES_DB_NAME", ""),
		port,
	)
	config := gorm.Config{}
	userSchema := os.Getenv("POSTGRES_USER_SCHEMA")
	if userSchema != "" {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix: fmt.Sprintf("%s.", userSchema),
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	return db, nil
}

func getenvOrDefault(key, default_ string) string {
	val := os.Getenv(key)
	if val == "" {
		return default_
	}

	return val
}

func buildFilters(c *gin.Context) ([]repositories.TransactionFilter, error) {
	var filters []repositories.TransactionFilter
	if stringTransactionId := c.DefaultQuery("transaction_id", ""); stringTransactionId != "" {
		id, err := strconv.ParseUint(stringTransactionId, 10, 64)
		if err != nil {
			return nil, err
		}

		filters = append(filters, repositories.FilterByTransactionId(id))
	}

	if terminalIds := c.QueryArray("terminal_id"); len(terminalIds) > 0 {
		var ids []uint64
		for _, stringId := range terminalIds {
			id, err := strconv.ParseUint(stringId, 10, 64)
			if err != nil {
				return nil, err
			}

			ids = append(ids, id)
		}
		filters = append(filters, repositories.FilterByTerminalId(ids))
	}

	if statusParam := c.DefaultQuery("status", ""); statusParam != "" {
		switch status := models.StatusType(statusParam); status {
		case models.ACCEPTED, models.DECLINED:
			filters = append(filters, repositories.FilterByStatus(status))
		default:
			return nil, fmt.Errorf(
				"value of \"status\" parameter should be either \"%v\" or \"%v\"",
				models.ACCEPTED,
				models.DECLINED,
			)
		}
	}

	if paymentTypeParam := c.DefaultQuery("payment_type", ""); paymentTypeParam != "" {
		switch paymentType := models.PaymentTypeType(paymentTypeParam); paymentType {
		case models.CASH, models.CARD:
			filters = append(filters, repositories.FilterByPaymentType(paymentType))
		default:
			return nil, fmt.Errorf(
				"value of \"payment_type\" parameter should be either \"%v\" or \"%v\"",
				models.CASH,
				models.CARD,
			)
		}
	}

	if fromString := c.DefaultQuery("date_post_from", ""); fromString != "" {
		toString := c.DefaultQuery("date_post_to", "")
		if toString == "" {
			return nil, errors.New("parameter \"date_post_to\" is required when using \"date_post_from\"")
		}

		from, err := time.Parse(models.TimeLayout, fromString)
		if err != nil {
			return nil, err
		}

		to, err := time.Parse(models.TimeLayout, toString)
		if err != nil {
			return nil, err
		}

		filters = append(filters, repositories.FilterByDatePostTimeRange(from, to))
	}

	if paymentNarrative := c.DefaultQuery("payment_narrative", ""); paymentNarrative != "" {
		filters = append(filters, repositories.ContainsTextInPaymentNarrative(paymentNarrative))
	}

	return filters, nil
}
