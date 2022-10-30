package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestApplication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	app := Application{}
	t.Run(
		"addRoutes", func(t *testing.T) {
			SubTestApplication_addRoutes(t, &app)
		},
	)
	t.Run(
		"sendErrorResponse", func(t *testing.T) {
			SubTestApplication_sendErrorResponse(t, &app)
		},
	)
	t.Run(
		"", func(t *testing.T) {

		},
	)
}

func SubTestApplication_addRoutes(t *testing.T, app *Application) {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	app.addRoutes(router)
	routes := router.Routes()
	if len(routes) != 3 {
		t.Errorf("expected routes count %d, actual %d", 3, len(routes))
	}

	sort.Slice(
		routes, func(i, j int) bool {
			return routes[i].Path < routes[j].Path
		},
	)

	addRoutesAssertPathAndMethod(t, routes[0], "/api/transactions/csv", "GET")
	addRoutesAssertPathAndMethod(t, routes[1], "/api/transactions/json", "GET")
	addRoutesAssertPathAndMethod(t, routes[2], "/api/transactions/upload", "POST")
}

func addRoutesAssertPathAndMethod(t *testing.T, route gin.RouteInfo, expectedPath, expectedMethod string) {
	if route.Path != expectedPath {
		t.Errorf("expected path \"%s\", actual path \"%s\"", expectedPath, route.Path)
	}

	if route.Method != expectedMethod {
		t.Errorf("expected method \"%s\", actual method \"%s\"", expectedMethod, route.Method)
	}
}

func SubTestApplication_sendErrorResponse(t *testing.T, app *Application) {
	t.Run(
		"InternalError", func(t *testing.T) {
			SubTestApplication_sendInternalError(t, app)
		},
	)
	t.Run(
		"BadRequest", func(t *testing.T) {
			SubTestApplication_sendBadRequest(t, app)
		},
	)
}

func SubTestApplication_sendInternalError(t *testing.T, app *Application) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	app.sendInternalError(c, "")
	sendErrorResponseAssertStatus(t, http.StatusInternalServerError, w.Code)
	sendErrorResponseAssertMessage(t, w.Body, "internal error")
}

func SubTestApplication_sendBadRequest(t *testing.T, app *Application) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	expectedMessage := "some error"
	app.sendBadRequest(c, expectedMessage)

	sendErrorResponseAssertStatus(t, http.StatusBadRequest, w.Code)
	sendErrorResponseAssertMessage(t, w.Body, expectedMessage)
}

func sendErrorResponseAssertMessage(t *testing.T, body *bytes.Buffer, expectedMessage string) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		t.Error(err)
	}

	resp := MessageResponseMock{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		t.Error(err)
	}

	if resp.Message != expectedMessage {
		t.Errorf("expected %s, actual %s", expectedMessage, resp.Message)
	}
}

func sendErrorResponseAssertStatus(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Errorf("expected %d, actual %d", expected, actual)
	}
}

type MessageResponseMock struct {
	Message string `json:"message"`
}

func Test_configureMaxMultipartMemory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	t.Run(
		"FromEnv", func(t *testing.T) {
			SubTest_configureMaxMultipartMemory_FromEnv(t, router)
		},
	)
	_, router = gin.CreateTestContext(w)
	t.Run(
		"SetDefaultDueToEmptyEnvVar", func(t *testing.T) {
			SubTest_configureMaxMultipartMemory_SetDefaultDueToEmptyEnvVar(t, router)
		},
	)
	_, router = gin.CreateTestContext(w)
	t.Run(
		"SetDefaultDueToNonIntValue", func(t *testing.T) {
			SubTest_configureMaxMultipartMemory_SetDefaultDueToNonIntValue(t, router)
		},
	)
}

func SubTest_configureMaxMultipartMemory_FromEnv(t *testing.T, router *gin.Engine) {
	envVar := "GIN_MAX_MULTIPART_MEMORY"
	envValue := int64(12345678912345)
	_ = os.Setenv(envVar, fmt.Sprintf("%d", envValue))
	setMaxMultipartMemoryOrDefault(router, envValue)
	if router.MaxMultipartMemory != envValue {
		t.Errorf("expected %d, actual %d", envValue, router.MaxMultipartMemory)
	}

	_ = os.Unsetenv(envVar)
}

func SubTest_configureMaxMultipartMemory_SetDefaultDueToEmptyEnvVar(t *testing.T, router *gin.Engine) {
	envVar := "GIN_MAX_MULTIPART_MEMORY"
	expectedValue := int64(12345)
	_ = os.Setenv(envVar, "")
	setMaxMultipartMemoryOrDefault(router, expectedValue)
	if router.MaxMultipartMemory != expectedValue {
		t.Errorf("expected %d, actual %d", expectedValue, router.MaxMultipartMemory)
	}

	_ = os.Unsetenv(envVar)
}

func SubTest_configureMaxMultipartMemory_SetDefaultDueToNonIntValue(t *testing.T, router *gin.Engine) {
	expectedValue := int64(12345)
	setMaxMultipartMemoryOrDefault(router, expectedValue)
	if router.MaxMultipartMemory != expectedValue {
		t.Errorf("expected %d, actual %d", expectedValue, router.MaxMultipartMemory)
	}
}
