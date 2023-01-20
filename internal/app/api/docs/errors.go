package docs

// Not Found
// swagger:response notFoundError
type NotFoundError struct {
	// in: body
	Body string
}

// Internal error
// Some internal error happend
// swagger:response internalError
type InternalError struct {
	// in: body
	Body string
}

// A ValidationError is an error that is used when the required input fails validation.
// swagger:response validationError
type ValidationError struct {
	// The error message
	// in: body
	Body struct {
		// The validation message
		//
		// Required: true
		// Example: []
		Messages []string `json:"messages"`
	}
}

// Payment required
// Returns when consumer have no enough balance for withdraw
// swagger:response paymentRequired
type PaymentRequired struct {
	// in: body
	Body string
}
