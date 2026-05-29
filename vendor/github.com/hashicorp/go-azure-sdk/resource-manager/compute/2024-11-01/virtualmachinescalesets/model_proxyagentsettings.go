package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProxyAgentSettings struct {
	Enabled          *bool                 `json:"enabled,omitempty"`
	Imds             *HostEndpointSettings `json:"imds,omitempty"`
	KeyIncarnationId *int64                `json:"keyIncarnationId,omitempty"`
	Mode             *Mode                 `json:"mode,omitempty"`
	WireServer       *HostEndpointSettings `json:"wireServer,omitempty"`
}
