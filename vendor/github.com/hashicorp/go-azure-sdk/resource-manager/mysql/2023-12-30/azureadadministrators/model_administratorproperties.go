package azureadadministrators

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdministratorProperties struct {
	AdministratorType  *AdministratorType `json:"administratorType,omitempty"`
	IdentityResourceId *string            `json:"identityResourceId,omitempty"`
	Login              *string            `json:"login,omitempty"`
	Sid                *string            `json:"sid,omitempty"`
	TenantId           *string            `json:"tenantId,omitempty"`
}
