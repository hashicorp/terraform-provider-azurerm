package healthbots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultProperties struct {
	KeyName      string  `json:"keyName"`
	KeyVaultUri  string  `json:"keyVaultUri"`
	KeyVersion   *string `json:"keyVersion,omitempty"`
	UserIdentity *string `json:"userIdentity,omitempty"`
}
