// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iotcentral_application_network_rule_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral.AppsClient
			var state IotCentralApplicationNetworkRuleSetModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, err := apps.ParseIotAppID(state.IotCentralApplicationId)
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			app, err := client.Get(ctx, *id)
			if err != nil || app.Model == nil {
				return fmt.Errorf("checking for the presence of existing %q: %+v", id, err)
			}

			model := app.Model

			// This resource is unique to the corresponding IoT Central Application.
			// It will be created automatically along with the IoT Central Application, therefore we check whether this resource is identical to a "deleted" one
			if property := model.Properties; property != nil {
				if property.NetworkRuleSets != nil {
					if !isNetworkRuleSetNullified(*property.NetworkRuleSets) {
						if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
							return tf.ImportAsExistsError(r.ResourceType(), id.ID())
						}
					}
				}
			}

			if model.Properties == nil {
				model.Properties = &apps.AppProperties{}
			}

			model.Properties.NetworkRuleSets = &apps.NetworkRuleSets{
				ApplyToDevices: pointer.To(state.ApplyToDevice),
				// ApplyToIoTCentral must be set to false explicitly
				ApplyToIoTCentral: pointer.To(false),
				DefaultAction:     &state.DefaultAction,
				IPRules:           expandIotCentralApplicationNetworkRuleSetIPRule(state.IPRule),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("creating Network Rule Set of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandIotCentralApplicationNetworkRuleSetIPRule(input []IPRule) *[]apps.NetworkRuleSetIPRule {
	results := make([]apps.NetworkRuleSetIPRule, 0)
	for _, item := range input {
		results = append(results, apps.NetworkRuleSetIPRule{
			FilterName: pointer.To(item.Name),
			IPMask:     pointer.To(item.IPMask),
		})
	}
	return &results
}
