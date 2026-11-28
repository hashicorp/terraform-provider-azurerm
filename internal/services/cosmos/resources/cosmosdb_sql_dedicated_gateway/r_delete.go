// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_dedicated_gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/sqldedicatedgateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
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
