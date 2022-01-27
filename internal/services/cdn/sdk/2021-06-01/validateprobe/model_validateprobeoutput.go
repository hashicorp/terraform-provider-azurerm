package validateprobe

type ValidateProbeOutput struct {
	ErrorCode *string `json:"errorCode,omitempty"`
	IsValid   *bool   `json:"isValid,omitempty"`
	Message   *string `json:"message,omitempty"`
}
