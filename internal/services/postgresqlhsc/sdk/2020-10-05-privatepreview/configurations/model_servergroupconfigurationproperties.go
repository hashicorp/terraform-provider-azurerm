package configurations

type ServerGroupConfigurationProperties struct {
	AllowedValues                 *string                        `json:"allowedValues,omitempty"`
	DataType                      *ConfigurationDataType         `json:"dataType,omitempty"`
	Description                   *string                        `json:"description,omitempty"`
	ServerRoleGroupConfigurations []ServerRoleGroupConfiguration `json:"serverRoleGroupConfigurations"`
}
