package trafficmanagers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HeatMapEndpoint struct {
	EndpointId *int64  `json:"endpointId,omitempty"`
	ResourceId *string `json:"resourceId,omitempty"`
}
