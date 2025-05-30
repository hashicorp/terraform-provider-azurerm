// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-09-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-09-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-09-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKubernetesClusterNodePool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKubernetesClusterNodePoolCreate,
		Read:   resourceKubernetesClusterNodePoolRead,
		Update: resourceKubernetesClusterNodePoolUpdate,
		Delete: resourceKubernetesClusterNodePoolDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NodePoolID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KubernetesClusterNodePoolV0ToV1{},
		}),

		Schema: resourceKubernetesClusterNodePoolSchema(),

		CustomizeDiff: pluginsdk.CustomDiffInSequence(
			pluginsdk.ForceNewIfChange("os_sku", func(ctx context.Context, old, new, meta interface{}) bool {
				// Ubuntu and AzureLinux are currently the only allowed Linux OSSKU Migration targets.
				if old != string(agentpools.OSSKUUbuntu) && old != string(agentpools.OSSKUAzureLinux) {
					return true
				}

				if new != string(agentpools.OSSKUUbuntu) && new != string(agentpools.OSSKUAzureLinux) {
					return true
				}

				return false
			}),
			// The behaviour of the API requires this, but this could be removed when https://github.com/Azure/azure-rest-api-specs/issues/27373 has been addressed
			pluginsdk.ForceNewIfChange("upgrade_settings.0.drain_timeout_in_minutes", func(ctx context.Context, old, new, meta interface{}) bool {
				return old != 0 && new == 0
			}),
		),
	}
}

func resourceKubernetesClusterNodePoolSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: containerValidate.KubernetesAgentPoolName,
		},

		"kubernetes_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: containerValidate.ClusterID,
		},

		"node_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(0, 1000),
		},

		"tags": commonschema.Tags(),

		"vm_size": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"host_group_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: computeValidate.HostGroupID,
		},

		// Optional
		"capacity_reservation_group_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
		},

		"eviction_policy": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.ScaleSetEvictionPolicyDelete),
				string(agentpools.ScaleSetEvictionPolicyDeallocate),
			}, false),
		},

		"kubelet_config": schemaNodePoolKubeletConfig(),

		"linux_os_config": schemaNodePoolLinuxOSConfig(),

		"fips_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"gpu_instance": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.GPUInstanceProfileMIGOneg),
				string(managedclusters.GPUInstanceProfileMIGTwog),
				string(managedclusters.GPUInstanceProfileMIGThreeg),
				string(managedclusters.GPUInstanceProfileMIGFourg),
				string(managedclusters.GPUInstanceProfileMIGSeveng),
			}, false),
		},

		"kubelet_disk_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.KubeletDiskTypeOS),
				string(agentpools.KubeletDiskTypeTemporary),
			}, false),
		},

		"max_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 1000),
		},

		"max_pods": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(agentpools.AgentPoolModeUser),
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.AgentPoolModeSystem),
				string(agentpools.AgentPoolModeUser),
			}, false),
		},

		"min_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: rather than setting `0` users should instead pass `null` here
			ValidateFunc: validation.IntBetween(0, 1000),
		},

		"node_network_profile": schemaNodePoolNetworkProfile(),

		"node_labels": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"node_public_ip_prefix_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"node_public_ip_enabled"},
		},

		// Node Taints control the behaviour of the Node Pool, as such they should not be computed and
		// must be specified/reconciled as required
		"node_taints": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"orchestrator_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"os_disk_size_gb": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"os_disk_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  agentpools.OSDiskTypeManaged,
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.OSDiskTypeEphemeral),
				string(agentpools.OSDiskTypeManaged),
			}, false),
		},

		"os_sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // defaults to Ubuntu if using Linux
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.OSSKUAzureLinux),
				string(agentpools.OSSKUUbuntu),
				string(agentpools.OSSKUWindowsTwoZeroOneNine),
				string(agentpools.OSSKUWindowsTwoZeroTwoTwo),
			}, false),
		},

		"os_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(agentpools.OSTypeLinux),
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.OSTypeLinux),
				string(agentpools.OSTypeWindows),
			}, false),
		},

		"pod_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"priority": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(agentpools.ScaleSetPriorityRegular),
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.ScaleSetPriorityRegular),
				string(agentpools.ScaleSetPrioritySpot),
			}, false),
		},

		"proximity_placement_group_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
		},

		"snapshot_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: snapshots.ValidateSnapshotID,
		},

		"spot_max_price": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			ForceNew:     true,
			Default:      -1.0,
			ValidateFunc: computeValidate.SpotMaxPrice,
		},

		"scale_down_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(agentpools.ScaleDownModeDelete),
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.ScaleDownModeDeallocate),
				string(agentpools.ScaleDownModeDelete),
			}, false),
		},

		"temporary_name_for_rotation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: containerValidate.KubernetesAgentPoolName,
		},

		"ultra_ssd_enabled": {
			Type:     pluginsdk.TypeBool,
			Default:  false,
			Optional: true,
		},

		"vnet_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"upgrade_settings": upgradeSettingsSchema(),

		"windows_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"outbound_nat_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  true,
					},
				},
			},
		},

		"workload_runtime": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.WorkloadRuntimeOCIContainer),
				string(agentpools.WorkloadRuntimeWasmWasi),
			}, false),
		},

		"zones": commonschema.ZonesMultipleOptional(),

		"auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"node_public_ip_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"host_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}

	return s
}

func resourceKubernetesClusterNodePoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	containersClient := meta.(*clients.Client).Containers
	clustersClient := containersClient.KubernetesClustersClient
	poolsClient := containersClient.AgentPoolsClient
	subnetClient := meta.(*clients.Client).Network.Client.Subnets
	vnetClient := meta.(*clients.Client).Network.VirtualNetworks

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId, err := commonids.ParseKubernetesClusterID(d.Get("kubernetes_cluster_id").(string))
	if err != nil {
		return err
	}

	var subnetID *commonids.SubnetId
	if subnetIDValue, ok := d.GetOk("vnet_subnet_id"); ok {
		subnetID, err = commonids.ParseSubnetID(subnetIDValue.(string))
		if err != nil {
			return err
		}
	}

	id := agentpools.NewAgentPoolID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.ManagedClusterName, d.Get("name").(string))

	log.Printf("[DEBUG] Retrieving %s...", *clusterId)
	cluster, err := clustersClient.Get(ctx, *clusterId)
	if err != nil {
		if response.WasNotFound(cluster.HttpResponse) {
			return fmt.Errorf("%s was not found", *clusterId)
		}

		return fmt.Errorf("retrieving %s: %+v", *clusterId, err)
	}

	// try to provide a more helpful error here
	defaultPoolIsVMSS := false
	if model := cluster.Model; model != nil && model.Properties != nil {
		props := model.Properties
		if pools := props.AgentPoolProfiles; pools != nil {
			for _, p := range *pools {
				if p.Type != nil && *p.Type == managedclusters.AgentPoolTypeVirtualMachineScaleSets {
					defaultPoolIsVMSS = true
					break
				}
			}
		}
	}
	if !defaultPoolIsVMSS {
		return fmt.Errorf("multiple node pools are only supported when the Default Node Pool uses a VMScaleSet (but %s doesn't)", *clusterId)
	}

	existing, err := poolsClient.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_kubernetes_cluster_node_pool", id.ID())
	}

	count := d.Get("node_count").(int)

	enableAutoScaling := d.Get("auto_scaling_enabled").(bool)
	hostEncryption := d.Get("host_encryption_enabled").(bool)
	nodeIp := d.Get("node_public_ip_enabled").(bool)

	evictionPolicy := d.Get("eviction_policy").(string)
	mode := agentpools.AgentPoolMode(d.Get("mode").(string))
	osType := d.Get("os_type").(string)
	priority := d.Get("priority").(string)
	spotMaxPrice := d.Get("spot_max_price").(float64)
	t := d.Get("tags").(map[string]interface{})

	profile := agentpools.ManagedClusterAgentPoolProfileProperties{
		OsType:                 pointer.To(agentpools.OSType(osType)),
		EnableAutoScaling:      pointer.To(enableAutoScaling),
		EnableFIPS:             pointer.To(d.Get("fips_enabled").(bool)),
		EnableEncryptionAtHost: pointer.To(hostEncryption),
		EnableUltraSSD:         pointer.To(d.Get("ultra_ssd_enabled").(bool)),
		EnableNodePublicIP:     pointer.To(nodeIp),
		KubeletDiskType:        pointer.To(agentpools.KubeletDiskType(d.Get("kubelet_disk_type").(string))),
		Mode:                   pointer.To(mode),
		ScaleSetPriority:       pointer.To(agentpools.ScaleSetPriority(d.Get("priority").(string))),
		Tags:                   tags.Expand(t),
		Type:                   pointer.To(agentpools.AgentPoolTypeVirtualMachineScaleSets),
		VMSize:                 pointer.To(d.Get("vm_size").(string)),
		UpgradeSettings:        expandAgentPoolUpgradeSettings(d.Get("upgrade_settings").([]interface{})),
		WindowsProfile:         expandAgentPoolWindowsProfile(d.Get("windows_profile").([]interface{})),

		// this must always be sent during creation, but is optional for auto-scaled clusters during update
		Count: utils.Int64(int64(count)),
	}

	if gpuInstanceProfile := d.Get("gpu_instance").(string); gpuInstanceProfile != "" {
		profile.GpuInstanceProfile = pointer.To(agentpools.GPUInstanceProfile(gpuInstanceProfile))
	}

	if osSku := d.Get("os_sku").(string); osSku != "" {
		profile.OsSKU = pointer.To(agentpools.OSSKU(osSku))
	}

	if scaleDownMode := d.Get("scale_down_mode").(string); scaleDownMode != "" {
		profile.ScaleDownMode = pointer.To(agentpools.ScaleDownMode(scaleDownMode))
	}

	if workloadRuntime := d.Get("workload_runtime").(string); workloadRuntime != "" {
		profile.WorkloadRuntime = pointer.To(agentpools.WorkloadRuntime(workloadRuntime))
	}

	if priority == string(managedclusters.ScaleSetPrioritySpot) {
		profile.ScaleSetEvictionPolicy = pointer.To(agentpools.ScaleSetEvictionPolicy(evictionPolicy))
		profile.SpotMaxPrice = utils.Float(spotMaxPrice)
	} else {
		if evictionPolicy != "" {
			return fmt.Errorf("`eviction_policy` can only be set when `priority` is set to `Spot`")
		}

		if spotMaxPrice != -1.0 {
			return fmt.Errorf("`spot_max_price` can only be set when `priority` is set to `Spot`")
		}
	}

	orchestratorVersion := d.Get("orchestrator_version").(string)
	if orchestratorVersion != "" {
		if err := validateNodePoolSupportsVersion(ctx, containersClient, "", id, orchestratorVersion); err != nil {
			return err
		}

		profile.OrchestratorVersion = utils.String(orchestratorVersion)
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		profile.AvailabilityZones = &zones
	}

	if maxPods := int64(d.Get("max_pods").(int)); maxPods > 0 {
		profile.MaxPods = utils.Int64(maxPods)
	}

	nodeLabelsRaw := d.Get("node_labels").(map[string]interface{})
	if nodeLabels := expandNodeLabels(nodeLabelsRaw); len(*nodeLabels) > 0 {
		profile.NodeLabels = nodeLabels
	}

	if nodePublicIPPrefixID := d.Get("node_public_ip_prefix_id").(string); nodePublicIPPrefixID != "" {
		profile.NodePublicIPPrefixID = utils.String(nodePublicIPPrefixID)
	}

	nodeTaintsRaw := d.Get("node_taints").([]interface{})
	if nodeTaints := utils.ExpandStringSlice(nodeTaintsRaw); len(*nodeTaints) > 0 {
		profile.NodeTaints = nodeTaints
	}

	if osDiskSizeGB := d.Get("os_disk_size_gb").(int); osDiskSizeGB > 0 {
		profile.OsDiskSizeGB = utils.Int64(int64(osDiskSizeGB))
	}

	proximityPlacementGroupId := d.Get("proximity_placement_group_id").(string)
	if proximityPlacementGroupId != "" {
		profile.ProximityPlacementGroupID = &proximityPlacementGroupId
	}

	if osDiskType := d.Get("os_disk_type").(string); osDiskType != "" {
		profile.OsDiskType = pointer.To(agentpools.OSDiskType(osDiskType))
	}

	if podSubnetID := d.Get("pod_subnet_id").(string); podSubnetID != "" {
		profile.PodSubnetID = utils.String(podSubnetID)
	}

	if subnetID != nil {
		profile.VnetSubnetID = utils.String(subnetID.ID())
	}

	if hostGroupID := d.Get("host_group_id").(string); hostGroupID != "" {
		profile.HostGroupID = utils.String(hostGroupID)
	}

	if capacityReservationGroupId := d.Get("capacity_reservation_group_id").(string); capacityReservationGroupId != "" {
		profile.CapacityReservationGroupID = utils.String(capacityReservationGroupId)
	}

	maxCount := d.Get("max_count").(int)
	minCount := d.Get("min_count").(int)

	if enableAutoScaling {
		// handle count being optional
		if count == 0 {
			profile.Count = utils.Int64(int64(minCount))
		}

		if maxCount >= 0 {
			profile.MaxCount = utils.Int64(int64(maxCount))
		} else {
			return fmt.Errorf("`max_count` must be configured when `auto_scaling_enabled` is set to `true`")
		}

		if minCount >= 0 {
			profile.MinCount = utils.Int64(int64(minCount))
		} else {
			return fmt.Errorf("`min_count` must be configured when `auto_scaling_enabled` is set to `true`")
		}

		if minCount > maxCount {
			return fmt.Errorf("`max_count` must be >= `min_count`")
		}
	} else if minCount > 0 || maxCount > 0 {
		return fmt.Errorf("`max_count` and `min_count` must be set to `null` when auto_scaling_enabled is set to `false`")
	}

	if kubeletConfig := d.Get("kubelet_config").([]interface{}); len(kubeletConfig) > 0 {
		profile.KubeletConfig = expandAgentPoolKubeletConfig(kubeletConfig)
	}

	if linuxOSConfig := d.Get("linux_os_config").([]interface{}); len(linuxOSConfig) > 0 {
		if osType != string(managedclusters.OSTypeLinux) {
			return fmt.Errorf("`linux_os_config` can only be configured when `os_type` is set to `linux`")
		}
		linuxOSConfig, err := expandAgentPoolLinuxOSConfig(linuxOSConfig)
		if err != nil {
			return err
		}
		profile.LinuxOSConfig = linuxOSConfig
	}

	if networkProfile := d.Get("node_network_profile").([]interface{}); len(networkProfile) > 0 {
		profile.NetworkProfile = expandAgentPoolNetworkProfile(networkProfile)
	}

	if snapshotId := d.Get("snapshot_id").(string); snapshotId != "" {
		profile.CreationData = &agentpools.CreationData{
			SourceResourceId: utils.String(snapshotId),
		}
	}

	parameters := agentpools.AgentPool{
		Name:       utils.String(id.AgentPoolName),
		Properties: &profile,
	}

	err = poolsClient.CreateOrUpdateThenPoll(ctx, id, parameters, agentpools.DefaultCreateOrUpdateOperationOptions())
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if subnetID != nil {
		// Wait for vnet to come back to Succeeded before releasing any locks
		timeout, ok := ctx.Deadline()
		if !ok {
			return fmt.Errorf("internal-error: context had no deadline")
		}

		// TODO: refactor this into a `custompoller` within the `network` package
		stateConf := &pluginsdk.StateChangeConf{
			Pending:    []string{string(subnets.ProvisioningStateUpdating)},
			Target:     []string{string(subnets.ProvisioningStateSucceeded)},
			Refresh:    network.SubnetProvisioningStateRefreshFunc(ctx, subnetClient, *subnetID),
			MinTimeout: 1 * time.Minute,
			Timeout:    time.Until(timeout),
		}
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for provisioning state of subnet for AKS Node Pool creation %s: %+v", *subnetID, err)
		}

		vnetId := commonids.NewVirtualNetworkID(subnetID.SubscriptionId, subnetID.ResourceGroupName, subnetID.VirtualNetworkName)
		vnetStateConf := &pluginsdk.StateChangeConf{
			Pending:    []string{string(subnets.ProvisioningStateUpdating)},
			Target:     []string{string(subnets.ProvisioningStateSucceeded)},
			Refresh:    network.VirtualNetworkProvisioningStateRefreshFunc(ctx, vnetClient, vnetId),
			MinTimeout: 1 * time.Minute,
			Timeout:    time.Until(timeout),
		}
		if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for provisioning state of virtual network for AKS Node Pool creation %s: %+v", vnetId, err)
		}
	}

	d.SetId(id.ID())
	return resourceKubernetesClusterNodePoolRead(d, meta)
}

func resourceKubernetesClusterNodePoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	containersClient := meta.(*clients.Client).Containers
	client := containersClient.AgentPoolsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agentpools.ParseAgentPoolID(d.Id())
	if err != nil {
		return err
	}

	d.Partial(true)

	log.Printf("[DEBUG] Retrieving existing %s..", *id)
	existing, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	props := existing.Model.Properties

	// store the existing value should the user have opted to ignore it
	enableAutoScaling := false
	if props.EnableAutoScaling != nil {
		enableAutoScaling = *props.EnableAutoScaling
	}

	log.Printf("[DEBUG] Determining delta for existing %s..", *id)

	// delta patching
	if d.HasChange("auto_scaling_enabled") {
		enableAutoScaling = d.Get("auto_scaling_enabled").(bool)
		props.EnableAutoScaling = utils.Bool(enableAutoScaling)
	}

	if d.HasChange("fips_enabled") {
		props.EnableFIPS = pointer.To(d.Get("fips_enabled").(bool))
	}

	if d.HasChange("host_encryption_enabled") {
		props.EnableEncryptionAtHost = pointer.To(d.Get("host_encryption_enabled").(bool))
	}

	if d.HasChange("kubelet_config") {
		kubeletConfigRaw := d.Get("kubelet_config").([]interface{})
		props.KubeletConfig = expandAgentPoolKubeletConfig(kubeletConfigRaw)
	}

	if d.HasChange("kubelet_disk_type") {
		props.KubeletDiskType = pointer.To(agentpools.KubeletDiskType(d.Get("kubelet_disk_type").(string)))
	}

	if d.HasChange("linux_os_config") {
		linuxOSConfigRaw := d.Get("linux_os_config").([]interface{})
		if d.Get("os_type").(string) != string(managedclusters.OSTypeLinux) {
			return fmt.Errorf("`linux_os_config` can only be configured when `os_type` is set to `linux`")
		}
		linuxOSConfig, err := expandAgentPoolLinuxOSConfig(linuxOSConfigRaw)
		if err != nil {
			return err
		}
		props.LinuxOSConfig = linuxOSConfig
	}

	if d.HasChange("max_count") || enableAutoScaling {
		props.MaxCount = utils.Int64(int64(d.Get("max_count").(int)))
	}

	if d.HasChange("max_pods") {
		props.MaxPods = pointer.To(int64(d.Get("max_pods").(int)))
	}

	if d.HasChange("mode") {
		mode := agentpools.AgentPoolMode(d.Get("mode").(string))
		props.Mode = &mode
	}

	if d.HasChange("min_count") || enableAutoScaling {
		props.MinCount = utils.Int64(int64(d.Get("min_count").(int)))
	}

	if d.HasChange("node_count") {
		props.Count = utils.Int64(int64(d.Get("node_count").(int)))
	}

	if d.HasChange("node_public_ip_enabled") {
		props.EnableNodePublicIP = pointer.To(d.Get("node_public_ip_enabled").(bool))
	}

	if d.HasChange("node_public_ip_prefix_id") {
		props.NodePublicIPPrefixID = pointer.To(d.Get("node_public_ip_prefix_id").(string))
	}

	if d.HasChange("orchestrator_version") {
		existingNodePoolResp, err := client.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving Node Pool %s: %+v", *id, err)
		}
		if existingNodePool := existingNodePoolResp.Model; existingNodePool != nil && existingNodePool.Properties != nil {
			orchestratorVersion := d.Get("orchestrator_version").(string)
			currentOrchestratorVersion := ""
			if v := existingNodePool.Properties.CurrentOrchestratorVersion; v != nil {
				currentOrchestratorVersion = *v
			}
			if err := validateNodePoolSupportsVersion(ctx, containersClient, currentOrchestratorVersion, *id, orchestratorVersion); err != nil {
				return err
			}

			props.OrchestratorVersion = utils.String(orchestratorVersion)
		}
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})
		props.Tags = tags.Expand(t)
	}

	if d.HasChange("os_disk_type") {
		props.OsDiskType = pointer.To(agentpools.OSDiskType(d.Get("os_disk_type").(string)))
	}

	if d.HasChange("os_disk_size_gb") {
		props.OsDiskSizeGB = pointer.To(int64(d.Get("os_disk_size_gb").(int)))
	}

	if d.HasChange("os_sku") {
		props.OsSKU = pointer.To(agentpools.OSSKU(d.Get("os_sku").(string)))
	}

	if d.HasChange("pod_subnet_id") {
		props.PodSubnetID = pointer.To(d.Get("pod_subnet_id").(string))
	}

	if d.HasChange("ultra_ssd_enabled") {
		props.EnableUltraSSD = pointer.To(d.Get("ultra_ssd_enabled").(bool))
	}

	if d.HasChange("upgrade_settings") {
		upgradeSettingsRaw := d.Get("upgrade_settings").([]interface{})
		props.UpgradeSettings = expandAgentPoolUpgradeSettings(upgradeSettingsRaw)
	}

	if d.HasChange("scale_down_mode") {
		mode := agentpools.ScaleDownMode(d.Get("scale_down_mode").(string))
		props.ScaleDownMode = &mode
	}

	if d.HasChange("snapshot_id") {
		props.CreationData = &agentpools.CreationData{
			SourceResourceId: pointer.To(d.Get("snapshot_id").(string)),
		}
	}

	if d.HasChange("vm_size") {
		props.VMSize = pointer.To(d.Get("vm_size").(string))
	}

	if d.HasChange("vnet_subnet_id") {
		if subnetIDValue, ok := d.GetOk("vnet_subnet_id"); ok {
			subnetID, err := commonids.ParseSubnetID(subnetIDValue.(string))
			if err != nil {
				return err
			}
			props.VnetSubnetID = pointer.To(subnetID.ID())
		}
	}

	if d.HasChange("workload_runtime") {
		runtime := agentpools.WorkloadRuntime(d.Get("workload_runtime").(string))
		props.WorkloadRuntime = &runtime
	}

	if d.HasChange("node_labels") {
		props.NodeLabels = expandNodeLabels(d.Get("node_labels").(map[string]interface{}))
	}

	if d.HasChange("node_taints") {
		props.NodeTaints = utils.ExpandStringSlice(d.Get("node_taints").([]interface{}))
	}

	if d.HasChange("node_network_profile") {
		props.NetworkProfile = expandAgentPoolNetworkProfile(d.Get("node_network_profile").([]interface{}))
	}

	if d.HasChange("zones") {
		zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
		props.AvailabilityZones = &zones
	}

	// validate the auto-scale fields are both set/unset to prevent a continual diff
	maxCount := 0
	if props.MaxCount != nil {
		maxCount = int(*props.MaxCount)
	}
	minCount := 0
	if props.MinCount != nil {
		minCount = int(*props.MinCount)
	}
	if enableAutoScaling {
		if maxCount == 0 {
			return fmt.Errorf("`max_count` must be configured when `auto_scaling_enabled` is set to `true`")
		}

		if minCount > maxCount {
			return fmt.Errorf("`max_count` must be >= `min_count`")
		}
	} else {
		if minCount > 0 || maxCount > 0 {
			return fmt.Errorf("`max_count` and `min_count` must be set to `nil` when `auto_scaling_enabled` is set to `false`")
		}

		// @tombuildsstuff: as of API version 2019-11-01 we need to explicitly nil these out
		props.MaxCount = nil
		props.MinCount = nil
	}

	// evaluate if the nodepool needs to be cycled
	cycleNodePoolProperties := []string{
		"fips_enabled",
		"host_encryption_enabled",
		"kubelet_config",
		"kubelet_disk_type",
		"linux_os_config",
		"max_pods",
		"node_public_ip_enabled",
		"os_disk_size_gb",
		"os_disk_type",
		"pod_subnet_id",
		"snapshot_id",
		"ultra_ssd_enabled",
		"vm_size",
		"vnet_subnet_id",
		"zones",
	}

	// if the node pool name has changed, it means the initial attempt at resizing failed
	cycleNodePool := d.HasChanges(cycleNodePoolProperties...)
	// os_sku can only be updated if the current and new os_sku are either Ubuntu or AzureLinux
	if d.HasChange("os_sku") {
		oldOsSkuRaw, newOsSkuRaw := d.GetChange("os_sku")
		oldOsSku := oldOsSkuRaw.(string)
		newOsSku := newOsSkuRaw.(string)
		if oldOsSku != string(managedclusters.OSSKUUbuntu) && oldOsSku != string(managedclusters.OSSKUAzureLinux) {
			cycleNodePool = true
		}
		if newOsSku != string(managedclusters.OSSKUUbuntu) && newOsSku != string(managedclusters.OSSKUAzureLinux) {
			cycleNodePool = true
		}
	}

	if cycleNodePool {
		log.Printf("[DEBUG] Cycling Node Pool..")
		// to provide a seamless updating experience for the node pool we need to cycle it by provisioning a temporary one,
		// tearing down the existing node pool and then bringing up the new one.

		if v := d.Get("temporary_name_for_rotation").(string); v == "" {
			return fmt.Errorf("`temporary_name_for_rotation` must be specified when updating any of the following properties %q", cycleNodePoolProperties)
		}

		temporaryNodePoolName := d.Get("temporary_name_for_rotation").(string)
		tempNodePoolId := agentpools.NewAgentPoolID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, temporaryNodePoolName)

		tempExisting, err := client.Get(ctx, tempNodePoolId)
		if !response.WasNotFound(tempExisting.HttpResponse) && err != nil {
			return fmt.Errorf("checking for existing temporary node pool %s: %+v", tempNodePoolId, err)
		}

		tempAgentProfile := *existing.Model
		tempAgentProfile.Name = &temporaryNodePoolName

		// if the temp node pool already exists due to a previous failure, don't bother spinning it up.
		// the temporary nodepool is created with the new values
		if tempExisting.Model == nil {
			if err := retryNodePoolCreation(ctx, client, tempNodePoolId, tempAgentProfile); err != nil {
				return fmt.Errorf("creating temporary %s: %+v", tempNodePoolId, err)
			}
		}

		// delete the old node pool if it exists
		if existing.Model != nil {
			if err := client.DeleteThenPoll(ctx, *id, agentpools.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting old %s: %+v", *id, err)
			}
		}

		// create the new node pool with the new data
		if err := retryNodePoolCreation(ctx, client, *id, *existing.Model); err != nil {
			log.Printf("[DEBUG] Creation of redefined node pool failed")
			return fmt.Errorf("creating default %s: %+v", *id, err)
		}

		if err := client.DeleteThenPoll(ctx, tempNodePoolId, agentpools.DefaultDeleteOperationOptions()); err != nil {
			return fmt.Errorf("deleting temporary %s: %+v", tempNodePoolId, err)
		}

		log.Printf("[DEBUG] Cycled Node Pool..")
	} else {
		log.Printf("[DEBUG] Updating existing %s..", *id)
		err = client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, agentpools.DefaultCreateOrUpdateOperationOptions())
		if err != nil {
			return fmt.Errorf("updating Node Pool %s: %+v", *id, err)
		}
	}

	d.Partial(false)

	return resourceKubernetesClusterNodePoolRead(d, meta)
}

func resourceKubernetesClusterNodePoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	poolsClient := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agentpools.ParseAgentPoolID(d.Id())
	if err != nil {
		return err
	}

	clusterId := commonids.NewKubernetesClusterID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName)

	resp, err := poolsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AgentPoolName)
	d.Set("kubernetes_cluster_id", clusterId.ID())

	if model := resp.Model; model != nil && model.Properties != nil {
		props := model.Properties
		d.Set("zones", zones.FlattenUntyped(props.AvailabilityZones))

		d.Set("auto_scaling_enabled", props.EnableAutoScaling)
		d.Set("node_public_ip_enabled", props.EnableNodePublicIP)
		d.Set("host_encryption_enabled", props.EnableEncryptionAtHost)
		d.Set("fips_enabled", props.EnableFIPS)
		d.Set("ultra_ssd_enabled", props.EnableUltraSSD)

		if v := props.KubeletDiskType; v != nil {
			d.Set("kubelet_disk_type", string(*v))
		}

		if v := props.GpuInstanceProfile; v != nil {
			d.Set("gpu_instance", string(*v))
		}

		if props.CreationData != nil {
			d.Set("snapshot_id", props.CreationData.SourceResourceId)
		}

		scaleDownMode := string(managedclusters.ScaleDownModeDelete)
		if v := props.ScaleDownMode; v != nil {
			scaleDownMode = string(*v)
		}
		d.Set("scale_down_mode", scaleDownMode)

		if v := props.WorkloadRuntime; v != nil {
			d.Set("workload_runtime", string(*v))
		}

		evictionPolicy := ""
		if v := props.ScaleSetEvictionPolicy; v != nil && *v != "" {
			evictionPolicy = string(*v)
		}
		d.Set("eviction_policy", evictionPolicy)

		if err := d.Set("kubelet_config", flattenAgentPoolKubeletConfig(props.KubeletConfig)); err != nil {
			return fmt.Errorf("setting `kubelet_config`: %+v", err)
		}

		linuxOSConfig, err := flattenAgentPoolLinuxOSConfig(props.LinuxOSConfig)
		if err != nil {
			return err
		}
		if err := d.Set("linux_os_config", linuxOSConfig); err != nil {
			return fmt.Errorf("setting `linux_os_config`: %+v", err)
		}

		maxCount := 0
		if props.MaxCount != nil {
			maxCount = int(*props.MaxCount)
		}
		d.Set("max_count", maxCount)

		maxPods := 0
		if props.MaxPods != nil {
			maxPods = int(*props.MaxPods)
		}
		d.Set("max_pods", maxPods)

		minCount := 0
		if props.MinCount != nil {
			minCount = int(*props.MinCount)
		}
		d.Set("min_count", minCount)

		mode := string(managedclusters.AgentPoolModeUser)
		if v := props.Mode; v != nil && *v != "" {
			mode = string(*props.Mode)
		}
		d.Set("mode", mode)

		count := 0
		if props.Count != nil {
			count = int(*props.Count)
		}
		d.Set("node_count", count)

		if err := d.Set("node_labels", props.NodeLabels); err != nil {
			return fmt.Errorf("setting `node_labels`: %+v", err)
		}

		d.Set("node_public_ip_prefix_id", props.NodePublicIPPrefixID)

		if err := d.Set("node_taints", utils.FlattenStringSlice(props.NodeTaints)); err != nil {
			return fmt.Errorf("setting `node_taints`: %+v", err)
		}

		// NOTE: workaround for migration from 2022-01-02-preview (<3.12.0) to 2022-03-02-preview (>=3.12.0). Before terraform apply is run against the new API, Azure will respond only with currentOrchestratorVersion, orchestratorVersion will be absent. More details: https://github.com/hashicorp/terraform-provider-azurerm/issues/17833#issuecomment-1227583353
		if props.OrchestratorVersion != nil {
			d.Set("orchestrator_version", props.OrchestratorVersion)
		} else {
			d.Set("orchestrator_version", props.CurrentOrchestratorVersion)
		}

		osDiskSizeGB := 0
		if props.OsDiskSizeGB != nil {
			osDiskSizeGB = int(*props.OsDiskSizeGB)
		}
		d.Set("os_disk_size_gb", osDiskSizeGB)

		osDiskType := agentpools.OSDiskTypeManaged
		if v := props.OsDiskType; v != nil && *v != "" {
			osDiskType = *v
		}
		d.Set("os_disk_type", osDiskType)

		if v := props.OsType; v != nil {
			d.Set("os_type", string(*v))
		}
		if v := props.OsSKU; v != nil {
			d.Set("os_sku", string(*v))
		}
		d.Set("pod_subnet_id", props.PodSubnetID)

		// not returned from the API if not Spot
		priority := string(managedclusters.ScaleSetPriorityRegular)
		if v := props.ScaleSetPriority; v != nil && *v != "" {
			priority = string(*v)
		}
		d.Set("priority", priority)

		d.Set("proximity_placement_group_id", props.ProximityPlacementGroupID)

		spotMaxPrice := -1.0
		if props.SpotMaxPrice != nil {
			spotMaxPrice = *props.SpotMaxPrice
		}
		d.Set("spot_max_price", spotMaxPrice)

		d.Set("vnet_subnet_id", props.VnetSubnetID)
		d.Set("vm_size", props.VMSize)
		d.Set("host_group_id", props.HostGroupID)
		d.Set("capacity_reservation_group_id", props.CapacityReservationGroupID)

		if err := d.Set("upgrade_settings", flattenAgentPoolUpgradeSettings(props.UpgradeSettings)); err != nil {
			return fmt.Errorf("setting `upgrade_settings`: %+v", err)
		}

		if err := d.Set("windows_profile", flattenAgentPoolWindowsProfile(props.WindowsProfile)); err != nil {
			return fmt.Errorf("setting `windows_profile`: %+v", err)
		}

		if err := d.Set("node_network_profile", flattenAgentPoolNetworkProfile(props.NetworkProfile)); err != nil {
			return fmt.Errorf("setting `node_network_profile`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Model.Properties.Tags)
}

func resourceKubernetesClusterNodePoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agentpools.ParseAgentPoolID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id, agentpools.DefaultDeleteOperationOptions())
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func upgradeSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_surge": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"drain_timeout_in_minutes": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},
				"node_soak_duration_in_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 30),
				},
			},
		},
	}
}

func upgradeSettingsForDataSourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_surge": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"drain_timeout_in_minutes": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
				"node_soak_duration_in_minutes": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func expandAgentPoolKubeletConfig(input []interface{}) *agentpools.KubeletConfig {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	result := &agentpools.KubeletConfig{
		CpuCfsQuota: utils.Bool(raw["cpu_cfs_quota_enabled"].(bool)),
		// must be false, otherwise the backend will report error: CustomKubeletConfig.FailSwapOn must be set to false to enable swap file on nodes.
		FailSwapOn:           utils.Bool(false),
		AllowedUnsafeSysctls: utils.ExpandStringSlice(raw["allowed_unsafe_sysctls"].(*pluginsdk.Set).List()),
	}

	if v := raw["cpu_manager_policy"].(string); v != "" {
		result.CpuManagerPolicy = utils.String(v)
	}
	if v := raw["cpu_cfs_quota_period"].(string); v != "" {
		result.CpuCfsQuotaPeriod = utils.String(v)
	}
	if v := raw["image_gc_high_threshold"].(int); v != 0 {
		result.ImageGcHighThreshold = utils.Int64(int64(v))
	}
	if v := raw["image_gc_low_threshold"].(int); v != 0 {
		result.ImageGcLowThreshold = utils.Int64(int64(v))
	}
	if v := raw["topology_manager_policy"].(string); v != "" {
		result.TopologyManagerPolicy = utils.String(v)
	}
	if v := raw["container_log_max_size_mb"].(int); v != 0 {
		result.ContainerLogMaxSizeMB = utils.Int64(int64(v))
	}
	if v := raw["container_log_max_line"].(int); v != 0 {
		result.ContainerLogMaxFiles = utils.Int64(int64(v))
	}
	if v := raw["pod_max_pid"].(int); v != 0 {
		result.PodMaxPids = utils.Int64(int64(v))
	}

	return result
}

func expandAgentPoolUpgradeSettings(input []interface{}) *agentpools.AgentPoolUpgradeSettings {
	setting := &agentpools.AgentPoolUpgradeSettings{}
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	if maxSurgeRaw := v["max_surge"].(string); maxSurgeRaw != "" {
		setting.MaxSurge = utils.String(maxSurgeRaw)
	}
	if drainTimeoutInMinutesRaw, ok := v["drain_timeout_in_minutes"].(int); ok {
		setting.DrainTimeoutInMinutes = pointer.To(int64(drainTimeoutInMinutesRaw))
	}
	if nodeSoakDurationInMinutesRaw, ok := v["node_soak_duration_in_minutes"].(int); ok {
		setting.NodeSoakDurationInMinutes = pointer.To(int64(nodeSoakDurationInMinutesRaw))
	}
	return setting
}

func flattenAgentPoolUpgradeSettings(input *agentpools.AgentPoolUpgradeSettings) []interface{} {
	// The API returns an empty upgrade settings object for spot node pools, so we need to explicitly check whether there's anything in it
	if input == nil || (input.MaxSurge == nil && input.DrainTimeoutInMinutes == nil && input.NodeSoakDurationInMinutes == nil) {
		return []interface{}{}
	}

	values := make(map[string]interface{})

	if input.MaxSurge != nil && *input.MaxSurge != "" {
		values["max_surge"] = *input.MaxSurge
	}

	if input.DrainTimeoutInMinutes != nil {
		values["drain_timeout_in_minutes"] = *input.DrainTimeoutInMinutes
	}

	if input.NodeSoakDurationInMinutes != nil {
		values["node_soak_duration_in_minutes"] = *input.NodeSoakDurationInMinutes
	}

	return []interface{}{values}
}

func expandNodeLabels(input map[string]interface{}) *map[string]string {
	result := make(map[string]string)
	for k, v := range input {
		result[k] = v.(string)
	}
	return &result
}

func expandAgentPoolLinuxOSConfig(input []interface{}) (*agentpools.LinuxOSConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	raw := input[0].(map[string]interface{})
	sysctlConfig, err := expandAgentPoolSysctlConfig(raw["sysctl_config"].([]interface{}))
	if err != nil {
		return nil, err
	}

	result := &agentpools.LinuxOSConfig{
		Sysctls: sysctlConfig,
	}
	if v := raw["transparent_huge_page_enabled"].(string); v != "" {
		result.TransparentHugePageEnabled = utils.String(v)
	}
	if v := raw["transparent_huge_page_defrag"].(string); v != "" {
		result.TransparentHugePageDefrag = utils.String(v)
	}
	if v := raw["swap_file_size_mb"].(int); v != 0 {
		result.SwapFileSizeMB = utils.Int64(int64(v))
	}
	return result, nil
}

func expandAgentPoolSysctlConfig(input []interface{}) (*agentpools.SysctlConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	raw := input[0].(map[string]interface{})
	result := &agentpools.SysctlConfig{
		NetIPv4TcpTwReuse: utils.Bool(raw["net_ipv4_tcp_tw_reuse"].(bool)),
	}
	if v := raw["net_core_somaxconn"].(int); v != 0 {
		result.NetCoreSomaxconn = utils.Int64(int64(v))
	}
	if v := raw["net_core_netdev_max_backlog"].(int); v != 0 {
		result.NetCoreNetdevMaxBacklog = utils.Int64(int64(v))
	}
	if v := raw["net_core_rmem_default"].(int); v != 0 {
		result.NetCoreRmemDefault = utils.Int64(int64(v))
	}
	if v := raw["net_core_rmem_max"].(int); v != 0 {
		result.NetCoreRmemMax = utils.Int64(int64(v))
	}
	if v := raw["net_core_wmem_default"].(int); v != 0 {
		result.NetCoreWmemDefault = utils.Int64(int64(v))
	}
	if v := raw["net_core_wmem_max"].(int); v != 0 {
		result.NetCoreWmemMax = utils.Int64(int64(v))
	}
	if v := raw["net_core_optmem_max"].(int); v != 0 {
		result.NetCoreOptmemMax = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_max_syn_backlog"].(int); v != 0 {
		result.NetIPv4TcpMaxSynBacklog = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_max_tw_buckets"].(int); v != 0 {
		result.NetIPv4TcpMaxTwBuckets = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_fin_timeout"].(int); v != 0 {
		result.NetIPv4TcpFinTimeout = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_time"].(int); v != 0 {
		result.NetIPv4TcpKeepaliveTime = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_probes"].(int); v != 0 {
		result.NetIPv4TcpKeepaliveProbes = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_intvl"].(int); v != 0 {
		result.NetIPv4TcpkeepaliveIntvl = utils.Int64(int64(v))
	}
	netIpv4IPLocalPortRangeMin := raw["net_ipv4_ip_local_port_range_min"].(int)
	netIpv4IPLocalPortRangeMax := raw["net_ipv4_ip_local_port_range_max"].(int)
	if (netIpv4IPLocalPortRangeMin != 0 && netIpv4IPLocalPortRangeMax == 0) || (netIpv4IPLocalPortRangeMin == 0 && netIpv4IPLocalPortRangeMax != 0) {
		return nil, fmt.Errorf("`net_ipv4_ip_local_port_range_min` and `net_ipv4_ip_local_port_range_max` should both be set or unset")
	}
	if netIpv4IPLocalPortRangeMin > netIpv4IPLocalPortRangeMax {
		return nil, fmt.Errorf("`net_ipv4_ip_local_port_range_min` should be no larger than `net_ipv4_ip_local_port_range_max`")
	}
	if netIpv4IPLocalPortRangeMin != 0 && netIpv4IPLocalPortRangeMax != 0 {
		result.NetIPv4IPLocalPortRange = utils.String(fmt.Sprintf("%d %d", netIpv4IPLocalPortRangeMin, netIpv4IPLocalPortRangeMax))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh1"].(int); v != 0 {
		result.NetIPv4NeighDefaultGcThresh1 = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh2"].(int); v != 0 {
		result.NetIPv4NeighDefaultGcThresh2 = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh3"].(int); v != 0 {
		result.NetIPv4NeighDefaultGcThresh3 = utils.Int64(int64(v))
	}
	if v := raw["net_netfilter_nf_conntrack_max"].(int); v != 0 {
		result.NetNetfilterNfConntrackMax = utils.Int64(int64(v))
	}
	if v := raw["net_netfilter_nf_conntrack_buckets"].(int); v != 0 {
		result.NetNetfilterNfConntrackBuckets = utils.Int64(int64(v))
	}
	if v := raw["fs_aio_max_nr"].(int); v != 0 {
		result.FsAioMaxNr = utils.Int64(int64(v))
	}
	if v := raw["fs_inotify_max_user_watches"].(int); v != 0 {
		result.FsInotifyMaxUserWatches = utils.Int64(int64(v))
	}
	if v := raw["fs_file_max"].(int); v != 0 {
		result.FsFileMax = utils.Int64(int64(v))
	}
	if v := raw["fs_nr_open"].(int); v != 0 {
		result.FsNrOpen = utils.Int64(int64(v))
	}
	if v := raw["kernel_threads_max"].(int); v != 0 {
		result.KernelThreadsMax = utils.Int64(int64(v))
	}
	if v := raw["vm_max_map_count"].(int); v != 0 {
		result.VMMaxMapCount = utils.Int64(int64(v))
	}
	if v := raw["vm_swappiness"].(int); v != 0 {
		result.VMSwappiness = utils.Int64(int64(v))
	}
	if v := raw["vm_vfs_cache_pressure"].(int); v != 0 {
		result.VMVfsCachePressure = utils.Int64(int64(v))
	}
	return result, nil
}

func flattenAgentPoolLinuxOSConfig(input *agentpools.LinuxOSConfig) ([]interface{}, error) {
	if input == nil {
		return make([]interface{}, 0), nil
	}

	var swapFileSizeMB int
	if input.SwapFileSizeMB != nil {
		swapFileSizeMB = int(*input.SwapFileSizeMB)
	}
	var transparentHugePageDefrag string
	if input.TransparentHugePageDefrag != nil {
		transparentHugePageDefrag = *input.TransparentHugePageDefrag
	}
	var transparentHugePageEnabled string
	if input.TransparentHugePageEnabled != nil {
		transparentHugePageEnabled = *input.TransparentHugePageEnabled
	}
	sysctlConfig, err := flattenAgentPoolSysctlConfig(input.Sysctls)
	if err != nil {
		return nil, err
	}
	return []interface{}{
		map[string]interface{}{
			"swap_file_size_mb":             swapFileSizeMB,
			"sysctl_config":                 sysctlConfig,
			"transparent_huge_page_defrag":  transparentHugePageDefrag,
			"transparent_huge_page_enabled": transparentHugePageEnabled,
		},
	}, nil
}

func flattenAgentPoolSysctlConfig(input *agentpools.SysctlConfig) ([]interface{}, error) {
	if input == nil {
		return make([]interface{}, 0), nil
	}

	var fsAioMaxNr int
	if input.FsAioMaxNr != nil {
		fsAioMaxNr = int(*input.FsAioMaxNr)
	}
	var fsFileMax int
	if input.FsFileMax != nil {
		fsFileMax = int(*input.FsFileMax)
	}
	var fsInotifyMaxUserWatches int
	if input.FsInotifyMaxUserWatches != nil {
		fsInotifyMaxUserWatches = int(*input.FsInotifyMaxUserWatches)
	}
	var fsNrOpen int
	if input.FsNrOpen != nil {
		fsNrOpen = int(*input.FsNrOpen)
	}
	var kernelThreadsMax int
	if input.KernelThreadsMax != nil {
		kernelThreadsMax = int(*input.KernelThreadsMax)
	}
	var netCoreNetdevMaxBacklog int
	if input.NetCoreNetdevMaxBacklog != nil {
		netCoreNetdevMaxBacklog = int(*input.NetCoreNetdevMaxBacklog)
	}
	var netCoreOptmemMax int
	if input.NetCoreOptmemMax != nil {
		netCoreOptmemMax = int(*input.NetCoreOptmemMax)
	}
	var netCoreRmemDefault int
	if input.NetCoreRmemDefault != nil {
		netCoreRmemDefault = int(*input.NetCoreRmemDefault)
	}
	var netCoreRmemMax int
	if input.NetCoreRmemMax != nil {
		netCoreRmemMax = int(*input.NetCoreRmemMax)
	}
	var netCoreSomaxconn int
	if input.NetCoreSomaxconn != nil {
		netCoreSomaxconn = int(*input.NetCoreSomaxconn)
	}
	var netCoreWmemDefault int
	if input.NetCoreWmemDefault != nil {
		netCoreWmemDefault = int(*input.NetCoreWmemDefault)
	}
	var netCoreWmemMax int
	if input.NetCoreWmemMax != nil {
		netCoreWmemMax = int(*input.NetCoreWmemMax)
	}
	var netIpv4IpLocalPortRangeMin, netIpv4IpLocalPortRangeMax int
	if input.NetIPv4IPLocalPortRange != nil {
		arr := regexp.MustCompile("[ \t]+").Split(*input.NetIPv4IPLocalPortRange, -1)
		if len(arr) != 2 {
			return nil, fmt.Errorf("parsing `NetIPv4IPLocalPortRange` %s", *input.NetIPv4IPLocalPortRange)
		}
		var err error
		netIpv4IpLocalPortRangeMin, err = strconv.Atoi(arr[0])
		if err != nil {
			return nil, err
		}
		netIpv4IpLocalPortRangeMax, err = strconv.Atoi(arr[1])
		if err != nil {
			return nil, err
		}
	}
	var netIpv4NeighDefaultGcThresh1 int
	if input.NetIPv4NeighDefaultGcThresh1 != nil {
		netIpv4NeighDefaultGcThresh1 = int(*input.NetIPv4NeighDefaultGcThresh1)
	}
	var netIpv4NeighDefaultGcThresh2 int
	if input.NetIPv4NeighDefaultGcThresh2 != nil {
		netIpv4NeighDefaultGcThresh2 = int(*input.NetIPv4NeighDefaultGcThresh2)
	}
	var netIpv4NeighDefaultGcThresh3 int
	if input.NetIPv4NeighDefaultGcThresh3 != nil {
		netIpv4NeighDefaultGcThresh3 = int(*input.NetIPv4NeighDefaultGcThresh3)
	}
	var netIpv4TcpFinTimeout int
	if input.NetIPv4TcpFinTimeout != nil {
		netIpv4TcpFinTimeout = int(*input.NetIPv4TcpFinTimeout)
	}
	var netIpv4TcpkeepaliveIntvl int
	if input.NetIPv4TcpkeepaliveIntvl != nil {
		netIpv4TcpkeepaliveIntvl = int(*input.NetIPv4TcpkeepaliveIntvl)
	}
	var netIpv4TcpKeepaliveProbes int
	if input.NetIPv4TcpKeepaliveProbes != nil {
		netIpv4TcpKeepaliveProbes = int(*input.NetIPv4TcpKeepaliveProbes)
	}
	var netIpv4TcpKeepaliveTime int
	if input.NetIPv4TcpKeepaliveTime != nil {
		netIpv4TcpKeepaliveTime = int(*input.NetIPv4TcpKeepaliveTime)
	}
	var netIpv4TcpMaxSynBacklog int
	if input.NetIPv4TcpMaxSynBacklog != nil {
		netIpv4TcpMaxSynBacklog = int(*input.NetIPv4TcpMaxSynBacklog)
	}
	var netIpv4TcpMaxTwBuckets int
	if input.NetIPv4TcpMaxTwBuckets != nil {
		netIpv4TcpMaxTwBuckets = int(*input.NetIPv4TcpMaxTwBuckets)
	}
	var netIpv4TcpTwReuse bool
	if input.NetIPv4TcpTwReuse != nil {
		netIpv4TcpTwReuse = *input.NetIPv4TcpTwReuse
	}
	var netNetfilterNfConntrackBuckets int
	if input.NetNetfilterNfConntrackBuckets != nil {
		netNetfilterNfConntrackBuckets = int(*input.NetNetfilterNfConntrackBuckets)
	}
	var netNetfilterNfConntrackMax int
	if input.NetNetfilterNfConntrackMax != nil {
		netNetfilterNfConntrackMax = int(*input.NetNetfilterNfConntrackMax)
	}
	var vmMaxMapCount int
	if input.VMMaxMapCount != nil {
		vmMaxMapCount = int(*input.VMMaxMapCount)
	}
	var vmSwappiness int
	if input.VMSwappiness != nil {
		vmSwappiness = int(*input.VMSwappiness)
	}
	var vmVfsCachePressure int
	if input.VMVfsCachePressure != nil {
		vmVfsCachePressure = int(*input.VMVfsCachePressure)
	}
	return []interface{}{
		map[string]interface{}{
			"fs_aio_max_nr":                      fsAioMaxNr,
			"fs_file_max":                        fsFileMax,
			"fs_inotify_max_user_watches":        fsInotifyMaxUserWatches,
			"fs_nr_open":                         fsNrOpen,
			"kernel_threads_max":                 kernelThreadsMax,
			"net_core_netdev_max_backlog":        netCoreNetdevMaxBacklog,
			"net_core_optmem_max":                netCoreOptmemMax,
			"net_core_rmem_default":              netCoreRmemDefault,
			"net_core_rmem_max":                  netCoreRmemMax,
			"net_core_somaxconn":                 netCoreSomaxconn,
			"net_core_wmem_default":              netCoreWmemDefault,
			"net_core_wmem_max":                  netCoreWmemMax,
			"net_ipv4_ip_local_port_range_min":   netIpv4IpLocalPortRangeMin,
			"net_ipv4_ip_local_port_range_max":   netIpv4IpLocalPortRangeMax,
			"net_ipv4_neigh_default_gc_thresh1":  netIpv4NeighDefaultGcThresh1,
			"net_ipv4_neigh_default_gc_thresh2":  netIpv4NeighDefaultGcThresh2,
			"net_ipv4_neigh_default_gc_thresh3":  netIpv4NeighDefaultGcThresh3,
			"net_ipv4_tcp_fin_timeout":           netIpv4TcpFinTimeout,
			"net_ipv4_tcp_keepalive_intvl":       netIpv4TcpkeepaliveIntvl,
			"net_ipv4_tcp_keepalive_probes":      netIpv4TcpKeepaliveProbes,
			"net_ipv4_tcp_keepalive_time":        netIpv4TcpKeepaliveTime,
			"net_ipv4_tcp_max_syn_backlog":       netIpv4TcpMaxSynBacklog,
			"net_ipv4_tcp_max_tw_buckets":        netIpv4TcpMaxTwBuckets,
			"net_ipv4_tcp_tw_reuse":              netIpv4TcpTwReuse,
			"net_netfilter_nf_conntrack_buckets": netNetfilterNfConntrackBuckets,
			"net_netfilter_nf_conntrack_max":     netNetfilterNfConntrackMax,
			"vm_max_map_count":                   vmMaxMapCount,
			"vm_swappiness":                      vmSwappiness,
			"vm_vfs_cache_pressure":              vmVfsCachePressure,
		},
	}, nil
}

func expandAgentPoolWindowsProfile(input []interface{}) *agentpools.AgentPoolWindowsProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	outboundNatEnabled := v["outbound_nat_enabled"].(bool)
	return &agentpools.AgentPoolWindowsProfile{
		DisableOutboundNat: utils.Bool(!outboundNatEnabled),
	}
}

func flattenAgentPoolWindowsProfile(input *agentpools.AgentPoolWindowsProfile) []interface{} {
	if input == nil || input.DisableOutboundNat == nil {
		return []interface{}{}
	}

	outboundNatEnabled := true
	if input.DisableOutboundNat != nil {
		outboundNatEnabled = !*input.DisableOutboundNat
	}

	return []interface{}{
		map[string]interface{}{
			"outbound_nat_enabled": outboundNatEnabled,
		},
	}
}

func expandAgentPoolNetworkProfile(input []interface{}) *agentpools.AgentPoolNetworkProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &agentpools.AgentPoolNetworkProfile{
		AllowedHostPorts:          expandAgentPoolNetworkProfileAllowedHostPorts(v["allowed_host_ports"].([]interface{})),
		ApplicationSecurityGroups: utils.ExpandStringSlice(v["application_security_group_ids"].([]interface{})),
		NodePublicIPTags:          expandAgentPoolNetworkProfileNodePublicIPTags(v["node_public_ip_tags"].(map[string]interface{})),
	}
}

func expandAgentPoolNetworkProfileAllowedHostPorts(input []interface{}) *[]agentpools.PortRange {
	if len(input) == 0 {
		return nil
	}
	out := make([]agentpools.PortRange, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		var portEnd, portStart int64
		var protocol agentpools.Protocol
		if raw["port_end"] != nil {
			portEnd = int64(raw["port_end"].(int))
		}
		if raw["port_start"] != nil {
			portStart = int64(raw["port_start"].(int))
		}
		if raw["protocol"] != nil {
			protocol = agentpools.Protocol(raw["protocol"].(string))
		}
		out = append(out, agentpools.PortRange{
			PortEnd:   pointer.To(portEnd),
			PortStart: pointer.To(portStart),
			Protocol:  pointer.To(protocol),
		})
	}
	return &out
}

func expandAgentPoolNetworkProfileNodePublicIPTags(input map[string]interface{}) *[]agentpools.IPTag {
	if len(input) == 0 {
		return nil
	}
	out := make([]agentpools.IPTag, 0)

	for key, val := range input {
		ipTag := agentpools.IPTag{
			IPTagType: utils.String(key),
			Tag:       utils.String(val.(string)),
		}
		out = append(out, ipTag)
	}
	return &out
}

func flattenAgentPoolNetworkProfile(input *agentpools.AgentPoolNetworkProfile) []interface{} {
	if input == nil || input.NodePublicIPTags == nil && input.AllowedHostPorts == nil && input.ApplicationSecurityGroups == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"allowed_host_ports":             flattenAgentPoolNetworkProfileAllowedHostPorts(input.AllowedHostPorts),
			"application_security_group_ids": utils.FlattenStringSlice(input.ApplicationSecurityGroups),
			"node_public_ip_tags":            flattenAgentPoolNetworkProfileNodePublicIPTags(input.NodePublicIPTags),
		},
	}
}

func flattenAgentPoolNetworkProfileAllowedHostPorts(input *[]agentpools.PortRange) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	out := make([]interface{}, 0)
	for _, portRange := range *input {
		out = append(out, map[string]interface{}{
			"port_end":   pointer.From(portRange.PortEnd),
			"port_start": pointer.From(portRange.PortStart),
			"protocol":   pointer.From(portRange.Protocol),
		})
	}
	return out
}

func flattenAgentPoolNetworkProfileNodePublicIPTags(input *[]agentpools.IPTag) map[string]interface{} {
	if input == nil {
		return map[string]interface{}{}
	}
	out := make(map[string]interface{})

	for _, tag := range *input {
		if tag.IPTagType != nil {
			out[*tag.IPTagType] = tag.Tag
		}
	}

	return out
}
