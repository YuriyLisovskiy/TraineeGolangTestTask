package cli

import (
	"fmt"
	"os"
	"testing"

	"TraineeGolangTestTask/app"
)

func Test_getPageSizeFromEnv_GetFromEnv(t *testing.T) {
	expected := 1
	_ = os.Setenv(app.EnvAppPageSize, fmt.Sprintf("%d", expected))
	actual := getPageSizeFromEnvOrDefault(50)
	if actual != expected {
		t.Errorf("expected %d, actual %d", expected, actual)
	}

	_ = os.Unsetenv(app.EnvAppPageSize)
}

func Test_getPageSizeFromEnv_GetDefaultDueToEmptyEnv(t *testing.T) {
	_ = os.Unsetenv(app.EnvAppPageSize)
	runTest_getPageSizeFromEnv_GetDefault(t)
}

func Test_getPageSizeFromEnv_GetDefaultDueToInvalidEnv(t *testing.T) {
	_ = os.Setenv(app.EnvAppPageSize, "non-int value")
	runTest_getPageSizeFromEnv_GetDefault(t)
	_ = os.Unsetenv(app.EnvAppPageSize)
}

func runTest_getPageSizeFromEnv_GetDefault(t *testing.T) {
	expected := 1
	actual := getPageSizeFromEnvOrDefault(expected)
	if actual != expected {
		t.Errorf("expected %d, actual %d", expected, actual)
	}
}
