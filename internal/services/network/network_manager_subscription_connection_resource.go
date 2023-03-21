package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
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
	return validate.NetworkManagerSubscriptionConnectionID
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
			ValidateFunc: validate.NetworkManagerID,
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

			client := metadata.Client.Network.ManagerSubscriptionConnectionsClient
			subscriptionId, err := commonids.ParseSubscriptionID(model.SubscriptionId)
			if err != nil {
				return err
			}

			id := parse.NewNetworkManagerSubscriptionConnectionID(subscriptionId.SubscriptionId, model.Name)
			existing, err := client.Get(ctx, id.NetworkManagerConnectionName)
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

			if _, err := client.CreateOrUpdate(ctx, *managerConnection, id.NetworkManagerConnectionName); err != nil {
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
			client := metadata.Client.Network.ManagerSubscriptionConnectionsClient

			id, err := parse.NetworkManagerSubscriptionConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerSubscriptionConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.NetworkManagerConnectionName)
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

			if _, err := client.CreateOrUpdate(ctx, existing, id.NetworkManagerConnectionName); err != nil {
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
			client := metadata.Client.Network.ManagerSubscriptionConnectionsClient

			id, err := parse.NetworkManagerSubscriptionConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.NetworkManagerConnectionName)
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

			state := ManagerSubscriptionConnectionModel{
				Name:           id.NetworkManagerConnectionName,
				SubscriptionId: commonids.NewSubscriptionID(id.SubscriptionId).ID(),
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

func (r ManagerSubscriptionConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerSubscriptionConnectionsClient

			id, err := parse.NetworkManagerSubscriptionConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.NetworkManagerConnectionName); err != nil {
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
					resp, err := client.Get(ctx, id.NetworkManagerConnectionName)
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
