package models

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Test_hasEnumType(t *testing.T) {
	db, err := openTestDb()
	if err != nil {
		t.Error("unable to open the database")
	}

	t.Run(
		"True", func(t *testing.T) {
			SubTest_hasEnumType_True(t, db)
		},
	)
	t.Run(
		"False", func(t *testing.T) {
			SubTest_hasEnumType_False(t, db)
		},
	)
}

func SubTest_hasEnumType_True(t *testing.T, db *gorm.DB) {
	exists, err := hasEnumType(db, "bool")
	if err != nil {
		t.Error(err)
	}

	if !exists {
		t.Errorf("expected %v, received %v", true, exists)
	}
}

func SubTest_hasEnumType_False(t *testing.T, db *gorm.DB) {
	exists, err := hasEnumType(db, "my_type")
	if err != nil {
		t.Error(err)
	}

	if exists {
		t.Errorf("expected %v, received %v", false, exists)
	}
}

func Test_createEnumIfNotExists(t *testing.T) {
	typeToCreate := "my_type"
	db, err := openTestDb()
	if err != nil {
		t.Error("unable to open the database")
	}

	tearDown := func(db *gorm.DB) {
		db.Exec("DROP TYPE " + typeToCreate)
	}

	t.Run(
		"Create", func(t *testing.T) {
			SubTest_createEnumIfNotExists_Create(t, db, typeToCreate)
			tearDown(db)
		},
	)

	t.Run(
		"AlreadyExists", func(t *testing.T) {
			SubTest_createEnumIfNotExists_AlreadyExists(t, db, typeToCreate)
			tearDown(db)
		},
	)

	t.Run(
		"InvalidTypeName", func(t *testing.T) {
			SubTest_createEnumIfNotExists_InvalidTypeName(t, db)
			tearDown(db)
		},
	)
}

func SubTest_createEnumIfNotExists_Create(t *testing.T, db *gorm.DB, typeToCreate string) {
	exists, _ := hasEnumType(db, typeToCreate)
	if exists {
		t.Errorf("type %s already exists", typeToCreate)
	}

	err := createEnumIfNotExists(db, typeToCreate, []string{"1", "2", "3"})
	if err != nil {
		t.Error(err)
	}

	exists, _ = hasEnumType(db, typeToCreate)
	if !exists {
		t.Errorf("created type %s does not exist", typeToCreate)
	}
}

func SubTest_createEnumIfNotExists_AlreadyExists(t *testing.T, db *gorm.DB, typeToCreate string) {
	err := db.Exec(fmt.Sprintf("CREATE TYPE %s AS ENUM ('1', '2', '3')", typeToCreate)).Error
	if err != nil {
		t.Error(err)
	}

	err = createEnumIfNotExists(db, typeToCreate, []string{"1", "2", "3"})
	if err != nil {
		t.Error(err)
	}
}

func SubTest_createEnumIfNotExists_InvalidTypeName(t *testing.T, db *gorm.DB) {
	err := createEnumIfNotExists(db, "1", []string{"1", "2", "3"})
	if err == nil {
		t.Error("error is nil")
	}
}

func Test_migrateAll(t *testing.T) {
	db, err := openTestDb()
	if err != nil {
		t.Error("unable to open the database")
	}

	db.Exec("CREATE SCHEMA rest_api")
	err = migrateAll(db)
	if err != nil {
		t.Error(err)
	}

	ensureEntityExists(t, db, fmt.Sprintf("SELECT count(*) FROM pg_type WHERE typname = '%s'", "status_type"))
	ensureEntityExists(t, db, fmt.Sprintf("SELECT count(*) FROM pg_type WHERE typname = '%s'", "payment_type_type"))
	ensureEntityExists(t, db, fmt.Sprintf("SELECT count(*) FROM pg_tables WHERE tablename = '%s'", "transactions"))

	log.Println(db.Exec("DROP TABLE rest_api.transactions").Error)
	log.Println(db.Exec("DROP TYPE rest_api.status_type").Error)
	log.Println(db.Exec("DROP TYPE rest_api.payment_type_type").Error)
	log.Println(db.Exec("DROP SCHEMA rest_api").Error)
}

func ensureEntityExists(t *testing.T, db *gorm.DB, sql string) {
	var exists bool
	err := db.Raw(sql).Row().Scan(&exists)
	if err != nil {
		t.Error(err)
	}

	if !exists {
		t.Error("entity does not exist")
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
