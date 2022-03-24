package connectiongateways

type ConnectionGatewayDefinitionProperties struct {
	BackendUri                    *string                     `json:"backendUri,omitempty"`
	ConnectionGatewayInstallation *ConnectionGatewayReference `json:"connectionGatewayInstallation,omitempty"`
	ContactInformation            *[]string                   `json:"contactInformation,omitempty"`
	Description                   *string                     `json:"description,omitempty"`
	DisplayName                   *string                     `json:"displayName,omitempty"`
	MachineName                   *string                     `json:"machineName,omitempty"`
	Status                        *interface{}                `json:"status,omitempty"`
}
