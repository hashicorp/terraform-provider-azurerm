package routes

type RouteProperties struct {
	CacheConfiguration  *AfdRouteCacheConfiguration   `json:"cacheConfiguration,omitempty"`
	CustomDomains       *[]ActivatedResourceReference `json:"customDomains,omitempty"`
	DeploymentStatus    *DeploymentStatus             `json:"deploymentStatus,omitempty"`
	EnabledState        *EnabledState                 `json:"enabledState,omitempty"`
	EndpointName        *string                       `json:"endpointName,omitempty"`
	ForwardingProtocol  *ForwardingProtocol           `json:"forwardingProtocol,omitempty"`
	HttpsRedirect       *HttpsRedirect                `json:"httpsRedirect,omitempty"`
	LinkToDefaultDomain *LinkToDefaultDomain          `json:"linkToDefaultDomain,omitempty"`
	OriginGroup         ResourceReference             `json:"originGroup"`
	OriginPath          *string                       `json:"originPath,omitempty"`
	PatternsToMatch     *[]string                     `json:"patternsToMatch,omitempty"`
	ProvisioningState   *AfdProvisioningState         `json:"provisioningState,omitempty"`
	RuleSets            *[]ResourceReference          `json:"ruleSets,omitempty"`
	SupportedProtocols  *[]AFDEndpointProtocols       `json:"supportedProtocols,omitempty"`
}
