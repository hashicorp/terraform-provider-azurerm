// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/managedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKubernetesClusterNodePoolUpdate(d *pluginsdk.ResourceData, meta any) error {
	containersClient := meta.(*clients.Client).Containers
	client := containersClient.AgentPoolsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agentpools.ParseAgentPoolID(d.Id())
	if err != nil {
		return err
	}

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
	enableAutoScaling := pointer.From(props.EnableAutoScaling)

	log.Printf("[DEBUG] Determining delta for existing %s..", *id)

	// delta patching
	if d.HasChange("auto_scaling_enabled") {
		enableAutoScaling = d.Get("auto_scaling_enabled").(bool)
		props.EnableAutoScaling = pointer.To(enableAutoScaling)
	}

	if d.HasChange("fips_enabled") {
		props.EnableFIPS = pointer.To(d.Get("fips_enabled").(bool))
	}

	if d.HasChange("host_encryption_enabled") {
		props.EnableEncryptionAtHost = pointer.To(d.Get("host_encryption_enabled").(bool))
	}

	if d.HasChange("kubelet_config") {
		kubeletConfigRaw := d.Get("kubelet_config").([]any)
		props.KubeletConfig = expandAgentPoolKubeletConfig(kubeletConfigRaw)
	}

	if d.HasChange("kubelet_disk_type") {
		props.KubeletDiskType = pointer.ToEnum[agentpools.KubeletDiskType](d.Get("kubelet_disk_type").(string))
	}

	if d.HasChange("linux_os_config") {
		linuxOSConfigRaw := d.Get("linux_os_config").([]any)
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
		props.MaxCount = pointer.To(int64(d.Get("max_count").(int)))
	}

	if d.HasChange("max_pods") {
		props.MaxPods = pointer.To(int64(d.Get("max_pods").(int)))
	}

	if d.HasChange("mode") {
		mode := agentpools.AgentPoolMode(d.Get("mode").(string))
		props.Mode = &mode
	}

	if d.HasChange("min_count") || enableAutoScaling {
		props.MinCount = pointer.To(int64(d.Get("min_count").(int)))
	}

	if d.HasChange("node_count") {
		props.Count = pointer.To(int64(d.Get("node_count").(int)))
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
			currentOrchestratorVersion := pointer.From(existingNodePool.Properties.CurrentOrchestratorVersion)
			if err := helpers.ValidateNodePoolSupportsVersion(ctx, containersClient, currentOrchestratorVersion, *id, orchestratorVersion); err != nil {
				return err
			}

			props.OrchestratorVersion = pointer.To(orchestratorVersion)
		}
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]any)
		props.Tags = tags.Expand(t)
	}

	if d.HasChange("os_disk_type") {
		props.OsDiskType = pointer.ToEnum[agentpools.OSDiskType](d.Get("os_disk_type").(string))
	}

	if d.HasChange("os_disk_size_gb") {
		props.OsDiskSizeGB = pointer.To(int64(d.Get("os_disk_size_gb").(int)))
	}

	if d.HasChange("os_sku") {
		props.OsSKU = pointer.ToEnum[agentpools.OSSKU](d.Get("os_sku").(string))
	}

	if d.HasChange("pod_subnet_id") {
		props.PodSubnetID = pointer.To(d.Get("pod_subnet_id").(string))
	}

	if d.HasChange("ultra_ssd_enabled") {
		props.EnableUltraSSD = pointer.To(d.Get("ultra_ssd_enabled").(bool))
	}

	if d.HasChange("upgrade_settings") {
		upgradeSettingsRaw := d.Get("upgrade_settings").([]any)
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
		props.NodeLabels = helpers.ExpandNodeLabels(d.Get("node_labels").(map[string]any))
	}

	if d.HasChange("node_taints") {
		props.NodeTaints = pluginsdk.ExpandStringSlice(d.Get("node_taints").([]any))
	}

	if d.HasChange("node_network_profile") {
		props.NetworkProfile = expandAgentPoolNetworkProfile(d.Get("node_network_profile").([]any))
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
		if !strings.HasPrefix(oldOsSku, string(managedclusters.OSSKUUbuntu)) && !strings.HasPrefix(oldOsSku, string(managedclusters.OSSKUAzureLinux)) {
			cycleNodePool = true
		}
		if !strings.HasPrefix(newOsSku, string(managedclusters.OSSKUUbuntu)) && !strings.HasPrefix(newOsSku, string(managedclusters.OSSKUAzureLinux)) {
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

		if tempAgentProfile.Properties != nil {
			tempAgentProfile.Properties.NodeImageVersion = nil
		}
		if existing.Model != nil && existing.Model.Properties != nil {
			existing.Model.Properties.NodeImageVersion = nil
		}

		// if the temp node pool already exists due to a previous failure, don't bother spinning it up.
		// the temporary nodepool is created with the new values
		if tempExisting.Model == nil {
			if err := helpers.RetryNodePoolCreation(ctx, client, tempNodePoolId, tempAgentProfile); err != nil {
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
		if err := helpers.RetryNodePoolCreation(ctx, client, *id, *existing.Model); err != nil {
			log.Printf("[DEBUG] Creation of redefined node pool failed")
			return fmt.Errorf("creating default %s: %+v", *id, err)
		}

		if err := client.DeleteThenPoll(ctx, tempNodePoolId, agentpools.DefaultDeleteOperationOptions()); err != nil {
			return fmt.Errorf("deleting temporary %s: %+v", tempNodePoolId, err)
		}

		log.Printf("[DEBUG] Cycled Node Pool..")
	} else {
		if err = client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, agentpools.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating Node Pool %s: %+v", *id, err)
		}
	}

	return resourceKubernetesClusterNodePoolRead(d, meta)
}
