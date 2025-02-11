package credentials

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicePrincipalCredentialTypeProperties struct {
	ServicePrincipalId  *string                       `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey *AzureKeyVaultSecretReference `json:"servicePrincipalKey,omitempty"`
	Tenant              *string                       `json:"tenant,omitempty"`
}
