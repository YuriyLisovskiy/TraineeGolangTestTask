package app

import (
	"github.com/gin-gonic/gin"
)

func (a *Application) addRoutes(r *gin.Engine) {
	apiTransactions := r.Group("/api/transactions")
	apiTransactions.GET("/json", a.handleTransactionsAsJson)
	apiTransactions.GET("/csv", a.handleTransactionsAsCsv)
	apiTransactions.POST("/upload", a.handleTransactionsUpload)
}
