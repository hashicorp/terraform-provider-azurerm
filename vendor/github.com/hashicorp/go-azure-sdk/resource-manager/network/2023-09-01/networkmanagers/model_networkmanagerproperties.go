package networkmanagers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkManagerProperties struct {
	Description                 *string                                      `json:"description,omitempty"`
	NetworkManagerScopeAccesses []ConfigurationType                          `json:"networkManagerScopeAccesses"`
	NetworkManagerScopes        NetworkManagerPropertiesNetworkManagerScopes `json:"networkManagerScopes"`
	ProvisioningState           *ProvisioningState                           `json:"provisioningState,omitempty"`
	ResourceGuid                *string                                      `json:"resourceGuid,omitempty"`
}
