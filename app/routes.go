package app

import (
	"github.com/gin-gonic/gin"
)

func (a *Application) addRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/transactions/json", a.handleTransactionsAsJson)
	api.GET("/transactions/csv", a.handleTransactionsAsCsv)
	api.POST("/upload", a.handleUpload)
}
