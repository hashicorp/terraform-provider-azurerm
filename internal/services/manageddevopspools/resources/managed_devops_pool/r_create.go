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
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-09-20/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedDevOpsPools.PoolsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagedDevOpsPoolModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := pools.NewPoolID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			expandedIdentity, err := identity.ExpandUserAssignedMapFromModel(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			var agentProfile pools.AgentProfile
			if _, ok := metadata.ResourceData.GetOk("stateful_agent"); ok {
				agentProfile = expandStatefulAgentModel(config.StatefulAgent)
			} else {
				agentProfile = expandStatelessAgentModel(config.StatelessAgent)
			}

			payload := pools.Pool{
				Name:     pointer.To(config.Name),
				Location: config.Location,
				Identity: expandedIdentity,
				Properties: &pools.PoolProperties{
					DevCenterProjectResourceId: config.DevCenterProjectId,
					MaximumConcurrency:         config.MaximumConcurrency,
					AgentProfile:               agentProfile,
					OrganizationProfile:        expandAzureDevOpsOrganizationModel(config.AzureDevOpsOrganization),
					FabricProfile:              expandVirtualMachineScaleSetFabricModel(config.VirtualMachineScaleSetFabric),
					RuntimeConfiguration: &pools.RuntimeConfiguration{
						WorkFolder: pointer.To(config.WorkFolder),
					},
				},
				Tags: pointer.To(config.Tags),
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
