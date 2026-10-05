package devcenters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomerManagedKeyEncryption struct {
	KeyEncryptionKeyIdentity *CustomerManagedKeyEncryptionKeyEncryptionKeyIdentity `json:"keyEncryptionKeyIdentity,omitempty"`
	KeyEncryptionKeyURL      *string                                               `json:"keyEncryptionKeyUrl,omitempty"`
}
