package blobinventorypolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobInventoryPolicySchema struct {
	Destination *string                   `json:"destination,omitempty"`
	Enabled     bool                      `json:"enabled"`
	Rules       []BlobInventoryPolicyRule `json:"rules"`
	Type        InventoryRuleType         `json:"type"`
}
