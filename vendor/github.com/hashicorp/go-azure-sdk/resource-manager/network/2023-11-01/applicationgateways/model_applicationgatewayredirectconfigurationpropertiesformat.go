package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayRedirectConfigurationPropertiesFormat struct {
	IncludePath         *bool                           `json:"includePath,omitempty"`
	IncludeQueryString  *bool                           `json:"includeQueryString,omitempty"`
	PathRules           *[]SubResource                  `json:"pathRules,omitempty"`
	RedirectType        *ApplicationGatewayRedirectType `json:"redirectType,omitempty"`
	RequestRoutingRules *[]SubResource                  `json:"requestRoutingRules,omitempty"`
	TargetListener      *SubResource                    `json:"targetListener,omitempty"`
	TargetURL           *string                         `json:"targetUrl,omitempty"`
	UrlPathMaps         *[]SubResource                  `json:"urlPathMaps,omitempty"`
}
