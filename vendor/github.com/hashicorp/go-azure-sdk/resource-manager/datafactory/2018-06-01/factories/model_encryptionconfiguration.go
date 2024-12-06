package factories

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionConfiguration struct {
	Identity     *CMKIdentityDefinition `json:"identity,omitempty"`
	KeyName      string                 `json:"keyName"`
	KeyVersion   *string                `json:"keyVersion,omitempty"`
	VaultBaseURL string                 `json:"vaultBaseUrl"`
}
