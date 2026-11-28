// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_database

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func mysqlFlexibleDatabaseSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the MySQL Database, which needs [to be a valid MySQL identifier](https://dev.mysql.com/doc/refman/8.4/en/identifiers.html).",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"server_name": {
			Description:  "The name of the MySQL Flexible Server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ServerName,
		},

		"charset": {
			Description:      "The Charset for the MySQL Database, which needs [to be a valid MySQL Charset](https://dev.mysql.com/doc/refman/8.4/en/charset-charsets.html).",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: resourceMysqlFlexibleDatabaseAliasDiffSuppress("utf8", "utf8mb3"),
		},

		"collation": {
			Description:      "The Collation for the MySQL Database, which needs [to be a valid MySQL Collation](https://dev.mysql.com/doc/refman/8.4/en/charset-mysql.html).",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: resourceMysqlFlexibleDatabaseAliasDiffSuppress("utf8_unicode_ci", "utf8mb3_unicode_ci"),
		},
	}
}
