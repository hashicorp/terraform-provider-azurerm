package configurations

type ServerRoleGroupConfiguration struct {
	DefaultValue *string    `json:"defaultValue,omitempty"`
	Role         ServerRole `json:"role"`
	Source       *string    `json:"source,omitempty"`
	Value        string     `json:"value"`
}
