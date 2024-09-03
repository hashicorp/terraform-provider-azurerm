package arckubernetes

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridazurekubernetesservice/2024-01-01/provisionedclusterinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridazurekubernetesservice/2024-01-01/virtualnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/arckubernetes/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/arckubernetes/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = ArcKubernetesProvisionedClusterResource{}
	_ sdk.ResourceWithUpdate = ArcKubernetesProvisionedClusterResource{}
)

type ArcKubernetesProvisionedClusterResource struct{}

func (ArcKubernetesProvisionedClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ArcKubernetesProvisionedClusterID
}

func (ArcKubernetesProvisionedClusterResource) ResourceType() string {
	return "azurerm_arc_kubernetes_provisioned_cluster"
}

func (ArcKubernetesProvisionedClusterResource) ModelObject() interface{} {
	return &ArcKubernetesProvisionedClusterResourceModel{}
}

type ArcKubernetesProvisionedClusterResourceModel struct {
	AgentPoolProfile      []ProvisionedClusterInstanceAgentPoolProfile       `tfschema:"agent_pool_profile"`
	CloudProviderProfile  []ProvisionedClusterInstanceCloudProviderProfile   `tfschema:"cloud_provider_profile"`
	ClusterID             string                                             `tfschema:"cluster_id"`
	ClusterVmAcessProfile []ProvisionedClusterInstanceClusterVmAccessProfile `tfschema:"cluster_vm_access_profile"`
	ControlPlaneProfile   []ProvisionedClusterInstanceControlPlaneProfile    `tfschema:"control_plane_profile"`
	CustomLocationId      string                                             `tfschema:"custom_location_id"`
	KubernetesVersion     string                                             `tfschema:"kubernetes_version"`
	LinuxProfile          []ProvisionedClusterInstanceLinuxProfile           `tfschema:"linux_profile"`
	LicenseProfile        []ProvisionedClusterInstanceLicenseProfile         `tfschema:"license_profile"`
	NetworkProfile        []ProvisionedClusterInstanceNetworkProfile         `tfschema:"network_profile"`
	StorageProfile        []ProvisionedClusterInstanceStorageProfile         `tfschema:"storage_profile"`
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

type ProvisionedClusterInstanceCloudProviderProfile struct {
	InfraNetworkProfile []ProvisionedClusterInstanceInfraNetworkProfile `tfschema:"infra_network_profile"`
}

type ProvisionedClusterInstanceInfraNetworkProfile struct {
	VnetSubnetIds []string `tfschema:"vnet_subnet_ids"`
}

type ProvisionedClusterInstanceClusterVmAccessProfile struct {
	AuthorizedIPRanges string `tfschema:"authorized_ip_ranges"`
}

type ProvisionedClusterInstanceControlPlaneProfile struct {
	Count  int64  `tfschema:"count"`
	HostIp string `tfschema:"host_ip"`
	VmSize string `tfschema:"vm_size"`
}

type ProvisionedClusterInstanceLinuxProfile struct {
	SshKey []string `tfschema:"ssh_key"`
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

func (ArcKubernetesProvisionedClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cluster_id": commonschema.ResourceIDReferenceRequiredForceNew(&connectedclusters.ConnectedClusterId{}),

		"custom_location_id": commonschema.ResourceIDReferenceRequiredForceNew(&customlocations.CustomLocationId{}),

		"agent_pool_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
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
					"os_sku": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(provisionedclusterinstances.OSSKUCBLMariner),
							string(provisionedclusterinstances.OSSKUWindowsTwoZeroOneNine),
							string(provisionedclusterinstances.OSSKUWindowsTwoZeroTwoTwo),
						}, false),
					},
					"os_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(provisionedclusterinstances.OsTypeLinux),
							string(provisionedclusterinstances.OsTypeWindows),
						}, false),
					},
					"vm_size": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"auto_scaling_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
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
						ValidateFunc: validation.IntAtLeast(1),
					},
					"max_pods": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
					"min_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
					"node_labels": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"node_taints": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
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
									Type:     pluginsdk.TypeList,
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
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"vm_size": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"host_ip": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ForceNew:      true,
						ValidateFunc:  validation.IsIPv4Address,
						ConflictsWith: []string{"network_profile.0.load_balancer_profile"},
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
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_policy": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(provisionedclusterinstances.NetworkPolicyCalico),
						}, false),
					},
					"pod_cidr": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsCIDR,
					},
					"load_balancer_profile": {
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ForceNew:      true,
						MaxItems:      1,
						ConflictsWith: []string{"control_plane_profile.0.host_ip"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntAtLeast(1),
								},
							},
						},
					},
				},
			},
		},

		"cluster_vm_access_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authorized_ip_ranges": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
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
						ForceNew: true,
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
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"azure_hybrid_benefit": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(connectedclusters.AzureHybridBenefitNotApplicable),
						ValidateFunc: validation.StringInSlice([]string{
							string(provisionedclusterinstances.AzureHybridBenefitNotApplicable),
							string(provisionedclusterinstances.AzureHybridBenefitFalse),
							string(provisionedclusterinstances.AzureHybridBenefitTrue),
						}, false),
					},
				},
			},
		},

		"storage_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"smb_csi_driver_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"nfs_csi_driver_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},
	}
}

func (ArcKubernetesProvisionedClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ArcKubernetesProvisionedClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ProvisionedClusterInstancesClient

			var config ArcKubernetesProvisionedClusterResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			connectedClusterId, err := connectedclusters.ParseConnectedClusterID(config.ClusterID)
			if err != nil {
				return err
			}

			scopeId := commonids.NewScopeID(connectedClusterId.ID())
			provisionedClusterInstanceId := parse.NewArcKubernetesProvisionedClusterID(connectedClusterId.SubscriptionId, connectedClusterId.ResourceGroupName, connectedClusterId.ConnectedClusterName, "default")

			existing, err := client.ProvisionedClusterInstancesGet(ctx, scopeId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", provisionedClusterInstanceId, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), provisionedClusterInstanceId)
			}

			payload := provisionedclusterinstances.ProvisionedCluster{
				ExtendedLocation: &provisionedclusterinstances.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(provisionedclusterinstances.ExtendedLocationTypesCustomLocation),
				},
				Properties: &provisionedclusterinstances.ProvisionedClusterProperties{
					AgentPoolProfiles:      expandProvisionedClusterAgentPoolProfiles(config.AgentPoolProfile),
					CloudProviderProfile:   expandProvisionedClusterCloudProviderProfile(config.CloudProviderProfile),
					ClusterVMAccessProfile: expandProvisionedClusterClusterVmAccessProfile(config.ClusterVmAcessProfile),
					ControlPlane:           expandProvisionedClusterControlPlaneProfile(config.ControlPlaneProfile),
					KubernetesVersion:      pointer.To(config.KubernetesVersion),
					LinuxProfile:           expandProvisionedClusterLinuxProfile(config.LinuxProfile),
					LicenseProfile:         expandProvisionedClusterLicenseProfile(config.LicenseProfile),
					NetworkProfile:         expandProvisionedClusterNetworkProfile(config.NetworkProfile),
					StorageProfile:         expandProvisionedClusterStorageProfile(config.StorageProfile),
				},
			}

			if err := client.ProvisionedClusterInstancesCreateOrUpdateThenPoll(ctx, scopeId, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", provisionedClusterInstanceId, err)
			}

			metadata.SetID(provisionedClusterInstanceId)

			return nil
		},
	}
}

func (r ArcKubernetesProvisionedClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ProvisionedClusterInstancesClient

			id, err := parse.ArcKubernetesProvisionedClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			connectedClusterId := connectedclusters.NewConnectedClusterID(id.SubscriptionId, id.ResourceGroup, id.ConnectedClusterName)
			scopeId := commonids.NewScopeID(connectedClusterId.ID())

			resp, err := client.ProvisionedClusterInstancesGet(ctx, scopeId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema := ArcKubernetesProvisionedClusterResourceModel{
				ClusterID: connectedClusterId.ID(),
			}

			if model := resp.Model; model != nil {
				if model.ExtendedLocation != nil && model.ExtendedLocation.Name != nil {
					customLocationId, err := customlocations.ParseCustomLocationIDInsensitively(*model.ExtendedLocation.Name)
					if err != nil {
						return err
					}

					schema.CustomLocationId = customLocationId.ID()
				}

				if props := model.Properties; props != nil {
					schema.CloudProviderProfile = flattenProvisionedClusterCloudProviderProfile(props.CloudProviderProfile)
					schema.AgentPoolProfile = flattenProvisionedClusterAgentPoolProfiles(props.AgentPoolProfiles)
					schema.ClusterVmAcessProfile = flattenProvisionedClusterClusterVmAccessProfile(props.ClusterVMAccessProfile)
					schema.ControlPlaneProfile = flattenProvisionedClusterControlPlaneProfile(props.ControlPlane)
					schema.KubernetesVersion = pointer.From(props.KubernetesVersion)
					schema.LinuxProfile = flattenProvisionedClusterLinuxProfile(props.LinuxProfile)
					schema.LicenseProfile = flattenProvisionedClusterLicenseProfile(props.LicenseProfile)
					schema.NetworkProfile = flattenProvisionedClusterNetworkProfile(props.NetworkProfile)
					schema.StorageProfile = flattenProvisionedClusterStorageProfile(props.StorageProfile)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ArcKubernetesProvisionedClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ProvisionedClusterInstancesClient

			var config ArcKubernetesProvisionedClusterResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := parse.ArcKubernetesProvisionedClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			connectedClusterId := connectedclusters.NewConnectedClusterID(id.SubscriptionId, id.ResourceGroup, id.ConnectedClusterName)
			scopeId := commonids.NewScopeID(connectedClusterId.ID())

			resp, err := client.ProvisionedClusterInstancesGet(ctx, scopeId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			parameters := resp.Model
			if parameters == nil || parameters.Properties == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if metadata.ResourceData.HasChange("agent_pool_profile") {
				parameters.Properties.AgentPoolProfiles = expandProvisionedClusterAgentPoolProfiles(config.AgentPoolProfile)
			}

			if metadata.ResourceData.HasChange("control_plane_profile") {
				parameters.Properties.ControlPlane = expandProvisionedClusterControlPlaneProfile(config.ControlPlaneProfile)
			}

			if metadata.ResourceData.HasChange("license_profile") {
				parameters.Properties.LicenseProfile = expandProvisionedClusterLicenseProfile(config.LicenseProfile)
			}

			if metadata.ResourceData.HasChange("storage_profile") {
				parameters.Properties.StorageProfile = expandProvisionedClusterStorageProfile(config.StorageProfile)
			}

			if err := client.ProvisionedClusterInstancesCreateOrUpdateThenPoll(ctx, scopeId, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r ArcKubernetesProvisionedClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ProvisionedClusterInstancesClient

			var config ArcKubernetesProvisionedClusterResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := parse.ArcKubernetesProvisionedClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			connectedClusterId := connectedclusters.NewConnectedClusterID(id.SubscriptionId, id.ResourceGroup, id.ConnectedClusterName)
			scopeId := commonids.NewScopeID(connectedClusterId.ID())

			if err := client.ProvisionedClusterInstancesDeleteThenPoll(ctx, scopeId); err != nil {
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
	if len(input) == 0 || len(input[0].InfraNetworkProfile) == 0 {
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
	if input == nil || input.InfraNetworkProfile == nil {
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

func expandProvisionedClusterClusterVmAccessProfile(input []ProvisionedClusterInstanceClusterVmAccessProfile) *provisionedclusterinstances.ClusterVMAccessProfile {
	if len(input) == 0 {
		return nil
	}

	return &provisionedclusterinstances.ClusterVMAccessProfile{
		AuthorizedIPRanges: pointer.To(input[0].AuthorizedIPRanges),
	}
}

func flattenProvisionedClusterClusterVmAccessProfile(input *provisionedclusterinstances.ClusterVMAccessProfile) []ProvisionedClusterInstanceClusterVmAccessProfile {
	if input == nil || input.AuthorizedIPRanges == nil {
		return make([]ProvisionedClusterInstanceClusterVmAccessProfile, 0)
	}

	return []ProvisionedClusterInstanceClusterVmAccessProfile{
		{
			AuthorizedIPRanges: pointer.From(input.AuthorizedIPRanges),
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
			KeyData: pointer.To(v),
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

	sshKey := make([]string, 0)
	for _, v := range *input.Ssh.PublicKeys {
		sshKey = append(sshKey, pointer.From(v.KeyData))
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
	if input.LoadBalancerProfile != nil && input.LoadBalancerProfile.Count != nil {
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
