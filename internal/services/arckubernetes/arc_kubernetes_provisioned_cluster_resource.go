package arckubernetes

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridazurekubernetesservice/2024-01-01/provisionedclusterinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const ArcKubernetesProvisionedClusterInstanceResourceIdSuffix = "/providers/Microsoft.HybridContainerService/provisionedClusterInstances/default"

var (
	_ sdk.Resource           = ArcKubernetesProvisionedClusterInstanceResource{}
	_ sdk.ResourceWithUpdate = ArcKubernetesProvisionedClusterInstanceResource{}
)

type ArcKubernetesProvisionedClusterInstanceResource struct{}

func (ArcKubernetesProvisionedClusterInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(val interface{}, key string) (warns []string, errs []error) {
		idRaw, ok := val.(string)
		if !ok {
			errs = append(errs, fmt.Errorf("expected `id` to be a string but got %+v", val))
			return
		}

		scopeId := strings.TrimRight(idRaw, ArcKubernetesProvisionedClusterInstanceResourceIdSuffix)
		if !strings.EqualFold(idRaw, scopeId+ArcKubernetesProvisionedClusterInstanceResourceIdSuffix) {
			errs = append(errs, fmt.Errorf("expect `id` ends with %s, but got %s", ArcKubernetesProvisionedClusterInstanceResourceIdSuffix, idRaw))
		}

		if _, err := connectedclusters.ParseConnectedClusterID(scopeId); err != nil {
			errs = append(errs, fmt.Errorf("parsing the scope of %q as a Connected Cluster ID: %+v", idRaw, err))
			return
		}

		return
	}
}

func (ArcKubernetesProvisionedClusterInstanceResource) ResourceType() string {
	return "azurerm_arc_kubernetes_provisioned_cluster_instance"
}

func (ArcKubernetesProvisionedClusterInstanceResource) ModelObject() interface{} {
	return &ArcKubernetesProvisionedClusterInstanceResourceModel{}
}

type ArcKubernetesProvisionedClusterInstanceResourceModel struct {
	Name              string                                        `tfschema:"name"`
	ResourceGroupName string                                        `tfschema:"resource_group_name"`
	Location          string                                        `tfschema:"location"`
	CustomLocationId  string                                        `tfschema:"custom_location_id"`
	AgentPoolProfile  []ProvisionedClusterInstanceAgentPoolProfile  `tfschema:"agent_pool_profile"`
	AutoScalerProfile []ProvisionedClusterInstanceAutoScalerProfile `tfschema:"auto_scaler_profile"`
	VirtualSwitchName string                                        `tfschema:"virtual_switch_name"`
	Tags              map[string]interface{}                        `tfschema:"tags"`
}

type ProvisionedClusterInstanceAgentPoolProfile struct {
	Count              string                 `tfschema:"count"`
	AutoScalingEnabled string                 `tfschema:"auto_scaling_enabled"`
	MaxCount           int64                  `tfschema:"max_count"`
	MaxPods            int64                  `tfschema:"max_pods"`
	MinCount           int64                  `tfschema:"min_count"`
	Name               string                 `tfschema:"name"`
	NodeLabel          map[string]interface{} `tfschema:"node_label"`
	NodeTaint          []string               `tfschema:"node_taint"`
	OsSku              string                 `tfschema:"os_sku"`
	OsType             string                 `tfschema:"os_type"`
	VmSize             string                 `tfschema:"vm_size"`
}

type ProvisionedClusterInstanceAutoScalerProfile struct {
	Balan
}

type StackHCIRouteModel struct {
	Name             string `tfschema:"name"`
	AddressPrefix    string `tfschema:"address_prefix"`
	NextHopIpAddress string `tfschema:"next_hop_ip_address"`
}

func (ArcKubernetesProvisionedClusterInstanceResource) Arguments() map[string]*pluginsdk.Schema {
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
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^[a-zA-Z0-9][\-\.\_a-zA-Z0-9]{0,78}[a-zA-Z0-9]$`),
										"name must be between 2 and 80 characters and can only contain alphanumberic characters, hyphen, dot and underline",
									),
								},

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

func (ArcKubernetesProvisionedClusterInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ArcKubernetesProvisionedClusterInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.LogicalNetworks

			var config ArcKubernetesProvisionedClusterInstanceResourceModel
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

			payload := provisionedclusterinstances.ID{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				ExtendedLocation: &logicalnetworks.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(logicalnetworks.ExtendedLocationTypesCustomLocation),
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

func (r ArcKubernetesProvisionedClusterInstanceResource) Read() sdk.ResourceFunc {
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

			schema := ArcKubernetesProvisionedClusterInstanceResourceModel{
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

func (r ArcKubernetesProvisionedClusterInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.LogicalNetworks

			id, err := logicalnetworks.ParseLogicalNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ArcKubernetesProvisionedClusterInstanceResourceModel
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

func (r ArcKubernetesProvisionedClusterInstanceResource) Delete() sdk.ResourceFunc {
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
		routes = append(routes, logicalnetworks.Route{
			Name: pointer.To(v.Name),
			Properties: &logicalnetworks.RoutePropertiesFormat{
				AddressPrefix:    pointer.To(v.AddressPrefix),
				NextHopIPAddress: pointer.To(v.NextHopIpAddress),
			},
		})
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
