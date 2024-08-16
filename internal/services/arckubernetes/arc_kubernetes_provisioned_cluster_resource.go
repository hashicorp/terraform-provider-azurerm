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
	AutoScalingEnabled bool              `tfschema:"auto_scaling_enabled"`
	Count              int64             `tfschema:"count"`
	MaxCount           int64             `tfschema:"max_count"`
	MaxPods            int64             `tfschema:"max_pods"`
	MinCount           int64             `tfschema:"min_count"`
	Name               string            `tfschema:"name"`
	NodeLabels         map[string]string `tfschema:"node_labels"`
	NodeTaints         []string          `tfschema:"node_taints"`
	OsSku              string            `tfschema:"os_sku"`
	OsType             string            `tfschema:"os_type"`
	VmSize             string            `tfschema:"vm_size"`
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

		"agent_pool_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
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
			Optional: true,
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

			connectedClusterId, err := connectedclusters.ParseConnectedClusterID(config.ConnectedClusterID)
			if err != nil {
				return err
			}

			scopeId := commonids.NewScopeID(connectedClusterId.ID())
			// TODO: the id should use resourceids.ResourceId
			provisionedClusterInstanceId := connectedClusterId.ID() + ArcKubernetesProvisionedClusterInstanceResourceIdSuffix

			existing, err := client.ProvisionedClusterInstancesGet(ctx, scopeId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", provisionedClusterInstanceId, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), scopeId)
			}

			payload := provisionedclusterinstances.ProvisionedCluster{
				ExtendedLocation: &provisionedclusterinstances.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(provisionedclusterinstances.ExtendedLocationTypesCustomLocation),
				},
				Properties: &provisionedclusterinstances.ProvisionedClusterProperties{
					AgentPoolProfiles:    expandProvisionedClusterAgentPoolProfiles(config.AgentPoolProfile),
					CloudProviderProfile: expandProvisionedClusterCloudProviderProfile(config.CloudProviderProfile),
					ControlPlane:         expandProvisionedClusterControlPlaneProfile(config.ControlPlaneProfile),
					KubernetesVersion:    pointer.To(config.KubernetesVersion),
					LinuxProfile:         expandProvisionedClusterLinuxProfile(config.LinuxProfile),
					LicenseProfile:       expandProvisionedClusterLicenseProfile(config.LicenseProfile),
					NetworkProfile:       expandProvisionedClusterNetworkProfile(config.NetworkProfile),
					StorageProfile:       expandProvisionedClusterStorageProfile(config.StorageProfile),
				},
			}

			if err := client.ProvisionedClusterInstancesCreateOrUpdateThenPoll(ctx, scopeId, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", connectedClusterId, err)
			}

			metadata.SetID(connectedClusterId)

			return nil
		},
	}
}

func (r ArcKubernetesProvisionedClusterInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ProvisionedClusterInstancesClient

			idRaw := metadata.ResourceData.Id()
			scopeId := strings.TrimRight(idRaw, ArcKubernetesProvisionedClusterInstanceResourceIdSuffix)
			if !strings.EqualFold(idRaw, scopeId+ArcKubernetesProvisionedClusterInstanceResourceIdSuffix) {
				return fmt.Errorf("expect `id` ends with %s, but got %s", ArcKubernetesProvisionedClusterInstanceResourceIdSuffix, idRaw)
			}

			if connectedClusterId, err := connectedclusters.ParseConnectedClusterID(scopeId); err != nil {
				return fmt.Errorf("parsing the scope of %q as a Connected Cluster ID: %+v", idRaw, err)
			}

			resp, err := client.ProvisionedClusterInstancesGet(ctx, commonids.scopeId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			schema := ArcKubernetesProvisionedClusterInstanceResourceModel{}

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

func expandProvisionedClusterAgentPoolProfiles(input []ProvisionedClusterInstanceAgentPoolProfile) *[]provisionedclusterinstances.NamedAgentPoolProfile {
	if len(input) == 0 {
		return nil
	}

	output := make([]provisionedclusterinstances.NamedAgentPoolProfile, 0)
	for _, v := range input {
		agentPool := provisionedclusterinstances.NamedAgentPoolProfile{
			Name:              pointer.To(v.Name),
			Count:             pointer.To(v.Count),
			EnableAutoScaling: pointer.To(v.AutoScalingEnabled),
			MaxCount:          pointer.To(v.MaxCount),
			MaxPods:           pointer.To(v.MaxPods),
			MinCount:          pointer.To(v.MinCount),
			NodeLabels:        pointer.To(v.NodeLabels),
			NodeTaints:        pointer.To(v.NodeTaints),
			OsSKU:             pointer.To(provisionedclusterinstances.OSSKU(v.OsSku)),
			OsType:            pointer.To(provisionedclusterinstances.OsType(v.OsType)),
			VMSize:            pointer.To(v.VmSize),
		}

		output = append(output, agentPool)
	}

	return &output
}

func flattenProvisionedClusterAgentPoolProfiles(input *[]provisionedclusterinstances.NamedAgentPoolProfile) []ProvisionedClusterInstanceAgentPoolProfile {
	if input == nil {
		return make([]ProvisionedClusterInstanceAgentPoolProfile, 0)
	}

	output := make([]ProvisionedClusterInstanceAgentPoolProfile, 0)
	for _, v := range *input {
		agentPool := ProvisionedClusterInstanceAgentPoolProfile{
			Name:               pointer.From(v.Name),
			Count:              pointer.From(v.Count),
			AutoScalingEnabled: pointer.From(v.EnableAutoScaling),
			MaxCount:           pointer.From(v.MaxCount),
			MaxPods:            pointer.From(v.MaxPods),
			MinCount:           pointer.From(v.MinCount),
			NodeLabels:         pointer.From(v.NodeLabels),
			NodeTaints:         pointer.From(v.NodeTaints),
			OsSku:              string(pointer.From(v.OsSKU)),
			OsType:             string(pointer.From(v.OsType)),
			VmSize:             pointer.From(v.VMSize),
		}

		output = append(output, agentPool)
	}

	return output
}

func expandProvisionedClusterCloudProviderProfile(input []ProvisionedClusterInstanceCloudProviderProfile) *provisionedclusterinstances.CloudProviderProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	return &provisionedclusterinstances.CloudProviderProfile{
		InfraNetworkProfile: &provisionedclusterinstances.CloudProviderProfileInfraNetworkProfile{
			VnetSubnetIds: pointer.To(v.InfraNetworkProfile[0].VnetSubnetIds),
		},
	}
}

func flattenProvisionedClusterCloudProviderProfile(input *provisionedclusterinstances.CloudProviderProfile) []ProvisionedClusterInstanceCloudProviderProfile {
	if input == nil {
		return make([]ProvisionedClusterInstanceCloudProviderProfile, 0)
	}

	return []ProvisionedClusterInstanceCloudProviderProfile{
		{
			InfraNetworkProfile: []ProvisionedClusterInstanceInfraNetworkProfile{
				{
					VnetSubnetIds: pointer.From(input.InfraNetworkProfile.VnetSubnetIds),
				},
			},
		},
	}
}

func expandProvisionedClusterControlPlaneProfile(input []ProvisionedClusterInstanceControlPlaneProfile) *provisionedclusterinstances.ControlPlaneProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	return &provisionedclusterinstances.ControlPlaneProfile{
		Count:  pointer.To(v.Count),
		VMSize: pointer.To(v.VmSize),
		ControlPlaneEndpoint: &provisionedclusterinstances.ControlPlaneProfileControlPlaneEndpoint{
			HostIP: pointer.To(v.HostIp),
		},
	}
}

func flattenProvisionedClusterControlPlaneProfile(input *provisionedclusterinstances.ControlPlaneProfile) []ProvisionedClusterInstanceControlPlaneProfile {
	if input == nil {
		return make([]ProvisionedClusterInstanceControlPlaneProfile, 0)
	}

	var hostIp string
	if input.ControlPlaneEndpoint != nil {
		hostIp = pointer.From(input.ControlPlaneEndpoint.HostIP)
	}

	return []ProvisionedClusterInstanceControlPlaneProfile{
		{
			Count:  pointer.From(input.Count),
			VmSize: pointer.From(input.VMSize),
			HostIp: hostIp,
		},
	}
}

func expandProvisionedClusterLinuxProfile(input []ProvisionedClusterInstanceLinuxProfile) *provisionedclusterinstances.LinuxProfileProperties {
	if len(input) == 0 {
		return nil
	}

	publicKeys := make([]provisionedclusterinstances.LinuxProfilePropertiesSshPublicKeysInlined, 0)
	for _, v := range input[0].SshKey {
		publicKeys = append(publicKeys, provisionedclusterinstances.LinuxProfilePropertiesSshPublicKeysInlined{
			KeyData: pointer.To(v.KeyData),
		})
	}

	return &provisionedclusterinstances.LinuxProfileProperties{
		Ssh: &provisionedclusterinstances.LinuxProfilePropertiesSsh{
			PublicKeys: &publicKeys,
		},
	}
}

func flattenProvisionedClusterLinuxProfile(input *provisionedclusterinstances.LinuxProfileProperties) []ProvisionedClusterInstanceLinuxProfile {
	if input == nil || input.Ssh == nil || input.Ssh.PublicKeys == nil {
		return nil
	}

	sshKey := make([]ProvisionedClusterInstanceSshKey, 0)
	for _, v := range *input.Ssh.PublicKeys {
		sshKey = append(sshKey, ProvisionedClusterInstanceSshKey{
			KeyData: pointer.From(v.KeyData),
		})
	}

	return []ProvisionedClusterInstanceLinuxProfile{
		{
			SshKey: sshKey,
		},
	}
}

func expandProvisionedClusterLicenseProfile(input []ProvisionedClusterInstanceLicenseProfile) *provisionedclusterinstances.ProvisionedClusterLicenseProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	return &provisionedclusterinstances.ProvisionedClusterLicenseProfile{
		AzureHybridBenefit: pointer.To(provisionedclusterinstances.AzureHybridBenefit(v.AzureHybridBenefit)),
	}
}

func flattenProvisionedClusterLicenseProfile(input *provisionedclusterinstances.ProvisionedClusterLicenseProfile) []ProvisionedClusterInstanceLicenseProfile {
	if input == nil || input.AzureHybridBenefit == nil {
		return make([]ProvisionedClusterInstanceLicenseProfile, 0)
	}

	return []ProvisionedClusterInstanceLicenseProfile{
		{
			AzureHybridBenefit: string(pointer.From(input.AzureHybridBenefit)),
		},
	}
}

func expandProvisionedClusterNetworkProfile(input []ProvisionedClusterInstanceNetworkProfile) *provisionedclusterinstances.NetworkProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	networkProfile := &provisionedclusterinstances.NetworkProfile{
		NetworkPolicy: pointer.To(provisionedclusterinstances.NetworkPolicy(v.NetworkPolicy)),
		PodCidr:       pointer.To(v.PodCidr),
	}

	if len(v.LoadBalancerProfile) > 0 && v.LoadBalancerProfile[0].Count != 0 {
		networkProfile.LoadBalancerProfile = &provisionedclusterinstances.NetworkProfileLoadBalancerProfile{
			Count: pointer.To(v.LoadBalancerProfile[0].Count),
		}
	}

	return networkProfile
}

func flattenProvisionedClusterNetworkProfile(input *provisionedclusterinstances.NetworkProfile) []ProvisionedClusterInstanceNetworkProfile {
	if input == nil {
		return make([]ProvisionedClusterInstanceNetworkProfile, 0)
	}

	loadBalancer := make([]ProvisionedClusterInstanceLoadBalancerProfile, 0)
	if input.LoadBalancerProfile != nil {
		loadBalancer = append(loadBalancer, ProvisionedClusterInstanceLoadBalancerProfile{
			Count: pointer.From(input.LoadBalancerProfile.Count),
		})
	}

	return []ProvisionedClusterInstanceNetworkProfile{
		{
			NetworkPolicy:       string(pointer.From(input.NetworkPolicy)),
			PodCidr:             pointer.From(input.PodCidr),
			LoadBalancerProfile: loadBalancer,
		},
	}
}

func expandProvisionedClusterStorageProfile(input []ProvisionedClusterInstanceStorageProfile) *provisionedclusterinstances.StorageProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	return &provisionedclusterinstances.StorageProfile{
		NfsCsiDriver: &provisionedclusterinstances.StorageProfileNfsCSIDriver{
			Enabled: pointer.To(v.NfsCsiDriverEnabled),
		},
		SmbCsiDriver: &provisionedclusterinstances.StorageProfileSmbCSIDriver{
			Enabled: pointer.To(v.SmbCsiDriverEnabled),
		},
	}
}

func flattenProvisionedClusterStorageProfile(input *provisionedclusterinstances.StorageProfile) []ProvisionedClusterInstanceStorageProfile {
	if input == nil {
		return make([]ProvisionedClusterInstanceStorageProfile, 0)
	}

	var nfsCsiDriverEnabled, smbCsiDriverEnabled bool
	if input.NfsCsiDriver != nil {
		nfsCsiDriverEnabled = pointer.From(input.NfsCsiDriver.Enabled)
	}
	if input.SmbCsiDriver != nil {
		smbCsiDriverEnabled = pointer.From(input.SmbCsiDriver.Enabled)
	}

	return []ProvisionedClusterInstanceStorageProfile{
		{
			NfsCsiDriverEnabled: nfsCsiDriverEnabled,
			SmbCsiDriverEnabled: smbCsiDriverEnabled,
		},
	}
}
