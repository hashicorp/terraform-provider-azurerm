package restorepointcollections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultSecretReference struct {
	SecretURL   string      `json:"secretUrl"`
	SourceVault SubResource `json:"sourceVault"`
}
