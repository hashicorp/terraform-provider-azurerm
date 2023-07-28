package batchendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchEndpoint struct {
	AuthMode          EndpointAuthMode           `json:"authMode"`
	Defaults          *BatchEndpointDefaults     `json:"defaults,omitempty"`
	Description       *string                    `json:"description,omitempty"`
	Keys              *EndpointAuthKeys          `json:"keys,omitempty"`
	Properties        *map[string]string         `json:"properties,omitempty"`
	ProvisioningState *EndpointProvisioningState `json:"provisioningState,omitempty"`
	ScoringUri        *string                    `json:"scoringUri,omitempty"`
	SwaggerUri        *string                    `json:"swaggerUri,omitempty"`
}
