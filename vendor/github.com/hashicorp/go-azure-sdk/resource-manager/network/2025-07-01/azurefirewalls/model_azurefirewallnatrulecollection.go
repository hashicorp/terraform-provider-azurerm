package azurefirewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFirewallNatRuleCollection struct {
	Etag       *string                                   `json:"etag,omitempty"`
	Id         *string                                   `json:"id,omitempty"`
	Name       *string                                   `json:"name,omitempty"`
	Properties *AzureFirewallNatRuleCollectionProperties `json:"properties,omitempty"`
}
