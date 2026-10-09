package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageMount struct {
	CredentialsKeyVaultReference *KeyVaultReferenceWithStatus `json:"credentialsKeyVaultReference,omitempty"`
	DestinationPath              *string                      `json:"destinationPath,omitempty"`
	Name                         *string                      `json:"name,omitempty"`
	Source                       *string                      `json:"source,omitempty"`
	Type                         *StorageMountType            `json:"type,omitempty"`
}
