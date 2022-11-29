package diskencryptionsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskEncryptionSetUpdateProperties struct {
	ActiveKey                         *KeyForDiskEncryptionSet `json:"activeKey,omitempty"`
	EncryptionType                    *DiskEncryptionSetType   `json:"encryptionType,omitempty"`
	FederatedClientId                 *string                  `json:"federatedClientId,omitempty"`
	RotationToLatestKeyVersionEnabled *bool                    `json:"rotationToLatestKeyVersionEnabled,omitempty"`
}
