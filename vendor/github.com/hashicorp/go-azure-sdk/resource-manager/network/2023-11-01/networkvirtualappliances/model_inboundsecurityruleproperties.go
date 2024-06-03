package networkvirtualappliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InboundSecurityRuleProperties struct {
	ProvisioningState *ProvisioningState       `json:"provisioningState,omitempty"`
	RuleType          *InboundSecurityRuleType `json:"ruleType,omitempty"`
	Rules             *[]InboundSecurityRules  `json:"rules,omitempty"`
}
