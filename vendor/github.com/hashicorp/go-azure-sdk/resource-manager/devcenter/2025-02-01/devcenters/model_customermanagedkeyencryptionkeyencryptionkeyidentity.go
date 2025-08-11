package devcenters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomerManagedKeyEncryptionKeyEncryptionKeyIdentity struct {
	DelegatedIdentityClientId      *string       `json:"delegatedIdentityClientId,omitempty"`
	IdentityType                   *IdentityType `json:"identityType,omitempty"`
	UserAssignedIdentityResourceId *string       `json:"userAssignedIdentityResourceId,omitempty"`
}
