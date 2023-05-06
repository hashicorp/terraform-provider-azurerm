package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	managementParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	managementValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
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
	return validate.NetworkManagerManagementGroupConnectionID
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
			ValidateFunc: validate.NetworkManagerID,
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

			client := metadata.Client.Network.ManagerManagementGroupConnectionsClient
			managementGroupId, err := managementParse.ManagementGroupID(model.ManagementGroupId)
			if err != nil {
				return err
			}

			id := parse.NewNetworkManagerManagementGroupConnectionID(managementGroupId.Name, model.Name)
			existing, err := client.Get(ctx, id.ManagementGroupName, id.NetworkManagerConnectionName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			managerConnection := &network.ManagerConnection{
				ManagerConnectionProperties: &network.ManagerConnectionProperties{},
			}

			if model.Description != "" {
				managerConnection.ManagerConnectionProperties.Description = &model.Description
			}

			if model.NetworkManagerId != "" {
				managerConnection.ManagerConnectionProperties.NetworkManagerID = &model.NetworkManagerId
			}

			if _, err := client.CreateOrUpdate(ctx, *managerConnection, id.ManagementGroupName, id.NetworkManagerConnectionName); err != nil {
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
			client := metadata.Client.Network.ManagerManagementGroupConnectionsClient

			id, err := parse.NetworkManagerManagementGroupConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerManagementGroupConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ManagementGroupName, id.NetworkManagerConnectionName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.ManagerConnectionProperties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Description = &model.Description
				}
			}

			if metadata.ResourceData.HasChange("network_manager_id") {
				if model.NetworkManagerId != "" {
					properties.NetworkManagerID = &model.NetworkManagerId
				}
			}

			if _, err := client.CreateOrUpdate(ctx, existing, id.ManagementGroupName, id.NetworkManagerConnectionName); err != nil {
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
			client := metadata.Client.Network.ManagerManagementGroupConnectionsClient

			id, err := parse.NetworkManagerManagementGroupConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ManagementGroupName, id.NetworkManagerConnectionName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.ManagerConnectionProperties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			state := ManagerManagementGroupConnectionModel{
				Name:              id.NetworkManagerConnectionName,
				ManagementGroupId: managementParse.NewManagementGroupId(id.ManagementGroupName).ID(),
			}

			state.ConnectionState = string(properties.ConnectionState)

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.NetworkManagerID != nil {
				state.NetworkManagerId = *properties.NetworkManagerID
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerManagementGroupConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerManagementGroupConnectionsClient

			id, err := parse.NetworkManagerManagementGroupConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ManagementGroupName, id.NetworkManagerConnectionName); err != nil {
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
					resp, err := client.Get(ctx, id.ManagementGroupName, id.NetworkManagerConnectionName)
					if err != nil {
						if utils.ResponseWasNotFound(resp.Response) {
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
