// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_dedicated_gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/sqldedicatedgateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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

			parameters := &sqldedicatedgateway.ServiceResourceCreateUpdateParameters{
				Properties: &sqldedicatedgateway.ServiceResourceCreateUpdateProperties{
					ServiceType:   pointer.To(sqldedicatedgateway.ServiceTypeSqlDedicatedGateway),
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
