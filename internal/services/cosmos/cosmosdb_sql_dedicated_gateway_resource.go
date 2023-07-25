// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/sqldedicatedgateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CosmosDbSqlDedicatedGatewayModel struct {
	CosmosDbAccountId string                          `tfschema:"cosmosdb_account_id"`
	InstanceCount     int64                           `tfschema:"instance_count"`
	InstanceSize      sqldedicatedgateway.ServiceSize `tfschema:"instance_size"`
}

type CosmosDbSqlDedicatedGatewayResource struct{}

var _ sdk.ResourceWithUpdate = CosmosDbSqlDedicatedGatewayResource{}

func (r CosmosDbSqlDedicatedGatewayResource) ResourceType() string {
	return "azurerm_cosmosdb_sql_dedicated_gateway"
}

func (r CosmosDbSqlDedicatedGatewayResource) ModelObject() interface{} {
	return &CosmosDbSqlDedicatedGatewayModel{}
}

func (r CosmosDbSqlDedicatedGatewayResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sqldedicatedgateway.ValidateServiceID
}

func (r CosmosDbSqlDedicatedGatewayResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cosmosdb_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cosmosdb.ValidateDatabaseAccountID,
		},

		"instance_size": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sqldedicatedgateway.ServiceSizeCosmosPointDFours),
				string(sqldedicatedgateway.ServiceSizeCosmosPointDEights),
				string(sqldedicatedgateway.ServiceSizeCosmosPointDOneSixs),
			}, false),
		},

		"instance_count": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 5),
		},
	}
}

func (r CosmosDbSqlDedicatedGatewayResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CosmosDbSqlDedicatedGatewayResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CosmosDbSqlDedicatedGatewayModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cosmos.SqlDedicatedGatewayClient
			cosmosdbAccountId, err := cosmosdb.ParseDatabaseAccountID(model.CosmosDbAccountId)
			if err != nil {
				return err
			}

			id := sqldedicatedgateway.NewServiceID(cosmosdbAccountId.SubscriptionId, cosmosdbAccountId.ResourceGroupName, cosmosdbAccountId.DatabaseAccountName, string(sqldedicatedgateway.ServiceTypeSqlDedicatedGateway))
			existing, err := client.ServiceGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			serviceType := sqldedicatedgateway.ServiceTypeSqlDedicatedGateway

			parameters := &sqldedicatedgateway.ServiceResourceCreateUpdateParameters{
				Properties: &sqldedicatedgateway.ServiceResourceCreateUpdateProperties{
					ServiceType:   &serviceType,
					InstanceCount: &model.InstanceCount,
					InstanceSize:  &model.InstanceSize,
				},
			}

			if err := client.ServiceCreateThenPoll(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CosmosDbSqlDedicatedGatewayResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.SqlDedicatedGatewayClient

			id, err := sqldedicatedgateway.ParseServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CosmosDbSqlDedicatedGatewayModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.ServiceGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			serviceType := sqldedicatedgateway.ServiceTypeSqlDedicatedGateway

			parameters := &sqldedicatedgateway.ServiceResourceCreateUpdateParameters{
				Properties: &sqldedicatedgateway.ServiceResourceCreateUpdateProperties{
					ServiceType:   &serviceType,
					InstanceCount: &model.InstanceCount,
					InstanceSize:  &model.InstanceSize,
				},
			}

			if err := client.ServiceCreateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CosmosDbSqlDedicatedGatewayResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.SqlDedicatedGatewayClient

			id, err := sqldedicatedgateway.ParseServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ServiceGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := CosmosDbSqlDedicatedGatewayModel{
				CosmosDbAccountId: cosmosdb.NewDatabaseAccountID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName).ID(),
			}

			if props := model.Properties; props != nil {
				existing := props.(sqldedicatedgateway.SqlDedicatedGatewayServiceResourceProperties)

				if existing.InstanceCount != nil {
					state.InstanceCount = *existing.InstanceCount
				}

				if existing.InstanceSize != nil {
					state.InstanceSize = *existing.InstanceSize
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CosmosDbSqlDedicatedGatewayResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.SqlDedicatedGatewayClient

			id, err := sqldedicatedgateway.ParseServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.ServiceDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
