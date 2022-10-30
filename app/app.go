package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"TraineeGolangTestTask/repositories"
	"github.com/gin-gonic/gin"
)

const (
	EnvAppPageSize           = "APP_PAGE_SIZE"
	EnvGinMaxMultipartMemory = "GIN_MAX_MULTIPART_MEMORY"
	EnvGinShutdownTimeout    = "GIN_SHUTDOWN_TIMEOUT"

	DefaultAppPageSize           = 30
	DefaultGinMaxMultipartMemory = 8 << 22 // 32 mb
	DefaultGinShutdownTimeout    = 5
)

type Application struct {
	PageSize              int
	TransactionRepository repositories.TransactionRepository
}

func (a *Application) Execute(addr string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := a.configureRouter(gin.Default())
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// listen for the interrupt signal
	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	shutdownTimeoutSec, err := strconv.Atoi(os.Getenv(EnvGinShutdownTimeout))
	if err != nil {
		shutdownTimeoutSec = DefaultGinShutdownTimeout
		// ignore error
	}

	// inform the server it has n seconds to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeoutSec)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return errors.New(fmt.Sprintf("Server forced to shut down: %v", err))
	}

	log.Println("Server exiting")
	return nil
}

func (a *Application) configureRouter(router *gin.Engine) *gin.Engine {
	setMaxMultipartMemoryOrDefault(router, DefaultGinMaxMultipartMemory)
	a.addRoutes(router)
	return router
}

func (a *Application) addRoutes(r *gin.Engine) {
	apiTransactions := r.Group("/api/transactions")
	apiTransactions.GET("/csv", a.handleTransactionsAsCsv)
	apiTransactions.GET("/json", a.handleTransactionsAsJson)
	apiTransactions.POST("/upload", a.handleTransactionsUpload)
}

func (a *Application) sendInternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
	if message != "" {
		log.Println(message)
	}
}

func (a *Application) sendBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"message": message})
	if message != "" {
		log.Println(message)
	}
}

func setMaxMultipartMemoryOrDefault(router *gin.Engine, defaultValue int64) {
	maxMultipartMemoryString := os.Getenv(EnvGinMaxMultipartMemory)
	if maxMultipartMemoryString != "" {
		maxMultipartMemory, err := strconv.ParseInt(maxMultipartMemoryString, 10, 64)
		if err == nil {
			router.MaxMultipartMemory = maxMultipartMemory
			return
		}

		log.Printf(
			"unable to parse %s to a 64-bit integer; using %d by default\n",
			maxMultipartMemoryString,
			defaultValue,
		)
	}

	router.MaxMultipartMemory = defaultValue
}
