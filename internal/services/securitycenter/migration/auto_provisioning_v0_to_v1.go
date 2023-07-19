// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutoProvisioningV0ToV1 struct{}

func (a AutoProvisioningV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"auto_provision": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (a AutoProvisioningV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)

		parsed, err := parse.AutoProvisioningSettingIDInsensitively(oldId)
		if err != nil {
			return nil, fmt.Errorf("parsing old ID %q: %+v", oldId, err)
		}

		// potentially overkill, but the name is fixed (we can't guarantee the casing, however)
		parsed.Name = "default"

		newId := parsed.ID()
		log.Printf("[DEBUG] Updating the ID from %q to %q..", oldId, newId)
		rawState["id"] = newId
		log.Printf("[DEBUG] Updated the ID from %q to %q.", oldId, newId)

		return rawState, nil
	}
}
