package customresourceprovider

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomRPResourceTypeRouteDefinition struct {
	Endpoint    string               `json:"endpoint"`
	Name        string               `json:"name"`
	RoutingType *ResourceTypeRouting `json:"routingType,omitempty"`
}
