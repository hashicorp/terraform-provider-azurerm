package customresourceprovider

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomRPActionRouteDefinition struct {
	Endpoint    string         `json:"endpoint"`
	Name        string         `json:"name"`
	RoutingType *ActionRouting `json:"routingType,omitempty"`
}
