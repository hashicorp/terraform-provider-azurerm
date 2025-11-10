package expressroutegateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteGatewayProperties struct {
	AllowNonVirtualWanTraffic *bool                                                `json:"allowNonVirtualWanTraffic,omitempty"`
	AutoScaleConfiguration    *ExpressRouteGatewayPropertiesAutoScaleConfiguration `json:"autoScaleConfiguration,omitempty"`
	ExpressRouteConnections   *[]ExpressRouteConnection                            `json:"expressRouteConnections,omitempty"`
	ProvisioningState         *ProvisioningState                                   `json:"provisioningState,omitempty"`
	VirtualHub                VirtualHubId                                         `json:"virtualHub"`
}
