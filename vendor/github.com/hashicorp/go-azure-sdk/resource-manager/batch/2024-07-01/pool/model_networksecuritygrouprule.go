package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityGroupRule struct {
	Access              NetworkSecurityGroupRuleAccess `json:"access"`
	Priority            int64                          `json:"priority"`
	SourceAddressPrefix string                         `json:"sourceAddressPrefix"`
	SourcePortRanges    *[]string                      `json:"sourcePortRanges,omitempty"`
}
