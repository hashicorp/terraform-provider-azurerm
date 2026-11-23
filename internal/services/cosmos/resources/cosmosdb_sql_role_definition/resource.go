// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_role_definition

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/rbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var CosmosDbAccountResourceName = "azurerm_cosmosdb_account"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbSQLRoleDefinitionCreate,
		Read:   resourceCosmosDbSQLRoleDefinitionRead,
		Update: resourceCosmosDbSQLRoleDefinitionUpdate,
		Delete: resourceCosmosDbSQLRoleDefinitionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := rbacs.ParseSqlRoleDefinitionID(id)
			return err
		}),

		Schema: cosmosdbSqlRoleDefinitionSchema(),
	}
}
