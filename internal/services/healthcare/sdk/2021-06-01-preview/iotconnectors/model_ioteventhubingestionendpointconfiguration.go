package iotconnectors

type IotEventHubIngestionEndpointConfiguration struct {
	ConsumerGroup                   *string `json:"consumerGroup,omitempty"`
	EventHubName                    *string `json:"eventHubName,omitempty"`
	FullyQualifiedEventHubNamespace *string `json:"fullyQualifiedEventHubNamespace,omitempty"`
}
