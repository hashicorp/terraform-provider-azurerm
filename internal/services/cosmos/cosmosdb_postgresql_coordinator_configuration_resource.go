// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CosmosDbPostgreSQLCoordinatorConfigurationModel struct {
	Name      string `tfschema:"name"`
	ClusterId string `tfschema:"cluster_id"`
	Value     string `tfschema:"value"`
}

type CosmosDbPostgreSQLCoordinatorConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = CosmosDbPostgreSQLCoordinatorConfigurationResource{}

func (r CosmosDbPostgreSQLCoordinatorConfigurationResource) ResourceType() string {
	return "azurerm_cosmosdb_postgresql_coordinator_configuration"
}

func (r CosmosDbPostgreSQLCoordinatorConfigurationResource) ModelObject() interface{} {
	return &CosmosDbPostgreSQLCoordinatorConfigurationResource{}
}

func (r CosmosDbPostgreSQLCoordinatorConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return configurations.ValidateCoordinatorConfigurationID
}

func (r CosmosDbPostgreSQLCoordinatorConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: configurations.ValidateServerGroupsv2ID,
		},

		"value": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r CosmosDbPostgreSQLCoordinatorConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CosmosDbPostgreSQLCoordinatorConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CosmosDbPostgreSQLCoordinatorConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cosmos.ConfigurationsClient
			clusterId, err := configurations.ParseServerGroupsv2ID(model.ClusterId)
			if err != nil {
				return err
			}

			id := configurations.NewCoordinatorConfigurationID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.ServerGroupsv2Name, model.Name)

			locks.ByName(id.ServerGroupsv2Name, CosmosDbPostgreSQLClusterResourceName)
			defer locks.UnlockByName(id.ServerGroupsv2Name, CosmosDbPostgreSQLClusterResourceName)

			parameters := configurations.ServerConfiguration{
				Properties: &configurations.ServerConfigurationProperties{
					Value: model.Value,
				},
			}
			if err := client.UpdateOnCoordinatorThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CosmosDbPostgreSQLCoordinatorConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.ConfigurationsClient

			id, err := configurations.ParseCoordinatorConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.ServerGroupsv2Name, CosmosDbPostgreSQLClusterResourceName)
			defer locks.UnlockByName(id.ServerGroupsv2Name, CosmosDbPostgreSQLClusterResourceName)

			var model CosmosDbPostgreSQLCoordinatorConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("value") {
				parameters := configurations.ServerConfiguration{
					Properties: &configurations.ServerConfigurationProperties{
						Value: model.Value,
					},
				}

				if err := client.UpdateOnCoordinatorThenPoll(ctx, *id, parameters); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r CosmosDbPostgreSQLCoordinatorConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.ConfigurationsClient

			id, err := configurations.ParseCoordinatorConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetCoordinator(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := CosmosDbPostgreSQLCoordinatorConfigurationModel{
				Name:      id.CoordinatorConfigurationName,
				ClusterId: configurations.NewServerGroupsv2ID(id.SubscriptionId, id.ResourceGroupName, id.ServerGroupsv2Name).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Value = props.Value
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CosmosDbPostgreSQLCoordinatorConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.ConfigurationsClient

			id, err := configurations.ParseCoordinatorConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.ServerGroupsv2Name, CosmosDbPostgreSQLClusterResourceName)
			defer locks.UnlockByName(id.ServerGroupsv2Name, CosmosDbPostgreSQLClusterResourceName)

			resp, err := client.GetCoordinator(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			defaultValue := ""
			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.DefaultValue != nil {
				defaultValue = *resp.Model.Properties.DefaultValue
			}

			parameters := configurations.ServerConfiguration{
				Properties: &configurations.ServerConfigurationProperties{
					Value: defaultValue,
				},
			}

			if err = client.UpdateOnCoordinatorThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
