package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataBoxEdgeDeviceExtendedInfoProperties struct {
	ChannelIntegrityKeyName    *string             `json:"channelIntegrityKeyName,omitempty"`
	ChannelIntegrityKeyVersion *string             `json:"channelIntegrityKeyVersion,omitempty"`
	ClientSecretStoreId        *string             `json:"clientSecretStoreId,omitempty"`
	ClientSecretStoreUrl       *string             `json:"clientSecretStoreUrl,omitempty"`
	DeviceSecrets              *DeviceSecrets      `json:"deviceSecrets,omitempty"`
	EncryptionKey              *string             `json:"encryptionKey,omitempty"`
	EncryptionKeyThumbprint    *string             `json:"encryptionKeyThumbprint,omitempty"`
	KeyVaultSyncStatus         *KeyVaultSyncStatus `json:"keyVaultSyncStatus,omitempty"`
	ResourceKey                *string             `json:"resourceKey,omitempty"`
}
