package integrationaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultKeyAttributes struct {
	Created *int64 `json:"created,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
	Updated *int64 `json:"updated,omitempty"`
}
