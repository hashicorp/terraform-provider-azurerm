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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppServiceEnvironmentV3DataSource struct{}

var _ sdk.DataSource = AppServiceEnvironmentV3DataSource{}

type AppServiceEnvironmentV3DataSourceModel struct {
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

func (r AppServiceEnvironmentV3DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AppServiceEnvironmentName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r AppServiceEnvironmentV3DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"allow_new_private_endpoint_connections": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"cluster_setting": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"value": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"dedicated_host_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

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

		"internal_load_balancing_mode": {
			Type:     pluginsdk.TypeString,
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

		"ip_ssl_address_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"remote_debugging_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"subnet_id": {
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

		"zone_redundant": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (r AppServiceEnvironmentV3DataSource) ModelObject() interface{} {
	return &AppServiceEnvironmentV3Model{}
}

func (r AppServiceEnvironmentV3DataSource) ResourceType() string {
	return "azurerm_app_service_environment_v3"
}

func (r AppServiceEnvironmentV3DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.AppServiceEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state AppServiceEnvironmentV3DataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := commonids.NewAppServiceEnvironmentID(subscriptionId, state.ResourceGroup, state.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if props := model.Properties; props != nil {
					state.SubnetId = props.VirtualNetwork.Id
					state.InternalLoadBalancingMode = string(pointer.From(props.InternalLoadBalancingMode))
					state.DedicatedHostCount = pointer.From(props.DedicatedHostCount)
					state.PricingTier = pointer.From(props.MultiSize)
					state.ClusterSetting = flattenClusterSettingsModel(props.ClusterSettings)
					state.DnsSuffix = utils.NormalizeNilableString(props.DnsSuffix)
					state.IpSSLAddressCount = pointer.From(props.IPsslAddressCount)
					state.ZoneRedundant = pointer.From(props.ZoneRedundant)
				}

				existingNetwork, err := client.GetAseV3NetworkingConfiguration(ctx, id)
				if err != nil {
					return fmt.Errorf("reading network configuration for %s: %+v", id, err)
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

				inboundNetworkDependencies, err := flattenInboundNetworkDependencies(ctx, client, &id)
				if err != nil {
					return err
				}
				state.InboundNetworkDependencies = *inboundNetworkDependencies

				state.Tags = pointer.From(model.Tags)

			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
