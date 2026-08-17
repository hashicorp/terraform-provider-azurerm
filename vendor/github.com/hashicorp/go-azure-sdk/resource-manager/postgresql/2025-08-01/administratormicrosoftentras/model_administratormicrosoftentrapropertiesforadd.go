package administratormicrosoftentras

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdministratorMicrosoftEntraPropertiesForAdd struct {
	PrincipalName *string        `json:"principalName,omitempty"`
	PrincipalType *PrincipalType `json:"principalType,omitempty"`
	TenantId      *string        `json:"tenantId,omitempty"`
}
