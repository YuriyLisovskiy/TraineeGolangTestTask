package cli

import (
	"fmt"
	"log"

	"TraineeGolangTestTask/app"
	"TraineeGolangTestTask/models"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply migrations to the database",
	RunE:  runMigrateCommand,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func runMigrateCommand(*cobra.Command, []string) error {
	db, err := app.ConnectToPostgreSQLWithEnv()
	if err != nil {
		return err
	}

	err = models.MigrateAll(db)
	if err != nil {
		return fmt.Errorf("failed to apply migration to the database: %v", err)
	}

	log.Println("Migrated successfully.")
	return nil
}
