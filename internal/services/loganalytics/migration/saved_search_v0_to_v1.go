// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/savedsearches"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

var _ pluginsdk.StateUpgrade = SavedSearchV0ToV1{}

type SavedSearchV0ToV1 struct{}

func (SavedSearchV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"log_analytics_workspace_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     savedsearches.ValidateWorkspaceID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"name": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"category": {
			Type:     pluginsdk.TypeString,
			ForceNew: true,
			Required: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"query": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"function_alias": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"function_parameters": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (SavedSearchV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)

		id, err := savedsearches.ParseSavedSearchID(fmt.Sprintf("/%s", strings.TrimPrefix(oldId, "/")))
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q..", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
