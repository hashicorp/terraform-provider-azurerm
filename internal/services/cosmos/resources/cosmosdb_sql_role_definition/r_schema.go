// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_role_definition

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/rbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cosmosdbSqlRoleDefinitionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_definition_id": {
			Description:  "The GUID as the name of the Cosmos DB SQL Role Definition - one will be generated if not specified.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Description:  "The name of the Cosmos DB Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosAccountName,
		},

		"type": {
			Description:  "The type of the Cosmos DB SQL Role Definition.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(rbacs.RoleDefinitionTypeCustomRole),
			ValidateFunc: validation.StringInSlice(rbacs.PossibleValuesForRoleDefinitionType(), false),
		},

		"assignable_scopes": {
			Description: "A list of fully qualified scopes at or below which Role Assignments may be created using this Cosmos DB SQL Role Definition. It will allow application of this Cosmos DB SQL Role Definition on the entire Database Account or any underlying Database/Collection. Scopes higher than Database Account are not enforceable as assignable scopes.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"name": {
			Description:  "A user-friendly name for the Cosmos DB SQL Role Definition which must be unique for the Database Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"permissions": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"data_actions": {
						Description: "A list of data actions that are allowed for the Cosmos DB SQL Role Definition.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}
