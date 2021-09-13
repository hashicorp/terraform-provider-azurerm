package servers

type AnalysisServicesServerUpdateParameters struct {
	Properties *AnalysisServicesServerMutableProperties `json:"properties,omitempty"`
	Sku        *ResourceSku                             `json:"sku,omitempty"`
	Tags       *map[string]string                       `json:"tags,omitempty"`
}
