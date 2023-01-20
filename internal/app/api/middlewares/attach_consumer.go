package middlewares

import (
	"context"
	"github.com/satont/test/internal/app/api"
	"net/http"
)

func AttachConsumer(next http.Handler, app *api.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("api-key")
		if apiKey == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		consumer, err := app.ConsumersService.GetByApiKey(apiKey)
		if err != nil {
			app.Logger.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if consumer == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "consumer", consumer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
