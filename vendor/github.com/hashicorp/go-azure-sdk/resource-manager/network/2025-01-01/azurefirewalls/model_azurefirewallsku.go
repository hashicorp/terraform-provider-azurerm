package azurefirewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFirewallSku struct {
	Name *AzureFirewallSkuName `json:"name,omitempty"`
	Tier *AzureFirewallSkuTier `json:"tier,omitempty"`
}
