// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package database_migration_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func databaseMigrationServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the database migration service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ServiceName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"subnet_id": {
			Description:  "The ID of the virtual subnet resource to which the database migration service should be joined.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"sku_name": {
			Description: "The SKU name of the database migration service.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				// No const defined in go sdk, the literal listed below is derived from the response of listskus endpoint.
				// See: https://learn.microsoft.com/rest/api/datamigration/resourceskus/listskus
				"Premium_4vCores",
				"Standard_1vCores",
				"Standard_2vCores",
				"Standard_4vCores",
			}, false),
		},

		"tags": commonschema.Tags(),

		"network_interface_id": {
			Description: "The ID of the Network Interface created by Azure for the Database Migration Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
