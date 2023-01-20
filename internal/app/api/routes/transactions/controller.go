package transactions

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/satont/test/internal/app/api"
	"github.com/satont/test/internal/app/api/api_errors"
	"github.com/satont/test/internal/db/models"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
)

// swagger:response transactionResponse
type TransactionResponse struct {
	// in: body
	Body struct {
		models.Transaction
	}
}

// swagger:response transactionsResponse
type TransactionsResponse struct {
	// in: body
	Body []struct {
		models.Transaction
	}
}

// swagger:route GET /consumer/transactions/{transactionId} transactions getSingleTransaction
//
// # Get transaction by id
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//
//	Security:
//	  api_key:
//
//	Parameters:
//	  + name: transactionId
//	    in: path
//	    description: id of transaction
//	    required: true
//	    type: string
//
//	Responses:
//	  200: transactionResponse
//	  404: notFoundError
//	  500: internalError
var Get = func(app *api.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		consumer := r.Context().Value("consumer").(*models.Consumer)
		transactionId := chi.URLParam(r, "transactionId")

		transaction, err := GetService(app, consumer.ID, transactionId)
		if err != nil {
			app.Logger.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if transaction == nil {
			http.NotFound(w, r)
			return
		}

		data, err := json.Marshal(transaction)
		if err != nil {
			app.Logger.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Write(data)
	}
}

// swagger:route GET /consumer/transactions transactions getManyTransactions
//
// # Get many transactions
//
//		Produces:
//		- application/json
//
//		Schemes: http, https
//
//
//		Security:
//		  api_key:
//
//		Parameters:
//		  + name: limit
//		    in: query
//		    description: limit
//		    required: false
//			schema:
//	           type: integer
//	           max: 100
//			   min: 1
//		  + name: offset
//		    in: query
//		    description: offset
//		    required: false
//			schema:
//	           type: integer
//	           max: 100
//			   min: 1
//		  + name: status
//		    in: query
//		    description: status
//		    required: false
//			schema:
//	           type: string
//
//		Responses:
//		  200: transactionsResponse
//		  400: validationError
//		  500: internalError
var GetMany = func(app *api.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		consumer := r.Context().Value("consumer").(*models.Consumer)
		var limit uint64 = 100
		var offset uint64 = 0
		var status models.TransactionStatus

		query := r.URL.Query()

		statusParam := query.Get("status")
		if len(statusParam) != 0 {
			switch models.TransactionStatus(statusParam) {
			case models.StatusCreated:
				status = models.StatusCreated
			case models.StatusProcessing:
				status = models.StatusProcessing
			case models.StatusProcessed:
				status = models.StatusProcessed
			case models.StatusCanceled:
				status = models.StatusCanceled
			default:
				status = ""
			}
		}

		limitParam := query.Get("limit")
		if len(limitParam) != 0 {
			newLimit, err := strconv.ParseUint(limitParam, 10, 64)
			if err != nil {
				response := api_errors.CreateBadRequestError([]string{"wrong limit"})
				w.WriteHeader(http.StatusBadRequest)
				w.Write(response)
				return
			}
			if newLimit > 100 {
				response := api_errors.CreateBadRequestError([]string{"limit cannot be higher than 100"})
				w.WriteHeader(http.StatusBadRequest)
				w.Write(response)
				return
			}
			limit = newLimit
		}

		offsetParam := query.Get("offset")
		if len(offsetParam) != 0 {
			newOffset, err := strconv.ParseUint(offsetParam, 10, 64)
			if err != nil {
				response := api_errors.CreateBadRequestError([]string{"wrong offset"})
				w.WriteHeader(http.StatusBadRequest)
				w.Write(response)
				return
			}
			if newOffset > 100 {
				response := api_errors.CreateBadRequestError([]string{"offset cannot be higher than 100"})
				w.WriteHeader(http.StatusBadRequest)
				w.Write(response)
				return
			}
			offset = newOffset
		}

		transactions, err := GetManyService(app, consumer.ID, limit, offset, status)
		if err != nil {
			app.Logger.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(transactions)
		if err != nil {
			app.Logger.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Write(data)
	}
}

// swagger:route POST /consumer/transactions/replenish transactions postReplenish
//
// # Create replenish transaction
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//
//	Security:
//	  api_key:
//
//	Responses:
//	  200: emptySuccessResponse
//	  400: validationError
//	  404: notFoundError
//	  500: internalError
var Replenish = func(app *api.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		consumer := r.Context().Value("consumer").(*models.Consumer)
		dto := r.Context().Value("body").(*ReplenishValidationDto)

		if dto.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
			response := api_errors.CreateBadRequestError([]string{"amount cannot be lower or equals 0"})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
			return
		}

		err := ReplenishService(app, consumer.ID, dto)
		if err != nil {
			app.Logger.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("ok"))
	}
}

// swagger:route POST /consumer/transactions/withdraw transactions postWithDraw
//
// # Create withdraw transaction
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Security:
//	  api_key:
//
//	Responses:
//	  200: emptySuccessResponse
//	  400: validationError
//	  402: paymentRequired
//	  404: notFoundError
//	  500: internalError
var WithDraw = func(app *api.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		consumer := r.Context().Value("consumer").(*models.Consumer)
		dto := r.Context().Value("body").(*WithDrawValidationDto)

		if dto.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
			response := api_errors.CreateBadRequestError([]string{"amount cannot be lower or equals 0"})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}

		err := WithDrawService(app, consumer.ID, dto)
		if err != nil {
			if err == NotEnoughBalance {
				http.Error(w, err.Error(), http.StatusPaymentRequired)
				return
			}

			app.Logger.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("ok"))
	}
}
