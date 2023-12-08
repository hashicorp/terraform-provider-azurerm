package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionProperty struct {
	Identity           *IdentityForCmk              `json:"identity,omitempty"`
	KeyVaultProperties EncryptionKeyVaultProperties `json:"keyVaultProperties"`
	Status             EncryptionStatus             `json:"status"`
}
