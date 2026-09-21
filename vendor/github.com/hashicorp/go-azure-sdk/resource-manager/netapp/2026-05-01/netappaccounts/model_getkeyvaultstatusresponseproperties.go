package netappaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetKeyVaultStatusResponseProperties struct {
	KeyName                  *string                    `json:"keyName,omitempty"`
	KeyVaultPrivateEndpoints *[]KeyVaultPrivateEndpoint `json:"keyVaultPrivateEndpoints,omitempty"`
	KeyVaultResourceId       *string                    `json:"keyVaultResourceId,omitempty"`
	KeyVaultUri              *string                    `json:"keyVaultUri,omitempty"`
}
