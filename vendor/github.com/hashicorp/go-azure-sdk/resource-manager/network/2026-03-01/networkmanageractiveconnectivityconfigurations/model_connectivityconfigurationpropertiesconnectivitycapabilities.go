package networkmanageractiveconnectivityconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityConfigurationPropertiesConnectivityCapabilities struct {
	ConnectedGroupAddressOverlap        ConnectedGroupAddressOverlap        `json:"connectedGroupAddressOverlap"`
	ConnectedGroupPrivateEndpointsScale ConnectedGroupPrivateEndpointsScale `json:"connectedGroupPrivateEndpointsScale"`
	PeeringEnforcement                  PeeringEnforcement                  `json:"peeringEnforcement"`
}
