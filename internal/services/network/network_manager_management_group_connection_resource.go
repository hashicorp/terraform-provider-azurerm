// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkmanagerconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	managementParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	managementValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagerManagementGroupConnectionModel struct {
	Name              string `tfschema:"name"`
	ManagementGroupId string `tfschema:"management_group_id"`
	ConnectionState   string `tfschema:"connection_state"`
	Description       string `tfschema:"description"`
	NetworkManagerId  string `tfschema:"network_manager_id"`
}

type ManagerManagementGroupConnectionResource struct{}

var _ sdk.ResourceWithUpdate = ManagerManagementGroupConnectionResource{}

func (r ManagerManagementGroupConnectionResource) ResourceType() string {
	return "azurerm_network_manager_management_group_connection"
}

func (r ManagerManagementGroupConnectionResource) ModelObject() interface{} {
	return &ManagerManagementGroupConnectionModel{}
}

func (r ManagerManagementGroupConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkmanagerconnections.ValidateProviders2NetworkManagerConnectionID
}

func (r ManagerManagementGroupConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"management_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managementValidate.ManagementGroupID,
		},

		"network_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkmanagers.ValidateNetworkManagerID,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ManagerManagementGroupConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"connection_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagerManagementGroupConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerManagementGroupConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.NetworkManagerConnections
			managementGroupId, err := managementParse.ManagementGroupID(model.ManagementGroupId)
			if err != nil {
				return err
			}

			id := networkmanagerconnections.NewProviders2NetworkManagerConnectionID(managementGroupId.Name, model.Name)
			existing, err := client.ManagementGroupNetworkManagerConnectionsGet(ctx, id)
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

			if _, err := client.ManagementGroupNetworkManagerConnectionsCreateOrUpdate(ctx, id, managerConnection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagerManagementGroupConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerConnections

			id, err := networkmanagerconnections.ParseProviders2NetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerManagementGroupConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.ManagementGroupNetworkManagerConnectionsGet(ctx, *id)
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

			if _, err := client.ManagementGroupNetworkManagerConnectionsCreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagerManagementGroupConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerConnections

			id, err := networkmanagerconnections.ParseProviders2NetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.ManagementGroupNetworkManagerConnectionsGet(ctx, *id)
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
			state := ManagerManagementGroupConnectionModel{
				Name:              id.NetworkManagerConnectionName,
				ManagementGroupId: managementParse.NewManagementGroupId(id.ManagementGroupId).ID(),
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

func (r ManagerManagementGroupConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerConnections

			id, err := networkmanagerconnections.ParseProviders2NetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.ManagementGroupNetworkManagerConnectionsDelete(ctx, *id); err != nil {
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
					resp, err := client.ManagementGroupNetworkManagerConnectionsGet(ctx, *id)
					if err != nil {
						if response.WasNotFound(resp.HttpResponse) {
							return "NotFound", "NotFound", nil
						}
						return "Error", "Error", err
					}
					return resp, "Exists", nil
				},
				PollInterval:              3 * time.Second,
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
