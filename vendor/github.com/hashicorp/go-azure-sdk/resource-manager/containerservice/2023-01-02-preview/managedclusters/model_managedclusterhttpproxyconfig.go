package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterHTTPProxyConfig struct {
	EffectiveNoProxy *[]string `json:"effectiveNoProxy,omitempty"`
	HTTPProxy        *string   `json:"httpProxy,omitempty"`
	HTTPSProxy       *string   `json:"httpsProxy,omitempty"`
	NoProxy          *[]string `json:"noProxy,omitempty"`
	TrustedCa        *string   `json:"trustedCa,omitempty"`
}
