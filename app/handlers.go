package app

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"TraineeGolangTestTask/models"
	"TraineeGolangTestTask/repositories"
	"github.com/gin-gonic/gin"
)

func makeFilters(c *gin.Context) ([]repositories.TransactionFilter, error) {
	var filters []repositories.TransactionFilter
	if stringId := c.DefaultQuery("id", ""); stringId != "" {
		id, err := strconv.ParseUint(stringId, 10, 64)
		if err != nil {
			return nil, err
		}

		filters = append(filters, repositories.FilterById(id))
	}

	if terminalIds := c.QueryArray("terminal_id"); len(terminalIds) > 0 {
		var ids []uint64
		for _, stringId := range terminalIds {
			id, err := strconv.ParseUint(stringId, 10, 64)
			if err != nil {
				return nil, err
			}

			ids = append(ids, id)
		}
		filters = append(filters, repositories.FilterByTerminalId(ids))
	}

	if status := c.DefaultQuery("status", ""); status != "" {
		switch status {
		case "accepted", "declined":
			filters = append(filters, repositories.FilterByStatus(status))
		default:
			return nil, errors.New("value of \"status\" parameter should be either \"accepted\" or \"declined\"")
		}
	}

	if paymentType := c.DefaultQuery("payment_type", ""); paymentType != "" {
		switch paymentType {
		case "cash", "card":
			filters = append(filters, repositories.FilterByPaymentType(paymentType))
		default:
			return nil, errors.New("value of \"payment_type\" parameter should be either \"cash\" or \"card\"")
		}
	}

	if fromString := c.DefaultQuery("from", ""); fromString != "" {
		toString := c.DefaultQuery("to", "")
		if toString == "" {
			return nil, errors.New("parameter \"to\" is required when using \"from\"")
		}

		from, err := time.Parse(models.TimeLayout, fromString)
		if err != nil {
			return nil, err
		}

		to, err := time.Parse(models.TimeLayout, toString)
		if err != nil {
			return nil, err
		}

		filters = append(filters, repositories.FilterByTimeRange(from, to))
	}

	if paymentNarrative := c.DefaultQuery("payment_narrative", ""); paymentNarrative != "" {
		filters = append(filters, repositories.ContainsTextInPaymentNarrative(paymentNarrative))
	}

	return filters, nil
}

func (a *Application) handleTransactionsAsJson(c *gin.Context) {
	pageParameter := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageParameter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "the \"page\" parameter is required to be an integer value"})
		log.Println(err)
		return
	}

	filters, err := makeFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		log.Println(err)
		return
	}

	transactions := a.TransactionRepository.Filter(filters, page, a.PageSize)
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
	filters, err := makeFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		log.Println(err)
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
		// TODO: check if it is indeed internal error
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		log.Println(err)
		return
	}

	err = a.TransactionRepository.ForEach(
		filters, func(model *models.Transaction) error {
			// TODO: handle error
			_, err := writer.Write([]byte(fmt.Sprintf("%s\n", model.ToCsvRow())))
			if err != nil {
				return err
			}

			flusher.Flush()

			// TODO: must be removed later
			// time.Sleep(time.Duration(1) * time.Second)
			return nil
		},
	)
	if err != nil {
		// TODO: check if it is indeed internal error
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		log.Println(err)
		return
	}

	flusher.Flush()
}

func (a *Application) handleUpload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		log.Println(err)
		return
	}

	log.Println(fileHeader.Filename)
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		log.Println(err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	uploadedRowsCount := 0
	scanner.Scan() // skips header of CSV file
	err = a.TransactionRepository.CreateBatch(
		func(repository repositories.TransactionRepository) error {
			for scanner.Scan() {
				transaction, err := models.NewTransactionFromCSVRow(scanner.Text())
				if err != nil {
					return errors.New(fmt.Sprintf("invalid file data: %v", err))
				}

				err = repository.Create(transaction)
				if err != nil {
					return err
				}

				uploadedRowsCount++
			}

			return nil
		},
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Println(err)
		return
	}

	if err = scanner.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"row_count": uploadedRowsCount})
}
