// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iotcentral_application_network_rule_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral.AppsClient
			id, err := apps.ParseIotAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.NetworkRuleSets == nil {
				return fmt.Errorf("reading Network Rule Set of %q: %+v", id, err)
			}

			networkRuleSet := resp.Model.Properties.NetworkRuleSets
			state := IotCentralApplicationNetworkRuleSetModel{
				IotCentralApplicationId: id.ID(),
				IPRule:                  flattenIotCentralApplicationNetworkRuleSetIPRule(networkRuleSet.IPRules),
			}

			applyToDevice := true
			if networkRuleSet.ApplyToDevices != nil {
				applyToDevice = *networkRuleSet.ApplyToDevices
			}
			state.ApplyToDevice = applyToDevice

			defaultAction := apps.NetworkActionDeny
			if networkRuleSet.DefaultAction != nil {
				defaultAction = *networkRuleSet.DefaultAction
			}
			state.DefaultAction = defaultAction

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func flattenIotCentralApplicationNetworkRuleSetIPRule(input *[]apps.NetworkRuleSetIPRule) []IPRule {
	if input == nil {
		return []IPRule{}
	}

	results := make([]IPRule, 0)
	for _, item := range *input {
		obj := IPRule{}

		if item.FilterName != nil {
			obj.Name = *item.FilterName
		}

		if item.IPMask != nil {
			obj.IPMask = *item.IPMask
		}

		results = append(results, obj)
	}
	return results
}

type IPRule struct {
	Name   string `tfschema:"name"`
	IPMask string `tfschema:"ip_mask"`
}

func isNetworkRuleSetNullified(networkRuleSet apps.NetworkRuleSets) bool {
	if networkRuleSet.ApplyToDevices != nil && *networkRuleSet.ApplyToDevices {
		return false
	}

	if networkRuleSet.DefaultAction != nil && *networkRuleSet.DefaultAction != apps.NetworkActionAllow {
		return false
	}

	if networkRuleSet.IPRules != nil && len(*networkRuleSet.IPRules) > 0 {
		return false
	}

	return true
}
