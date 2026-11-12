package goldengateconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GoldenGateConnectionUpdateProperties struct {
	ConnectionType   *ConnectionType `json:"connectionType,omitempty"`
	DisplayName      *string         `json:"displayName,omitempty"`
	DoesUseSecretIds *bool           `json:"doesUseSecretIds,omitempty"`
	KeyId            *string         `json:"keyId,omitempty"`
	RoutingMethod    *RoutingMethod  `json:"routingMethod,omitempty"`
	VaultId          *string         `json:"vaultId,omitempty"`
}
