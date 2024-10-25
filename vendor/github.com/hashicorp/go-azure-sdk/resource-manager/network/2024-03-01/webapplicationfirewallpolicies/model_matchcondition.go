package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MatchCondition struct {
	MatchValues      []string                           `json:"matchValues"`
	MatchVariables   []MatchVariable                    `json:"matchVariables"`
	NegationConditon *bool                              `json:"negationConditon,omitempty"`
	Operator         WebApplicationFirewallOperator     `json:"operator"`
	Transforms       *[]WebApplicationFirewallTransform `json:"transforms,omitempty"`
}
