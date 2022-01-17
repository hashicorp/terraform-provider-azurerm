package configurations

type ServerConfigurationProperties struct {
	AllowedValues *string                `json:"allowedValues,omitempty"`
	DataType      *ConfigurationDataType `json:"dataType,omitempty"`
	DefaultValue  *string                `json:"defaultValue,omitempty"`
	Description   *string                `json:"description,omitempty"`
	Source        *string                `json:"source,omitempty"`
	Value         string                 `json:"value"`
}
