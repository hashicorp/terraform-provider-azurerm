// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_database

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

func postgresqlFlexibleServerDatabaseSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the PostgreSQL Database, which needs [to be a valid PostgreSQL identifier](https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIERS).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FlexibleServerDatabaseName,
		},

		"server_id": {
			Description:  "The ID of the Azure PostgreSQL Flexible Server from which to create this PostgreSQL Flexible Server Database.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: databases.ValidateFlexibleServerID,
		},

		"charset": {
			Description:      "The Charset for the Azure PostgreSQL Flexible Server Database, which needs [to be a valid PostgreSQL Charset](https://www.postgresql.org/docs/current/static/multibyte.html).",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validate.DatabaseCharset,
			Default:          "UTF8",
		},

		"collation": {
			Description:  "The Collation for the Azure PostgreSQL Flexible Server Database, which needs [to be a valid PostgreSQL Collation](https://www.postgresql.org/docs/current/static/collation.html).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.PostgresqlFlexibleServerDatabaseCollation,
			Default:      "en_US.utf8",
		},
	}
}
