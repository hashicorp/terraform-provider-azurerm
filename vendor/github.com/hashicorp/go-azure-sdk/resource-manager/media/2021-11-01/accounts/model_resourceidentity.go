package accounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceIdentity struct {
	UseSystemAssignedIdentity bool    `json:"useSystemAssignedIdentity"`
	UserAssignedIdentity      *string `json:"userAssignedIdentity,omitempty"`
}
