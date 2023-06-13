package prometheusrulegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrometheusRuleGroupProperties struct {
	ClusterName *string          `json:"clusterName,omitempty"`
	Description *string          `json:"description,omitempty"`
	Enabled     *bool            `json:"enabled,omitempty"`
	Interval    *string          `json:"interval,omitempty"`
	Rules       []PrometheusRule `json:"rules"`
	Scopes      []string         `json:"scopes"`
}
