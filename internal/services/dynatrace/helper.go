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

func FlattenDynatraceUserInfo(input []interface{}) []UserInfo {
	if len(input) == 0 {
		return []UserInfo{}
	}

	v := input[0].(map[string]interface{})
	return []UserInfo{
		{
			Country:      v["country"].(string),
			EmailAddress: v["email"].(string),
			FirstName:    v["first_name"].(string),
			LastName:     v["last_name"].(string),
			PhoneNumber:  v["phone_number"].(string),
		},
	}
}
