package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicySku struct {
	Tier *FirewallPolicySkuTier `json:"tier,omitempty"`
}
