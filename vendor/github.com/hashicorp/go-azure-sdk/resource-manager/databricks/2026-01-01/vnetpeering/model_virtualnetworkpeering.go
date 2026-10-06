package vnetpeering

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkPeering struct {
	Id         *string                               `json:"id,omitempty"`
	Name       *string                               `json:"name,omitempty"`
	Properties VirtualNetworkPeeringPropertiesFormat `json:"properties"`
	Type       *string                               `json:"type,omitempty"`
}
