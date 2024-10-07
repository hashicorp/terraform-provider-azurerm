// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appserviceenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const KindASEV3 = "ASEV3"

type ClusterSettingModel struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

type AppServiceEnvironmentV3Model struct {
	Name                               string                            `tfschema:"name"`
	ResourceGroup                      string                            `tfschema:"resource_group_name"`
	SubnetId                           string                            `tfschema:"subnet_id"`
	AllowNewPrivateEndpointConnections bool                              `tfschema:"allow_new_private_endpoint_connections"`
	ClusterSetting                     []ClusterSettingModel             `tfschema:"cluster_setting"`
	DedicatedHostCount                 int64                             `tfschema:"dedicated_host_count"`
	InternalLoadBalancingMode          string                            `tfschema:"internal_load_balancing_mode"`
	RemoteDebuggingEnabled             bool                              `tfschema:"remote_debugging_enabled"`
	ZoneRedundant                      bool                              `tfschema:"zone_redundant"`
	Tags                               map[string]string                 `tfschema:"tags"`
	DnsSuffix                          string                            `tfschema:"dns_suffix"`
	ExternalInboundIPAddresses         []string                          `tfschema:"external_inbound_ip_addresses"`
	InboundNetworkDependencies         []AppServiceV3InboundDependencies `tfschema:"inbound_network_dependencies"`
	InternalInboundIPAddresses         []string                          `tfschema:"internal_inbound_ip_addresses"`
	IpSSLAddressCount                  int64                             `tfschema:"ip_ssl_address_count"`
	LinuxOutboundIPAddresses           []string                          `tfschema:"linux_outbound_ip_addresses"`
	Location                           string                            `tfschema:"location"`
	PricingTier                        string                            `tfschema:"pricing_tier"`
	WindowsOutboundIPAddresses         []string                          `tfschema:"windows_outbound_ip_addresses"`
}

type AppServiceV3InboundDependencies struct {
	Description string   `tfschema:"description"`
	IPAddresses []string `tfschema:"ip_addresses"`
	Ports       []string `tfschema:"ports"`
}

// (@jackofallops) - Important property missing from the SDK / Swagger that will need to be added later: `upgrade_preference` https://docs.microsoft.com/en-us/azure/app-service/environment/using#upgrade-preference

type AppServiceEnvironmentV3Resource struct{}

var (
	_ sdk.Resource           = AppServiceEnvironmentV3Resource{}
	_ sdk.ResourceWithUpdate = AppServiceEnvironmentV3Resource{}
)

func (r AppServiceEnvironmentV3Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AppServiceEnvironmentName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"subnet_id": { // (@jackofallops) - This _should_ be updatable via `ChangeVnet`, but the service returns Code="NotImplemented" Message="The requested method is not implemented."
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"allow_new_private_endpoint_connections": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"cluster_setting": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"value": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"dedicated_host_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(2, 2), // Docs suggest is limited to 2 physical hosts at this time
			ConflictsWith: []string{
				"zone_redundant",
			},
		},

		"internal_load_balancing_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(appserviceenvironments.LoadBalancingModeNone),
			ValidateFunc: validation.StringInSlice([]string{
				string(appserviceenvironments.LoadBalancingModeNone),
				string(appserviceenvironments.LoadBalancingModeWebPublishing),
			}, false),
		},

		"remote_debugging_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"zone_redundant": {
			Type:     pluginsdk.TypeBool,
			ForceNew: true,
			Optional: true,
			Default:  false,
			ConflictsWith: []string{
				"dedicated_host_count",
			},
		},

		"tags": tags.Schema(),
	}
}

func (r AppServiceEnvironmentV3Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dns_suffix": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"external_inbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"inbound_network_dependencies": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"ip_addresses": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"ports": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"internal_inbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"ip_ssl_address_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"linux_outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"pricing_tier": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"windows_outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r AppServiceEnvironmentV3Resource) ModelObject() interface{} {
	return &AppServiceEnvironmentV3Model{}
}

func (r AppServiceEnvironmentV3Resource) ResourceType() string {
	return "azurerm_app_service_environment_v3"
}

func (r AppServiceEnvironmentV3Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.AppServiceEnvironmentClient
			networksClient := metadata.Client.Network.VirtualNetworks
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AppServiceEnvironmentV3Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			subnet, err := commonids.ParseSubnetID(model.SubnetId)
			if err != nil {
				return err
			}

			vnetId := commonids.NewVirtualNetworkID(subnet.SubscriptionId, subnet.ResourceGroupName, subnet.VirtualNetworkName)

			vnet, err := networksClient.Get(ctx, vnetId, virtualnetworks.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving Virtual Network %q (Resource Group %q): %+v", subnet.VirtualNetworkName, subnet.ResourceGroupName, err)
			}
			if vnet.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", subnet)
			}

			vnetLoc := location.NormalizeNilable(vnet.Model.Location)
			if vnetLoc == "" {
				return fmt.Errorf("determining Location from Virtual Network %q (Resource Group %q): `location` was missing", subnet.VirtualNetworkName, subnet.ResourceGroupName)
			}

			id := commonids.NewAppServiceEnvironmentID(subscriptionId, model.ResourceGroup, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			envelope := appserviceenvironments.AppServiceEnvironmentResource{
				Kind:     pointer.To(KindASEV3),
				Location: location.Normalize(vnetLoc),
				Properties: &appserviceenvironments.AppServiceEnvironment{
					DedicatedHostCount:        pointer.To(model.DedicatedHostCount),
					ClusterSettings:           expandClusterSettingsModel(model.ClusterSetting),
					InternalLoadBalancingMode: pointer.To(appserviceenvironments.LoadBalancingMode(model.InternalLoadBalancingMode)),
					VirtualNetwork: appserviceenvironments.VirtualNetworkProfile{
						Id: model.SubnetId,
					},
					ZoneRedundant: pointer.To(model.ZoneRedundant),
				},
				Tags: pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, envelope); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// Networking config cannot be sent in the initial create and must be updated post-creation.
			aseNetworkConfig := appserviceenvironments.AseV3NetworkingConfiguration{
				Properties: &appserviceenvironments.AseV3NetworkingConfigurationProperties{
					AllowNewPrivateEndpointConnections: pointer.To(model.AllowNewPrivateEndpointConnections),
					RemoteDebugEnabled:                 pointer.To(model.RemoteDebuggingEnabled),
				},
			}

			if _, err := client.UpdateAseNetworkingConfiguration(ctx, id, aseNetworkConfig); err != nil {
				return fmt.Errorf("setting Allow New Private Endpoint Connections on %s: %+v", id, err)
			}

			// Updating Network Config returns quickly, but is actually async on some properties.
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("the Network Configuration Update request context had no deadline")
			}

			updateWait := &pluginsdk.StateChangeConf{
				Pending:      []string{"Pending"},
				Target:       []string{"Succeeded"},
				PollInterval: 10 * time.Second,
				Delay:        10 * time.Second,
				Timeout:      time.Until(deadline),
				Refresh:      checkNetworkConfigUpdate(ctx, client, id, *aseNetworkConfig.Properties),
			}

			if _, err := updateWait.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for Network Update for %s to complete: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AppServiceEnvironmentV3Resource) Read() sdk.ResourceFunc {
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
					state.InternalLoadBalancingMode = string(pointer.From(props.InternalLoadBalancingMode))
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

func (r AppServiceEnvironmentV3Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.AppServiceEnvironmentClient

			id, err := commonids.ParseAppServiceEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			deleteOpts := appserviceenvironments.DeleteOperationOptions{
				ForceDelete: pointer.To(false),
			}

			if err := client.DeleteThenPoll(ctx, *id, deleteOpts); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r AppServiceEnvironmentV3Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AppServiceEnvironmentID
}

func (r AppServiceEnvironmentV3Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.AppServiceEnvironmentClient

			id, err := commonids.ParseAppServiceEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Info("Decoding state...")
			var state AppServiceEnvironmentV3Model
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := existing.Model
			if model == nil {
				return fmt.Errorf("reading %s for update: model was nil", *id)
			}

			metadata.Logger.Infof("updating %s", id)

			if metadata.ResourceData.HasChange("cluster_setting") {
				model.Properties.ClusterSettings = expandClusterSettingsModel(state.ClusterSetting)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(state.Tags)
			}

			aseNetworkConfig := appserviceenvironments.AseV3NetworkingConfiguration{
				Properties: &appserviceenvironments.AseV3NetworkingConfigurationProperties{
					AllowNewPrivateEndpointConnections: pointer.To(state.AllowNewPrivateEndpointConnections),
					RemoteDebugEnabled:                 pointer.To(state.RemoteDebuggingEnabled),
				},
			}

			if _, err := client.UpdateAseNetworkingConfiguration(ctx, *id, aseNetworkConfig); err != nil {
				return fmt.Errorf("setting Allow New Private Endpoint Connections on %s: %+v", id, err)
			}

			// Updating Network Config returns quickly, but is actually async on some properties. e.g. `RemoteDebuggingEnabled`
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("the Network Configuration Update request context had no deadline")
			}

			updateWait := &pluginsdk.StateChangeConf{
				Pending:      []string{"Pending"},
				Target:       []string{"Succeeded"},
				PollInterval: 10 * time.Second,
				Delay:        10 * time.Second,
				Timeout:      time.Until(deadline),
				Refresh:      checkNetworkConfigUpdate(ctx, client, *id, *aseNetworkConfig.Properties),
			}

			if _, err := updateWait.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for Network Update for %s to complete: %+v", *id, err)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func flattenClusterSettingsModel(input *[]appserviceenvironments.NameValuePair) []ClusterSettingModel {
	var output []ClusterSettingModel
	if input == nil || len(*input) == 0 {
		return output
	}

	for _, v := range *input {
		if v.Name == nil {
			continue
		}

		output = append(output, ClusterSettingModel{
			Name:  *v.Name,
			Value: utils.NormalizeNilableString(v.Value),
		})
	}
	return output
}

func expandClusterSettingsModel(input []ClusterSettingModel) *[]appserviceenvironments.NameValuePair {
	var clusterSettings []appserviceenvironments.NameValuePair
	if input == nil {
		return &clusterSettings
	}

	for _, v := range input {
		clusterSettings = append(clusterSettings, appserviceenvironments.NameValuePair{
			Name:  utils.String(v.Name),
			Value: utils.String(v.Value),
		})
	}

	return &clusterSettings
}

func flattenInboundNetworkDependencies(ctx context.Context, client *appserviceenvironments.AppServiceEnvironmentsClient, id *commonids.AppServiceEnvironmentId) (*[]AppServiceV3InboundDependencies, error) {
	var results []AppServiceV3InboundDependencies
	inboundNetworking, err := client.GetInboundNetworkDependenciesEndpointsComplete(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading paged results for Inbound Network Dependencies for %s: %+v", id, err)
	}
	for _, v := range inboundNetworking.Items {
		if err != nil {
			return nil, fmt.Errorf("reading Inbound Network dependencies for %s: %+v", id, err)
		}
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

func checkNetworkConfigUpdate(ctx context.Context, client *appserviceenvironments.AppServiceEnvironmentsClient, id commonids.AppServiceEnvironmentId, values appserviceenvironments.AseV3NetworkingConfigurationProperties) pluginsdk.StateRefreshFunc {
	return func() (result interface{}, state string, err error) {
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
