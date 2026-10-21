// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_assignment_dynamic_scope

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Maintenance.ConfigurationAssignmentsClient

			var state MaintenanceDynamicScopeModel
			id, err := configurationassignments.ParseConfigurationAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ForSubscriptionsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return err
			}

			if model := resp.Model; model != nil {
				state.Name = id.ConfigurationAssignmentName

				if properties := model.Properties; properties != nil {
					if properties.MaintenanceConfigurationId != nil {
						maintenanceConfigurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationIDInsensitively(pointer.From(properties.MaintenanceConfigurationId))
						if err != nil {
							return fmt.Errorf("parsing %q: %+v", pointer.From(properties.MaintenanceConfigurationId), err)
						}
						state.MaintenanceConfigurationId = maintenanceConfigurationId.ID()
					}

					if filter := properties.Filter; filter != nil {
						filterProp := make([]Filter, 0)
						tagsListProp := make([]Tag, 0)
						tagFilterProp := ""
						if tags := filter.TagSettings; tags != nil {
							tagFilterProp = pointer.FromEnum(tags.FilterOperator)
							for k, v := range pointer.From(tags.Tags) {
								tagsListProp = append(tagsListProp, Tag{
									Tag:    k,
									Values: v,
								})
							}
						}
						filterProp = append(filterProp, Filter{
							Locations:      pointer.From(filter.Locations),
							OsTypes:        pointer.From(filter.OsTypes),
							ResourceGroups: pointer.From(filter.ResourceGroups),
							ResourceTypes:  pointer.From(filter.ResourceTypes),
							Tags:           tagsListProp,
							TagFilter:      tagFilterProp,
						})
						state.Filter = filterProp
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}

type Tag struct {
	Tag    string   `tfschema:"tag"`
	Values []string `tfschema:"values"`
}

type Filter struct {
	Locations      []string `tfschema:"locations"`
	OsTypes        []string `tfschema:"os_types"`
	ResourceGroups []string `tfschema:"resource_groups"`
	ResourceTypes  []string `tfschema:"resource_types"`
	Tags           []Tag    `tfschema:"tags"`
	TagFilter      string   `tfschema:"tag_filter"`
}
