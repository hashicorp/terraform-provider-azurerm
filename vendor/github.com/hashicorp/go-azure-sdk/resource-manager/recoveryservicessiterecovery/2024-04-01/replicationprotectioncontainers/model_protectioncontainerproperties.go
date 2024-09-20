package replicationprotectioncontainers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectionContainerProperties struct {
	FabricFriendlyName    *string                                   `json:"fabricFriendlyName,omitempty"`
	FabricSpecificDetails *ProtectionContainerFabricSpecificDetails `json:"fabricSpecificDetails,omitempty"`
	FabricType            *string                                   `json:"fabricType,omitempty"`
	FriendlyName          *string                                   `json:"friendlyName,omitempty"`
	PairingStatus         *string                                   `json:"pairingStatus,omitempty"`
	ProtectedItemCount    *int64                                    `json:"protectedItemCount,omitempty"`
	Role                  *string                                   `json:"role,omitempty"`
}
