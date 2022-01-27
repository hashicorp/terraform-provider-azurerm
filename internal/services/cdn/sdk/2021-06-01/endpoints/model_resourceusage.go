package endpoints

type ResourceUsage struct {
	CurrentValue *int64  `json:"currentValue,omitempty"`
	Limit        *int64  `json:"limit,omitempty"`
	ResourceType *string `json:"resourceType,omitempty"`
	Unit         *string `json:"unit,omitempty"`
}
