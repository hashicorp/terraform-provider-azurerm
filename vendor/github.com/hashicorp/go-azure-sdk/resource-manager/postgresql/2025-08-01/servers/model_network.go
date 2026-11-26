package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Network struct {
	DelegatedSubnetResourceId   *string                         `json:"delegatedSubnetResourceId,omitempty"`
	PrivateDnsZoneArmResourceId *string                         `json:"privateDnsZoneArmResourceId,omitempty"`
	PublicNetworkAccess         *ServerPublicNetworkAccessState `json:"publicNetworkAccess,omitempty"`
}
