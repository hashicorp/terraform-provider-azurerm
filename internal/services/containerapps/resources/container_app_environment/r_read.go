// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient
			id, err := managedenvironments.ParseManagedEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var existingState ContainerAppEnvironmentModel
			if err := metadata.Decode(&existingState); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			var state ContainerAppEnvironmentModel

			if model := existing.Model; model != nil {
				state.Name = id.ManagedEnvironmentName
				state.ResourceGroup = id.ResourceGroupName
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)
				if model.Identity != nil {
					ident, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
					if err != nil {
						return fmt.Errorf("flattening identity: %+v", err)
					}
					state.Identity = ident
				}

				if props := model.Properties; props != nil {
					if vnet := props.VnetConfiguration; vnet != nil {
						state.InfrastructureSubnetId = pointer.From(vnet.InfrastructureSubnetId)
						state.InternalLoadBalancerEnabled = pointer.From(vnet.Internal)
						state.DockerBridgeCidr = pointer.From(vnet.DockerBridgeCidr)
						state.PlatformReservedCidr = pointer.From(vnet.PlatformReservedCidr)
						state.PlatformReservedDnsIP = pointer.From(vnet.PlatformReservedDnsIP)
					}

					if appLogsConfig := props.AppLogsConfiguration; appLogsConfig != nil {
						state.LogsDestination = pointer.From(appLogsConfig.Destination)
						if appLogsConfig.LogAnalyticsConfiguration != nil && appLogsConfig.LogAnalyticsConfiguration.CustomerId != nil {
							workspaceId, err := findWorkspaceResourceIDFromCustomerID(ctx, metadata, *appLogsConfig.LogAnalyticsConfiguration.CustomerId)
							// During refreshing stage, `GetRawConfig()` may return null value.

							if err == nil {
								if workspaceId != nil {
									state.LogAnalyticsWorkspaceId = workspaceId.ID()
								} else {
									state.LogAnalyticsWorkspaceId = existingState.LogAnalyticsWorkspaceId
								}
							}
						}
					}

					state.PublicNetworkAccess = pointer.FromEnum(props.PublicNetworkAccess)
					state.ZoneRedundant = pointer.From(props.ZoneRedundant)
					state.StaticIP = pointer.From(props.StaticIP)
					state.DefaultDomain = pointer.From(props.DefaultDomain)
					state.WorkloadProfiles = FlattenWorkloadProfiles(props.WorkloadProfiles)
					state.InfrastructureResourceGroup = pointer.From(props.InfrastructureResourceGroup)

					if props.CustomDomainConfiguration != nil {
						state.CustomDomainVerificationId = pointer.From(props.CustomDomainConfiguration.CustomDomainVerificationId)
					}

					if props.PeerAuthentication != nil && props.PeerAuthentication.Mtls != nil {
						state.Mtls = pointer.From(props.PeerAuthentication.Mtls.Enabled)
					}
				}
			}

			// `dapr_application_insights_connection_string` is sensitive and not returned by API
			if v := metadata.ResourceData.Get("dapr_application_insights_connection_string").(string); v != "" {
				state.DaprApplicationInsightsConnectionString = v
			}

			if err = metadata.Encode(&state); err != nil {
				return fmt.Errorf("encoding: %+v", err)
			}

			return nil
		},
	}
}

const (
	LogsDestinationLogAnalytics string = "log-analytics"
	LogsDestinationAzureMonitor string = "azure-monitor"
	LogsDestinationNone         string = ""
)

func findWorkspaceResourceIDFromCustomerID(ctx context.Context, meta sdk.ResourceMetaData, customerID string) (*workspaces.WorkspaceId, error) {
	client := meta.Client.LogAnalytics.WorkspaceClient

	subscriptionId := commonids.NewSubscriptionID(meta.Client.Account.SubscriptionId)

	list, err := client.List(ctx, subscriptionId)
	if err != nil {
		return nil, err
	}

	model := list.Model
	if model == nil {
		return nil, fmt.Errorf("could not resolve Log Analytics Workspace ID for %s, list model was nil", customerID)
	}

	if len(*model) == 0 {
		return nil, fmt.Errorf("could not resolve Log Analytics Workspace ID for %s, no Log Analytics Workspaces found in %s", customerID, subscriptionId)
	}

	for _, v := range *model {
		if v.Properties != nil && v.Properties.CustomerId != nil && strings.EqualFold(*v.Properties.CustomerId, customerID) {
			result, err := workspaces.ParseWorkspaceIDInsensitively(pointer.From(v.Id))
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}

	return nil, nil
}

func getSharedKeyForWorkspace(ctx context.Context, meta sdk.ResourceMetaData, workspaceID string) (*string, *string, error) {
	logAnalyticsClient := meta.Client.LogAnalytics.SharedKeyWorkspacesClient

	logAnalyticsId, err := workspaces.ParseWorkspaceID(workspaceID)
	if err != nil {
		return nil, nil, err
	}

	workspace, err := logAnalyticsClient.Get(ctx, *logAnalyticsId)
	if err != nil {
		return nil, nil, fmt.Errorf("retrieving %s: %+v", logAnalyticsId, err)
	}

	if workspace.Model == nil || workspace.Model.Properties == nil {
		return nil, nil, fmt.Errorf("reading customer ID from %s", logAnalyticsId)
	}

	if workspace.Model.Properties.CustomerId == nil {
		return nil, nil, fmt.Errorf("reading customer ID from %s, `customer_id` is nil", logAnalyticsId)
	}

	keys, err := logAnalyticsClient.SharedKeysGetSharedKeys(ctx, *logAnalyticsId)
	if err != nil {
		return nil, nil, fmt.Errorf("retrieving access keys to %s: %+v", logAnalyticsId, err)
	}
	if keys.Model.PrimarySharedKey == nil {
		return nil, nil, fmt.Errorf("reading shared key for %s", logAnalyticsId)
	}

	return workspace.Model.Properties.CustomerId, keys.Model.PrimarySharedKey, nil
}
