// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_prometheus_rule_group

type AlertPrometheusRuleGroupResourceModel struct {
	Name              string                `tfschema:"name"`
	Location          string                `tfschema:"location"`
	ResourceGroupName string                `tfschema:"resource_group_name"`
	Rule              []PrometheusRuleModel `tfschema:"rule"`
	Scopes            []string              `tfschema:"scopes"`
	ClusterName       string                `tfschema:"cluster_name"`
	Description       string                `tfschema:"description"`
	Interval          string                `tfschema:"interval"`
	RuleGroupEnabled  bool                  `tfschema:"rule_group_enabled"`
	Tags              map[string]string     `tfschema:"tags"`
}

type PrometheusRuleModel struct {
	Action          []PrometheusRuleGroupActionModel     `tfschema:"action"`
	Alert           string                               `tfschema:"alert"`
	Annotations     map[string]string                    `tfschema:"annotations"`
	Enabled         bool                                 `tfschema:"enabled"`
	Expression      string                               `tfschema:"expression"`
	For             string                               `tfschema:"for"`
	Labels          map[string]string                    `tfschema:"labels"`
	Record          string                               `tfschema:"record"`
	AlertResolution []PrometheusRuleAlertResolutionModel `tfschema:"alert_resolution"`
	Severity        int64                                `tfschema:"severity"`
}

type PrometheusRuleGroupActionModel struct {
	ActionGroupId    string            `tfschema:"action_group_id"`
	ActionProperties map[string]string `tfschema:"action_properties"`
}

type PrometheusRuleAlertResolutionModel struct {
	AutoResolved  bool   `tfschema:"auto_resolved"`
	TimeToResolve string `tfschema:"time_to_resolve"`
}

func (r Resource) ModelObject() any {
	return &AlertPrometheusRuleGroupResourceModel{}
}
