package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CmkKekIdentity struct {
	UseSystemAssignedIdentity *bool   `json:"useSystemAssignedIdentity,omitempty"`
	UserAssignedIdentity      *string `json:"userAssignedIdentity,omitempty"`
}
