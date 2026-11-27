// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_dedicated_gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/sqldedicatedgateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ServiceGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := &sqldedicatedgateway.ServiceResourceCreateUpdateParameters{
				Properties: &sqldedicatedgateway.ServiceResourceCreateUpdateProperties{
					ServiceType:   pointer.To(sqldedicatedgateway.ServiceTypeSqlDedicatedGateway),
					InstanceCount: &model.InstanceCount,
					InstanceSize:  &model.InstanceSize,
				},
			}

			if err := client.ServiceCreateCallbackThenPoll(ctx, id, *parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
