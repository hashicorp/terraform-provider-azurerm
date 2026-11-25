// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Kubernetes Cluster Node Pool.",
		Read:        dataSourceKubernetesClusterNodePoolRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of this Kubernetes Cluster Node Pool.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.KubernetesAgentPoolName,
			},

			"kubernetes_cluster_name": {
				Description:  "The Name of the Kubernetes Cluster where this Node Pool is located.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"auto_scaling_enabled": {
				Description: "Whether this Node Pool has Auto-Scaling enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"eviction_policy": {
				Description: "The eviction policy used for Virtual Machines in the Virtual Machine Scale Set, when `priority` is set to `Spot`.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"gpu_driver": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"max_count": {
				Description: "The maximum number of Nodes allowed when auto-scaling is enabled.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"max_pods": {
				Description: "The maximum number of Pods allowed on each Node in this Node Pool.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"mode": {
				Description: "The Mode for this Node Pool, specifying how these Nodes should be used (for either System or User resources).",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"min_count": {
				Description: "The minimum number of Nodes allowed when auto-scaling is enabled.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"node_count": {
				Description: "The current number of Nodes in the Node Pool.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"node_labels": {
				Description: "A map of Kubernetes Labels applied to each Node in this Node Pool.",
				Type:        pluginsdk.TypeMap,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"node_public_ip_enabled": {
				Description: "Whether nodes in this Node Pool have a Public IP Address.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"node_public_ip_prefix_id": {
				Description: "Resource ID for the Public IP Addresses Prefix for the nodes in this Agent Pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"node_taints": {
				Description: "A map of Kubernetes Taints applied to each Node in this Node Pool.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"node_image_version": {
				Description: "The current node image version running on this Node Pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"orchestrator_version": {
				Description: "The version of Kubernetes configured on each Node in this Node Pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"os_disk_size_gb": {
				Description: "The size of the OS Disk on each Node in this Node Pool.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"os_disk_type": {
				Description: "The type of the OS Disk on each Node in this Node Pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"os_type": {
				Description: "The operating system used on each Node in this Node Pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"priority": {
				Description: "The priority of the Virtual Machines in the Virtual Machine Scale Set backing this Node Pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"proximity_placement_group_id": {
				Description: "The ID of the Proximity Placement Group where the Virtual Machine Scale Set backing this Node Pool will be placed.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"spot_max_price": {
				Description: "The maximum price being paid for Virtual Machines in this Scale Set. `-1` means the current on-demand price for a Virtual Machine.",
				Type:        pluginsdk.TypeFloat,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),

			"upgrade_settings": helpers.UpgradeSettingsForDataSourceSchema(),

			"vm_size": {
				Description: "The size of the Virtual Machines used in the Virtual Machine Scale Set backing this Node Pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"vnet_subnet_id": {
				Description: "The ID of the Subnet in which this Node Pool exists.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"zones": commonschema.ZonesMultipleComputed(),
		},
	}
}

func dataSourceKubernetesClusterNodePoolRead(d *pluginsdk.ResourceData, meta any) error {
	clustersClient := meta.(*clients.Client).Containers.KubernetesClustersClient
	poolsClient := meta.(*clients.Client).Containers.AgentPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId := commonids.NewKubernetesClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("kubernetes_cluster_name").(string))

	// if the parent cluster doesn't exist then the node pool won't
	cluster, err := clustersClient.Get(ctx, clusterId)
	if err != nil {
		if response.WasNotFound(cluster.HttpResponse) {
			return fmt.Errorf("%s was not found", clusterId)
		}

		return fmt.Errorf("retrieving %s: %+v", clusterId, err)
	}

	id := agentpools.NewAgentPoolID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.ManagedClusterName, d.Get("name").(string))
	resp, err := poolsClient.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.AgentPoolName)
	d.Set("kubernetes_cluster_name", id.ManagedClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil && model.Properties != nil {
		props := model.Properties
		d.Set("zones", zones.FlattenUntyped(props.AvailabilityZones))

		d.Set("auto_scaling_enabled", props.EnableAutoScaling)
		d.Set("node_public_ip_enabled", props.EnableNodePublicIP)

		evictionPolicy := ""
		if props.ScaleSetEvictionPolicy != nil && *props.ScaleSetEvictionPolicy != "" {
			evictionPolicy = string(*props.ScaleSetEvictionPolicy)
		}
		d.Set("eviction_policy", evictionPolicy)

		gpuDriver := ""
		if props.GpuProfile != nil {
			gpuDriver = pointer.FromEnum(props.GpuProfile.Driver)
		}
		d.Set("gpu_driver", gpuDriver)

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

		mode := string(agentpools.AgentPoolModeUser)
		if props.Mode != nil && *props.Mode != "" {
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

		if err := d.Set("node_taints", pluginsdk.FlattenSlice(props.NodeTaints)); err != nil {
			return fmt.Errorf("setting `node_taints`: %+v", err)
		}

		d.Set("node_image_version", props.NodeImageVersion)

		d.Set("orchestrator_version", props.OrchestratorVersion)
		osDiskSizeGB := 0
		if props.OsDiskSizeGB != nil {
			osDiskSizeGB = int(*props.OsDiskSizeGB)
		}
		d.Set("os_disk_size_gb", osDiskSizeGB)

		osDiskType := agentpools.OSDiskTypeManaged
		if props.OsDiskType != nil && *props.OsDiskType != "" {
			osDiskType = *props.OsDiskType
		}
		d.Set("os_disk_type", string(osDiskType))

		if props.OsDiskType != nil {
			d.Set("os_type", string(*props.OsType))
		}

		// not returned from the API if not Spot
		priority := string(agentpools.ScaleSetPriorityRegular)
		if props.ScaleSetPriority != nil && *props.ScaleSetPriority != "" {
			priority = string(*props.ScaleSetPriority)
		}
		d.Set("priority", priority)

		d.Set("proximity_placement_group_id", pointer.From(props.ProximityPlacementGroupID))

		spotMaxPrice := -1.0
		if props.SpotMaxPrice != nil {
			spotMaxPrice = *props.SpotMaxPrice
		}
		d.Set("spot_max_price", spotMaxPrice)

		if err := d.Set("upgrade_settings", flattenAgentPoolUpgradeSettings(props.UpgradeSettings)); err != nil {
			return fmt.Errorf("setting `upgrade_settings`: %+v", err)
		}

		d.Set("vnet_subnet_id", props.VnetSubnetID)
		d.Set("vm_size", props.VMSize)

		d.Set("tags", tags.Flatten(props.Tags))
	}

	return nil
}
