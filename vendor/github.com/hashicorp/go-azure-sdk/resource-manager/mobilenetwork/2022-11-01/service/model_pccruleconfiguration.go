package service

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PccRuleConfiguration struct {
	RuleName                 string                    `json:"ruleName"`
	RulePrecedence           int64                     `json:"rulePrecedence"`
	RuleQosPolicy            *PccRuleQosPolicy         `json:"ruleQosPolicy,omitempty"`
	ServiceDataFlowTemplates []ServiceDataFlowTemplate `json:"serviceDataFlowTemplates"`
	TrafficControl           *TrafficControlPermission `json:"trafficControl,omitempty"`
}
