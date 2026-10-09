package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VnetRouteProperties struct {
	EndAddress   *string    `json:"endAddress,omitempty"`
	RouteType    *RouteType `json:"routeType,omitempty"`
	StartAddress *string    `json:"startAddress,omitempty"`
}
