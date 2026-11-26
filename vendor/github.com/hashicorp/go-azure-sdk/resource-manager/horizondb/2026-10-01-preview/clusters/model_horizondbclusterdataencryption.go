package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HorizonDbClusterDataEncryption struct {
	PrimaryEncryptionKeyCurrentVersion *string                              `json:"primaryEncryptionKeyCurrentVersion,omitempty"`
	PrimaryEncryptionKeyStatus         *HorizonDbClusterEncryptionKeyStatus `json:"primaryEncryptionKeyStatus,omitempty"`
	PrimaryKeyUri                      *string                              `json:"primaryKeyUri,omitempty"`
	PrimaryUserAssignedIdentityId      *string                              `json:"primaryUserAssignedIdentityId,omitempty"`
	Type                               *HorizonDbClusterDataEncryptionType  `json:"type,omitempty"`
}
