package collection

type ServiceAuthenticationConfigurationInfo struct {
	Audience          *string `json:"audience,omitempty"`
	Authority         *string `json:"authority,omitempty"`
	SmartProxyEnabled *bool   `json:"smartProxyEnabled,omitempty"`
}
