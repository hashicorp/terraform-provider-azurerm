package backupvaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionSettings struct {
	InfrastructureEncryption *InfrastructureEncryptionState `json:"infrastructureEncryption,omitempty"`
	KekIdentity              *CmkKekIdentity                `json:"kekIdentity,omitempty"`
	KeyVaultProperties       *CmkKeyVaultProperties         `json:"keyVaultProperties,omitempty"`
	State                    *EncryptionState               `json:"state,omitempty"`
}
