package links

type ValidationResultItem struct {
	Description  *string                 `json:"description,omitempty"`
	ErrorCode    *string                 `json:"errorCode,omitempty"`
	ErrorMessage *string                 `json:"errorMessage,omitempty"`
	Name         *string                 `json:"name,omitempty"`
	Result       *ValidationResultStatus `json:"result,omitempty"`
}
