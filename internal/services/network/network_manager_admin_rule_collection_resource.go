// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networkgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/adminrulecollections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerAdminRuleCollectionModel struct {
	Name                         string   `tfschema:"name"`
	SecurityAdminConfigurationId string   `tfschema:"security_admin_configuration_id"`
	NetworkGroupIds              []string `tfschema:"network_group_ids"`
	Description                  string   `tfschema:"description"`
}

type ManagerAdminRuleCollectionResource struct{}

var _ sdk.ResourceWithUpdate = ManagerAdminRuleCollectionResource{}

func (r ManagerAdminRuleCollectionResource) ResourceType() string {
	return "azurerm_network_manager_admin_rule_collection"
}

func (r ManagerAdminRuleCollectionResource) ModelObject() interface{} {
	return &ManagerAdminRuleCollectionModel{}
}

func (r ManagerAdminRuleCollectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return adminrulecollections.ValidateRuleCollectionID
}

func (r ManagerAdminRuleCollectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"security_admin_configuration_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: adminrulecollections.ValidateSecurityAdminConfigurationID,
		},

		"network_group_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: networkgroups.ValidateNetworkGroupID,
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r ManagerAdminRuleCollectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerAdminRuleCollectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerAdminRuleCollectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.AdminRuleCollections
			configurationId, err := adminrulecollections.ParseSecurityAdminConfigurationID(model.SecurityAdminConfigurationId)
			if err != nil {
				return err
			}

			id := adminrulecollections.NewRuleCollectionID(configurationId.SubscriptionId, configurationId.ResourceGroupName,
				configurationId.NetworkManagerName, configurationId.SecurityAdminConfigurationName, model.Name)
			existing, err := client.Get(ctx, id)

			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			adminRuleCollection := adminrulecollections.AdminRuleCollection{
				Properties: &adminrulecollections.AdminRuleCollectionPropertiesFormat{
					AppliesToGroups: expandNetworkManagerNetworkGroupIds(model.NetworkGroupIds),
				},
			}

			if model.Description != "" {
				adminRuleCollection.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, adminRuleCollection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagerAdminRuleCollectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.AdminRuleCollections

			id, err := adminrulecollections.ParseRuleCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerAdminRuleCollectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
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

			if metadata.ResourceData.HasChange("network_group_ids") {
				properties.AppliesToGroups = expandNetworkManagerNetworkGroupIds(model.NetworkGroupIds)
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

func (r ManagerAdminRuleCollectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.AdminRuleCollections

			id, err := adminrulecollections.ParseRuleCollectionID(metadata.ResourceData.Id())
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

			state := ManagerAdminRuleCollectionModel{
				Name:                         id.RuleCollectionName,
				SecurityAdminConfigurationId: adminrulecollections.NewSecurityAdminConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.SecurityAdminConfigurationName).ID(),
				NetworkGroupIds:              flattenNetworkManagerNetworkGroupIds(properties.AppliesToGroups),
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerAdminRuleCollectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.AdminRuleCollections

			id, err := adminrulecollections.ParseRuleCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.DeleteThenPoll(ctx, *id, adminrulecollections.DeleteOperationOptions{
				Force: utils.Bool(true),
			})
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandNetworkManagerNetworkGroupIds(inputList []string) []adminrulecollections.NetworkManagerSecurityGroupItem {
	var outputList []adminrulecollections.NetworkManagerSecurityGroupItem
	for _, v := range inputList {
		input := v
		output := adminrulecollections.NetworkManagerSecurityGroupItem{
			NetworkGroupId: input,
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenNetworkManagerNetworkGroupIds(inputList []adminrulecollections.NetworkManagerSecurityGroupItem) []string {
	var outputList []string

	for _, input := range inputList {
		outputList = append(outputList, input.NetworkGroupId)
	}

	return outputList
}
