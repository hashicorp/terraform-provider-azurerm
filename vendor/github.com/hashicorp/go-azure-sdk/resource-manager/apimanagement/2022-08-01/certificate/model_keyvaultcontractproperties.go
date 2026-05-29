package certificate

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultContractProperties struct {
	IdentityClientId *string                                     `json:"identityClientId,omitempty"`
	LastStatus       *KeyVaultLastAccessStatusContractProperties `json:"lastStatus,omitempty"`
	SecretIdentifier *string                                     `json:"secretIdentifier,omitempty"`
}
