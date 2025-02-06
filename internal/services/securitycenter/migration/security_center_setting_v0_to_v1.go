// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = SecurityCenterSettingsV0ToV1{}

type SecurityCenterSettingsV0ToV1 struct{}

func (SecurityCenterSettingsV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"enabled": {
			Required: true,
			Type:     pluginsdk.TypeBool,
		},
		"setting_name": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (SecurityCenterSettingsV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Println("[DEBUG] Migrating Security Center Settings from v0 to v1 format")
		oldId := strings.Split(rawState["id"].(string), "/")
		if oldId[len(oldId)-1] == "SENTINEL" {
			oldId[len(oldId)-1] = "Sentinel"
		}
		newId := strings.Join(oldId, "/")

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId

		return rawState, nil
	}
}
