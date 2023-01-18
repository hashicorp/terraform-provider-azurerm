package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VaultProperties struct {
	Encryption                          *VaultPropertiesEncryption                  `json:"encryption,omitempty"`
	MoveDetails                         *VaultPropertiesMoveDetails                 `json:"moveDetails,omitempty"`
	MoveState                           *ResourceMoveState                          `json:"moveState,omitempty"`
	PrivateEndpointConnections          *[]PrivateEndpointConnectionVaultProperties `json:"privateEndpointConnections,omitempty"`
	PrivateEndpointStateForBackup       *VaultPrivateEndpointState                  `json:"privateEndpointStateForBackup,omitempty"`
	PrivateEndpointStateForSiteRecovery *VaultPrivateEndpointState                  `json:"privateEndpointStateForSiteRecovery,omitempty"`
	ProvisioningState                   *string                                     `json:"provisioningState,omitempty"`
	UpgradeDetails                      *UpgradeDetails                             `json:"upgradeDetails,omitempty"`
}
