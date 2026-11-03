// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package new_relic_tag_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/tagrules"
)

type NewRelicTagRuleModel struct {
	NewRelicMonitorId      string              `tfschema:"monitor_id"`
	AadLogEnabled          bool                `tfschema:"azure_active_directory_log_enabled"`
	ActivityLogEnabled     bool                `tfschema:"activity_log_enabled"`
	LogTagFilter           []FilteringTagModel `tfschema:"log_tag_filter"`
	MetricEnabled          bool                `tfschema:"metric_enabled"`
	MetricTagFilter        []FilteringTagModel `tfschema:"metric_tag_filter"`
	SubscriptionLogEnabled bool                `tfschema:"subscription_log_enabled"`
}

type FilteringTagModel struct {
	Action tagrules.TagAction `tfschema:"action"`
	Name   string             `tfschema:"name"`
	Value  string             `tfschema:"value"`
}

func (r Resource) ModelObject() any {
	return &NewRelicTagRuleModel{}
}
