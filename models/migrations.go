package models

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"
)

func MigrateAll(db *gorm.DB) error {
	return db.Transaction(migrateAll)
}

func migrateAll(tx *gorm.DB) error {
	err := createEnumIfNotExists(tx, "status_type", []string{string(ACCEPTED), string(DECLINED)})
	if err != nil {
		return err
	}

	err = createEnumIfNotExists(tx, "payment_type_type", []string{string(CASH), string(CARD)})
	if err != nil {
		return err
	}

	return tx.AutoMigrate(&Transaction{})
}

func createEnumIfNotExists(tx *gorm.DB, typeName string, values []string) error {
	exists, err := hasEnumType(tx, typeName)
	if err != nil {
		return err
	}

	if !exists {
		valuesAsString := fmt.Sprintf("'%s'", strings.Join(values, "', '"))
		schema := os.Getenv("POSTGRES_USER_SCHEMA")
		if schema != "" {
			schema += "."
		}

		sql := fmt.Sprintf("CREATE TYPE %s%s AS ENUM (%s)", schema, typeName, valuesAsString)
		return tx.Raw(sql).Row().Err()
	}

	return nil
}

func hasEnumType(tx *gorm.DB, typeName string) (bool, error) {
	var exists bool
	err := tx.Raw("SELECT count(*) > 0 FROM pg_type WHERE typname = ?", typeName).Row().Scan(&exists)
	return exists, err
}
