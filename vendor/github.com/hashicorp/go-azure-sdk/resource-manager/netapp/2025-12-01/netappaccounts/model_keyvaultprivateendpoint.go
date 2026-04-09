package netappaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultPrivateEndpoint struct {
	PrivateEndpointId *string `json:"privateEndpointId,omitempty"`
	VirtualNetworkId  *string `json:"virtualNetworkId,omitempty"`
}
