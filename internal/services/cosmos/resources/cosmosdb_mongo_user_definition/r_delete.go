// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_user_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/mongorbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.MongoRBACClient

			id, err := mongorbacs.ParseMongodbUserDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
			defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

			if err := client.MongoDBResourcesDeleteMongoUserDefinitionThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
