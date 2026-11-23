// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_role_definition

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/rbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLRoleDefinitionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.RbacsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rbacs.ParseSqlRoleDefinitionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SqlResourcesGetSqlRoleDefinition(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("role_definition_id", id.RoleDefinitionId)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.DatabaseAccountName)

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil {
			d.Set("assignable_scopes", pluginsdk.FlattenSlice(props.AssignableScopes))
			d.Set("name", props.RoleName)
			d.Set("type", pointer.FromEnum(props.Type))

			if err := d.Set("permissions", flattenSqlRoleDefinitionPermissions(props.Permissions)); err != nil {
				return fmt.Errorf("setting `permissions`: %+v", err)
			}
		}
	}

	return nil
}

func flattenSqlRoleDefinitionPermissions(input *[]rbacs.Permission) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"data_actions": pluginsdk.FlattenSlice(item.DataActions),
		})
	}

	return results
}
