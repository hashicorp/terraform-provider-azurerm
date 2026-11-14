// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_database

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

func postgresqlFlexibleServerDatabaseSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FlexibleServerDatabaseName,
		},

		"server_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: databases.ValidateFlexibleServerID,
		},

		"charset": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validate.DatabaseCharset,
			Default:          "UTF8",
		},

		"collation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.PostgresqlFlexibleServerDatabaseCollation,
			Default:      "en_US.utf8",
		},
	}
}
