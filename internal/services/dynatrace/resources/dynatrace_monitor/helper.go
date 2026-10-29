// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dynatrace_monitor

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
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

	userInfo := monitors.UserInfo{
		EmailAddress: pointer.To(v.EmailAddress),
		FirstName:    pointer.To(v.FirstName),
		LastName:     pointer.To(v.LastName),
	}

	if v.Country != "" {
		userInfo.Country = pointer.To(v.Country)
	}

	if v.PhoneNumber != "" {
		userInfo.PhoneNumber = pointer.To(v.PhoneNumber)
	}

	return pointer.To(userInfo)
}

func ExpandDynatraceEnvironmentProperties(input []EnvironmentProperties) *monitors.DynatraceEnvironmentProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0]

	var environmentInfo monitors.EnvironmentInfo
	if len(v.EnvironmentInfo) > 0 {
		environmentInfo = monitors.EnvironmentInfo{
			EnvironmentId: pointer.To(v.EnvironmentInfo[0].EnvironmentId),
		}
	}

	return pointer.To(monitors.DynatraceEnvironmentProperties{
		EnvironmentInfo: &environmentInfo,
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
