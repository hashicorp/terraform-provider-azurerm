// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_environment_v3

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appserviceenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.AppServiceEnvironmentClient
			id, err := commonids.ParseAppServiceEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := AppServiceEnvironmentV3Model{
				Name:          id.HostingEnvironmentName,
				ResourceGroup: id.ResourceGroupName,
			}

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if props := model.Properties; props != nil {
					state.SubnetId = props.VirtualNetwork.Id
					state.InternalLoadBalancingMode = pointer.FromEnum(props.InternalLoadBalancingMode)
					state.DedicatedHostCount = pointer.From(props.DedicatedHostCount)
					state.PricingTier = pointer.From(props.MultiSize)
					state.ClusterSetting = flattenClusterSettingsModel(props.ClusterSettings)
					state.DnsSuffix = pointer.From(props.DnsSuffix)
					state.IpSSLAddressCount = pointer.From(props.IPsslAddressCount)
					state.ZoneRedundant = pointer.From(props.ZoneRedundant)
				}

				existingNetwork, err := client.GetAseV3NetworkingConfiguration(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading network configuration for %s: %+v", *id, err)
				}

				if networkModel := existingNetwork.Model; networkModel != nil {
					if props := networkModel.Properties; props != nil {
						state.WindowsOutboundIPAddresses = pointer.From(props.WindowsOutboundIPAddresses)
						state.LinuxOutboundIPAddresses = pointer.From(props.LinuxOutboundIPAddresses)
						state.InternalInboundIPAddresses = pointer.From(props.InternalInboundIPAddresses)
						state.ExternalInboundIPAddresses = pointer.From(props.ExternalInboundIPAddresses)
						state.AllowNewPrivateEndpointConnections = pointer.From(props.AllowNewPrivateEndpointConnections)
						state.RemoteDebuggingEnabled = pointer.From(props.RemoteDebugEnabled)
					}
				}
				inboundNetworkDependencies, err := flattenInboundNetworkDependencies(ctx, client, id)
				if err != nil {
					return err
				}

				state.InboundNetworkDependencies = *inboundNetworkDependencies

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenClusterSettingsModel(input *[]appserviceenvironments.NameValuePair) []ClusterSettingModel {
	if input == nil || len(*input) == 0 {
		return []ClusterSettingModel{}
	}

	output := make([]ClusterSettingModel, 0, len(*input))
	for _, v := range *input {
		if v.Name == nil {
			continue
		}

		output = append(output, ClusterSettingModel{
			Name:  *v.Name,
			Value: pointer.From(v.Value),
		})
	}
	return output
}

func flattenInboundNetworkDependencies(ctx context.Context, client *appserviceenvironments.AppServiceEnvironmentsClient, id *commonids.AppServiceEnvironmentId) (*[]AppServiceV3InboundDependencies, error) {
	inboundNetworking, err := client.GetInboundNetworkDependenciesEndpointsComplete(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading paged results for Inbound Network Dependencies for %s: %+v", id, err)
	}

	results := make([]AppServiceV3InboundDependencies, 0, len(inboundNetworking.Items))
	for _, v := range inboundNetworking.Items {
		result := AppServiceV3InboundDependencies{
			Description: pointer.From(v.Description),
		}

		if v.Endpoints != nil {
			result.IPAddresses = *v.Endpoints
		}

		if v.Ports != nil {
			result.Ports = *v.Ports
		}

		results = append(results, result)
	}

	return &results, nil
}

const KindASEV3 = "ASEV3"

type AppServiceV3InboundDependencies struct {
	Description string   `tfschema:"description"`
	IPAddresses []string `tfschema:"ip_addresses"`
	Ports       []string `tfschema:"ports"`
}

// resolvePreflightVnetLocation determines the Virtual Network location used for preflight
// validation. It returns skip=true when validation should be gracefully skipped (e.g. the
// subnet_id is not yet known, or the Virtual Network doesn't exist and no location fallback
// is configured).
func resolvePreflightVnetLocation(ctx context.Context, metadata sdk.ResourceMetaData, model AppServiceEnvironmentV3Model) (vnetLoc string, skip bool) {
	if model.SubnetId == "" {
		return "", false
	}

	subnet, err := commonids.ParseSubnetID(model.SubnetId)
	if err != nil {
		metadata.Logger.Info(fmt.Sprintf("skipping preflight validation for %q: 'subnet_id' unknown", model.Name))
		return "", true
	}

	vnetId := commonids.NewVirtualNetworkID(subnet.SubscriptionId, subnet.ResourceGroupName, subnet.VirtualNetworkName)

	vnet, err := metadata.Client.Network.VirtualNetworks.Get(ctx, vnetId, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(vnet.HttpResponse) {
			return "", true
		}

		// The VNet doesn't exist yet, so we can't determine its location. Rely on the configured
		// fallback if one is available, otherwise gracefully skip preflight validation.
		fallback := metadata.Client.Features.EnhancedValidation.LocationFallback
		if fallback == nil {
			metadata.Logger.Info(fmt.Sprintf("skipping preflight validation for %q: Virtual Network %q not found and no location fallback configured", model.Name, vnetId.ID()))
			return "", true
		}

		metadata.Logger.Info(fmt.Sprintf("Virtual Network %q not found, relying on location fallback for preflight validation", vnetId.ID()))
		return *fallback, false
	}

	if vnet.Model == nil || vnet.Model.Location == nil {
		return "", true
	}

	return *vnet.Model.Location, false
}

func checkNetworkConfigUpdate(ctx context.Context, client *appserviceenvironments.AppServiceEnvironmentsClient, id commonids.AppServiceEnvironmentId, values appserviceenvironments.AseV3NetworkingConfigurationProperties) pluginsdk.StateRefreshFunc {
	return func() (result any, state string, err error) {
		resp, err := client.GetAseV3NetworkingConfiguration(ctx, id)
		if err != nil || resp.Model == nil || resp.Model.Properties == nil {
			return nil, "", err
		}

		props := *resp.Model.Properties

		debugEnabledReq := pointer.From(values.RemoteDebugEnabled)
		newPECReq := pointer.From(values.AllowNewPrivateEndpointConnections)

		debugEnableResp := pointer.From(props.RemoteDebugEnabled)
		newPECResp := pointer.From(props.AllowNewPrivateEndpointConnections)

		if debugEnableResp != debugEnabledReq || newPECResp != newPECReq {
			return props, "Pending", nil
		}

		return props, "Succeeded", nil
	}
}
