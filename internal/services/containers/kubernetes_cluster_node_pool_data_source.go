package containers

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerservice/mgmt/2022-03-02-preview/containerservice"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceKubernetesClusterNodePool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKubernetesClusterNodePoolRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.KubernetesAgentPoolName,
			},

			"kubernetes_cluster_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_auto_scaling": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_node_public_ip": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"eviction_policy": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"max_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"max_pods": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"mode": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"min_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"node_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"node_labels": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"node_public_ip_prefix_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"node_taints": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"orchestrator_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"os_disk_size_gb": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"os_disk_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"os_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"priority": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"proximity_placement_group_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"spot_max_price": {
				Type:     pluginsdk.TypeFloat,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),

			"upgrade_settings": upgradeSettingsForDataSourceSchema(),

			"vm_size": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"vnet_subnet_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"zones": commonschema.ZonesMultipleComputed(),
		},
	}
}

func dataSourceKubernetesClusterNodePoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).Containers.KubernetesClustersClient
	poolsClient := meta.(*clients.Client).Containers.AgentPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId := parse.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("kubernetes_cluster_name").(string))

	// if the parent cluster doesn't exist then the node pool won't
	cluster, err := clustersClient.Get(ctx, clusterId.ResourceGroup, clusterId.ManagedClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			return fmt.Errorf("%s was not found", clusterId)
		}

		return fmt.Errorf("retrieving %s: %+v", clusterId, err)
	}

	id := parse.NewNodePoolID(clusterId.SubscriptionId, clusterId.ResourceGroup, clusterId.ManagedClusterName, d.Get("name").(string))
	resp, err := poolsClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.AgentPoolName)
	d.Set("kubernetes_cluster_name", id.ManagedClusterName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.ManagedClusterAgentPoolProfileProperties; props != nil {
		d.Set("zones", zones.Flatten(props.AvailabilityZones))

		d.Set("enable_auto_scaling", props.EnableAutoScaling)
		d.Set("enable_node_public_ip", props.EnableNodePublicIP)

		evictionPolicy := ""
		if props.ScaleSetEvictionPolicy != "" {
			evictionPolicy = string(props.ScaleSetEvictionPolicy)
		}
		d.Set("eviction_policy", evictionPolicy)

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

		mode := string(containerservice.AgentPoolModeUser)
		if props.Mode != "" {
			mode = string(props.Mode)
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

		d.Set("orchestrator_version", props.OrchestratorVersion)
		osDiskSizeGB := 0
		if props.OsDiskSizeGB != nil {
			osDiskSizeGB = int(*props.OsDiskSizeGB)
		}
		d.Set("os_disk_size_gb", osDiskSizeGB)

		osDiskType := containerservice.OSDiskTypeManaged
		if props.OsDiskType != "" {
			osDiskType = props.OsDiskType
		}
		d.Set("os_disk_type", string(osDiskType))
		d.Set("os_type", string(props.OsType))

		// not returned from the API if not Spot
		priority := string(containerservice.ScaleSetPriorityRegular)
		if props.ScaleSetPriority != "" {
			priority = string(props.ScaleSetPriority)
		}
		d.Set("priority", priority)

		proximityPlacementGroupId := ""
		if props.ProximityPlacementGroupID != nil {
			proximityPlacementGroupId = *props.ProximityPlacementGroupID
		}
		d.Set("proximity_placement_group_id", proximityPlacementGroupId)

		spotMaxPrice := -1.0
		if props.SpotMaxPrice != nil {
			spotMaxPrice = *props.SpotMaxPrice
		}
		d.Set("spot_max_price", spotMaxPrice)

		if err := d.Set("upgrade_settings", flattenUpgradeSettings(props.UpgradeSettings)); err != nil {
			return fmt.Errorf("setting `upgrade_settings`: %+v", err)
		}

		d.Set("vnet_subnet_id", props.VnetSubnetID)
		d.Set("vm_size", props.VMSize)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
