package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthConfig struct {
	ActiveDirectoryAuth *MicrosoftEntraAuth `json:"activeDirectoryAuth,omitempty"`
	PasswordAuth        *PasswordBasedAuth  `json:"passwordAuth,omitempty"`
	TenantId            *string             `json:"tenantId,omitempty"`
}
