package endpoints

type EndpointPropertiesUpdateParameters struct {
	ContentTypesToCompress           *[]string                                                           `json:"contentTypesToCompress,omitempty"`
	DefaultOriginGroup               *ResourceReference                                                  `json:"defaultOriginGroup,omitempty"`
	DeliveryPolicy                   *EndpointPropertiesUpdateParametersDeliveryPolicy                   `json:"deliveryPolicy,omitempty"`
	GeoFilters                       *[]GeoFilter                                                        `json:"geoFilters,omitempty"`
	IsCompressionEnabled             *bool                                                               `json:"isCompressionEnabled,omitempty"`
	IsHttpAllowed                    *bool                                                               `json:"isHttpAllowed,omitempty"`
	IsHttpsAllowed                   *bool                                                               `json:"isHttpsAllowed,omitempty"`
	OptimizationType                 *OptimizationType                                                   `json:"optimizationType,omitempty"`
	OriginHostHeader                 *string                                                             `json:"originHostHeader,omitempty"`
	OriginPath                       *string                                                             `json:"originPath,omitempty"`
	ProbePath                        *string                                                             `json:"probePath,omitempty"`
	QueryStringCachingBehavior       *QueryStringCachingBehavior                                         `json:"queryStringCachingBehavior,omitempty"`
	UrlSigningKeys                   *[]UrlSigningKey                                                    `json:"urlSigningKeys,omitempty"`
	WebApplicationFirewallPolicyLink *EndpointPropertiesUpdateParametersWebApplicationFirewallPolicyLink `json:"webApplicationFirewallPolicyLink,omitempty"`
}
