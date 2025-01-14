// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynatrace

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

	environmentInfo := FlattenDynatraceEnvironmentInfo(input.EnvironmentInfo)

	return []EnvironmentProperties{
		{
			EnvironmentInfo: environmentInfo,
		},
	}
}

func FlattenDynatraceEnvironmentInfo(input *monitors.EnvironmentInfo) []EnvironmentInfo {
	if input == nil {
		return []EnvironmentInfo{}
	}

	return []EnvironmentInfo{
		{
			EnvironmentId: pointer.From(input.EnvironmentId),
		},
	}
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
