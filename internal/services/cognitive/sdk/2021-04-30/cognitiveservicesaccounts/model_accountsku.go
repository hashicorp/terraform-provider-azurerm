package cognitiveservicesaccounts

type AccountSku struct {
	ResourceType *string `json:"resourceType,omitempty"`
	Sku          *Sku    `json:"sku,omitempty"`
}
