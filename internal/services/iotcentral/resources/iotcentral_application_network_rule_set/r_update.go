// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iotcentral_application_network_rule_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral.AppsClient
			id, err := apps.ParseIotAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state IotCentralApplicationNetworkRuleSetModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model.Properties == nil {
				existing.Model.Properties = &apps.AppProperties{}
			}

			if existing.Model.Properties.NetworkRuleSets == nil {
				existing.Model.Properties.NetworkRuleSets = &apps.NetworkRuleSets{}
			}

			// ApplyToIoTCentral must be set to false explicitly
			existing.Model.Properties.NetworkRuleSets.ApplyToIoTCentral = pointer.To(false)

			if metadata.ResourceData.HasChange("apply_to_device") {
				existing.Model.Properties.NetworkRuleSets.ApplyToDevices = pointer.To(state.ApplyToDevice)
			}

			if metadata.ResourceData.HasChange("default_action") {
				existing.Model.Properties.NetworkRuleSets.DefaultAction = &state.DefaultAction
			}

			if metadata.ResourceData.HasChange("ip_rule") {
				existing.Model.Properties.NetworkRuleSets.IPRules = expandIotCentralApplicationNetworkRuleSetIPRule(state.IPRule)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating Network Rule Set of %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
