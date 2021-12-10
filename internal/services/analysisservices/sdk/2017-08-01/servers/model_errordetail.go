package servers

type ErrorDetail struct {
	AdditionalInfo *[]ErrorAdditionalInfo `json:"additionalInfo,omitempty"`
	Code           *string                `json:"code,omitempty"`
	Details        *[]ErrorDetail         `json:"details,omitempty"`
	HttpStatusCode *int64                 `json:"httpStatusCode,omitempty"`
	Message        *string                `json:"message,omitempty"`
	SubCode        *int64                 `json:"subCode,omitempty"`
	Target         *string                `json:"target,omitempty"`
	TimeStamp      *string                `json:"timeStamp,omitempty"`
}
