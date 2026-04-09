package privateclouds

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Circuit struct {
	ExpressRouteID               *string `json:"expressRouteID,omitempty"`
	ExpressRoutePrivatePeeringID *string `json:"expressRoutePrivatePeeringID,omitempty"`
	PrimarySubnet                *string `json:"primarySubnet,omitempty"`
	SecondarySubnet              *string `json:"secondarySubnet,omitempty"`
}
