package collection

type ServiceCorsConfigurationInfo struct {
	AllowCredentials *bool     `json:"allowCredentials,omitempty"`
	Headers          *[]string `json:"headers,omitempty"`
	MaxAge           *int64    `json:"maxAge,omitempty"`
	Methods          *[]string `json:"methods,omitempty"`
	Origins          *[]string `json:"origins,omitempty"`
}
