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

type Application struct {
	PageSize              int
	TransactionRepository repositories.TransactionRepository
}

func (a *Application) Configure() error {
	mode := gin.ReleaseMode
	debug, err := strconv.ParseBool(os.Getenv("GIN_DEBUG"))
	if err != nil {
		return err
	}

	if debug {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)
	return nil
}

func (a *Application) buildRouter() *gin.Engine {
	router := gin.Default()
	maxMultipartMemoryString := os.Getenv("GIN_MAX_MULTIPART_MEMORY")
	if maxMultipartMemoryString != "" {
		var err error
		router.MaxMultipartMemory, err = strconv.ParseInt(maxMultipartMemoryString, 10, 64)
		if err != nil {

		}
	} else {
		router.MaxMultipartMemory = 8 << 22 // 32 mb
	}

	a.addRoutes(router)
	return router
}

func (a *Application) Execute(addr string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := a.buildRouter()
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	shutdownTimeoutSec, err := strconv.Atoi(os.Getenv("GIN_SHUTDOWN_TIMEOUT"))
	if err != nil {
		shutdownTimeoutSec = 5
		// ignore error
	}

	// The context is used to inform the server it has n seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeoutSec)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return errors.New(fmt.Sprintf("Server forced to shut down: %v", err))
	}

	log.Println("Server exiting")
	return nil
}
