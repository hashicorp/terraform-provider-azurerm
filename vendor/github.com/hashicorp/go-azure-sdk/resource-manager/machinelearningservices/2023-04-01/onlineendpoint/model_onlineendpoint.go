package onlineendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OnlineEndpoint struct {
	AuthMode            EndpointAuthMode           `json:"authMode"`
	Compute             *string                    `json:"compute,omitempty"`
	Description         *string                    `json:"description,omitempty"`
	Keys                *EndpointAuthKeys          `json:"keys,omitempty"`
	MirrorTraffic       *map[string]int64          `json:"mirrorTraffic,omitempty"`
	Properties          *map[string]string         `json:"properties,omitempty"`
	ProvisioningState   *EndpointProvisioningState `json:"provisioningState,omitempty"`
	PublicNetworkAccess *PublicNetworkAccessType   `json:"publicNetworkAccess,omitempty"`
	ScoringUri          *string                    `json:"scoringUri,omitempty"`
	SwaggerUri          *string                    `json:"swaggerUri,omitempty"`
	Traffic             *map[string]int64          `json:"traffic,omitempty"`
}
