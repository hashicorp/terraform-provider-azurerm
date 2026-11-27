// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_role_definition

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/rbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLRoleDefinitionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.RbacsClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rbacs.ParseSqlRoleDefinitionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
	defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

	if err := client.SqlResourcesDeleteSqlRoleDefinitionThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
