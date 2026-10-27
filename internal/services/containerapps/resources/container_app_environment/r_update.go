// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/azuresdkhacks"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient
			workaroundClient := azuresdkhacks.NewManagedEnvironmentWorkaroundClient(client)
			id, err := managedenvironments.ParseManagedEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ContainerAppEnvironmentModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			payload := azuresdkhacks.ManagedEnvironment{
				Name:       pointer.To(state.Name),
				Location:   state.Location,
				Properties: &azuresdkhacks.ManagedEnvironmentProperties{},
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(state.Tags)
			}

			if metadata.ResourceData.HasChange("identity") {
				ident, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = ident
			}

			if metadata.ResourceData.HasChange("workload_profile") {
				payload.Properties.WorkloadProfiles = ExpandWorkloadProfiles(state.WorkloadProfiles)
			}

			if metadata.ResourceData.HasChange("mutual_tls_enabled") {
				payload.Properties.PeerAuthentication = &managedenvironments.ManagedEnvironmentPropertiesPeerAuthentication{
					Mtls: &managedenvironments.Mtls{
						Enabled: pointer.To(state.Mtls),
					},
				}
				payload.Properties.PeerTrafficConfiguration = &managedenvironments.ManagedEnvironmentPropertiesPeerTrafficConfiguration{
					Encryption: &managedenvironments.ManagedEnvironmentPropertiesPeerTrafficConfigurationEncryption{
						Enabled: pointer.To(state.Mtls),
					},
				}
			}

			if metadata.ResourceData.HasChanges("logs_destination", "log_analytics_workspace_id") {
				switch state.LogsDestination {
				case LogsDestinationAzureMonitor:
					payload.Properties.AppLogsConfiguration = &azuresdkhacks.AppLogsConfiguration{
						Destination: pointer.To(LogsDestinationAzureMonitor),
					}
				case LogsDestinationLogAnalytics:
					if state.LogAnalyticsWorkspaceId != "" {
						customerId, sharedKey, err := getSharedKeyForWorkspace(ctx, metadata, state.LogAnalyticsWorkspaceId)
						if err != nil {
							return fmt.Errorf("retrieving access keys to Log Analytics Workspace for %s: %+v", id, err)
						}

						payload.Properties.AppLogsConfiguration = &azuresdkhacks.AppLogsConfiguration{
							Destination: pointer.To(LogsDestinationLogAnalytics),
							LogAnalyticsConfiguration: &managedenvironments.LogAnalyticsConfiguration{
								CustomerId: customerId,
								SharedKey:  sharedKey,
							},
						}
						// Set the `log_analytics_workspace_id` during creation, in case the workspace is created on another subscription.
						if state.LogAnalyticsWorkspaceId != "" {
							metadata.ResourceData.Set("log_analytics_workspace_id", state.LogAnalyticsWorkspaceId)
						}
					}
				default:
					payload.Properties.AppLogsConfiguration = &azuresdkhacks.AppLogsConfiguration{
						Destination: pointer.To(LogsDestinationNone),
					}
				}
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				payload.Properties.PublicNetworkAccess = pointer.ToEnum[managedenvironments.PublicNetworkAccess](state.PublicNetworkAccess)
			}

			if err := workaroundClient.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
