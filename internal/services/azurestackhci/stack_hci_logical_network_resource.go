// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

var (
	_ sdk.Resource           = StackHCILogicalNetworkResource{}
	_ sdk.ResourceWithUpdate = StackHCILogicalNetworkResource{}
)

type StackHCILogicalNetworkResource struct{}

func (StackHCILogicalNetworkResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return logicalnetworks.ValidateLogicalNetworkID
}

func (StackHCILogicalNetworkResource) ResourceType() string {
	return "azurerm_stack_hci_logical_network"
}

func (StackHCILogicalNetworkResource) ModelObject() interface{} {
	return &StackHCILogicalNetworkResourceModel{}
}

type StackHCILogicalNetworkResourceModel struct {
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Location          string                 `tfschema:"location"`
	CustomLocationId  string                 `tfschema:"custom_location_id"`
	DNSServers        []string               `tfschema:"dns_servers"`
	Subnet            []StackHCISubnetModel  `tfschema:"subnet"`
	VirtualSwitchName string                 `tfschema:"virtual_switch_name"`
	Tags              map[string]interface{} `tfschema:"tags"`
}

type StackHCISubnetModel struct {
	AddressPrefix      string                `tfschema:"address_prefix"`
	IpAllocationMethod string                `tfschema:"ip_allocation_method"`
	IpPool             []StackHCIIPPoolModel `tfschema:"ip_pool"`
	Route              []StackHCIRouteModel  `tfschema:"route"`
	VlanId             int64                 `tfschema:"vlan_id"`
}

type StackHCIIPPoolModel struct {
	Start string `tfschema:"start"`
	End   string `tfschema:"end"`
}

type StackHCIRouteModel struct {
	Name             string `tfschema:"name"`
	AddressPrefix    string `tfschema:"address_prefix"`
	NextHopIpAddress string `tfschema:"next_hop_ip_address"`
}

func (StackHCILogicalNetworkResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][\-\.\_a-zA-Z0-9]{0,62}[a-zA-Z0-9]$`),
				"name must be between 2 and 64 characters and can only contain alphanumberic characters, hyphen, dot and underline",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: customlocations.ValidateCustomLocationID,
		},

		"virtual_switch_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dns_servers": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsIPv4Address,
			},
		},

		"subnet": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_allocation_method": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(logicalnetworks.PossibleValuesForIPAllocationMethodEnum(), false),
					},

					"address_prefix": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsCIDR,
					},

					"ip_pool": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"start": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},
								"end": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},
							},
						},
					},

					"route": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"address_prefix": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsCIDR,
								},

								"next_hop_ip_address": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},

								"name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^[a-zA-Z0-9][\-\.\_a-zA-Z0-9]{0,78}[a-zA-Z0-9]$`),
										"name must be between 2 and 80 characters and can only contain alphanumberic characters, hyphen, dot and underline",
									),
								},
							},
						},
					},

					"vlan_id": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (StackHCILogicalNetworkResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StackHCILogicalNetworkResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.LogicalNetworks

			var config StackHCILogicalNetworkResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := logicalnetworks.NewLogicalNetworkID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := logicalnetworks.LogicalNetworks{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				ExtendedLocation: &logicalnetworks.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(logicalnetworks.ExtendedLocationTypesCustomLocation),
				},
				Properties: &logicalnetworks.LogicalNetworkProperties{
					VMSwitchName: pointer.To(config.VirtualSwitchName),
					Subnets:      expandStackHCILogicalNetworkSubnet(config.Subnet),
					DhcpOptions: &logicalnetworks.LogicalNetworkPropertiesDhcpOptions{
						DnsServers: pointer.To(config.DNSServers),
					},
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("performing create %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r StackHCILogicalNetworkResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.LogicalNetworks

			id, err := logicalnetworks.ParseLogicalNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			schema := StackHCILogicalNetworkResourceModel{
				Name:              id.LogicalNetworkName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				schema.Location = location.Normalize(model.Location)
				schema.Tags = tags.Flatten(model.Tags)

				if model.ExtendedLocation != nil && model.ExtendedLocation.Name != nil {
					customLocationId, err := customlocations.ParseCustomLocationIDInsensitively(*model.ExtendedLocation.Name)
					if err != nil {
						return err
					}

					schema.CustomLocationId = customLocationId.ID()
				}

				if props := model.Properties; props != nil {
					schema.Subnet = flattenStackHCILogicalNetworkSubnet(props.Subnets)
					schema.VirtualSwitchName = pointer.From(props.VMSwitchName)

					if props.DhcpOptions != nil {
						schema.DNSServers = pointer.From(props.DhcpOptions.DnsServers)
					}
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r StackHCILogicalNetworkResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.LogicalNetworks

			id, err := logicalnetworks.ParseLogicalNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StackHCILogicalNetworkResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			parameters := resp.Model
			if parameters == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = tags.Expand(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r StackHCILogicalNetworkResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.LogicalNetworks

			id, err := logicalnetworks.ParseLogicalNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandStackHCILogicalNetworkSubnet(input []StackHCISubnetModel) *[]logicalnetworks.Subnet {
	if len(input) == 0 {
		return nil
	}

	results := make([]logicalnetworks.Subnet, 0)
	for _, v := range input {
		results = append(results, logicalnetworks.Subnet{
			Properties: &logicalnetworks.SubnetPropertiesFormat{
				AddressPrefix:      pointer.To(v.AddressPrefix),
				IPAllocationMethod: pointer.To(logicalnetworks.IPAllocationMethodEnum(v.IpAllocationMethod)),
				IPPools:            expandStackHCILogicalNetworkIPPool(v.IpPool),
				RouteTable:         expandStackHCILogicalNetworkRouteTable(v.Route),
				Vlan:               pointer.To(v.VlanId),
			},
		})
	}

	return &results
}

func flattenStackHCILogicalNetworkSubnet(input *[]logicalnetworks.Subnet) []StackHCISubnetModel {
	if input == nil {
		return make([]StackHCISubnetModel, 0)
	}

	results := make([]StackHCISubnetModel, 0)
	for _, v := range *input {
		if v.Properties != nil {
			results = append(results, StackHCISubnetModel{
				AddressPrefix:      pointer.From(v.Properties.AddressPrefix),
				IpAllocationMethod: string(pointer.From(v.Properties.IPAllocationMethod)),
				IpPool:             flattenStackHCILogicalNetworkIPPool(v.Properties.IPPools),
				Route:              flattenStackHCILogicalNetworkRouteTable(v.Properties.RouteTable),
				VlanId:             pointer.From(v.Properties.Vlan),
			})
		}
	}

	return results
}

func expandStackHCILogicalNetworkIPPool(input []StackHCIIPPoolModel) *[]logicalnetworks.IPPool {
	if len(input) == 0 {
		return nil
	}

	results := make([]logicalnetworks.IPPool, 0)
	for _, v := range input {
		results = append(results, logicalnetworks.IPPool{
			Start: pointer.To(v.Start),
			End:   pointer.To(v.End),
		})
	}

	return &results
}

func flattenStackHCILogicalNetworkIPPool(input *[]logicalnetworks.IPPool) []StackHCIIPPoolModel {
	if input == nil {
		return make([]StackHCIIPPoolModel, 0)
	}

	results := make([]StackHCIIPPoolModel, 0)
	for _, v := range *input {
		results = append(results, StackHCIIPPoolModel{
			Start: pointer.From(v.Start),
			End:   pointer.From(v.End),
		})
	}

	return results
}

func expandStackHCILogicalNetworkRouteTable(input []StackHCIRouteModel) *logicalnetworks.RouteTable {
	if len(input) == 0 {
		return nil
	}

	routes := make([]logicalnetworks.Route, 0)
	for _, v := range input {
		route := logicalnetworks.Route{
			Properties: &logicalnetworks.RoutePropertiesFormat{
				AddressPrefix:    pointer.To(v.AddressPrefix),
				NextHopIPAddress: pointer.To(v.NextHopIpAddress),
			},
		}

		if v.Name != "" {
			route.Name = pointer.To(v.Name)
		}

		routes = append(routes, route)
	}

	return &logicalnetworks.RouteTable{
		Properties: &logicalnetworks.RouteTablePropertiesFormat{
			Routes: pointer.To(routes),
		},
	}
}

func flattenStackHCILogicalNetworkRouteTable(input *logicalnetworks.RouteTable) []StackHCIRouteModel {
	if input == nil || input.Properties == nil || input.Properties.Routes == nil {
		return make([]StackHCIRouteModel, 0)
	}

	results := make([]StackHCIRouteModel, 0)
	for _, v := range *input.Properties.Routes {
		route := StackHCIRouteModel{
			Name: pointer.From(v.Name),
		}

		if v.Properties != nil {
			route.AddressPrefix = pointer.From(v.Properties.AddressPrefix)
			route.NextHopIpAddress = pointer.From(v.Properties.NextHopIPAddress)
		}
		results = append(results, route)
	}

	return results
}
