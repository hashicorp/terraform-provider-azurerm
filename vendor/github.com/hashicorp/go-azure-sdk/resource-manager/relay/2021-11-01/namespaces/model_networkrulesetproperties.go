package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkRuleSetProperties struct {
	DefaultAction       *DefaultAction       `json:"defaultAction,omitempty"`
	IPRules             *[]NWRuleSetIPRules  `json:"ipRules,omitempty"`
	PublicNetworkAccess *PublicNetworkAccess `json:"publicNetworkAccess,omitempty"`
}
