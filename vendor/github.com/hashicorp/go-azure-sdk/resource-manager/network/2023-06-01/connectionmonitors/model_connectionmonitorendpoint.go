package connectionmonitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionMonitorEndpoint struct {
	Address       *string                          `json:"address,omitempty"`
	CoverageLevel *CoverageLevel                   `json:"coverageLevel,omitempty"`
	Filter        *ConnectionMonitorEndpointFilter `json:"filter,omitempty"`
	Name          string                           `json:"name"`
	ResourceId    *string                          `json:"resourceId,omitempty"`
	Scope         *ConnectionMonitorEndpointScope  `json:"scope,omitempty"`
	Type          *EndpointType                    `json:"type,omitempty"`
}
