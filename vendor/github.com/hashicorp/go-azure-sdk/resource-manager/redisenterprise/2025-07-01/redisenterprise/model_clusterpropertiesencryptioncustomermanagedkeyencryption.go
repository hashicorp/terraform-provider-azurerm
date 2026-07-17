package redisenterprise

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPropertiesEncryptionCustomerManagedKeyEncryption struct {
	KeyEncryptionKeyIdentity *ClusterPropertiesEncryptionCustomerManagedKeyEncryptionKeyEncryptionKeyIdentity `json:"keyEncryptionKeyIdentity,omitempty"`
	KeyEncryptionKeyURL      *string                                                                          `json:"keyEncryptionKeyUrl,omitempty"`
}
