package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerExternalAdministrator struct {
	AdministratorType         *AdministratorType `json:"administratorType,omitempty"`
	AzureADOnlyAuthentication *bool              `json:"azureADOnlyAuthentication,omitempty"`
	Login                     *string            `json:"login,omitempty"`
	PrincipalType             *PrincipalType     `json:"principalType,omitempty"`
	Sid                       *string            `json:"sid,omitempty"`
	TenantId                  *string            `json:"tenantId,omitempty"`
}
