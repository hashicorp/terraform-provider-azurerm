package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayRouteConfigProperties struct {
	AppResourceId     *string                              `json:"appResourceId,omitempty"`
	Filters           *[]string                            `json:"filters,omitempty"`
	OpenApi           *GatewayRouteConfigOpenApiProperties `json:"openApi,omitempty"`
	Predicates        *[]string                            `json:"predicates,omitempty"`
	Protocol          *GatewayRouteConfigProtocol          `json:"protocol,omitempty"`
	ProvisioningState *GatewayProvisioningState            `json:"provisioningState,omitempty"`
	Routes            *[]GatewayApiRoute                   `json:"routes,omitempty"`
	SsoEnabled        *bool                                `json:"ssoEnabled,omitempty"`
}
