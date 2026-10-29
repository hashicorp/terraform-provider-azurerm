package hcpopenshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomerManagedEncryptionProfile struct {
	EncryptionType *CustomerManagedEncryptionType `json:"encryptionType,omitempty"`
	Kms            *KmsEncryptionProfile          `json:"kms,omitempty"`
}
