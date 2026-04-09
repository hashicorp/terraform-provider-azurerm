package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebApplicationFirewallCustomRule struct {
	Action             WebApplicationFirewallAction                 `json:"action"`
	Etag               *string                                      `json:"etag,omitempty"`
	GroupByUserSession *[]GroupByUserSession                        `json:"groupByUserSession,omitempty"`
	MatchConditions    []MatchCondition                             `json:"matchConditions"`
	Name               *string                                      `json:"name,omitempty"`
	Priority           int64                                        `json:"priority"`
	RateLimitDuration  *ApplicationGatewayFirewallRateLimitDuration `json:"rateLimitDuration,omitempty"`
	RateLimitThreshold *int64                                       `json:"rateLimitThreshold,omitempty"`
	RuleType           WebApplicationFirewallRuleType               `json:"ruleType"`
	State              *WebApplicationFirewallState                 `json:"state,omitempty"`
}
