// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Default:  false for `auto_rotation_enabled`

var _ pluginsdk.StateUpgrade = MsSqlTransparentDataEncryptionV0ToV1{}

type MsSqlTransparentDataEncryptionV0ToV1 struct{}

func (d MsSqlTransparentDataEncryptionV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"server_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"key_vault_key_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (d MsSqlTransparentDataEncryptionV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[DEBUG] Upgrading from Transparent Data Encryption V0 to V1..")
		existing := rawState["auto_rotation_enabled"]
		if existing == nil {
			log.Printf("[DEBUG] Setting `auto_rotation_enabled` to `false`")
			rawState["auto_rotation_enabled"] = false
		}

		log.Printf("[DEBUG] Upgraded from Transparent Data Encryption V0 to V1..")
		return rawState, nil
	}
}
