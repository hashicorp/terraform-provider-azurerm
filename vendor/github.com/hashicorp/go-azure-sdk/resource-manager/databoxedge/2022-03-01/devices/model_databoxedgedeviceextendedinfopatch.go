package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataBoxEdgeDeviceExtendedInfoPatch struct {
	ChannelIntegrityKeyName    *string             `json:"channelIntegrityKeyName,omitempty"`
	ChannelIntegrityKeyVersion *string             `json:"channelIntegrityKeyVersion,omitempty"`
	ClientSecretStoreId        *string             `json:"clientSecretStoreId,omitempty"`
	ClientSecretStoreURL       *string             `json:"clientSecretStoreUrl,omitempty"`
	SyncStatus                 *KeyVaultSyncStatus `json:"syncStatus,omitempty"`
}
