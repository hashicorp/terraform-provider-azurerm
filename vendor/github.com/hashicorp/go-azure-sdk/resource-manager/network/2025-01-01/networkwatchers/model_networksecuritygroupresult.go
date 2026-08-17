package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityGroupResult struct {
	EvaluatedNetworkSecurityGroups *[]EvaluatedNetworkSecurityGroup `json:"evaluatedNetworkSecurityGroups,omitempty"`
	SecurityRuleAccessResult       *SecurityRuleAccess              `json:"securityRuleAccessResult,omitempty"`
}
