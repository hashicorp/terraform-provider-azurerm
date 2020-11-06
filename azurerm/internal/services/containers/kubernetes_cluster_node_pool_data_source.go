package containers

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-09-01/containerservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceKubernetesClusterNodePool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKubernetesClusterNodePoolRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.KubernetesAgentPoolName,
			},

			"kubernetes_cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			// Computed
			"availability_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"enable_auto_scaling": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enable_node_public_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"eviction_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"max_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_pods": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"min_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"node_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"node_taints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"orchestrator_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"os_disk_size_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"os_disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"priority": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"proximity_placement_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"spot_max_price": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),

			"vm_size": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vnet_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKubernetesClusterNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).Containers.KubernetesClustersClient
	poolsClient := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	nodePoolName := d.Get("name").(string)
	clusterName := d.Get("kubernetes_cluster_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	// if the parent cluster doesn't exist then the node pool won't
	cluster, err := clustersClient.Get(ctx, resourceGroup, clusterName)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			return fmt.Errorf("Kubernetes Cluster %q was not found in Resource Group %q", clusterName, resourceGroup)
		}

		return fmt.Errorf("retrieving Managed Kubernetes Cluster %q (Resource Group %q): %+v", clusterName, resourceGroup, err)
	}

	resp, err := poolsClient.Get(ctx, resourceGroup, clusterName, nodePoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Node Pool %q was not found in Managed Kubernetes Cluster %q / Resource Group %q", nodePoolName, clusterName, resourceGroup)
		}

		return fmt.Errorf("retrieving Node Pool %q (Managed Kubernetes Cluster %q / Resource Group %q): %+v", nodePoolName, clusterName, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("retrieving Node Pool %q (Managed Kubernetes Cluster %q / Resource Group %q): `id` was nil", nodePoolName, clusterName, resourceGroup)
	}

	d.SetId(*resp.ID)
	d.Set("name", nodePoolName)
	d.Set("kubernetes_cluster_name", clusterName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.ManagedClusterAgentPoolProfileProperties; props != nil {
		if err := d.Set("availability_zones", utils.FlattenStringSlice(props.AvailabilityZones)); err != nil {
			return fmt.Errorf("setting `availability_zones`: %+v", err)
		}

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

		mode := string(containerservice.User)
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

		if err := d.Set("node_taints", utils.FlattenStringSlice(props.NodeTaints)); err != nil {
			return fmt.Errorf("setting `node_taints`: %+v", err)
		}

		d.Set("orchestrator_version", props.OrchestratorVersion)
		osDiskSizeGB := 0
		if props.OsDiskSizeGB != nil {
			osDiskSizeGB = int(*props.OsDiskSizeGB)
		}
		d.Set("os_disk_size_gb", osDiskSizeGB)
		d.Set("os_disk_type", string(props.OsDiskType))
		d.Set("os_type", string(props.OsType))

		// not returned from the API if not Spot
		priority := string(containerservice.Regular)
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

		d.Set("vnet_subnet_id", props.VnetSubnetID)
		d.Set("vm_size", string(props.VMSize))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
