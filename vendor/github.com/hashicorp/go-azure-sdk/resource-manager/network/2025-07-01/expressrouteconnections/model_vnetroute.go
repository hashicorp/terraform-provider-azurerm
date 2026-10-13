package expressrouteconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VnetRoute struct {
	BgpConnections     *[]SubResource      `json:"bgpConnections,omitempty"`
	StaticRoutes       *[]StaticRoute      `json:"staticRoutes,omitempty"`
	StaticRoutesConfig *StaticRoutesConfig `json:"staticRoutesConfig,omitempty"`
}
