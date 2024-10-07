package expressrouteconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticRoutesConfig struct {
	PropagateStaticRoutes          *bool                           `json:"propagateStaticRoutes,omitempty"`
	VnetLocalRouteOverrideCriteria *VnetLocalRouteOverrideCriteria `json:"vnetLocalRouteOverrideCriteria,omitempty"`
}
