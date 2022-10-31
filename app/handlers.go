package app

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"strings"

	"TraineeGolangTestTask/models"
	"TraineeGolangTestTask/repositories"
	"github.com/gin-gonic/gin"
)

func (a *Application) handleTransactionsAsJson(c *gin.Context) {
	pageParameter := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageParameter)
	if err != nil {
		a.sendBadRequest(c, "the \"page\" parameter is required to be an integer number")
		log.Println(err)
		return
	}

	if page <= 0 {
		a.sendBadRequest(c, "the \"page\" parameter is required to be a positive integer number")
		return
	}

	filterBuilder := a.TransactionRepository.NewFilterBuilder()
	err = parseParameters(c, filterBuilder)
	if err != nil {
		a.sendBadRequest(c, err.Error())
		return
	}

	transactions := a.TransactionRepository.Filter(filterBuilder.GetFilters(), page, a.PageSize)
	var (
		previousPage *int
		nextPage     *int
	)
	if page-1 > 0 {
		previousPage = new(int)
		*previousPage = page - 1
	}

	transactionsLen := len(transactions)
	if transactionsLen == a.PageSize {
		nextPage = new(int)
		*nextPage = page + 1
	}

	c.JSON(
		http.StatusOK, gin.H{
			"count":         transactionsLen,
			"next_page":     nextPage,
			"previous_page": previousPage,
			"results":       transactions,
		},
	)
}

func (a *Application) handleTransactionsAsCsv(c *gin.Context) {
	filterBuilder := a.TransactionRepository.NewFilterBuilder()
	err := parseParameters(c, filterBuilder)
	if err != nil {
		a.sendBadRequest(c, err.Error())
		return
	}

	writer := c.Writer
	header := writer.Header()
	header.Set("Transfer-Encoding", "chunked")
	header.Set("Content-Type", "text/csv")
	writer.WriteHeader(http.StatusOK)
	flusher := writer.(http.Flusher)
	flusher.Flush()

	_, err = writer.Write([]byte("TransactionId,RequestId,TerminalId,PartnerObjectId,AmountTotal,AmountOriginal,CommissionPS,CommissionClient,CommissionProvider,DateInput,DatePost,Status,PaymentType,PaymentNumber,ServiceId,Service,PayeeId,PayeeName,PayeeBankMfo,PayeeBankAccount,PaymentNarrative\n"))
	if err != nil {
		// return from the handler to trigger closing the connection
		return
	}

	err = a.TransactionRepository.ForEach(
		filterBuilder.GetFilters(), func(model models.Transaction) error {
			_, err := writer.Write([]byte(fmt.Sprintf("%s\n", model.ToCsvRow())))
			if err != nil {
				return err
			}

			flusher.Flush()
			return nil
		},
	)
	if err != nil {
		return
	}

	flusher.Flush()
}

func (a *Application) handleTransactionsUpload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		a.sendBadRequest(c, err.Error())
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		a.sendInternalError(c, err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	// skip header of CSV file
	scanner.Scan()

	uploadedRowsCount := 0
	err = a.TransactionRepository.UseTransaction(
		func(repository repositories.TransactionRepository) error {
			i := MaxRowsPerDbCreateRequest
			for i == MaxRowsPerDbCreateRequest {
				i = 0
				var transactions []models.Transaction
				for scanner.Scan() && i < MaxRowsPerDbCreateRequest {
					text := scanner.Text()
					if text == "" {
						continue
					}

					transaction, err := models.NewTransactionFromCSVRow(strings.Split(text, ","))
					if err != nil {
						return errors.New(fmt.Sprintf("invalid file data: %v", err))
					}

					transactions = append(transactions, transaction)
					i++
				}

				transactionCount := len(transactions)
				if transactionCount > 0 {
					err = repository.CreateBatch(transactions)
					if err != nil {
						return err
					}

					uploadedRowsCount += transactionCount
				}

				transactions = []models.Transaction{}
			}

			return nil
		},
	)

	if err != nil {
		a.sendBadRequest(c, err.Error())
		return
	}

	// if err = scanner.Err(); err != nil {
	// 	a.sendInternalError(c, err.Error())
	// 	return
	// }

	log.Printf("File %s was uploaded.\n", fileHeader.Filename)
	c.JSON(http.StatusCreated, gin.H{"row_count": uploadedRowsCount})
}
