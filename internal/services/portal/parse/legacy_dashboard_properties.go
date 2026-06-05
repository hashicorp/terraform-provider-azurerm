// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"encoding/json"
	"maps"
	"slices"

	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2026-04-01/dashboards"
)

// LegacyDashboardProperties converts legacy dashboard properties payloads to the
// 2020-09-01-preview dashboard properties model. Should be removed after 5.0.
func LegacyDashboardProperties(v string) (*dashboards.DashboardPropertiesWithProvisioningState, bool) {

	var dashboardProperties dashboard.DashboardProperties
	var dashboardPropertiesWithProvisioningState dashboards.DashboardPropertiesWithProvisioningState

	if err := json.Unmarshal([]byte(v), &dashboardProperties); err == nil {
		if dashboardProperties.Lenses == nil {
			return nil, false
		}

		lenses := make([]dashboards.DashboardLens, 0, len(*dashboardProperties.Lenses))
		lensKeys := slices.Sorted(maps.Keys(*dashboardProperties.Lenses))
		for _, lensKey := range lensKeys {
			lens := (*dashboardProperties.Lenses)[lensKey]
			newLens := dashboards.DashboardLens{
				Order: lens.Order,
			}
			if lens.Metadata != nil {
				m := interface{}(*lens.Metadata)
				newLens.Metadata = &m
			}
			partKeys := slices.Sorted(maps.Keys(lens.Parts))
			parts := make([]dashboards.DashboardParts, 0, len(lens.Parts))
			for _, partKey := range partKeys {
				part := lens.Parts[partKey]
				p := dashboards.DashboardParts{
					Position: dashboards.DashboardPartsPosition{
						X:       part.Position.X,
						Y:       part.Position.Y,
						RowSpan: part.Position.RowSpan,
						ColSpan: part.Position.ColSpan,
					},
				}
				if part.Metadata != nil {
					if metaData, ok := (*part.Metadata).(map[string]interface{}); ok {
						if settings, ok := metaData["settings"].(map[string]interface{}); ok {
							if content, ok := settings["content"].(map[string]interface{}); ok {
								if inner, ok := content["settings"]; ok {
									settings["content"] = inner
									metaData["settings"] = settings
								}
							}
						}

						p.Metadata = rawDashboardPartMetadata(metaData)
					}
				}
				parts = append(parts, p)
			}
			newLens.Parts = parts
			lenses = append(lenses, newLens)
		}

		dashboardPropertiesWithProvisioningState.Lenses = &lenses

		if dashboardProperties.Metadata != nil {
			m := interface{}(*dashboardProperties.Metadata)
			dashboardPropertiesWithProvisioningState.Metadata = &m
		}

		return &dashboardPropertiesWithProvisioningState, true

	}

	return nil, false
}
