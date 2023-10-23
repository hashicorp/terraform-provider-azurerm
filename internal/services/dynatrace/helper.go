package dynatrace

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/tagrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SchemaPlanData() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
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
					Optional:     true,
					ValidateFunc: validation.IsRFC3339Time,
				},

				"plan": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
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
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"country": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"email": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"first_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"last_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"phone_number": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaFilteringTag() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"action": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"Include",
						"Exclude",
					}, false),
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"value": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaLogRule() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"filtering_tag": SchemaFilteringTag(),

				"send_aad_logs": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"send_activity_logs": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"send_subscription_logs": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func SchemaMetricRules() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"filtering_tag": SchemaFilteringTag(),
			},
		},
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
			Action: &action,
			Name:   &v.Name,
			Value:  &v.Value,
		},
	}
}

func ExpandLogRule(input []LogRule) *tagrules.LogRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	sendAadLogs := tagrules.SendAadLogsStatusDisabled
	sendActivityLogs := tagrules.SendActivityLogsStatusDisabled
	sendSubscriptionLogs := tagrules.SendSubscriptionLogsStatusDisabled

	if v.SendAadLogs {
		sendAadLogs = tagrules.SendAadLogsStatusEnabled
	}
	if v.SendActivityLogs {
		sendActivityLogs = tagrules.SendActivityLogsStatusEnabled
	}
	if v.SendSubscriptionLogs {
		sendSubscriptionLogs = tagrules.SendSubscriptionLogsStatusEnabled
	}

	return &tagrules.LogRules{
		FilteringTags:        ExpandFilteringTag(v.FilteringTags),
		SendAadLogs:          &sendAadLogs,
		SendActivityLogs:     &sendActivityLogs,
		SendSubscriptionLogs: &sendSubscriptionLogs,
	}
}

func ExpandMetricRules(input []MetricRule) *tagrules.MetricRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0]

	return &tagrules.MetricRules{
		FilteringTags: ExpandFilteringTag(v.FilteringTags),
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
	if input == nil || len(input) == 0 {
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

func FlattenFilteringTags(input *[]tagrules.FilteringTag) []FilteringTag {
	if input == nil || len(*input) == 0 {
		return []FilteringTag{}
	}

	var name string
	var value string
	var action string
	tags := *input
	v := tags[0]

	if v.Name != nil {
		name = *v.Name
	}

	if v.Value != nil {
		value = *v.Value
	}

	if v.Action != nil {
		action = string(*v.Action)
	}

	return []FilteringTag{
		{
			Name:   name,
			Value:  value,
			Action: action,
		},
	}
}

func FlattenLogRules(input *tagrules.LogRules) []LogRule {
	if input == nil {
		return []LogRule{}
	}

	var filteringTags []FilteringTag
	sendAadLogs := false
	sendActivityLogs := false
	sendSubscriptionLogs := false

	if input.FilteringTags != nil {
		filteringTags = FlattenFilteringTags(input.FilteringTags)
	}

	if input.SendActivityLogs != nil && *input.SendActivityLogs == tagrules.SendActivityLogsStatusEnabled {
		sendActivityLogs = true
	}

	if input.SendAadLogs != nil && *input.SendAadLogs == tagrules.SendAadLogsStatusEnabled {
		sendAadLogs = true
	}

	if input.SendSubscriptionLogs != nil && *input.SendSubscriptionLogs == tagrules.SendSubscriptionLogsStatusEnabled {
		sendSubscriptionLogs = true
	}

	return []LogRule{
		{
			FilteringTags:        filteringTags,
			SendAadLogs:          sendAadLogs,
			SendActivityLogs:     sendActivityLogs,
			SendSubscriptionLogs: sendSubscriptionLogs,
		},
	}
}

func FlattenMetricRules(input *tagrules.MetricRules) []MetricRule {
	if input == nil {
		return []MetricRule{}
	}

	var filteringTags []FilteringTag

	if input.FilteringTags != nil {
		filteringTags = FlattenFilteringTags(input.FilteringTags)
	}

	return []MetricRule{
		{
			FilteringTags: filteringTags,
		},
	}
}
