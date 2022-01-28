package iotconnectors

type IotConnectorProperties struct {
	DeviceMapping                  *IotMappingProperties                      `json:"deviceMapping,omitempty"`
	IngestionEndpointConfiguration *IotEventHubIngestionEndpointConfiguration `json:"ingestionEndpointConfiguration,omitempty"`
	ProvisioningState              *ProvisioningState                         `json:"provisioningState,omitempty"`
}
