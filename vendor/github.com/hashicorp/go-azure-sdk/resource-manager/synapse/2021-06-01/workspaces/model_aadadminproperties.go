package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AadAdminProperties struct {
	AdministratorType *string `json:"administratorType,omitempty"`
	Login             *string `json:"login,omitempty"`
	Sid               *string `json:"sid,omitempty"`
	TenantId          *string `json:"tenantId,omitempty"`
}
