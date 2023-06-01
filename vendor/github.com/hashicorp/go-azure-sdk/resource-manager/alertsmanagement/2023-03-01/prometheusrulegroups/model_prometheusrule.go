package prometheusrulegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrometheusRule struct {
	Actions              *[]PrometheusRuleGroupAction        `json:"actions,omitempty"`
	Alert                *string                             `json:"alert,omitempty"`
	Annotations          *map[string]string                  `json:"annotations,omitempty"`
	Enabled              *bool                               `json:"enabled,omitempty"`
	Expression           string                              `json:"expression"`
	For                  *string                             `json:"for,omitempty"`
	Labels               *map[string]string                  `json:"labels,omitempty"`
	Record               *string                             `json:"record,omitempty"`
	ResolveConfiguration *PrometheusRuleResolveConfiguration `json:"resolveConfiguration,omitempty"`
	Severity             *int64                              `json:"severity,omitempty"`
}
