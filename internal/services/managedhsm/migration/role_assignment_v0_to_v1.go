// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ManagedHSMRoleAssignmentV0ToV1{}

type ManagedHSMRoleAssignmentV0ToV1 struct {
}

func (m ManagedHSMRoleAssignmentV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"vault_base_url": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"scope": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"role_definition_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"principal_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m ManagedHSMRoleAssignmentV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := parseLegacyV0RoleAssignmentId(oldIdRaw)
		if err != nil {
			return rawState, fmt.Errorf("parsing the old Role Assignment ID %q: %+v", oldId, err)
		}

		newId := parse.NewManagedHSMDataPlaneRoleAssignmentID(oldId.managedHSMName, oldId.domainSuffix, oldId.scope, oldId.roleAssignmentName).ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldIdRaw, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
