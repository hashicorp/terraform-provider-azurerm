package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VaultCreateOrUpdateParameters struct {
	Location   string             `json:"location"`
	Properties VaultProperties    `json:"properties"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
