package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkFeaturesProperties struct {
	HybridConnections        *[]RelayServiceConnectionEntity `json:"hybridConnections,omitempty"`
	HybridConnectionsV2      *[]HybridConnection             `json:"hybridConnectionsV2,omitempty"`
	VirtualNetworkConnection *VnetInfo                       `json:"virtualNetworkConnection,omitempty"`
	VirtualNetworkName       *string                         `json:"virtualNetworkName,omitempty"`
}
