// Package classification API.
//
//		Schemes: http, https
//		Host: localhost:8080
//	 Version: 0.0.1
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//		Security:
//		- api_key:
//
//		SecurityDefinitions:
//		api_key:
//		     type: apiKey
//		     name: api-key
//		     in: header
//
// swagger:meta
package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/satont/test/internal/app/api"
	"github.com/satont/test/internal/app/api/middlewares"
	"github.com/satont/test/internal/app/api/routes/transactions"
	"net/http"
)

func InitRoutes(app *api.App) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/consumer/transactions", func(r chi.Router) {
		r.Use(func(handler http.Handler) http.Handler {
			return middlewares.AttachConsumer(handler, app)
		})

		r.Get("/", transactions.GetMany(app))
		r.Get("/{transactionId}", transactions.Get(app))
		r.With(func(handler http.Handler) http.Handler {
			return middlewares.ValidateAndAttachBody(handler, app, &transactions.ReplenishValidationDto{})
		}).Post("/replenish", transactions.Replenish(app))
		r.With(func(handler http.Handler) http.Handler {
			return middlewares.ValidateAndAttachBody(handler, app, &transactions.WithDrawValidationDto{})
		}).Post("/withdraw", transactions.WithDraw(app))
	})

	return router
}
