package connectiongateways

type ConnectionGatewayInstallationDefinitionProperties struct {
	BackendUri         *string                     `json:"backendUri,omitempty"`
	ConnectionGateway  *ConnectionGatewayReference `json:"connectionGateway,omitempty"`
	ContactInformation *[]string                   `json:"contactInformation,omitempty"`
	Description        *string                     `json:"description,omitempty"`
	DisplayName        *string                     `json:"displayName,omitempty"`
	MachineName        *string                     `json:"machineName,omitempty"`
	Status             *interface{}                `json:"status,omitempty"`
}
