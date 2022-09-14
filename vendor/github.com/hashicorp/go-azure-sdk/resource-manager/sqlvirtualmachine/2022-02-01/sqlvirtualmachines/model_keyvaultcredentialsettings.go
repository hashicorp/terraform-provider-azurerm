package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultCredentialSettings struct {
	AzureKeyVaultUrl       *string `json:"azureKeyVaultUrl,omitempty"`
	CredentialName         *string `json:"credentialName,omitempty"`
	Enable                 *bool   `json:"enable,omitempty"`
	ServicePrincipalName   *string `json:"servicePrincipalName,omitempty"`
	ServicePrincipalSecret *string `json:"servicePrincipalSecret,omitempty"`
}
