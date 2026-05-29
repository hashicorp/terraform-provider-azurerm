package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceIdentity struct {
	PrincipalId *string          `json:"principalId,omitempty"`
	TenantId    *string          `json:"tenantId,omitempty"`
	Type        *MsiIdentityType `json:"type,omitempty"`
}
