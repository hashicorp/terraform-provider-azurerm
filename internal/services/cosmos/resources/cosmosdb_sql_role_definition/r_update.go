// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_role_definition

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/rbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLRoleDefinitionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.RbacsClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rbacs.ParseSqlRoleDefinitionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
	defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

	existing, err := client.SqlResourcesGetSqlRoleDefinition(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %w", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", id)
	}

	parameters := rbacs.SqlRoleDefinitionCreateUpdateParameters{
		Properties: existing.Model.Properties,
	}

	if d.HasChange("assignable_scopes") {
		parameters.Properties.AssignableScopes = pluginsdk.ExpandStringSlice(d.Get("assignable_scopes").(*pluginsdk.Set).List())
	}

	if d.HasChange("name") {
		parameters.Properties.RoleName = pointer.To(d.Get("name").(string))
	}

	if d.HasChange("permissions") {
		parameters.Properties.Permissions = expandSqlRoleDefinitionPermissions(d.Get("permissions").(*pluginsdk.Set).List())
	}

	if err := client.SqlResourcesCreateUpdateSqlRoleDefinitionThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceCosmosDbSQLRoleDefinitionRead(d, meta)
}
