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
	pageSize              int
	TransactionRepository repositories.TransactionRepository
}

func (a *Application) Configure() error {
	mode := gin.ReleaseMode
	debug, err := strconv.ParseBool(os.Getenv("GIN_DEBUG"))
	if err != nil {
		return err
	}

	pageSizeString := os.Getenv("APP_PAGE_SIZE")
	if pageSizeString != "" {
		a.pageSize, err = strconv.Atoi(pageSizeString)
		if err != nil {
			return err
		}
	} else {
		a.pageSize = 30
	}

	if debug {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)
	return nil
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

	// listen for the interrupt signal
	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	shutdownTimeoutSec, err := strconv.Atoi(os.Getenv("GIN_SHUTDOWN_TIMEOUT"))
	if err != nil {
		shutdownTimeoutSec = 5
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

func (a *Application) buildRouter() *gin.Engine {
	router := gin.Default()
	configureMaxMultipartMemory(router)
	a.addRoutes(router)
	return router
}

func (a *Application) sendInternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
	if message != "" {
		log.Println(message)
	}
}

func (a *Application) sendBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{"message": message})
	if message != "" {
		log.Println(message)
	}
}

func configureMaxMultipartMemory(router *gin.Engine) {
	maxMultipartMemoryString := os.Getenv("GIN_MAX_MULTIPART_MEMORY")
	if maxMultipartMemoryString != "" {
		maxMultipartMemory, err := strconv.ParseInt(maxMultipartMemoryString, 10, 64)
		if err == nil {
			router.MaxMultipartMemory = maxMultipartMemory
			return
		}

		log.Printf("unable to parse %s to a 64-bit integer; using 32 mb by default\n", maxMultipartMemoryString)
	}

	router.MaxMultipartMemory = 8 << 22 // 32 mb by default
}
