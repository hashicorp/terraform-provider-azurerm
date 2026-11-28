// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_role_definition

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cosmos_mongo_database_id": {
			Description:  "The resource ID of the Mongo DB.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cosmosdb.ValidateMongodbDatabaseID,
		},

		"role_name": {
			Description:  "The user-friendly name for the Mongo Role Definition.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"inherited_role_names": {
			Description: "A list of Mongo Roles which are inherited to the Mongo Role Definition.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"privilege": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"actions": {
						Description: "A list of actions that are allowed.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"resource": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"collection_name": {
									Description:  "The name of the Mongo DB Collection that the Role Definition is applied to.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validate.CosmosEntityName,
								},

								"db_name": {
									Description:  "The name of the Mongo DB that the Role Definition is applied to.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validate.CosmosEntityName,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
