package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
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
	return validate.NetworkManagerScopeConnectionID
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
			ValidateFunc: validate.NetworkManagerID,
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

			client := metadata.Client.Network.ManagerScopeConnectionsClient
			networkManagerId, err := parse.NetworkManagerID(model.NetworkManagerId)
			if err != nil {
				return err
			}

			id := parse.NewNetworkManagerScopeConnectionID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroup, networkManagerId.Name, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.ScopeConnectionName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			scopeConnection := &network.ScopeConnection{
				ScopeConnectionProperties: &network.ScopeConnectionProperties{},
			}

			if model.Description != "" {
				scopeConnection.ScopeConnectionProperties.Description = &model.Description
			}

			if model.ResourceId != "" {
				scopeConnection.ScopeConnectionProperties.ResourceID = &model.ResourceId
			}

			if model.TenantId != "" {
				scopeConnection.ScopeConnectionProperties.TenantID = &model.TenantId
			}

			if _, err := client.CreateOrUpdate(ctx, *scopeConnection, id.ResourceGroup, id.NetworkManagerName, id.ScopeConnectionName); err != nil {
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
			client := metadata.Client.Network.ManagerScopeConnectionsClient

			id, err := parse.NetworkManagerScopeConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerScopeConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.ScopeConnectionName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.ScopeConnectionProperties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Description = &model.Description
				}
			}

			if metadata.ResourceData.HasChange("target_scope_id") {
				if model.ResourceId != "" {
					properties.ResourceID = &model.ResourceId
				}
			}

			if metadata.ResourceData.HasChange("tenant_id") {
				if model.TenantId != "" {
					properties.TenantID = &model.TenantId
				}
			}

			if _, err := client.CreateOrUpdate(ctx, existing, id.ResourceGroup, id.NetworkManagerName, id.ScopeConnectionName); err != nil {
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
			client := metadata.Client.Network.ManagerScopeConnectionsClient

			id, err := parse.NetworkManagerScopeConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.ScopeConnectionName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.ScopeConnectionProperties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			state := ManagerScopeConnectionModel{
				Name:             id.ScopeConnectionName,
				NetworkManagerId: parse.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName).ID(),
			}

			state.ConnectionState = string(properties.ConnectionState)

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.ResourceID != nil {
				state.ResourceId = *properties.ResourceID
			}

			if properties.TenantID != nil {
				state.TenantId = *properties.TenantID
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerScopeConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerScopeConnectionsClient

			id, err := parse.NetworkManagerScopeConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.NetworkManagerName, id.ScopeConnectionName); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("context had no deadline")
			}

			// https://github.com/Azure/azure-rest-api-specs/issues/23188
			// confirm the connection is fully deleted
			stateChangeConf := &pluginsdk.StateChangeConf{
				Pending: []string{"Exists"},
				Target:  []string{"NotFound"},
				Refresh: func() (result interface{}, state string, err error) {
					resp, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.ScopeConnectionName)
					if err != nil {
						if utils.ResponseWasNotFound(resp.Response) {
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
