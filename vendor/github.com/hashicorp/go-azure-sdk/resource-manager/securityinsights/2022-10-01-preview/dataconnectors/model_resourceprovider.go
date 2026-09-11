package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceProvider struct {
	PermissionsDisplayText *string                  `json:"permissionsDisplayText,omitempty"`
	Provider               *ProviderName            `json:"provider,omitempty"`
	ProviderDisplayName    *string                  `json:"providerDisplayName,omitempty"`
	RequiredPermissions    *RequiredPermissions     `json:"requiredPermissions,omitempty"`
	Scope                  *PermissionProviderScope `json:"scope,omitempty"`
}
