package operationsstatus

type ErrorResponse struct {
	Error *ErrorDetail `json:"error,omitempty"`
}
