package outboundnetworkdependenciesendpoints

type OutboundEnvironmentEndpoint struct {
	Category  *string               `json:"category,omitempty"`
	Endpoints *[]EndpointDependency `json:"endpoints,omitempty"`
}
