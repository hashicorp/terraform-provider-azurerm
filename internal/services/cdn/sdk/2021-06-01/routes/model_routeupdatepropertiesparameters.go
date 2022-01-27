package routes

type RouteUpdatePropertiesParameters struct {
	CacheConfiguration  *AfdRouteCacheConfiguration   `json:"cacheConfiguration,omitempty"`
	CustomDomains       *[]ActivatedResourceReference `json:"customDomains,omitempty"`
	EnabledState        *EnabledState                 `json:"enabledState,omitempty"`
	EndpointName        *string                       `json:"endpointName,omitempty"`
	ForwardingProtocol  *ForwardingProtocol           `json:"forwardingProtocol,omitempty"`
	HttpsRedirect       *HttpsRedirect                `json:"httpsRedirect,omitempty"`
	LinkToDefaultDomain *LinkToDefaultDomain          `json:"linkToDefaultDomain,omitempty"`
	OriginGroup         *ResourceReference            `json:"originGroup,omitempty"`
	OriginPath          *string                       `json:"originPath,omitempty"`
	PatternsToMatch     *[]string                     `json:"patternsToMatch,omitempty"`
	RuleSets            *[]ResourceReference          `json:"ruleSets,omitempty"`
	SupportedProtocols  *[]AFDEndpointProtocols       `json:"supportedProtocols,omitempty"`
}
