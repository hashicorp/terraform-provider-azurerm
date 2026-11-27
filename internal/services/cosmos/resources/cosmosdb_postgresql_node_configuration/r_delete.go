// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_node_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.ConfigurationsClient

			id, err := configurations.ParseNodeConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.ServerGroupsv2Name, CosmosDbPostgreSQLClusterResourceName)
			defer locks.UnlockByName(id.ServerGroupsv2Name, CosmosDbPostgreSQLClusterResourceName)

			resp, err := client.GetNode(ctx, *id)
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

			if err = client.UpdateOnNodeThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
