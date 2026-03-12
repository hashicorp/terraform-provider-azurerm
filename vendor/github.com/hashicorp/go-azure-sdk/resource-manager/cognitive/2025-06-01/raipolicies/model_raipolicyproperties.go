package raipolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RaiPolicyProperties struct {
	BasePolicyName   *string                   `json:"basePolicyName,omitempty"`
	ContentFilters   *[]RaiPolicyContentFilter `json:"contentFilters,omitempty"`
	CustomBlocklists *[]CustomBlocklistConfig  `json:"customBlocklists,omitempty"`
	Mode             *RaiPolicyMode            `json:"mode,omitempty"`
	Type             *RaiPolicyType            `json:"type,omitempty"`
}
