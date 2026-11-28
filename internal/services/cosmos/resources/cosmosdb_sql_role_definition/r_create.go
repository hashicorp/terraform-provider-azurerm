// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_role_definition

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/rbacs"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLRoleDefinitionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.RbacsClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	roleDefinitionId := d.Get("role_definition_id").(string)
	if roleDefinitionId == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("generating UUID for Cosmos DB SQL Role Definition: %+v", err)
		}

		roleDefinitionId = uuid
	}

	id := rbacs.NewSqlRoleDefinitionID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), roleDefinitionId)

	locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
	defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.SqlResourcesGetSqlRoleDefinition(ctx, id)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			return tf.ImportAsExistsError("azurerm_cosmosdb_sql_role_definition", id.ID())
		}
	}

	parameters := rbacs.SqlRoleDefinitionCreateUpdateParameters{
		Properties: &rbacs.SqlRoleDefinitionResource{
			RoleName:         pointer.To(d.Get("name").(string)),
			AssignableScopes: pluginsdk.ExpandStringSlice(d.Get("assignable_scopes").(*pluginsdk.Set).List()),
			Permissions:      expandSqlRoleDefinitionPermissions(d.Get("permissions").(*pluginsdk.Set).List()),
			Type:             pointer.ToEnum[rbacs.RoleDefinitionType](d.Get("type").(string)),
		},
	}

	if err := client.SqlResourcesCreateUpdateSqlRoleDefinitionCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbSQLRoleDefinitionRead(d, meta)
}

func expandSqlRoleDefinitionPermissions(input []any) *[]rbacs.Permission {
	results := make([]rbacs.Permission, 0)

	for _, item := range input {
		v := item.(map[string]any)

		results = append(results, rbacs.Permission{
			DataActions: pluginsdk.ExpandStringSlice(v["data_actions"].(*pluginsdk.Set).List()),
		})
	}

	return &results
}
