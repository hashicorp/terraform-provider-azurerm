package containers

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-10-01/containerservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SchemaDefaultNodePool() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// Required
				"name": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validate.KubernetesAgentPoolName,
				},

				"type": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Default:  string(containerservice.VirtualMachineScaleSets),
					ValidateFunc: validation.StringInSlice([]string{
						string(containerservice.AvailabilitySet),
						string(containerservice.VirtualMachineScaleSets),
					}, false),
				},

				"vm_size": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validate.NoEmptyStrings,
				},

				// Optional
				"availability_zones": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"enable_auto_scaling": {
					Type:     schema.TypeBool,
					Optional: true,
				},

				"enable_node_public_ip": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"max_count": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 100),
				},

				"max_pods": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},

				"min_count": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 100),
				},

				"node_count": {
					Type:         schema.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 100),
				},

				"node_taints": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},

				"os_disk_size_gb": {
					Type:         schema.TypeInt,
					Optional:     true,
					ForceNew:     true,
					Computed:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},

				"vnet_subnet_id": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
			},
		},
	}
}

func ConvertDefaultNodePoolToAgentPool(input *[]containerservice.ManagedClusterAgentPoolProfile) containerservice.AgentPool {
	defaultCluster := (*input)[0]
	return containerservice.AgentPool{
		Name: defaultCluster.Name,
		ManagedClusterAgentPoolProfileProperties: &containerservice.ManagedClusterAgentPoolProfileProperties{
			Count:                  defaultCluster.Count,
			VMSize:                 defaultCluster.VMSize,
			OsDiskSizeGB:           defaultCluster.OsDiskSizeGB,
			VnetSubnetID:           defaultCluster.VnetSubnetID,
			MaxPods:                defaultCluster.MaxPods,
			OsType:                 defaultCluster.OsType,
			MaxCount:               defaultCluster.MaxCount,
			MinCount:               defaultCluster.MinCount,
			EnableAutoScaling:      defaultCluster.EnableAutoScaling,
			Type:                   defaultCluster.Type,
			OrchestratorVersion:    defaultCluster.OrchestratorVersion,
			AvailabilityZones:      defaultCluster.AvailabilityZones,
			EnableNodePublicIP:     defaultCluster.EnableNodePublicIP,
			ScaleSetPriority:       defaultCluster.ScaleSetPriority,
			ScaleSetEvictionPolicy: defaultCluster.ScaleSetEvictionPolicy,
			NodeTaints:             defaultCluster.NodeTaints,
		},
	}
}

func ExpandDefaultNodePool(d *schema.ResourceData) (*[]containerservice.ManagedClusterAgentPoolProfile, error) {
	input := d.Get("default_node_pool").([]interface{})
	// TODO: in 2.0 make this Required
	// this exists to allow users to migrate to default_node_pool
	if len(input) == 0 {
		return nil, nil
	}

	raw := input[0].(map[string]interface{})
	enableAutoScaling := raw["enable_auto_scaling"].(bool)
	nodeTaintsRaw := raw["node_taints"].([]interface{})
	nodeTaints := utils.ExpandStringSlice(nodeTaintsRaw)

	profile := containerservice.ManagedClusterAgentPoolProfile{
		EnableAutoScaling:  utils.Bool(enableAutoScaling),
		EnableNodePublicIP: utils.Bool(raw["enable_node_public_ip"].(bool)),
		Name:               utils.String(raw["name"].(string)),
		NodeTaints:         nodeTaints,
		Type:               containerservice.AgentPoolType(raw["type"].(string)),
		VMSize:             containerservice.VMSizeTypes(raw["vm_size"].(string)),

		// at this time the default node pool has to be Linux or the AKS cluster fails to provision with:
		// Pods not in Running status: coredns-7fc597cc45-v5z7x,coredns-autoscaler-7ccc76bfbd-djl7j,metrics-server-cbd95f966-5rl97,tunnelfront-7d9884977b-wpbvn
		// Windows agents can be configured via the separate node pool resource
		OsType: containerservice.Linux,

		//// TODO: support these in time
		// OrchestratorVersion:    nil,
		// ScaleSetEvictionPolicy: "",
		// ScaleSetPriority:       "",
	}

	availabilityZonesRaw := raw["availability_zones"].([]interface{})
	availabilityZones := utils.ExpandStringSlice(availabilityZonesRaw)

	// otherwise: Standard Load Balancer is required for availability zone.
	if len(*availabilityZones) > 0 {
		profile.AvailabilityZones = availabilityZones
	}

	if maxPods := int32(raw["max_pods"].(int)); maxPods > 0 {
		profile.MaxPods = utils.Int32(maxPods)
	}

	if osDiskSizeGB := int32(raw["os_disk_size_gb"].(int)); osDiskSizeGB > 0 {
		profile.OsDiskSizeGB = utils.Int32(osDiskSizeGB)
	}

	if vnetSubnetID := raw["vnet_subnet_id"].(string); vnetSubnetID != "" {
		profile.VnetSubnetID = utils.String(vnetSubnetID)
	}

	count := raw["node_count"].(int)
	maxCount := raw["max_count"].(int)
	minCount := raw["min_count"].(int)

	// Count must be set for the initial creation when using AutoScaling but cannot be updated
	autoScaledCluster := enableAutoScaling && d.IsNewResource()

	// however it must always be sent for manually scaled clusters
	manuallyScaledCluster := !enableAutoScaling && (d.IsNewResource() || d.HasChange("default_node_pool.0.node_count"))

	if autoScaledCluster || manuallyScaledCluster {
		// users creating an auto-scaled cluster may not set the `node_count` field - if so use `min_count`
		if count == 0 && autoScaledCluster {
			count = minCount
		}

		profile.Count = utils.Int32(int32(count))
	}

	if enableAutoScaling {
		if maxCount > 0 {
			profile.MaxCount = utils.Int32(int32(maxCount))
		} else {
			return nil, fmt.Errorf("`max_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount > 0 {
			profile.MinCount = utils.Int32(int32(minCount))
		} else {
			return nil, fmt.Errorf("`min_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount > maxCount {
			return nil, fmt.Errorf("`max_count` must be >= `min_count`")
		}
	} else if minCount > 0 || maxCount > 0 {
		return nil, fmt.Errorf("`max_count` and `min_count` must be set to `0` when enable_auto_scaling is set to `false`")
	}

	return &[]containerservice.ManagedClusterAgentPoolProfile{
		profile,
	}, nil
}

func FlattenDefaultNodePool(input *[]containerservice.ManagedClusterAgentPoolProfile, d *schema.ResourceData) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	agentPool, err := findDefaultNodePool(input, d)
	if err != nil {
		return nil, err
	}

	var availabilityZones []string
	if agentPool.AvailabilityZones != nil {
		availabilityZones = *agentPool.AvailabilityZones
	}

	count := 0
	if agentPool.Count != nil {
		count = int(*agentPool.Count)
	}

	enableAutoScaling := false
	if agentPool.EnableAutoScaling != nil {
		enableAutoScaling = *agentPool.EnableAutoScaling
	}

	enableNodePublicIP := false
	if agentPool.EnableNodePublicIP != nil {
		enableNodePublicIP = *agentPool.EnableNodePublicIP
	}

	maxCount := 0
	if agentPool.MaxCount != nil {
		maxCount = int(*agentPool.MaxCount)
	}

	maxPods := 0
	if agentPool.MaxPods != nil {
		maxPods = int(*agentPool.MaxPods)
	}

	minCount := 0
	if agentPool.MinCount != nil {
		minCount = int(*agentPool.MinCount)
	}

	name := ""
	if agentPool.Name != nil {
		name = *agentPool.Name
	}

	var nodeTaints []string
	if agentPool.NodeTaints != nil {
		nodeTaints = *agentPool.NodeTaints
	}

	osDiskSizeGB := 0
	if agentPool.OsDiskSizeGB != nil {
		osDiskSizeGB = int(*agentPool.OsDiskSizeGB)
	}

	vnetSubnetId := ""
	if agentPool.VnetSubnetID != nil {
		vnetSubnetId = *agentPool.VnetSubnetID
	}

	return &[]interface{}{
		map[string]interface{}{
			"availability_zones":    availabilityZones,
			"enable_auto_scaling":   enableAutoScaling,
			"enable_node_public_ip": enableNodePublicIP,
			"max_count":             maxCount,
			"max_pods":              maxPods,
			"min_count":             minCount,
			"name":                  name,
			"node_count":            count,
			"node_taints":           nodeTaints,
			"os_disk_size_gb":       osDiskSizeGB,
			"type":                  string(agentPool.Type),
			"vm_size":               string(agentPool.VMSize),
			"vnet_subnet_id":        vnetSubnetId,
		},
	}, nil
}

func findDefaultNodePool(input *[]containerservice.ManagedClusterAgentPoolProfile, d *schema.ResourceData) (*containerservice.ManagedClusterAgentPoolProfile, error) {
	// first try loading this from the Resource Data if possible (e.g. when Created)
	defaultNodePoolName := d.Get("default_node_pool.0.name")

	var agentPool *containerservice.ManagedClusterAgentPoolProfile
	if defaultNodePoolName != "" {
		// find it
		for _, v := range *input {
			if v.Name != nil && *v.Name == defaultNodePoolName {
				agentPool = &v
				break
			}
		}
	} else {
		// otherwise we need to fall back to the name of the first agent pool
		for _, v := range *input {
			if v.Name == nil {
				continue
			}

			defaultNodePoolName = *v.Name
			agentPool = &v
			break
		}

		if defaultNodePoolName == nil {
			return nil, fmt.Errorf("Unable to Determine Default Agent Pool")
		}
	}

	if agentPool == nil {
		return nil, fmt.Errorf("The Default Agent Pool %q was not found", defaultNodePoolName)
	}

	return agentPool, nil
}
