// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SQLAdministratorV0ToV1 struct{}

func (s SQLAdministratorV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"server_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"login": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"object_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"tenant_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"azuread_authentication_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},
	}
}

func (s SQLAdministratorV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.AzureActiveDirectoryAdministratorIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
