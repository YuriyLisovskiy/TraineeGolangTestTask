package app

import (
	"fmt"
	"os"
	"strconv"

	"TraineeGolangTestTask/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ConnectToPostgreSQLWithEnv() (*gorm.DB, error) {
	port, err := strconv.Atoi(getEnvOrDefault("POSTGRES_PORT", "5432"))
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		getEnvOrDefault("POSTGRES_HOST", "localhost"),
		getEnvOrDefault("POSTGRES_USER", "postgres"),
		getEnvOrDefault("POSTGRES_PASSWORD", ""),
		getEnvOrDefault("POSTGRES_DB_NAME", ""),
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

func getEnvOrDefault(key, default_ string) string {
	val := os.Getenv(key)
	if val == "" {
		return default_
	}

	return val
}

func parseParameters(c *gin.Context, builder repositories.TransactionFilterBuilder) error {
	err := builder.AddTransactionId(c.DefaultQuery("transaction_id", ""))
	if err != nil {
		return err
	}

	err = builder.AddTerminalIds(c.QueryArray("terminal_id"))
	if err != nil {
		return err
	}

	err = builder.AddStatus(c.DefaultQuery("status", ""))
	if err != nil {
		return err
	}

	err = builder.AddPaymentType(c.DefaultQuery("payment_type", ""))
	if err != nil {
		return err
	}

	err = builder.AddDatePostRange(c.DefaultQuery("date_post_from", ""), c.DefaultQuery("date_post_to", ""))
	if err != nil {
		return err
	}

	return builder.AddPaymentNarrative(c.DefaultQuery("payment_narrative", ""))
}
