package outboundnetworkdependenciesendpoints

type EndpointDependency struct {
	DomainName      *string           `json:"domainName,omitempty"`
	EndpointDetails *[]EndpointDetail `json:"endpointDetails,omitempty"`
}
