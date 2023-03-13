package logz

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/logz/mgmt/2020-10-01/logz" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const PlanId100gb14days = "100gb14days"
const PlanDetails100gb14days = "100gb14days@TIDgmz7xq9ge3py"

func schemaTagFilter() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 10,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"action": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(logz.TagActionInclude),
						string(logz.TagActionExclude),
					}, false),
				},

				"value": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func getPlanDetails(plan string) (string, error) {
	if plan == PlanId100gb14days {
		return PlanDetails100gb14days, nil
	}

	return "", fmt.Errorf("cannot find plan details for id: %s", plan)
}

func getPlanId(planDetails string) (string, error) {
	if planDetails == PlanDetails100gb14days {
		return PlanId100gb14days, nil
	}

	return "", fmt.Errorf("cannot find plan id for details: %s", planDetails)
}
