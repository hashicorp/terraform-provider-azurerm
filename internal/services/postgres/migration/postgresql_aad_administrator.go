// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/serveradministrators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ pluginsdk.StateUpgrade = PostgresqlAADAdministratorV0ToV1{}

type PostgresqlAADAdministratorV0ToV1 struct{}

func (PostgresqlAADAdministratorV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"server_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"login": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AdminUsernames,
		},

		"object_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}

func (PostgresqlAADAdministratorV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		//  /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}/administrators/activeDirectory
		// new:
		//  /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}
		// summary:
		// Remove administrators chunk
		oldId := rawState["id"].(string)

		// This appeared to be using the wrong AAD Admin ID parser from the deprecated sql package which was removed post 4.0
		// but since this is a state migration this cannot be changed or reapplied retroactively
		id, err := parse.SqlAzureActiveDirectoryAdministratorID(oldId)
		if err != nil {
			return rawState, err
		}

		newId := serveradministrators.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId.ID()

		return rawState, nil
	}
}
