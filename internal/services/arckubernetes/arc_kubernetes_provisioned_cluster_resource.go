package arckubernetes

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridazurekubernetesservice/2024-01-01/provisionedclusterinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridazurekubernetesservice/2024-01-01/virtualnetworks"
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
	AgentPoolProfile     []ProvisionedClusterInstanceAgentPoolProfile     `tfschema:"agent_pool_profile"`
	AutoScalerProfile    []ProvisionedClusterInstanceAutoScalerProfile    `tfschema:"auto_scaler_profile"`
	CloudProviderProfile []ProvisionedClusterInstanceCloudProviderProfile `tfschema:"cloud_provider_profile"`
	ConnectedClusterID   string                                           `tfschema:"connected_cluster_id"`
	ControlPlaneProfile  []ProvisionedClusterInstanceControlPlaneProfile  `tfschema:"control_plane_profile"`
	CustomLocationId     string                                           `tfschema:"custom_location_id"`
	KubernetesVersion    string                                           `tfschema:"kubernetes_version"`
	LinuxProfile         []ProvisionedClusterInstanceLinuxProfile         `tfschema:"linux_profile"`
	LicenseProfile       []ProvisionedClusterInstanceLicenseProfile       `tfschema:"license_profile"`
	NetworkProfile       []ProvisionedClusterInstanceNetworkProfile       `tfschema:"network_profile"`
	StorageProfile       []ProvisionedClusterInstanceStorageProfile       `tfschema:"storage_profile"`
}

type ProvisionedClusterInstanceAgentPoolProfile struct {
	AutoScalingEnabled bool                   `tfschema:"auto_scaling_enabled"`
	Count              int64                  `tfschema:"count"`
	MaxCount           int64                  `tfschema:"max_count"`
	MaxPods            int64                  `tfschema:"max_pods"`
	MinCount           int64                  `tfschema:"min_count"`
	Name               string                 `tfschema:"name"`
	NodeLabels         map[string]interface{} `tfschema:"node_labels"`
	NodeTaints         []string               `tfschema:"node_taints"`
	OsSku              string                 `tfschema:"os_sku"`
	OsType             string                 `tfschema:"os_type"`
	VmSize             string                 `tfschema:"vm_size"`
}

type ProvisionedClusterInstanceAutoScalerProfile struct {
	BalanceSimilarNodeGroups      string `tfschema:"balance_similar_node_groups"`
	Expander                      string `tfschema:"expander"`
	MaxEmptyBulkDelete            string `tfschema:"max_empty_bulk_delete"`
	MaxGracefulTerminationSec     string `tfschema:"max_graceful_termination_sec"`
	MaxNodeProvisionTime          string `tfschema:"max_node_provision_time"`
	MaxTotalUnreadyPercentage     string `tfschema:"max_total_unready_percentage"`
	NewPodScaleUpDelay            string `tfschema:"new_pod_scale_up_delay"`
	OkTotalUnreadyCount           string `tfschema:"ok_total_unready_count"`
	ScanInterval                  string `tfschema:"scan_interval"`
	ScaleDownDelayAfterAdd        string `tfschema:"scale_down_delay_after_add"`
	ScaleDownDelayAfterDelete     string `tfschema:"scale_down_delay_after_delete"`
	ScaleDownDelayAfterFailure    string `tfschema:"scale_down_delay_after_failure"`
	ScaleDownUnneededTime         string `tfschema:"scale_down_unneeded_time"`
	ScaleDownUnreadyTime          string `tfschema:"scale_down_unready_time"`
	ScaleDownUtilizationThreshold string `tfschema:"scale_down_utilization_threshold"`
	SkipNodesWithLocalStorage     string `tfschema:"skip_nodes_with_local_storage"`
	SkipNodesWithSystemPods       string `tfschema:"skip_nodes_with_system_pods"`
}

type ProvisionedClusterInstanceCloudProviderProfile struct {
	InfraNetworkProfile []ProvisionedClusterInstanceInfraNetworkProfile `tfschema:"infra_network_profile"`
}

type ProvisionedClusterInstanceInfraNetworkProfile struct {
	VnetSubnetIds []string `tfschema:"vnet_subnet_ids"`
}

type ProvisionedClusterInstanceControlPlaneProfile struct {
	Count  int64  `tfschema:"count"`
	HostIp string `tfschema:"host_ip"`
	VmSize string `tfschema:"vm_size"`
}

type ProvisionedClusterInstanceLinuxProfile struct {
	SshKey []ProvisionedClusterInstanceSshKey `tfschema:"ssh_key"`
}

type ProvisionedClusterInstanceSshKey struct {
	KeyData string `tfschema:"key_data"`
}

type ProvisionedClusterInstanceLicenseProfile struct {
	AzureHybridBenefit string `tfschema:"azure_hybrid_benefit"`
}

type ProvisionedClusterInstanceNetworkProfile struct {
	LoadBalancerProfile []ProvisionedClusterInstanceLoadBalancerProfile `tfschema:"load_balancer_profile"`
	NetworkPolicy       string                                          `tfschema:"network_policy"`
	PodCidr             string                                          `tfschema:"pod_cidr"`
}

type ProvisionedClusterInstanceLoadBalancerProfile struct {
	Count int64 `tfschema:"count"`
}

type ProvisionedClusterInstanceStorageProfile struct {
	SmbCsiDriverEnabled bool `tfschema:"smb_csi_driver_enabled"`
	NfsCsiDriverEnabled bool `tfschema:"nfs_csi_driver_enabled"`
}

func (ArcKubernetesProvisionedClusterInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"connected_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: connectedclusters.ValidateConnectedClusterID,
		},

		"custom_location_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: customlocations.ValidateCustomLocationID,
		},

		"auto_scaler_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"balance_similar_node_groups": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"expander": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(provisionedclusterinstances.ExpanderLeastNegativewaste),
							string(provisionedclusterinstances.ExpanderMostNegativepods),
							string(provisionedclusterinstances.ExpanderPriority),
							string(provisionedclusterinstances.ExpanderRandom),
						}, false),
					},
					"max_empty_bulk_delete": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"max_graceful_termination_sec": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"max_node_provision_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"max_total_unready_percentage": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"new_pod_scale_up_delay": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"ok_total_unready_count": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"scan_interval": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"scale_down_delay_after_add": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"scale_down_delay_after_delete": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"scale_down_delay_after_failure": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"scale_down_unneeded_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"scale_down_unready_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"scale_down_utilization_threshold": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"skip_nodes_with_local_storage": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"skip_nodes_with_system_pods": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"agent_pool_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-z][a-z0-9]{2,11}$`),
							"`name` must start with lower cases characters and can only contains 3-12 lowercase alphanumberic characters"),
					},
					"vm_size": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"auto_scaling_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},
					"count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
						Default:      1,
					},
					"max_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
					"max_pods": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
					"min_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
					"node_labels": {
						Type:     pluginsdk.TypeMap,
						ForceNew: true,
						Optional: true,
					},
					"node_taints": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsIPv4Address,
						},
					},
					"os_sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(provisionedclusterinstances.OSSKUCBLMariner),
							string(provisionedclusterinstances.OSSKUWindowsTwoZeroOneNine),
							string(provisionedclusterinstances.OSSKUWindowsTwoZeroTwoTwo),
						}, false),
					},
					"os_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(provisionedclusterinstances.OsTypeLinux),
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(provisionedclusterinstances.OsTypeLinux),
							string(provisionedclusterinstances.OsTypeWindows),
						}, false),
					},
				},
			},
		},

		"cloud_provider_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"infra_network_profile": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"vnet_subnet_ids": {
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.Any(
											logicalnetworks.ValidateLogicalNetworkID,
											virtualnetworks.ValidateVirtualNetworkID,
										),
									},
								},
							},
						},
					},
				},
			},
		},

		"control_plane_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host_ip": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsIPv4Address,
					},
					"vm_size": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"kubernetes_version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"linux_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ssh_key": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"license_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"azure_hybrid_benefit": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(provisionedclusterinstances.AzureHybridBenefitNotApplicable),
						ValidateFunc: validation.StringInSlice([]string{
							string(provisionedclusterinstances.AzureHybridBenefitNotApplicable),
							string(provisionedclusterinstances.AzureHybridBenefitFalse),
							string(provisionedclusterinstances.AzureHybridBenefitTrue),
						}, false),
					},
				},
			},
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_policy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(provisionedclusterinstances.NetworkPolicyCalico),
						}, false),
					},
					"pod_cidr": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsCIDR,
					},
					"load_balancer_profile": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntAtLeast(1),
								},
							},
						},
					},
				},
			},
		},

		"storage_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"smb_csi_driver_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  true,
					},
					"nfs_csi_driver_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  true,
					},
				},
			},
		},
	}
}

func (ArcKubernetesProvisionedClusterInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ArcKubernetesProvisionedClusterInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ProvisionedClusterInstancesClient

			var config ArcKubernetesProvisionedClusterInstanceResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := connectedclusters.ParseConnectedClusterID(config.ConnectedClusterID)
			if err != nil {
				return err
			}

			scopeId := commonids.NewScopeID(id.ID())

			existing, err := client.ProvisionedClusterInstancesGet(ctx, scopeId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := provisionedclusterinstances.ProvisionedCluster{
				ExtendedLocation: &provisionedclusterinstances.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(provisionedclusterinstances.ExtendedLocationTypesCustomLocation),
				},
			}

			if err := client.ProvisionedClusterInstancesCreateOrUpdateThenPoll(ctx, scopeId, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
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
