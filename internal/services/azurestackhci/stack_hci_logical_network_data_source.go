package azurestackhci

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StackHCILogicalNetworkDataSource struct{}

var _ sdk.DataSource = StackHCILogicalNetworkDataSource{}

type StackHCILogicalNetworkDataSourceModel struct {
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Location          string                 `tfschema:"location"`
	CustomLocationId  string                 `tfschema:"custom_location_id"`
	DNSServers        []string               `tfschema:"dns_servers"`
	Subnet            []StackHCISubnetModel  `tfschema:"subnet"`
	VirtualSwitchName string                 `tfschema:"virtual_switch_name"`
	Tags              map[string]interface{} `tfschema:"tags"`
}

func (r StackHCILogicalNetworkDataSource) ResourceType() string {
	return "azurerm_stack_hci_logical_network"
}

func (r StackHCILogicalNetworkDataSource) ModelObject() interface{} {
	return &StackHCILogicalNetworkDataSourceModel{}
}

func (r StackHCILogicalNetworkDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][\-\.\_a-zA-Z0-9]{0,62}[a-zA-Z0-9]$`),
				"name must begin and end with an alphanumeric character, be between 2 and 64 characters in length and can only contain alphanumeric characters, hyphens, periods or underscores.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (r StackHCILogicalNetworkDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"custom_location_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"virtual_switch_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"dns_servers": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"subnet": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_allocation_method": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"address_prefix": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"ip_pool": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"start": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"end": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"route": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"address_prefix": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"next_hop_ip_address": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"vlan_id": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r StackHCILogicalNetworkDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.LogicalNetworks

			var state StackHCILogicalNetworkDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := logicalnetworks.NewLogicalNetworkID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s does not exist", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)

				if model.ExtendedLocation != nil && model.ExtendedLocation.Name != nil {
					customLocationId, err := customlocations.ParseCustomLocationIDInsensitively(*model.ExtendedLocation.Name)
					if err != nil {
						return err
					}

					state.CustomLocationId = customLocationId.ID()
				}

				if props := model.Properties; props != nil {
					state.Subnet = flattenStackHCILogicalNetworkSubnet(props.Subnets)
					state.VirtualSwitchName = pointer.From(props.VMSwitchName)

					if props.DhcpOptions != nil {
						state.DNSServers = pointer.From(props.DhcpOptions.DnsServers)
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
