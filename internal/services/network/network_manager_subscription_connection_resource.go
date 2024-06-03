// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networkmanagers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkmanagerconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagerSubscriptionConnectionModel struct {
	Name             string `tfschema:"name"`
	SubscriptionId   string `tfschema:"subscription_id"`
	ConnectionState  string `tfschema:"connection_state"`
	Description      string `tfschema:"description"`
	NetworkManagerId string `tfschema:"network_manager_id"`
}

type ManagerSubscriptionConnectionResource struct{}

var _ sdk.ResourceWithUpdate = ManagerSubscriptionConnectionResource{}

func (r ManagerSubscriptionConnectionResource) ResourceType() string {
	return "azurerm_network_manager_subscription_connection"
}

func (r ManagerSubscriptionConnectionResource) ModelObject() interface{} {
	return &ManagerSubscriptionConnectionModel{}
}

func (r ManagerSubscriptionConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkmanagerconnections.ValidateNetworkManagerConnectionID
}

func (r ManagerSubscriptionConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},

		"network_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: networkmanagers.ValidateNetworkManagerID,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ManagerSubscriptionConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"connection_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagerSubscriptionConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerSubscriptionConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.NetworkManagerConnections
			subscriptionId, err := commonids.ParseSubscriptionID(model.SubscriptionId)
			if err != nil {
				return err
			}

			id := networkmanagerconnections.NewNetworkManagerConnectionID(subscriptionId.SubscriptionId, model.Name)
			existing, err := client.SubscriptionNetworkManagerConnectionsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			managerConnection := networkmanagerconnections.NetworkManagerConnection{
				Properties: &networkmanagerconnections.NetworkManagerConnectionProperties{},
			}

			if model.Description != "" {
				managerConnection.Properties.Description = &model.Description
			}

			if model.NetworkManagerId != "" {
				managerConnection.Properties.NetworkManagerId = &model.NetworkManagerId
			}

			if _, err := client.SubscriptionNetworkManagerConnectionsCreateOrUpdate(ctx, id, managerConnection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagerSubscriptionConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerConnections

			id, err := networkmanagerconnections.ParseNetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerSubscriptionConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.SubscriptionNetworkManagerConnectionsGet(ctx, *id)
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
			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Description = &model.Description
				}
			}

			if metadata.ResourceData.HasChange("network_manager_id") {
				if model.NetworkManagerId != "" {
					properties.NetworkManagerId = &model.NetworkManagerId
				}
			}

			if _, err := client.SubscriptionNetworkManagerConnectionsCreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagerSubscriptionConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerConnections

			id, err := networkmanagerconnections.ParseNetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.SubscriptionNetworkManagerConnectionsGet(ctx, *id)
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

			state := ManagerSubscriptionConnectionModel{
				Name:           id.NetworkManagerConnectionName,
				SubscriptionId: commonids.NewSubscriptionID(id.SubscriptionId).ID(),
			}

			if properties.ConnectionState != nil {
				state.ConnectionState = string(*properties.ConnectionState)
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.NetworkManagerId != nil {
				state.NetworkManagerId = *properties.NetworkManagerId
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerSubscriptionConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerConnections

			id, err := networkmanagerconnections.ParseNetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.SubscriptionNetworkManagerConnectionsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			// https://github.com/Azure/azure-rest-api-specs/issues/23188
			// confirm the connection is fully deleted
			stateChangeConf := &pluginsdk.StateChangeConf{
				Pending: []string{"Exists"},
				Target:  []string{"NotFound"},
				Refresh: func() (result interface{}, state string, err error) {
					resp, err := client.SubscriptionNetworkManagerConnectionsGet(ctx, *id)
					if err != nil {
						if response.WasNotFound(resp.HttpResponse) {
							return "NotFound", "NotFound", nil
						}
						return "Error", "Error", err
					}
					return resp, "Exists", nil
				},
				MinTimeout:                3 * time.Second,
				ContinuousTargetOccurence: 3,
				Timeout:                   time.Until(deadline),
			}

			if _, err = stateChangeConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
			}

			return nil
		},
	}
}
