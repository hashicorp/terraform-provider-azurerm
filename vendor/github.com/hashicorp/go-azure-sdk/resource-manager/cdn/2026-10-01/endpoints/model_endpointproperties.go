package endpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointProperties struct {
	ContentTypesToCompress           *[]string                                                           `json:"contentTypesToCompress,omitempty"`
	CustomDomains                    *[]DeepCreatedCustomDomain                                          `json:"customDomains,omitempty"`
	DefaultOriginGroup               *ResourceReference                                                  `json:"defaultOriginGroup,omitempty"`
	DeliveryPolicy                   *EndpointPropertiesUpdateParametersDeliveryPolicy                   `json:"deliveryPolicy,omitempty"`
	GeoFilters                       *[]GeoFilter                                                        `json:"geoFilters,omitempty"`
	HostName                         *string                                                             `json:"hostName,omitempty"`
	IsCompressionEnabled             *bool                                                               `json:"isCompressionEnabled,omitempty"`
	IsHTTPAllowed                    *bool                                                               `json:"isHttpAllowed,omitempty"`
	IsHTTPSAllowed                   *bool                                                               `json:"isHttpsAllowed,omitempty"`
	OptimizationType                 *OptimizationType                                                   `json:"optimizationType,omitempty"`
	OriginGroups                     *[]DeepCreatedOriginGroup                                           `json:"originGroups,omitempty"`
	OriginHostHeader                 *string                                                             `json:"originHostHeader,omitempty"`
	OriginPath                       *string                                                             `json:"originPath,omitempty"`
	Origins                          []DeepCreatedOrigin                                                 `json:"origins"`
	ProbePath                        *string                                                             `json:"probePath,omitempty"`
	ProvisioningState                *EndpointProvisioningState                                          `json:"provisioningState,omitempty"`
	QueryStringCachingBehavior       *QueryStringCachingBehavior                                         `json:"queryStringCachingBehavior,omitempty"`
	ResourceState                    *EndpointResourceState                                              `json:"resourceState,omitempty"`
	UrlSigningKeys                   *[]URLSigningKey                                                    `json:"urlSigningKeys,omitempty"`
	WebApplicationFirewallPolicyLink *EndpointPropertiesUpdateParametersWebApplicationFirewallPolicyLink `json:"webApplicationFirewallPolicyLink,omitempty"`
}
