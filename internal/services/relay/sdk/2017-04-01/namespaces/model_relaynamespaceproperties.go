package namespaces

type RelayNamespaceProperties struct {
	CreatedAt          *string                `json:"createdAt,omitempty"`
	MetricId           *string                `json:"metricId,omitempty"`
	ProvisioningState  *ProvisioningStateEnum `json:"provisioningState,omitempty"`
	ServiceBusEndpoint *string                `json:"serviceBusEndpoint,omitempty"`
	UpdatedAt          *string                `json:"updatedAt,omitempty"`
}
