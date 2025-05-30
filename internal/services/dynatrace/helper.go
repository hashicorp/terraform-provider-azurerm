// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynatrace

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/tagrules"
)

func ExpandDynatracePlanData(input []PlanData) *monitors.PlanData {
	if len(input) == 0 {
		return nil
	}
	v := input[0]

	return pointer.To(monitors.PlanData{
		BillingCycle: &v.BillingCycle,
		PlanDetails:  &v.PlanDetails,
		UsageType:    &v.UsageType,
	})
}

func ExpandDynatraceUserInfo(input []UserInfo) *monitors.UserInfo {
	if len(input) == 0 {
		return nil
	}
	v := input[0]

	return pointer.To(monitors.UserInfo{
		Country:      pointer.To(v.Country),
		EmailAddress: pointer.To(v.EmailAddress),
		FirstName:    pointer.To(v.FirstName),
		LastName:     pointer.To(v.LastName),
		PhoneNumber:  pointer.To(v.PhoneNumber),
	})
}

func FlattenDynatraceEnvironmentProperties(input *monitors.DynatraceEnvironmentProperties) []EnvironmentProperties {
	if input == nil {
		return []EnvironmentProperties{}
	}
	var properties EnvironmentProperties
	if input.EnvironmentInfo.EnvironmentId != nil {
		environmentId := pointer.From(input.EnvironmentInfo.EnvironmentId)
		properties.EnvironmentInfo = []EnvironmentInfo{
			{
				EnvironmentId: environmentId,
			},
		}
	}
	return []EnvironmentProperties{properties}
}

func FlattenDynatracePlanData(input *monitors.PlanData) []PlanData {
	if input == nil {
		return []PlanData{}
	}

	var billingCycle string
	var effectiveDate string
	var planDetails string
	var usageType string

	if input.BillingCycle != nil {
		billingCycle = pointer.From(input.BillingCycle)
	}

	if input.EffectiveDate != nil {
		effectiveDate = pointer.From(input.EffectiveDate)
	}

	if input.PlanDetails != nil {
		planDetails = pointer.From(input.PlanDetails)
	}

	if input.UsageType != nil {
		usageType = pointer.From(input.UsageType)
	}

	return []PlanData{
		{
			BillingCycle:  billingCycle,
			EffectiveDate: effectiveDate,
			PlanDetails:   planDetails,
			UsageType:     usageType,
		},
	}
}

func FlattenDynatraceUserInfo(input *monitors.UserInfo) []UserInfo {
	if input == nil {
		return []UserInfo{}
	}

	return []UserInfo{
		{
			Country:      pointer.From(input.Country),
			EmailAddress: pointer.From(input.EmailAddress),
			FirstName:    pointer.From(input.FirstName),
			LastName:     pointer.From(input.LastName),
			PhoneNumber:  pointer.From(input.PhoneNumber),
		},
	}
}

func FlattenLogRules(input *tagrules.LogRules) []LogRule {
	if input == nil {
		return []LogRule{}
	}

	var logRule LogRule
	var sendAadLogs bool
	var sendActivityLogs bool
	var sendSubscriptionLogs bool

	if input.FilteringTags != nil {
		filteringTags := FlattenFilteringTags(input.FilteringTags)
		logRule.FilteringTags = filteringTags
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
			Action: string(pointer.From(tags.Action)),
		},
	}
}

func FlattenMetricRules(input *tagrules.MetricRules) []MetricRule {
	if input == nil {
		return []MetricRule{}
	}

	var metricRule MetricRule

	if input.FilteringTags != nil {
		filteringTags := FlattenFilteringTags(input.FilteringTags)
		metricRule.FilteringTags = filteringTags
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
