// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotcentral

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IotCentralApplicationNetworkRuleSetResource struct{}

var (
	_ sdk.ResourceWithUpdate = IotCentralApplicationNetworkRuleSetResource{}
)

type IotCentralApplicationNetworkRuleSetModel struct {
	IotCentralApplicationId string             `tfschema:"iotcentral_application_id"`
	ApplyToDevice           bool               `tfschema:"apply_to_device"`
	DefaultAction           apps.NetworkAction `tfschema:"default_action"`
	IPRule                  []IPRule           `tfschema:"ip_rule"`
}

type IPRule struct {
	Name   string `tfschema:"name"`
	IPMask string `tfschema:"ip_mask"`
}

func (r IotCentralApplicationNetworkRuleSetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"iotcentral_application_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apps.ValidateIotAppID,
		},

		"apply_to_device": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"default_action": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(apps.NetworkActionDeny),
			ValidateFunc: validation.StringInSlice([]string{
				string(apps.NetworkActionAllow),
				string(apps.NetworkActionDeny),
			}, false),
		},

		"ip_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: iothubValidate.IoTHubIpRuleName,
					},
					"ip_mask": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.CIDR,
					},
				},
			},
		},
	}
}

func (r IotCentralApplicationNetworkRuleSetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r IotCentralApplicationNetworkRuleSetResource) ResourceType() string {
	return "azurerm_iotcentral_application_network_rule_set"
}

func (r IotCentralApplicationNetworkRuleSetResource) ModelObject() interface{} {
	return &IotCentralApplicationNetworkRuleSetModel{}
}

func (r IotCentralApplicationNetworkRuleSetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return apps.ValidateIotAppID
}

func (r IotCentralApplicationNetworkRuleSetResource) Create() sdk.ResourceFunc {
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
						return tf.ImportAsExistsError(r.ResourceType(), id.ID())
					}
				}
			}

			if model.Properties == nil {
				model.Properties = &apps.AppProperties{}
			}

			model.Properties.NetworkRuleSets = &apps.NetworkRuleSets{
				ApplyToDevices: utils.Bool(state.ApplyToDevice),
				// ApplyToIoTCentral must be set to false explicitly
				ApplyToIoTCentral: utils.Bool(false),
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

func (r IotCentralApplicationNetworkRuleSetResource) Read() sdk.ResourceFunc {
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

func (r IotCentralApplicationNetworkRuleSetResource) Update() sdk.ResourceFunc {
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
			existing.Model.Properties.NetworkRuleSets.ApplyToIoTCentral = utils.Bool(false)

			if metadata.ResourceData.HasChange("apply_to_device") {
				existing.Model.Properties.NetworkRuleSets.ApplyToDevices = utils.Bool(state.ApplyToDevice)
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

func (r IotCentralApplicationNetworkRuleSetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral.AppsClient
			id, err := apps.ParseIotAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			app, err := client.Get(ctx, *id)
			if err != nil || app.Model == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if app.Model.Properties == nil {
				app.Model.Properties = &apps.AppProperties{}
			}

			app.Model.Properties.NetworkRuleSets = nil

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *app.Model); err != nil {
				return fmt.Errorf("deleting Network Rule Set of %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandIotCentralApplicationNetworkRuleSetIPRule(input []IPRule) *[]apps.NetworkRuleSetIPRule {
	results := make([]apps.NetworkRuleSetIPRule, 0)
	for _, item := range input {
		results = append(results, apps.NetworkRuleSetIPRule{
			FilterName: utils.String(item.Name),
			IPMask:     utils.String(item.IPMask),
		})
	}
	return &results
}

func flattenIotCentralApplicationNetworkRuleSetIPRule(input *[]apps.NetworkRuleSetIPRule) []IPRule {
	if input == nil {
		return nil
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
