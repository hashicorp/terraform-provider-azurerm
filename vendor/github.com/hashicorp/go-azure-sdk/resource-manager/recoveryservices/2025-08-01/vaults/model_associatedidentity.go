package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssociatedIdentity struct {
	OperationIdentityType *IdentityType `json:"operationIdentityType,omitempty"`
	UserAssignedIdentity  *string       `json:"userAssignedIdentity,omitempty"`
}
