package middlewares

import (
	"context"
	"encoding/json"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/satont/test/internal/app/api"
	"github.com/satont/test/internal/app/api/api_errors"
	"io"
	"net/http"
	"reflect"
	"strings"
)

var (
	validatorInstance = validator.New()
	en                = en_US.New()
	uni               = ut.New(en, en)
	transEN, _        = uni.GetTranslator("en_US")
)

var TagNameFunc = func(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

func ValidateAndAttachBody(next http.Handler, app *api.App, dto any) http.Handler {
	validatorInstance.RegisterTagNameFunc(TagNameFunc)
	enTranslations.RegisterDefaultTranslations(validatorInstance, transEN)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			app.Logger.Error(err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, dto); err != nil {
			app.Logger.Error(err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		err = validatorInstance.Struct(dto)
		if err != nil {
			castedErr := err.(validator.ValidationErrors)
			errors := []string{}
			for _, e := range castedErr {
				errors = append(errors, e.Translate(transEN))
			}

			response := api_errors.CreateBadRequestError(errors)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}

		ctx := context.WithValue(r.Context(), "body", dto)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
