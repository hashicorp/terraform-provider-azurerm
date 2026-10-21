// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_devops_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-09-20/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedDevOpsPools.PoolsClient

			var config ManagedDevOpsPoolModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("dev_center_project_id") {
				payload.Properties.DevCenterProjectResourceId = config.DevCenterProjectId
			}

			if metadata.ResourceData.HasChange("maximum_concurrency") {
				payload.Properties.MaximumConcurrency = config.MaximumConcurrency
			}

			if metadata.ResourceData.HasChanges("stateful_agent", "stateless_agent") {
				var agentProfile pools.AgentProfile

				if _, ok := metadata.ResourceData.GetOk("stateful_agent"); ok {
					agentProfile = expandStatefulAgentModel(config.StatefulAgent)
				} else {
					agentProfile = expandStatelessAgentModel(config.StatelessAgent)
				}

				payload.Properties.AgentProfile = agentProfile
			}

			if metadata.ResourceData.HasChange("azure_devops_organization") {
				payload.Properties.OrganizationProfile = expandAzureDevOpsOrganizationModel(config.AzureDevOpsOrganization)
			}

			if metadata.ResourceData.HasChange("virtual_machine_scale_set_fabric") {
				payload.Properties.FabricProfile = expandVirtualMachineScaleSetFabricModel(config.VirtualMachineScaleSetFabric)
			}

			if metadata.ResourceData.HasChange("work_folder") {
				payload.Properties.RuntimeConfiguration = &pools.RuntimeConfiguration{
					WorkFolder: pointer.To(config.WorkFolder),
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
