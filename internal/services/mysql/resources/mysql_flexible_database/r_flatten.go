// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_database

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceMySqlFlexibleDatabaseFlatten(d *pluginsdk.ResourceData, id *databases.DatabaseId, database *databases.Database) error {
	d.Set("name", id.DatabaseName)
	d.Set("server_name", id.FlexibleServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if database != nil {
		if props := database.Properties; props != nil {
			d.Set("charset", props.Charset)
			d.Set("collation", props.Collation)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
