package providers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderPermission struct {
	ApplicationId                     *string                            `json:"applicationId,omitempty"`
	ManagedByRoleDefinition           *RoleDefinition                    `json:"managedByRoleDefinition,omitempty"`
	ProviderAuthorizationConsentState *ProviderAuthorizationConsentState `json:"providerAuthorizationConsentState,omitempty"`
	RoleDefinition                    *RoleDefinition                    `json:"roleDefinition,omitempty"`
}
