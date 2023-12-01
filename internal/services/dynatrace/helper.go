package dynatrace

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SchemaPlanData() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"billing_cycle": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"MONTHLY",
						"WEEKLY",
					}, false),
				},

				"effective_date": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsRFC3339Time,
				},

				"plan": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"usage_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"PAYG",
						"COMMITTED",
					}, false),
				},
			},
		},
	}
}

func SchemaUserInfo() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"country": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"email": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"first_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"last_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"phone_number": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func ExpandDynatracePlanData(input []PlanData) *monitors.PlanData {
	if len(input) == 0 {
		return nil
	}
	v := input[0]

	return &monitors.PlanData{
		BillingCycle:  &v.BillingCycle,
		EffectiveDate: &v.EffectiveDate,
		PlanDetails:   &v.PlanDetails,
		UsageType:     &v.UsageType,
	}
}

func ExpandDynatraceUserInfo(input []UserInfo) *monitors.UserInfo {
	if len(input) == 0 {
		return nil
	}
	v := input[0]

	return &monitors.UserInfo{
		Country:      &v.Country,
		EmailAddress: &v.EmailAddress,
		FirstName:    &v.FirstName,
		LastName:     &v.LastName,
		PhoneNumber:  &v.PhoneNumber,
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
		billingCycle = *input.BillingCycle
	}

	if input.EffectiveDate != nil {
		effectiveDate = *input.EffectiveDate
	}

	if input.PlanDetails != nil {
		planDetails = *input.PlanDetails
	}

	if input.UsageType != nil {
		usageType = *input.UsageType
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

	country := v["country"].(string)
	emailAddress := v["email"].(string)
	firstName := v["first_name"].(string)
	lastName := v["last_name"].(string)
	phoneNumber := v["phone_number"].(string)

	return []UserInfo{
		{
			Country:      country,
			EmailAddress: emailAddress,
			FirstName:    firstName,
			LastName:     lastName,
			PhoneNumber:  phoneNumber,
		},
	}
}
