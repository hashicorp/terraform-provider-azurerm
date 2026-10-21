// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_devops_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-09-20/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/manageddevopspools/helpers"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedDevOpsPools.PoolsClient

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ManagedDevOpsPoolModel{
				Name:              id.PoolName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if model.Identity != nil {
					flattenedIdentity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
					if err != nil {
						return fmt.Errorf("flattening `identity`: %+v", err)
					}
					state.Identity = *flattenedIdentity
				}

				if props := model.Properties; props != nil {
					devCenterProjectId, err := projects.ParseProjectID(props.DevCenterProjectResourceId)
					if err != nil {
						return fmt.Errorf("parsing `dev_center_project_id`: %+v", err)
					}

					state.DevCenterProjectId = devCenterProjectId.ID()
					state.MaximumConcurrency = props.MaximumConcurrency

					if agentProfile := props.AgentProfile; agentProfile != nil {
						if stateful, ok := agentProfile.(pools.Stateful); ok {
							state.StatefulAgent = flattenStatefulAgentToModel(stateful)
						} else if stateless, ok := agentProfile.(pools.StatelessAgentProfile); ok {
							state.StatelessAgent = flattenStatelessAgentToModel(stateless)
						}
					}

					if organizationProfile := props.OrganizationProfile; organizationProfile != nil {
						if azureDevOpsOrganization, ok := organizationProfile.(pools.AzureDevOpsOrganizationProfile); ok {
							state.AzureDevOpsOrganization = flattenAzureDevOpsOrganizationToModel(azureDevOpsOrganization)
						}
					}

					if fabricProfile := props.FabricProfile; fabricProfile != nil {
						if vmssFabric, ok := fabricProfile.(pools.VMSSFabricProfile); ok {
							state.VirtualMachineScaleSetFabric = flattenVirtualMachineScaleSetFabricToModel(vmssFabric)
						}
					}

					if runtimeConfig := props.RuntimeConfiguration; runtimeConfig != nil {
						state.WorkFolder = pointer.From(runtimeConfig.WorkFolder)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func validateVirtualMachineScaleSetFabricImages(metadata sdk.ResourceMetaData, vmssFabrics []helpers.VirtualMachineScaleSetFabricModel) error {
	rawConfig := metadata.ResourceDiff.GetRawConfig().AsValueMap()

	vmssFabricValue, exists := rawConfig["virtual_machine_scale_set_fabric"]
	if !exists || vmssFabricValue.IsNull() {
		return nil
	}

	if !vmssFabricValue.IsWhollyKnown() {
		return nil
	}

	for _, vmssFabric := range vmssFabrics {
		for i, image := range vmssFabric.Images {
			haveResourceId := image.Id != ""
			haveWellKnownImageName := image.WellKnownImageName != ""

			if !haveResourceId && !haveWellKnownImageName {
				return fmt.Errorf("one of `id` or `well_known_image_name` must be specified for image %d in `virtual_machine_scale_set_fabric`", i)
			}

			if haveResourceId && haveWellKnownImageName {
				return fmt.Errorf("only one of `id` or `well_known_image_name` can be specified for image %d in `virtual_machine_scale_set_fabric`", i)
			}
		}
	}

	return nil
}

func validateManualAgentCounts(profileType string, manualPredictions helpers.ManualResourcePredictionModel, maxConcurrency int64) error {
	hasAllWeek := manualPredictions.AllWeekSchedule > 0

	daySchedules := []struct {
		name     string
		schedule []helpers.DayScheduleModel
	}{
		{"sunday_schedule", manualPredictions.SundaySchedule},
		{"monday_schedule", manualPredictions.MondaySchedule},
		{"tuesday_schedule", manualPredictions.TuesdaySchedule},
		{"wednesday_schedule", manualPredictions.WednesdaySchedule},
		{"thursday_schedule", manualPredictions.ThursdaySchedule},
		{"friday_schedule", manualPredictions.FridaySchedule},
		{"saturday_schedule", manualPredictions.SaturdaySchedule},
	}

	hasDaySchedule := false
	for _, sched := range daySchedules {
		if len(sched.schedule) > 0 {
			hasDaySchedule = true
			break
		}
	}

	if !hasAllWeek && !hasDaySchedule {
		return fmt.Errorf("%s `manual_resource_prediction` must specify either `all_week_schedule` or at least one day schedule", profileType)
	}

	if hasAllWeek {
		if manualPredictions.AllWeekSchedule > maxConcurrency {
			return fmt.Errorf("%s `all_week_schedule` agent count (%d) cannot exceed `maximum_concurrency` (%d)",
				profileType, manualPredictions.AllWeekSchedule, maxConcurrency)
		}
	}

	for _, sched := range daySchedules {
		if err := validateScheduleAgentCounts(profileType, sched.name, sched.schedule, maxConcurrency); err != nil {
			return err
		}
	}

	return nil
}

func validateScheduleAgentCounts(profileType, scheduleName string, schedule []helpers.DayScheduleModel, maxConcurrency int64) error {
	seen := make(map[string]bool, len(schedule))
	for _, entry := range schedule {
		if seen[entry.Time] {
			return fmt.Errorf("%s %s has duplicate time slot %q - entries with the same time must be merged into a single block",
				profileType, scheduleName, entry.Time)
		}
		seen[entry.Time] = true
		if entry.Count > maxConcurrency {
			return fmt.Errorf("%s %s time slot %s has agent count (%d) that exceeds `maximum_concurrency` (%d)",
				profileType, scheduleName, entry.Time, entry.Count, maxConcurrency)
		}
	}
	return nil
}
