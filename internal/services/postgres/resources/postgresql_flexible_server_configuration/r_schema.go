// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func postgresqlFlexibleServerConfigurationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the PostgreSQL Configuration, which needs [to be a valid PostgreSQL configuration name](https://www.postgresql.org/docs/current/static/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIER).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"server_id": {
			Description:  "The ID of the PostgreSQL Flexible Server where we want to change configuration.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: servers.ValidateFlexibleServerID,
		},

		"value": {
			Description:  "The value of the PostgreSQL Configuration. See the PostgreSQL documentation for valid values.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
