package web

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppServiceEnvironmentV3DataSource struct{}

var _ sdk.DataSource = AppServiceEnvironmentV3DataSource{}

func (r AppServiceEnvironmentV3DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AppServiceEnvironmentName,
		},

		"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
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

		"tags": tags.SchemaDataSource(),
	}
}

func (r AppServiceEnvironmentV3DataSource) ModelObject() interface{} {
	return AppServiceEnvironmentV3Model{}
}

func (r AppServiceEnvironmentV3DataSource) ResourceType() string {
	return "azurerm_app_service_environment_v3"
}

func (r AppServiceEnvironmentV3DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Web.AppServiceEnvironmentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var appServiceEnvironmentV3 AppServiceEnvironmentV3Model
			if err := metadata.Decode(&appServiceEnvironmentV3); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := parse.NewAppServiceEnvironmentID(subscriptionId, appServiceEnvironmentV3.ResourceGroup, appServiceEnvironmentV3.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.HostingEnvironmentName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := AppServiceEnvironmentV3Model{
				Name:          id.HostingEnvironmentName,
				ResourceGroup: id.ResourceGroup,
				Location:      location.NormalizeNilable(existing.Location),
			}

			if props := existing.AppServiceEnvironment; props != nil {
				if props.VirtualNetwork != nil {
					model.SubnetId = utils.NormalizeNilableString(props.VirtualNetwork.ID)
				}
				model.InternalLoadBalancingMode = string(props.InternalLoadBalancingMode)
				model.DedicatedHostCount = int(utils.NormaliseNilableInt32(props.DedicatedHostCount))
				model.PricingTier = utils.NormalizeNilableString(props.MultiSize)
				model.ClusterSetting = flattenClusterSettingsModel(props.ClusterSettings)
				model.DnsSuffix = utils.NormalizeNilableString(props.DNSSuffix)
				model.IpSSLAddressCount = int(utils.NormaliseNilableInt32(existing.IpsslAddressCount))
				// model.ZoneRedundant = *props.ZoneRedundant
			}

			existingNetwork, err := client.GetAseV3NetworkingConfiguration(ctx, id.ResourceGroup, id.HostingEnvironmentName)
			if err != nil {
				return fmt.Errorf("reading network configuration for %s: %+v", id, err)
			}

			if props := existingNetwork.AseV3NetworkingConfigurationProperties; props != nil {
				model.WindowsOutboundIPAddresses = *props.WindowsOutboundIPAddresses
				model.LinuxOutboundIPAddresses = *props.LinuxOutboundIPAddresses
				model.AllowNewPrivateEndpointConnections = *props.AllowNewPrivateEndpointConnections
			}

			inboundNetworkDependencies := &[]AppServiceV3InboundDependencies{}
			inboundNetworkDependencies, err = flattenInboundNetworkDependencies(ctx, client, &id)
			if err != nil {
				return err
			}
			model.InboundNetworkDependencies = *inboundNetworkDependencies

			model.Tags = tags.Flatten(existing.Tags)

			metadata.SetID(id)
			return metadata.Encode(&model)
		},
	}
}
