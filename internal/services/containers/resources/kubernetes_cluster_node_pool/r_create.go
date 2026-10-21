// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKubernetesClusterNodePoolCreate(d *pluginsdk.ResourceData, meta any) error {
	containersClient := meta.(*clients.Client).Containers
	clustersClient := containersClient.KubernetesClustersClient
	poolsClient := containersClient.AgentPoolsClient
	subnetClient := meta.(*clients.Client).Network.Subnets
	vnetClient := meta.(*clients.Client).Network.VirtualNetworks

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId, err := commonids.ParseKubernetesClusterID(d.Get("kubernetes_cluster_id").(string))
	if err != nil {
		return err
	}

	parseOptionalSubnetID := func(d *pluginsdk.ResourceData, key string) (*commonids.SubnetId, error) {
		if value, ok := d.GetOk(key); ok {
			return commonids.ParseSubnetID(value.(string))
		}
		return nil, nil
	}

	nodeSubnetID, err := parseOptionalSubnetID(d, "vnet_subnet_id")
	if err != nil {
		return err
	}

	podSubnetID, err := parseOptionalSubnetID(d, "pod_subnet_id")
	if err != nil {
		return err
	}

	id := agentpools.NewAgentPoolID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.ManagedClusterName, d.Get("name").(string))

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

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := poolsClient.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError(azureKubernetesClusterNodePoolResourceName, id.ID())
		}
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
	t := d.Get("tags").(map[string]any)

	profile := agentpools.ManagedClusterAgentPoolProfileProperties{
		OsType:                 pointer.ToEnum[agentpools.OSType](osType),
		EnableAutoScaling:      pointer.To(enableAutoScaling),
		EnableFIPS:             pointer.To(d.Get("fips_enabled").(bool)),
		EnableEncryptionAtHost: pointer.To(hostEncryption),
		EnableUltraSSD:         pointer.To(d.Get("ultra_ssd_enabled").(bool)),
		EnableNodePublicIP:     pointer.To(nodeIp),
		KubeletDiskType:        pointer.ToEnum[agentpools.KubeletDiskType](d.Get("kubelet_disk_type").(string)),
		Mode:                   pointer.To(mode),
		ScaleSetPriority:       pointer.ToEnum[agentpools.ScaleSetPriority](d.Get("priority").(string)),
		Tags:                   tags.Expand(t),
		Type:                   pointer.To(agentpools.AgentPoolTypeVirtualMachineScaleSets),
		VMSize:                 pointer.To(d.Get("vm_size").(string)),
		UpgradeSettings:        expandAgentPoolUpgradeSettings(d.Get("upgrade_settings").([]any)),
		WindowsProfile:         expandAgentPoolWindowsProfile(d.Get("windows_profile").([]any)),

		// this must always be sent during creation, but is optional for auto-scaled clusters during update
		Count: pointer.To(int64(count)),
	}

	if gpuInstanceProfile := d.Get("gpu_instance").(string); gpuInstanceProfile != "" {
		profile.GpuInstanceProfile = pointer.ToEnum[agentpools.GPUInstanceProfile](gpuInstanceProfile)
	}

	if gpuDriver := d.Get("gpu_driver").(string); gpuDriver != "" {
		profile.GpuProfile = &agentpools.GPUProfile{
			Driver: pointer.ToEnum[agentpools.GPUDriver](gpuDriver),
		}
	}

	if osSku := d.Get("os_sku").(string); osSku != "" {
		profile.OsSKU = pointer.ToEnum[agentpools.OSSKU](osSku)
	}

	if scaleDownMode := d.Get("scale_down_mode").(string); scaleDownMode != "" {
		profile.ScaleDownMode = pointer.ToEnum[agentpools.ScaleDownMode](scaleDownMode)
	}

	if workloadRuntime := d.Get("workload_runtime").(string); workloadRuntime != "" {
		profile.WorkloadRuntime = pointer.ToEnum[agentpools.WorkloadRuntime](workloadRuntime)
	}

	if priority == string(managedclusters.ScaleSetPrioritySpot) {
		profile.ScaleSetEvictionPolicy = pointer.ToEnum[agentpools.ScaleSetEvictionPolicy](evictionPolicy)
		profile.SpotMaxPrice = pointer.To(spotMaxPrice)
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
		if err := helpers.ValidateNodePoolSupportsVersion(ctx, containersClient, "", id, orchestratorVersion); err != nil {
			return err
		}

		profile.OrchestratorVersion = pointer.To(orchestratorVersion)
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		profile.AvailabilityZones = &zones
	}

	if maxPods := int64(d.Get("max_pods").(int)); maxPods > 0 {
		profile.MaxPods = pointer.To(maxPods)
	}

	if nodeLabels := helpers.ExpandNodeLabels(d.Get("node_labels").(map[string]any)); len(*nodeLabels) > 0 {
		profile.NodeLabels = nodeLabels
	}

	if nodePublicIPPrefixID := d.Get("node_public_ip_prefix_id").(string); nodePublicIPPrefixID != "" {
		profile.NodePublicIPPrefixID = pointer.To(nodePublicIPPrefixID)
	}

	if nodeTaints := pluginsdk.ExpandStringSlice(d.Get("node_taints").([]any)); len(*nodeTaints) > 0 {
		profile.NodeTaints = nodeTaints
	}

	if osDiskSizeGB := d.Get("os_disk_size_gb").(int); osDiskSizeGB > 0 {
		profile.OsDiskSizeGB = pointer.To(int64(osDiskSizeGB))
	}

	proximityPlacementGroupId := d.Get("proximity_placement_group_id").(string)
	if proximityPlacementGroupId != "" {
		profile.ProximityPlacementGroupID = &proximityPlacementGroupId
	}

	if osDiskType := d.Get("os_disk_type").(string); osDiskType != "" {
		profile.OsDiskType = pointer.ToEnum[agentpools.OSDiskType](osDiskType)
	}

	resourceIDsToLock := make([]string, 0)
	if podSubnetID != nil {
		// Lock the pod subnet and its vnet to avoid a race condition with AKS setting
		// vnet ownership across node pools that share the same virtual network.
		profile.PodSubnetID = pointer.To(podSubnetID.ID())
		resourceIDsToLock = append(resourceIDsToLock, commonids.NewVirtualNetworkID(podSubnetID.SubscriptionId, podSubnetID.ResourceGroupName, podSubnetID.VirtualNetworkName).ID())
		resourceIDsToLock = append(resourceIDsToLock, podSubnetID.ID())
	}

	if nodeSubnetID != nil {
		// Lock the node subnet and its vnet to avoid a race condition with AKS setting
		// vnet ownership across node pools that share the same virtual network.
		profile.VnetSubnetID = pointer.To(nodeSubnetID.ID())
		resourceIDsToLock = append(resourceIDsToLock, commonids.NewVirtualNetworkID(nodeSubnetID.SubscriptionId, nodeSubnetID.ResourceGroupName, nodeSubnetID.VirtualNetworkName).ID())
		resourceIDsToLock = append(resourceIDsToLock, nodeSubnetID.ID())
	}
	locks.MultipleByID(&resourceIDsToLock)
	defer locks.UnlockMultipleByID(&resourceIDsToLock)

	if hostGroupID := d.Get("host_group_id").(string); hostGroupID != "" {
		profile.HostGroupID = pointer.To(hostGroupID)
	}

	if capacityReservationGroupId := d.Get("capacity_reservation_group_id").(string); capacityReservationGroupId != "" {
		profile.CapacityReservationGroupID = pointer.To(capacityReservationGroupId)
	}

	maxCount := d.Get("max_count").(int)
	minCount := d.Get("min_count").(int)

	if enableAutoScaling {
		// handle count being optional
		if count == 0 {
			profile.Count = pointer.To(int64(minCount))
		}

		if maxCount >= 0 {
			profile.MaxCount = pointer.To(int64(maxCount))
		} else {
			return fmt.Errorf("`max_count` must be configured when `auto_scaling_enabled` is set to `true`")
		}

		if minCount >= 0 {
			profile.MinCount = pointer.To(int64(minCount))
		} else {
			return fmt.Errorf("`min_count` must be configured when `auto_scaling_enabled` is set to `true`")
		}

		if minCount > maxCount {
			return fmt.Errorf("`max_count` must be >= `min_count`")
		}
	} else if minCount > 0 || maxCount > 0 {
		return fmt.Errorf("`max_count` and `min_count` must be set to `null` when auto_scaling_enabled is set to `false`")
	}

	if kubeletConfig := d.Get("kubelet_config").([]any); len(kubeletConfig) > 0 {
		profile.KubeletConfig = expandAgentPoolKubeletConfig(kubeletConfig)
	}

	if linuxOSConfig := d.Get("linux_os_config").([]any); len(linuxOSConfig) > 0 {
		if osType != string(managedclusters.OSTypeLinux) {
			return fmt.Errorf("`linux_os_config` can only be configured when `os_type` is set to `linux`")
		}
		linuxOSConfig, err := expandAgentPoolLinuxOSConfig(linuxOSConfig)
		if err != nil {
			return err
		}
		profile.LinuxOSConfig = linuxOSConfig
	}

	if networkProfile := d.Get("node_network_profile").([]any); len(networkProfile) > 0 {
		profile.NetworkProfile = expandAgentPoolNetworkProfile(networkProfile)
	}

	if snapshotId := d.Get("snapshot_id").(string); snapshotId != "" {
		profile.CreationData = &agentpools.CreationData{
			SourceResourceId: pointer.To(snapshotId),
		}
	}
	parameters := agentpools.AgentPool{
		Name:       pointer.To(id.AgentPoolName),
		Properties: &profile,
	}

	if err := poolsClient.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, agentpools.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	// Wait for vnet and node subnet to come back to Succeeded before releasing any locks
	timeout, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	if nodeSubnetID != nil {
		// Wait for vnet and node subnet to come back to Succeeded before releasing any locks
		if err = network.NewSubnetAndVnetPoller(subnetClient, vnetClient, nodeSubnetID, timeout).Poll(ctx); err != nil {
			return fmt.Errorf("waiting for provisioning state of subnet for AKS Node Pool creation %s: %+v", *nodeSubnetID, err)
		}
	}

	if podSubnetID != nil {
		// Wait for vnet and pod subnet to come back to Succeeded before releasing any locks
		if err = network.NewSubnetAndVnetPoller(subnetClient, vnetClient, podSubnetID, timeout).Poll(ctx); err != nil {
			return fmt.Errorf("waiting for provisioning state of the pod subnet for AKS Node Pool creation %s: %+v", *podSubnetID, err)
		}
	}

	return resourceKubernetesClusterNodePoolRead(d, meta)
}
