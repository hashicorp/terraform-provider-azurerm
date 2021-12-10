package servers

type ErrorAdditionalInfo struct {
	Info *interface{} `json:"info,omitempty"`
	Type *string      `json:"type,omitempty"`
}
