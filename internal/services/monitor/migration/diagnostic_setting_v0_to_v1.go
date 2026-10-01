// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DiagnosticSettingV0ToV1{}

type DiagnosticSettingV0ToV1 struct{}

func (DiagnosticSettingV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"target_resource_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"eventhub_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"eventhub_authorization_rule_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"log_analytics_workspace_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"storage_account_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"partner_solution_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"log_analytics_destination_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"enabled_log": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"category_group": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"enabled_metric": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

func (DiagnosticSettingV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
		oldId, ok := rawState["id"].(string)
		if !ok || oldId == "" {
			return rawState, nil
		}

		if strings.Contains(oldId, "|") {
			parts := strings.Split(oldId, "|")
			if len(parts) != 2 {
				return nil, fmt.Errorf("expected the Monitor Diagnostics ID to be in the format `{resourceId}|{name}` but got %d segments", len(parts))
			}

			newId := diagnosticsettings.NewScopedDiagnosticSettingID(parts[0], parts[1])
			log.Printf("[DEBUG] Upgrading Monitor Diagnostic Setting ID from %q to %q", oldId, newId.ID())
			rawState["id"] = newId.ID()
		}

		return rawState, nil
	}
}
