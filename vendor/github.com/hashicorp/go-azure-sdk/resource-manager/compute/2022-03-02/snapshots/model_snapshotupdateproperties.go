package snapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotUpdateProperties struct {
	DataAccessAuthMode           *DataAccessAuthMode           `json:"dataAccessAuthMode,omitempty"`
	DiskAccessId                 *string                       `json:"diskAccessId,omitempty"`
	DiskSizeGB                   *int64                        `json:"diskSizeGB,omitempty"`
	Encryption                   *Encryption                   `json:"encryption,omitempty"`
	EncryptionSettingsCollection *EncryptionSettingsCollection `json:"encryptionSettingsCollection,omitempty"`
	NetworkAccessPolicy          *NetworkAccessPolicy          `json:"networkAccessPolicy,omitempty"`
	OsType                       *OperatingSystemTypes         `json:"osType,omitempty"`
	PublicNetworkAccess          *PublicNetworkAccess          `json:"publicNetworkAccess,omitempty"`
	SupportedCapabilities        *SupportedCapabilities        `json:"supportedCapabilities,omitempty"`
	SupportsHibernation          *bool                         `json:"supportsHibernation,omitempty"`
}
