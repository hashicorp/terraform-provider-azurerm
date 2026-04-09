package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicySNAT struct {
	AutoLearnPrivateRanges *AutoLearnPrivateRangesMode `json:"autoLearnPrivateRanges,omitempty"`
	PrivateRanges          *[]string                   `json:"privateRanges,omitempty"`
}
