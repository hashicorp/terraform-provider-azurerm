package raipolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomBlocklistConfig struct {
	Blocking      *bool                   `json:"blocking,omitempty"`
	BlocklistName *string                 `json:"blocklistName,omitempty"`
	Source        *RaiPolicyContentSource `json:"source,omitempty"`
}
