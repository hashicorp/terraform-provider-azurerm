package endpoints

type EndpointProperties struct {
	ContentTypesToCompress           *[]string                                                           `json:"contentTypesToCompress,omitempty"`
	CustomDomains                    *[]CustomDomain                                                     `json:"customDomains,omitempty"`
	DefaultOriginGroup               *ResourceReference                                                  `json:"defaultOriginGroup,omitempty"`
	DeliveryPolicy                   *EndpointPropertiesUpdateParametersDeliveryPolicy                   `json:"deliveryPolicy,omitempty"`
	GeoFilters                       *[]GeoFilter                                                        `json:"geoFilters,omitempty"`
	HostName                         *string                                                             `json:"hostName,omitempty"`
	IsCompressionEnabled             *bool                                                               `json:"isCompressionEnabled,omitempty"`
	IsHttpAllowed                    *bool                                                               `json:"isHttpAllowed,omitempty"`
	IsHttpsAllowed                   *bool                                                               `json:"isHttpsAllowed,omitempty"`
	OptimizationType                 *OptimizationType                                                   `json:"optimizationType,omitempty"`
	OriginGroups                     *[]DeepCreatedOriginGroup                                           `json:"originGroups,omitempty"`
	OriginHostHeader                 *string                                                             `json:"originHostHeader,omitempty"`
	OriginPath                       *string                                                             `json:"originPath,omitempty"`
	Origins                          []DeepCreatedOrigin                                                 `json:"origins"`
	ProbePath                        *string                                                             `json:"probePath,omitempty"`
	ProvisioningState                *string                                                             `json:"provisioningState,omitempty"`
	QueryStringCachingBehavior       *QueryStringCachingBehavior                                         `json:"queryStringCachingBehavior,omitempty"`
	ResourceState                    *EndpointResourceState                                              `json:"resourceState,omitempty"`
	UrlSigningKeys                   *[]UrlSigningKey                                                    `json:"urlSigningKeys,omitempty"`
	WebApplicationFirewallPolicyLink *EndpointPropertiesUpdateParametersWebApplicationFirewallPolicyLink `json:"webApplicationFirewallPolicyLink,omitempty"`
}
