package docs

import "github.com/shopspring/decimal"

// swagger:parameters postReplenish
type ReplenishValidationDto struct {
	// in:body
	Body struct {
		Amount decimal.Decimal `json:"amount"`
	}
}

// swagger:parameters postWithDraw
type WithDrawValidationDto struct {
	// in:body
	Body struct {
		Amount decimal.Decimal `json:"amount"`
	}
}
