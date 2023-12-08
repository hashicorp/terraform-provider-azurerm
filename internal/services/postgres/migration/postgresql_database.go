// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

var _ pluginsdk.StateUpgrade = PostgresqlDatabaseV0ToV1{}

type PostgresqlDatabaseV0ToV1 struct{}

func (PostgresqlDatabaseV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"server_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ServerName,
		},

		"charset": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ForceNew:         true,
		},

		"collation": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PostgresqlDatabaseCollation,
		},
	}
}

func (PostgresqlDatabaseV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		//  /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBForPostgreSQL/servers/{serverName}/databases/{databaseName}
		// new:
		//  /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}/databases/{databaseName}
		// summary:
		// Check for `For` and swap to `for`
		oldId := rawState["id"].(string)
		if strings.Contains(oldId, "Microsoft.DBForPostgreSQL") {
			modifiedId := strings.ReplaceAll(oldId, "Microsoft.DBForPostgreSQL", "Microsoft.DBforPostgreSQL")

			newId, err := databases.ParseDatabaseID(modifiedId)
			if err != nil {
				return rawState, err
			}
			log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
			rawState["id"] = newId.ID()
		}

		return rawState, nil
	}
}
