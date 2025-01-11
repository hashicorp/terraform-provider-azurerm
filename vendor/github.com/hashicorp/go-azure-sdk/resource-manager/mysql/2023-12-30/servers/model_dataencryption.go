package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataEncryption struct {
	GeoBackupKeyURI                 *string             `json:"geoBackupKeyURI,omitempty"`
	GeoBackupUserAssignedIdentityId *string             `json:"geoBackupUserAssignedIdentityId,omitempty"`
	PrimaryKeyURI                   *string             `json:"primaryKeyURI,omitempty"`
	PrimaryUserAssignedIdentityId   *string             `json:"primaryUserAssignedIdentityId,omitempty"`
	Type                            *DataEncryptionType `json:"type,omitempty"`
}
