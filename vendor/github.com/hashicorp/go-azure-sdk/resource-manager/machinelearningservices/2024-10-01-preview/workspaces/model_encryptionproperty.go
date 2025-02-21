package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionProperty struct {
	CosmosDbResourceId       *string            `json:"cosmosDbResourceId,omitempty"`
	Identity                 *IdentityForCmk    `json:"identity,omitempty"`
	KeyVaultProperties       KeyVaultProperties `json:"keyVaultProperties"`
	SearchAccountResourceId  *string            `json:"searchAccountResourceId,omitempty"`
	Status                   EncryptionStatus   `json:"status"`
	StorageAccountResourceId *string            `json:"storageAccountResourceId,omitempty"`
}
