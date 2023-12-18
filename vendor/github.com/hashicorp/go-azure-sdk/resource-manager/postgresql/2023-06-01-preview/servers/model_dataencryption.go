package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataEncryption struct {
	GeoBackupEncryptionKeyStatus    *KeyStatusEnum    `json:"geoBackupEncryptionKeyStatus,omitempty"`
	GeoBackupKeyURI                 *string           `json:"geoBackupKeyURI,omitempty"`
	GeoBackupUserAssignedIdentityId *string           `json:"geoBackupUserAssignedIdentityId,omitempty"`
	PrimaryEncryptionKeyStatus      *KeyStatusEnum    `json:"primaryEncryptionKeyStatus,omitempty"`
	PrimaryKeyURI                   *string           `json:"primaryKeyURI,omitempty"`
	PrimaryUserAssignedIdentityId   *string           `json:"primaryUserAssignedIdentityId,omitempty"`
	Type                            *ArmServerKeyType `json:"type,omitempty"`
}
