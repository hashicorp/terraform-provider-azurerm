// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var containerAppEnvironment ContainerAppEnvironmentModel

			if err := metadata.Decode(&containerAppEnvironment); err != nil {
				return err
			}

			id := managedenvironments.NewManagedEnvironmentID(subscriptionId, containerAppEnvironment.ResourceGroup, containerAppEnvironment.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			managedEnvironment := managedenvironments.ManagedEnvironment{
				Location: containerAppEnvironment.Location,
				Name:     pointer.To(containerAppEnvironment.Name),
				Properties: &managedenvironments.ManagedEnvironmentProperties{
					AppLogsConfiguration: &managedenvironments.AppLogsConfiguration{
						Destination: pointer.To(containerAppEnvironment.LogsDestination),
					},
					VnetConfiguration: &managedenvironments.VnetConfiguration{},
					ZoneRedundant:     pointer.To(containerAppEnvironment.ZoneRedundant),
					PeerAuthentication: &managedenvironments.ManagedEnvironmentPropertiesPeerAuthentication{
						Mtls: &managedenvironments.Mtls{
							Enabled: pointer.To(containerAppEnvironment.Mtls),
						},
					},
					PeerTrafficConfiguration: &managedenvironments.ManagedEnvironmentPropertiesPeerTrafficConfiguration{
						Encryption: &managedenvironments.ManagedEnvironmentPropertiesPeerTrafficConfigurationEncryption{
							Enabled: pointer.To(containerAppEnvironment.Mtls),
						},
					},
				},
				Tags: tags.Expand(containerAppEnvironment.Tags),
			}

			if containerAppEnvironment.PublicNetworkAccess != "" {
				managedEnvironment.Properties.PublicNetworkAccess = pointer.ToEnum[managedenvironments.PublicNetworkAccess](containerAppEnvironment.PublicNetworkAccess)
			}

			if containerAppEnvironment.DaprApplicationInsightsConnectionString != "" {
				managedEnvironment.Properties.DaprAIConnectionString = pointer.To(containerAppEnvironment.DaprApplicationInsightsConnectionString)
			}

			if containerAppEnvironment.InfrastructureResourceGroup != "" {
				managedEnvironment.Properties.InfrastructureResourceGroup = pointer.To(containerAppEnvironment.InfrastructureResourceGroup)
			}

			if containerAppEnvironment.LogAnalyticsWorkspaceId != "" {
				customerId, sharedKey, err := getSharedKeyForWorkspace(ctx, metadata, containerAppEnvironment.LogAnalyticsWorkspaceId)
				if err != nil {
					return fmt.Errorf("retrieving access keys to Log Analytics Workspace for %s: %+v", id, err)
				}

				managedEnvironment.Properties.AppLogsConfiguration.Destination = pointer.To(LogsDestinationLogAnalytics)

				managedEnvironment.Properties.AppLogsConfiguration.LogAnalyticsConfiguration = &managedenvironments.LogAnalyticsConfiguration{
					CustomerId: customerId,
					SharedKey:  sharedKey,
				}
			}

			if containerAppEnvironment.InfrastructureSubnetId != "" {
				managedEnvironment.Properties.VnetConfiguration.InfrastructureSubnetId = pointer.To(containerAppEnvironment.InfrastructureSubnetId)
				managedEnvironment.Properties.VnetConfiguration.Internal = pointer.To(containerAppEnvironment.InternalLoadBalancerEnabled)
			}

			ident, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(containerAppEnvironment.Identity)
			if err != nil {
				return fmt.Errorf("expanding identity: %+v", err)
			}
			managedEnvironment.Identity = ident

			managedEnvironment.Properties.WorkloadProfiles = ExpandWorkloadProfiles(containerAppEnvironment.WorkloadProfiles)

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, managedEnvironment, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			// Set the `log_analytics_workspace_id` during creation, in case the workspace is created on another subscription.
			if containerAppEnvironment.LogAnalyticsWorkspaceId != "" {
				metadata.ResourceData.Set("log_analytics_workspace_id", containerAppEnvironment.LogAnalyticsWorkspaceId)
			}
			return nil
		},
	}
}
