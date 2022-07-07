package managedidentity

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SystemAssignedIdentityProperties struct {
	ClientId        *string `json:"clientId,omitempty"`
	ClientSecretUrl *string `json:"clientSecretUrl,omitempty"`
	PrincipalId     *string `json:"principalId,omitempty"`
	TenantId        *string `json:"tenantId,omitempty"`
}
