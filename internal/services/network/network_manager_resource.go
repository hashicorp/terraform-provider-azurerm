// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkmanagers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	managementGroupValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerModel struct {
	CrossTenantScopes []ManagerCrossTenantScopeModel `tfschema:"cross_tenant_scopes"`
	Scope             []ManagerScopeModel            `tfschema:"scope"`
	ScopeAccesses     []string                       `tfschema:"scope_accesses"`
	Description       string                         `tfschema:"description"`
	Name              string                         `tfschema:"name"`
	Location          string                         `tfschema:"location"`
	ResourceGroupName string                         `tfschema:"resource_group_name"`
	Tags              map[string]interface{}         `tfschema:"tags"`
}

type ManagerScopeModel struct {
	ManagementGroups []string `tfschema:"management_group_ids"`
	Subscriptions    []string `tfschema:"subscription_ids"`
}

type ManagerCrossTenantScopeModel struct {
	TenantId         string   `tfschema:"tenant_id"`
	ManagementGroups []string `tfschema:"management_groups"`
	Subscriptions    []string `tfschema:"subscriptions"`
}

type ManagerResource struct{}

func (r ManagerResource) ResourceType() string {
	return "azurerm_network_manager"
}

func (r ManagerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkmanagers.ValidateNetworkManagerID
}

func (r ManagerResource) ModelObject() interface{} {
	return &ManagerModel{}
}

func (r ManagerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"scope": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"management_group_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: managementGroupValidate.ManagementGroupID,
						},
						AtLeastOneOf: []string{"scope.0.management_group_ids", "scope.0.subscription_ids"},
					},
					"subscription_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: commonids.ValidateSubscriptionID,
						},
						AtLeastOneOf: []string{"scope.0.management_group_ids", "scope.0.subscription_ids"},
					},
				},
			},
		},

		"scope_accesses": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(networkmanagers.ConfigurationTypeConnectivity),
					string(networkmanagers.ConfigurationTypeSecurityAdmin),
				}, false),
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r ManagerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cross_tenant_scopes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"subscriptions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"management_groups": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},
	}
}

func (r ManagerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var state ManagerModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Network.NetworkManagers
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := networkmanagers.NewNetworkManagerID(subscriptionId, state.ResourceGroupName, state.Name)
			metadata.Logger.Infof("creating %s", id)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := networkmanagers.NetworkManager{
				Location: utils.String(azure.NormalizeLocation(state.Location)),
				Name:     utils.String(state.Name),
				Properties: &networkmanagers.NetworkManagerProperties{
					Description:                 utils.String(state.Description),
					NetworkManagerScopes:        expandNetworkManagerScope(state.Scope),
					NetworkManagerScopeAccesses: expandNetworkManagerScopeAccesses(state.ScopeAccesses),
				},
				Tags: utils.ExpandPtrMapStringString(state.Tags),
			}

			if _, err := client.CreateOrUpdate(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ManagerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagers
			id, err := networkmanagers.ParseNetworkManagerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving %s", *id)
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := resp.Model.Properties
			var description string
			var scope []ManagerScopeModel
			var ScopeAccesses []string
			if properties.Description != nil {
				description = *properties.Description
			}
			scope = flattenNetworkManagerScope(properties.NetworkManagerScopes)
			ScopeAccesses = flattenNetworkManagerScopeAccesses(properties.NetworkManagerScopeAccesses)

			return metadata.Encode(&ManagerModel{
				Description:       description,
				Location:          location.NormalizeNilable(resp.Model.Location),
				Name:              id.NetworkManagerName,
				ResourceGroupName: id.ResourceGroupName,
				ScopeAccesses:     ScopeAccesses,
				Scope:             scope,
				Tags:              utils.FlattenPtrMapStringString(resp.Model.Tags),
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ManagerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := networkmanagers.ParseNetworkManagerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s..", *id)
			client := metadata.Client.Network.NetworkManagers
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

			var state ManagerModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("description") {
				existing.Model.Properties.Description = utils.String(state.Description)
			}

			if metadata.ResourceData.HasChange("scope") {
				existing.Model.Properties.NetworkManagerScopes = expandNetworkManagerScope(state.Scope)
			}

			if metadata.ResourceData.HasChange("scope_accesses") {
				existing.Model.Properties.NetworkManagerScopeAccesses = expandNetworkManagerScopeAccesses(state.ScopeAccesses)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = utils.ExpandPtrMapStringString(state.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ManagerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagers
			id, err := networkmanagers.ParseNetworkManagerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s..", *id)
			err = client.DeleteThenPoll(ctx, *id, networkmanagers.DeleteOperationOptions{
				Force: utils.Bool(true),
			})
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func stringSlice(input []string) *[]string {
	return &input
}

func expandNetworkManagerScope(input []ManagerScopeModel) networkmanagers.NetworkManagerPropertiesNetworkManagerScopes {
	return networkmanagers.NetworkManagerPropertiesNetworkManagerScopes{
		ManagementGroups: stringSlice(input[0].ManagementGroups),
		Subscriptions:    stringSlice(input[0].Subscriptions),
	}
}

func expandNetworkManagerScopeAccesses(input []string) []networkmanagers.ConfigurationType {
	result := make([]networkmanagers.ConfigurationType, 0)
	for _, v := range input {
		result = append(result, networkmanagers.ConfigurationType(v))
	}
	return result
}

func flattenStringSlicePtr(input *[]string) []string {
	if input == nil {
		return make([]string, 0)
	}
	return *input
}

func flattenNetworkManagerScope(input networkmanagers.NetworkManagerPropertiesNetworkManagerScopes) []ManagerScopeModel {
	return []ManagerScopeModel{{
		ManagementGroups: flattenStringSlicePtr(input.ManagementGroups),
		Subscriptions:    flattenStringSlicePtr(input.Subscriptions),
	}}
}

func flattenNetworkManagerScopeAccesses(input []networkmanagers.ConfigurationType) []string {
	var result []string
	for _, v := range input {
		result = append(result, string(v))
	}
	return result
}
