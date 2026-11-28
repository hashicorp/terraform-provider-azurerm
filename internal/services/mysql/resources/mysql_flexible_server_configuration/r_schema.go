// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_configuration

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func mysqlFlexibleServerConfigurationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the MySQL Flexible Server Configuration, which needs [to be a valid MySQL configuration name](https://dev.mysql.com/doc/refman/5.7/en/server-configuration.html).",
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

		"value": {
			Description: "The value of the MySQL Flexible Server Configuration. See the MySQL documentation for valid values.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},
	}
}
