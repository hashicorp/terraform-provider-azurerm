package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Secret struct {
	Identity    *string `json:"identity,omitempty"`
	KeyVaultURL *string `json:"keyVaultUrl,omitempty"`
	Name        *string `json:"name,omitempty"`
	Value       *string `json:"value,omitempty"`
}
