// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/networkmanagerroutingconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = ManagerRoutingConfigurationResource{}

type ManagerRoutingConfigurationResource struct{}

func (ManagerRoutingConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkmanagerroutingconfigurations.ValidateRoutingConfigurationID
}

func (ManagerRoutingConfigurationResource) ResourceType() string {
	return "azurerm_network_manager_routing_configuration"
}

func (ManagerRoutingConfigurationResource) ModelObject() interface{} {
	return &ManagerRoutingConfigurationResourceModel{}
}

type ManagerRoutingConfigurationResourceModel struct {
	Description      string `tfschema:"description"`
	Name             string `tfschema:"name"`
	NetworkManagerId string `tfschema:"network_manager_id"`
}

func (ManagerRoutingConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]{1,64}$`),
				"`name` must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"network_manager_id": commonschema.ResourceIDReferenceRequiredForceNew(&networkmanagerroutingconfigurations.NetworkManagerId{}),

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (ManagerRoutingConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerRoutingConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerRoutingConfigurations
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagerRoutingConfigurationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			networkManagerId, err := networkmanagerroutingconfigurations.ParseNetworkManagerID(config.NetworkManagerId)
			if err != nil {
				return err
			}

			id := networkmanagerroutingconfigurations.NewRoutingConfigurationID(subscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := networkmanagerroutingconfigurations.NetworkManagerRoutingConfiguration{
				Name: pointer.To(config.Name),
				Properties: &networkmanagerroutingconfigurations.NetworkManagerRoutingConfigurationPropertiesFormat{
					Description: pointer.To(config.Description),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ManagerRoutingConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerRoutingConfigurations

			id, err := networkmanagerroutingconfigurations.ParseRoutingConfigurationID(metadata.ResourceData.Id())
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

			networkManagerId := networkmanagerroutingconfigurations.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID()
			schema := ManagerRoutingConfigurationResourceModel{
				Name:             id.RoutingConfigurationName,
				NetworkManagerId: networkManagerId,
			}

			if model := resp.Model; model != nil && model.Properties != nil {
				schema.Description = pointer.From(model.Properties.Description)
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ManagerRoutingConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerRoutingConfigurations

			id, err := networkmanagerroutingconfigurations.ParseRoutingConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerRoutingConfigurationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			parameters := resp.Model

			if metadata.ResourceData.HasChange("description") {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r ManagerRoutingConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerRoutingConfigurations

			id, err := networkmanagerroutingconfigurations.ParseRoutingConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, networkmanagerroutingconfigurations.DeleteOperationOptions{
				Force: pointer.To(true),
			}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
