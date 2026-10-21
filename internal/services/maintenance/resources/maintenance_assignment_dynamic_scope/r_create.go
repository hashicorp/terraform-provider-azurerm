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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Maintenance.ConfigurationAssignmentsClient

			var model MaintenanceDynamicScopeModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			maintenanceConfigurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(model.MaintenanceConfigurationId)
			if err != nil {
				return err
			}

			id := configurationassignments.NewConfigurationAssignmentID(metadata.Client.Account.SubscriptionId, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ForSubscriptionsGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			configurationAssignment := configurationassignments.ConfigurationAssignment{
				Properties: &configurationassignments.ConfigurationAssignmentProperties{
					MaintenanceConfigurationId: pointer.To(maintenanceConfigurationId.ID()),
				},
			}

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
				configurationAssignment.Properties.Filter = pointer.To(filterProperties)
			}

			if _, err = client.ForSubscriptionsCreateOrUpdate(ctx, id, configurationAssignment); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
