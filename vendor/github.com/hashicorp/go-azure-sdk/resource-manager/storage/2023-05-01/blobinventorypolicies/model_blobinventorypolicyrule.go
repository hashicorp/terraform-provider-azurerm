package blobinventorypolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobInventoryPolicyRule struct {
	Definition  BlobInventoryPolicyDefinition `json:"definition"`
	Destination string                        `json:"destination"`
	Enabled     bool                          `json:"enabled"`
	Name        string                        `json:"name"`
}
