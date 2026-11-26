package endpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointPropertiesUpdateParameters struct {
	ContentTypesToCompress           *[]string                                                           `json:"contentTypesToCompress,omitempty"`
	DefaultOriginGroup               *ResourceReference                                                  `json:"defaultOriginGroup,omitempty"`
	DeliveryPolicy                   *EndpointPropertiesUpdateParametersDeliveryPolicy                   `json:"deliveryPolicy,omitempty"`
	GeoFilters                       *[]GeoFilter                                                        `json:"geoFilters,omitempty"`
	IsCompressionEnabled             *bool                                                               `json:"isCompressionEnabled,omitempty"`
	IsHTTPAllowed                    *bool                                                               `json:"isHttpAllowed,omitempty"`
	IsHTTPSAllowed                   *bool                                                               `json:"isHttpsAllowed,omitempty"`
	OptimizationType                 *OptimizationType                                                   `json:"optimizationType,omitempty"`
	OriginHostHeader                 *string                                                             `json:"originHostHeader,omitempty"`
	OriginPath                       *string                                                             `json:"originPath,omitempty"`
	ProbePath                        *string                                                             `json:"probePath,omitempty"`
	QueryStringCachingBehavior       *QueryStringCachingBehavior                                         `json:"queryStringCachingBehavior,omitempty"`
	UrlSigningKeys                   *[]URLSigningKey                                                    `json:"urlSigningKeys,omitempty"`
	WebApplicationFirewallPolicyLink *EndpointPropertiesUpdateParametersWebApplicationFirewallPolicyLink `json:"webApplicationFirewallPolicyLink,omitempty"`
}
