// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dynatrace_tag_rules

type TagRulesResourceModel struct {
	Name        string       `tfschema:"name"`
	Monitor     string       `tfschema:"monitor_id"`
	LogRules    []LogRule    `tfschema:"log_rule"`
	MetricRules []MetricRule `tfschema:"metric_rule"`
}

func (r Resource) ModelObject() any {
	return &TagRulesResourceModel{}
}
