package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsSettings struct {
	EnableProxy                 *bool     `json:"enableProxy,omitempty"`
	RequireProxyForNetworkRules *bool     `json:"requireProxyForNetworkRules,omitempty"`
	Servers                     *[]string `json:"servers,omitempty"`
}
