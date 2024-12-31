// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/firewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PlanData struct {
	BillingCycle string `json:"billing_cycle,omitempty"`
	PlanId       string `json:"plan_id,omitempty"`
	UsageType    string `json:"usage_type,omitempty"`
}

func PlanDataSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"plan_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      "panw-cloud-ngfw-payg",
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"billing_cycle": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      firewalls.BillingCycleMONTHLY,
					ValidateFunc: validation.StringInSlice(firewalls.PossibleValuesForBillingCycle(), false),
				},

				"usage_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      firewalls.UsageTypePAYG,
					ValidateFunc: validation.StringInSlice(firewalls.PossibleValuesForUsageType(), false),
				},
			},
		},
	}
}

func ExpandPlanData(input []PlanData) firewalls.PlanData {
	result := firewalls.PlanData{
		BillingCycle: firewalls.BillingCycleMONTHLY,
		UsageType:    pointer.To(firewalls.UsageTypePAYG),
		PlanId:       "panw-cloud-ngfw-payg",
	}

	if len(input) == 1 {
		p := input[0]

		if p.BillingCycle != "" {
			result.BillingCycle = firewalls.BillingCycle(p.BillingCycle)
		}

		if p.PlanId != "" {
			result.PlanId = p.PlanId
		}

		if p.UsageType != "" {
			result.UsageType = (*firewalls.UsageType)(&p.UsageType)
		}
	}

	return result
}

func FlattenPlanData(input firewalls.PlanData) []PlanData {
	result := PlanData{}

	result.BillingCycle = string(input.BillingCycle)
	result.PlanId = input.PlanId
	result.UsageType = string(*input.UsageType)

	return []PlanData{result}
}
