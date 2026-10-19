// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dynatrace_tag_rules

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/tagrules"
)

func FlattenLogRules(input *tagrules.LogRules) []LogRule {
	if input == nil {
		return []LogRule{}
	}

	var logRule LogRule
	var sendAadLogs bool
	var sendActivityLogs bool
	var sendSubscriptionLogs bool

	if input.FilteringTags != nil {
		logRule.FilteringTags = FlattenFilteringTags(input.FilteringTags)
	}

	if input.SendActivityLogs != nil {
		if pointer.From(input.SendActivityLogs) == tagrules.SendActivityLogsStatusEnabled {
			sendActivityLogs = true
		} else {
			sendActivityLogs = false
		}
		logRule.SendActivityLogs = sendActivityLogs
	}

	if input.SendAadLogs != nil {
		if pointer.From(input.SendAadLogs) == tagrules.SendAadLogsStatusEnabled {
			sendAadLogs = true
		} else {
			sendAadLogs = false
		}
		logRule.SendAadLogs = sendAadLogs
	}

	if input.SendSubscriptionLogs != nil {
		if pointer.From(input.SendSubscriptionLogs) == tagrules.SendSubscriptionLogsStatusEnabled {
			sendSubscriptionLogs = true
		} else {
			sendSubscriptionLogs = false
		}
		logRule.SendSubscriptionLogs = sendSubscriptionLogs
	}

	return []LogRule{logRule}
}

func FlattenFilteringTags(input *[]tagrules.FilteringTag) []FilteringTag {
	if input == nil || len(*input) == 0 {
		return []FilteringTag{}
	}

	tags := pointer.From(input)[0]

	return []FilteringTag{
		{
			Name:   pointer.From(tags.Name),
			Value:  pointer.From(tags.Value),
			Action: pointer.FromEnum(tags.Action),
		},
	}
}

func FlattenMetricRules(input *tagrules.MetricRules) []MetricRule {
	if input == nil {
		return []MetricRule{}
	}

	var metricRule MetricRule

	if input.FilteringTags != nil {
		metricRule.FilteringTags = FlattenFilteringTags(input.FilteringTags)
	}

	if input.SendingMetrics != nil {
		if pointer.From(input.SendingMetrics) == tagrules.SendingMetricsStatusEnabled {
			metricRule.SendingMetrics = true
		} else {
			metricRule.SendingMetrics = false
		}
	}

	return []MetricRule{metricRule}
}

func ExpandMetricRules(input []MetricRule) *tagrules.MetricRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	var sendingMetrics tagrules.SendingMetricsStatus

	if v.SendingMetrics {
		sendingMetrics = tagrules.SendingMetricsStatusEnabled
	} else {
		sendingMetrics = tagrules.SendingMetricsStatusDisabled
	}

	return &tagrules.MetricRules{
		FilteringTags:  ExpandFilteringTag(v.FilteringTags),
		SendingMetrics: pointer.To(sendingMetrics),
	}
}

func ExpandLogRule(input []LogRule) *tagrules.LogRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	var sendAadLogs tagrules.SendAadLogsStatus
	var sendActivityLogs tagrules.SendActivityLogsStatus
	var sendSubscriptionLogs tagrules.SendSubscriptionLogsStatus

	if v.SendAadLogs {
		sendAadLogs = tagrules.SendAadLogsStatusEnabled
	} else {
		sendAadLogs = tagrules.SendAadLogsStatusDisabled
	}
	if v.SendActivityLogs {
		sendActivityLogs = tagrules.SendActivityLogsStatusEnabled
	} else {
		sendActivityLogs = tagrules.SendActivityLogsStatusDisabled
	}
	if v.SendSubscriptionLogs {
		sendSubscriptionLogs = tagrules.SendSubscriptionLogsStatusEnabled
	} else {
		sendSubscriptionLogs = tagrules.SendSubscriptionLogsStatusDisabled
	}

	return &tagrules.LogRules{
		FilteringTags:        ExpandFilteringTag(v.FilteringTags),
		SendAadLogs:          pointer.To(sendAadLogs),
		SendActivityLogs:     pointer.To(sendActivityLogs),
		SendSubscriptionLogs: pointer.To(sendSubscriptionLogs),
	}
}

func ExpandFilteringTag(input []FilteringTag) *[]tagrules.FilteringTag {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	action := tagrules.TagAction(v.Action)

	return &[]tagrules.FilteringTag{
		{
			Action: pointer.To(action),
			Name:   pointer.To(v.Name),
			Value:  pointer.To(v.Value),
		},
	}
}
