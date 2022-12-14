package cli

import (
	"log"
	"os"
	"strconv"

	"TraineeGolangTestTask/app"
	"TraineeGolangTestTask/repositories"
	"github.com/spf13/cobra"
)

var (
	addressArg string

	rootCmd = &cobra.Command{
		Use:  "rest-api-app",
		RunE: runRootCommand,
	}
)

func RunCLI() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(
		&addressArg, "bind", "b", "127.0.0.1:8000", "bind address",
	)
}

func runRootCommand(*cobra.Command, []string) error {
	db, err := app.ConnectToPostgreSQLWithEnv()
	if err != nil {
		return err
	}

	application := app.Application{
		PageSize:              getPageSizeFromEnvOrDefault(app.DefaultAppPageSize),
		TransactionRepository: repositories.NewTransactionRepository(db),
	}

	log.Printf("Serving at %s\n", addressArg)
	return application.Execute(addressArg)
}

func getPageSizeFromEnvOrDefault(defaultValue int) int {
	pageSizeString := os.Getenv(app.EnvAppPageSize)
	if pageSizeString != "" {
		parsedPageSize, err := strconv.Atoi(pageSizeString)
		if err == nil {
			return parsedPageSize
		}
	}

	return defaultValue
}
