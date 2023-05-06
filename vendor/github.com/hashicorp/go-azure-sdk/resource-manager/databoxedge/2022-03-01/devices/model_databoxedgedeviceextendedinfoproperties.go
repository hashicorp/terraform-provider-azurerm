package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataBoxEdgeDeviceExtendedInfoProperties struct {
	ChannelIntegrityKeyName        *string             `json:"channelIntegrityKeyName,omitempty"`
	ChannelIntegrityKeyVersion     *string             `json:"channelIntegrityKeyVersion,omitempty"`
	ClientSecretStoreId            *string             `json:"clientSecretStoreId,omitempty"`
	ClientSecretStoreUrl           *string             `json:"clientSecretStoreUrl,omitempty"`
	CloudWitnessContainerName      *string             `json:"cloudWitnessContainerName,omitempty"`
	CloudWitnessStorageAccountName *string             `json:"cloudWitnessStorageAccountName,omitempty"`
	CloudWitnessStorageEndpoint    *string             `json:"cloudWitnessStorageEndpoint,omitempty"`
	ClusterWitnessType             *ClusterWitnessType `json:"clusterWitnessType,omitempty"`
	DeviceSecrets                  *map[string]Secret  `json:"deviceSecrets,omitempty"`
	EncryptionKey                  *string             `json:"encryptionKey,omitempty"`
	EncryptionKeyThumbprint        *string             `json:"encryptionKeyThumbprint,omitempty"`
	FileShareWitnessLocation       *string             `json:"fileShareWitnessLocation,omitempty"`
	FileShareWitnessUsername       *string             `json:"fileShareWitnessUsername,omitempty"`
	KeyVaultSyncStatus             *KeyVaultSyncStatus `json:"keyVaultSyncStatus,omitempty"`
	ResourceKey                    *string             `json:"resourceKey,omitempty"`
}
