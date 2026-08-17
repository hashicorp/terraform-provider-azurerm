package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyThreatIntelWhitelist struct {
	Fqdns       *[]string `json:"fqdns,omitempty"`
	IPAddresses *[]string `json:"ipAddresses,omitempty"`
}
