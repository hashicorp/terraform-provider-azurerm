package databases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabasePrincipal struct {
	AppId      *string               `json:"appId,omitempty"`
	Email      *string               `json:"email,omitempty"`
	Fqn        *string               `json:"fqn,omitempty"`
	Name       string                `json:"name"`
	Role       DatabasePrincipalRole `json:"role"`
	TenantName *string               `json:"tenantName,omitempty"`
	Type       DatabasePrincipalType `json:"type"`
}
