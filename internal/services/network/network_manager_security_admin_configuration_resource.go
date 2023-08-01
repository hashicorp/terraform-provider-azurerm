// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/securityadminconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerSecurityAdminConfigurationModel struct {
	Name                                    string   `tfschema:"name"`
	NetworkManagerId                        string   `tfschema:"network_manager_id"`
	ApplyOnNetworkIntentPolicyBasedServices []string `tfschema:"apply_on_network_intent_policy_based_services"`
	Description                             string   `tfschema:"description"`
}

type ManagerSecurityAdminConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = ManagerSecurityAdminConfigurationResource{}

func (r ManagerSecurityAdminConfigurationResource) ResourceType() string {
	return "azurerm_network_manager_security_admin_configuration"
}

func (r ManagerSecurityAdminConfigurationResource) ModelObject() interface{} {
	return &ManagerSecurityAdminConfigurationModel{}
}

func (r ManagerSecurityAdminConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return securityadminconfigurations.ValidateSecurityAdminConfigurationID
}

func (r ManagerSecurityAdminConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: securityadminconfigurations.ValidateNetworkManagerID,
		},

		"apply_on_network_intent_policy_based_services": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(securityadminconfigurations.NetworkIntentPolicyBasedServiceAll),
					string(securityadminconfigurations.NetworkIntentPolicyBasedServiceAllowRulesOnly),
					string(securityadminconfigurations.NetworkIntentPolicyBasedServiceNone),
				}, false),
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ManagerSecurityAdminConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerSecurityAdminConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerSecurityAdminConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.SecurityAdminConfigurations
			networkManagerId, err := securityadminconfigurations.ParseNetworkManagerID(model.NetworkManagerId)
			if err != nil {
				return err
			}

			id := securityadminconfigurations.NewSecurityAdminConfigurationID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			conf := securityadminconfigurations.SecurityAdminConfiguration{
				Properties: &securityadminconfigurations.SecurityAdminConfigurationPropertiesFormat{
					ApplyOnNetworkIntentPolicyBasedServices: expandNetworkIntentPolicyBasedServiceModel(model.ApplyOnNetworkIntentPolicyBasedServices),
				},
			}

			if model.Description != "" {
				conf.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, conf); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagerSecurityAdminConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.SecurityAdminConfigurations

			id, err := securityadminconfigurations.ParseSecurityAdminConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerSecurityAdminConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties

			if metadata.ResourceData.HasChange("apply_on_network_intent_policy_based_services") {
				properties.ApplyOnNetworkIntentPolicyBasedServices = expandNetworkIntentPolicyBasedServiceModel(model.ApplyOnNetworkIntentPolicyBasedServices)
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Description = utils.String(model.Description)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagerSecurityAdminConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.SecurityAdminConfigurations

			id, err := securityadminconfigurations.ParseSecurityAdminConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties

			state := ManagerSecurityAdminConfigurationModel{
				Name:                                    id.SecurityAdminConfigurationName,
				NetworkManagerId:                        securityadminconfigurations.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID(),
				ApplyOnNetworkIntentPolicyBasedServices: flattenNetworkIntentPolicyBasedServiceModel(properties.ApplyOnNetworkIntentPolicyBasedServices),
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerSecurityAdminConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.SecurityAdminConfigurations

			id, err := securityadminconfigurations.ParseSecurityAdminConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.DeleteThenPoll(ctx, *id, securityadminconfigurations.DeleteOperationOptions{
				Force: utils.Bool(true),
			})
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandNetworkIntentPolicyBasedServiceModel(inputList []string) *[]securityadminconfigurations.NetworkIntentPolicyBasedService {
	var outputList []securityadminconfigurations.NetworkIntentPolicyBasedService
	for _, input := range inputList {
		output := securityadminconfigurations.NetworkIntentPolicyBasedService(input)

		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenNetworkIntentPolicyBasedServiceModel(inputList *[]securityadminconfigurations.NetworkIntentPolicyBasedService) []string {
	var outputList []string
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		outputList = append(outputList, string(input))
	}

	return outputList
}
