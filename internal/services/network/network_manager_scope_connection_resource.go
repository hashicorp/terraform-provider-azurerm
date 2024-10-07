// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/scopeconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagerScopeConnectionModel struct {
	Name             string `tfschema:"name"`
	NetworkManagerId string `tfschema:"network_manager_id"`
	ConnectionState  string `tfschema:"connection_state"`
	Description      string `tfschema:"description"`
	ResourceId       string `tfschema:"target_scope_id"`
	TenantId         string `tfschema:"tenant_id"`
}

type ManagerScopeConnectionResource struct{}

var _ sdk.ResourceWithUpdate = ManagerScopeConnectionResource{}

func (r ManagerScopeConnectionResource) ResourceType() string {
	return "azurerm_network_manager_scope_connection"
}

func (r ManagerScopeConnectionResource) ModelObject() interface{} {
	return &ManagerScopeConnectionModel{}
}

func (r ManagerScopeConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return scopeconnections.ValidateScopeConnectionID
}

func (r ManagerScopeConnectionResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: scopeconnections.ValidateNetworkManagerID,
		},

		"target_scope_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ManagerScopeConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"connection_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagerScopeConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerScopeConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.ScopeConnections
			networkManagerId, err := scopeconnections.ParseNetworkManagerID(model.NetworkManagerId)
			if err != nil {
				return err
			}

			id := scopeconnections.NewScopeConnectionID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			scopeConnection := scopeconnections.ScopeConnection{
				Properties: &scopeconnections.ScopeConnectionProperties{},
			}

			if model.Description != "" {
				scopeConnection.Properties.Description = &model.Description
			}

			if model.ResourceId != "" {
				scopeConnection.Properties.ResourceId = &model.ResourceId
			}

			if model.TenantId != "" {
				scopeConnection.Properties.TenantId = &model.TenantId
			}

			if _, err := client.CreateOrUpdate(ctx, id, scopeConnection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagerScopeConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ScopeConnections

			id, err := scopeconnections.ParseScopeConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerScopeConnectionModel
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

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Description = &model.Description
				}
			}

			if metadata.ResourceData.HasChange("target_scope_id") {
				if model.ResourceId != "" {
					properties.ResourceId = &model.ResourceId
				}
			}

			if metadata.ResourceData.HasChange("tenant_id") {
				if model.TenantId != "" {
					properties.TenantId = &model.TenantId
				}
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagerScopeConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ScopeConnections

			id, err := scopeconnections.ParseScopeConnectionID(metadata.ResourceData.Id())
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

			state := ManagerScopeConnectionModel{
				Name:             id.ScopeConnectionName,
				NetworkManagerId: scopeconnections.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID(),
			}

			if properties.ConnectionState != nil {
				state.ConnectionState = string(*properties.ConnectionState)
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.ResourceId != nil {
				state.ResourceId = *properties.ResourceId
			}

			if properties.TenantId != nil {
				state.TenantId = *properties.TenantId
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerScopeConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ScopeConnections

			id, err := scopeconnections.ParseScopeConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
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
					resp, err := client.Get(ctx, *id)
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
