package caches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultKeyReference struct {
	KeyUrl      string                          `json:"keyUrl"`
	SourceVault KeyVaultKeyReferenceSourceVault `json:"sourceVault"`
}
