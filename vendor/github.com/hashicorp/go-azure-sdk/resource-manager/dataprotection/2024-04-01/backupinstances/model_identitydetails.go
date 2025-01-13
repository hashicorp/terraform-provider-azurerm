package backupinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IdentityDetails struct {
	UseSystemAssignedIdentity  *bool   `json:"useSystemAssignedIdentity,omitempty"`
	UserAssignedIdentityArmURL *string `json:"userAssignedIdentityArmUrl,omitempty"`
}
