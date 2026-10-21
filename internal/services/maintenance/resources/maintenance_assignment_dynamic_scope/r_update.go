// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_assignment_dynamic_scope

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Maintenance.ConfigurationAssignmentsClient

			id, err := configurationassignments.ParseConfigurationAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model MaintenanceDynamicScopeModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.ForSubscriptionsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			existing := resp.Model

			if metadata.ResourceData.HasChange("filter") {
				if len(model.Filter) > 0 {
					filter := model.Filter[0]
					filterProperties := configurationassignments.ConfigurationAssignmentFilterProperties{}

					if len(filter.Locations) > 0 {
						filterProperties.Locations = pointer.To(filter.Locations)
					}

					if len(filter.OsTypes) > 0 {
						filterProperties.OsTypes = pointer.To(filter.OsTypes)
					}

					if len(filter.ResourceGroups) > 0 {
						filterProperties.ResourceGroups = pointer.To(filter.ResourceGroups)
					}

					if len(filter.ResourceTypes) > 0 {
						filterProperties.ResourceTypes = pointer.To(filter.ResourceTypes)
					}

					if len(filter.Tags) > 0 || filter.TagFilter != "" {
						tags := make(map[string][]string)
						for _, tag := range filter.Tags {
							tags[tag.Tag] = tag.Values
						}

						filterProperties.TagSettings = &configurationassignments.TagSettingsProperties{
							FilterOperator: pointer.ToEnum[configurationassignments.TagOperators](filter.TagFilter),
							Tags:           pointer.To(tags),
						}
					}

					if pointer.To(filterProperties) != nil {
						existing.Properties.Filter = pointer.To(filterProperties)
					}
				} else {
					existing.Properties.Filter = &configurationassignments.ConfigurationAssignmentFilterProperties{}
				}
			}

			if _, err = client.ForSubscriptionsCreateOrUpdate(ctx, pointer.From(id), *existing); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}
			return nil
		},
	}
}
