package managedprivateendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedPrivateEndpoint struct {
	ConnectionState       *ConnectionStateProperties `json:"connectionState,omitempty"`
	Fqdns                 *[]string                  `json:"fqdns,omitempty"`
	GroupId               *string                    `json:"groupId,omitempty"`
	IsReserved            *bool                      `json:"isReserved,omitempty"`
	PrivateLinkResourceId *string                    `json:"privateLinkResourceId,omitempty"`
	ProvisioningState     *string                    `json:"provisioningState,omitempty"`
}
