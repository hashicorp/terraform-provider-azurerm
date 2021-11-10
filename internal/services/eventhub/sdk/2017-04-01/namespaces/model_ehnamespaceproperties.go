package namespaces

type EHNamespaceProperties struct {
	CreatedAt              *string `json:"createdAt,omitempty"`
	IsAutoInflateEnabled   *bool   `json:"isAutoInflateEnabled,omitempty"`
	KafkaEnabled           *bool   `json:"kafkaEnabled,omitempty"`
	MaximumThroughputUnits *int64  `json:"maximumThroughputUnits,omitempty"`
	MetricId               *string `json:"metricId,omitempty"`
	ProvisioningState      *string `json:"provisioningState,omitempty"`
	ServiceBusEndpoint     *string `json:"serviceBusEndpoint,omitempty"`
	UpdatedAt              *string `json:"updatedAt,omitempty"`
}
