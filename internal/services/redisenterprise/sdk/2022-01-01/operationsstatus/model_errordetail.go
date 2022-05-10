package operationsstatus

type ErrorDetail struct {
	AdditionalInfo *[]ErrorAdditionalInfo `json:"additionalInfo,omitempty"`
	Code           *string                `json:"code,omitempty"`
	Details        *[]ErrorDetail         `json:"details,omitempty"`
	Message        *string                `json:"message,omitempty"`
	Target         *string                `json:"target,omitempty"`
}
