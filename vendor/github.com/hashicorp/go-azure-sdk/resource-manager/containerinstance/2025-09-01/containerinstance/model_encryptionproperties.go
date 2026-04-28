package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionProperties struct {
	Identity     *string `json:"identity,omitempty"`
	KeyName      string  `json:"keyName"`
	KeyVersion   string  `json:"keyVersion"`
	VaultBaseURL string  `json:"vaultBaseUrl"`
}
