package containers

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-09-01/containerservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SchemaDefaultNodePool() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
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
					ValidateFunc: validation.StringIsNotEmpty,
				},

				// Optional
				"availability_zones": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
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
					Type:     schema.TypeInt,
					Optional: true,
					// NOTE: rather than setting `0` users should instead pass `null` here
					ValidateFunc: validation.IntBetween(1, 100),
				},

				"max_pods": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},

				"min_count": {
					Type:     schema.TypeInt,
					Optional: true,
					// NOTE: rather than setting `0` users should instead pass `null` here
					ValidateFunc: validation.IntBetween(1, 100),
				},

				"node_count": {
					Type:         schema.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 100),
				},

				"node_labels": {
					Type:     schema.TypeMap,
					ForceNew: true,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"node_taints": {
					Type:     schema.TypeList,
					ForceNew: true,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"tags": tags.Schema(),

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
				"orchestrator_version": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringIsNotEmpty,
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
			SpotMaxPrice:           defaultCluster.SpotMaxPrice,
			Mode:                   defaultCluster.Mode,
			NodeLabels:             defaultCluster.NodeLabels,
			NodeTaints:             defaultCluster.NodeTaints,
			Tags:                   defaultCluster.Tags,
		},
	}
}

func ExpandDefaultNodePool(d *schema.ResourceData) (*[]containerservice.ManagedClusterAgentPoolProfile, error) {
	input := d.Get("default_node_pool").([]interface{})

	raw := input[0].(map[string]interface{})
	enableAutoScaling := raw["enable_auto_scaling"].(bool)
	nodeLabelsRaw := raw["node_labels"].(map[string]interface{})
	nodeLabels := utils.ExpandMapStringPtrString(nodeLabelsRaw)
	nodeTaintsRaw := raw["node_taints"].([]interface{})
	nodeTaints := utils.ExpandStringSlice(nodeTaintsRaw)

	if len(*nodeTaints) != 0 {
		return nil, fmt.Errorf("The AKS API has removed support for tainting all nodes in the default node pool and it is no longer possible to configure this. To taint a node pool, create a separate one")
	}

	t := raw["tags"].(map[string]interface{})

	profile := containerservice.ManagedClusterAgentPoolProfile{
		EnableAutoScaling:  utils.Bool(enableAutoScaling),
		EnableNodePublicIP: utils.Bool(raw["enable_node_public_ip"].(bool)),
		Name:               utils.String(raw["name"].(string)),
		NodeLabels:         nodeLabels,
		Tags:               tags.Expand(t),
		Type:               containerservice.AgentPoolType(raw["type"].(string)),
		VMSize:             containerservice.VMSizeTypes(raw["vm_size"].(string)),

		// at this time the default node pool has to be Linux or the AKS cluster fails to provision with:
		// Pods not in Running status: coredns-7fc597cc45-v5z7x,coredns-autoscaler-7ccc76bfbd-djl7j,metrics-server-cbd95f966-5rl97,tunnelfront-7d9884977b-wpbvn
		// Windows agents can be configured via the separate node pool resource
		OsType: containerservice.Linux,

		// without this set the API returns:
		// Code="MustDefineAtLeastOneSystemPool" Message="Must define at least one system pool."
		// since this is the "default" node pool we can assume this is a system node pool
		Mode: containerservice.System,

		// // TODO: support these in time
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

	if orchestratorVersion := raw["orchestrator_version"].(string); orchestratorVersion != "" {
		profile.OrchestratorVersion = utils.String(orchestratorVersion)
	}

	count := raw["node_count"].(int)
	maxCount := raw["max_count"].(int)
	minCount := raw["min_count"].(int)

	// Count must always be set (see #6094), RP behaviour has changed
	// since the API version upgrade in v2.1.0 making Count required
	// for all create/update requests
	profile.Count = utils.Int32(int32(count))

	if enableAutoScaling {
		// if Count has not been set use min count
		if count == 0 {
			count = minCount
			profile.Count = utils.Int32(int32(count))
		}

		// Count must be set for the initial creation when using AutoScaling but cannot be updated
		if d.HasChange("default_node_pool.0.node_count") && !d.IsNewResource() {
			return nil, fmt.Errorf("cannot change `node_count` when `enable_auto_scaling` is set to `true`")
		}

		if maxCount > 0 {
			profile.MaxCount = utils.Int32(int32(maxCount))
			if maxCount < count {
				return nil, fmt.Errorf("`node_count`(%d) must be equal to or less than `max_count`(%d) when `enable_auto_scaling` is set to `true`", count, maxCount)
			}
		} else {
			return nil, fmt.Errorf("`max_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount > 0 {
			profile.MinCount = utils.Int32(int32(minCount))

			if minCount > count && d.IsNewResource() {
				return nil, fmt.Errorf("`node_count`(%d) must be equal to or greater than `min_count`(%d) when `enable_auto_scaling` is set to `true`", count, minCount)
			}
		} else {
			return nil, fmt.Errorf("`min_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount > maxCount {
			return nil, fmt.Errorf("`max_count` must be >= `min_count`")
		}
	} else if minCount > 0 || maxCount > 0 {
		return nil, fmt.Errorf("`max_count`(%d) and `min_count`(%d) must be set to `null` when `enable_auto_scaling` is set to `false`", maxCount, minCount)
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

	var nodeLabels map[string]string
	if agentPool.NodeLabels != nil {
		nodeLabels = make(map[string]string)
		for k, v := range agentPool.NodeLabels {
			nodeLabels[k] = *v
		}
	}

	osDiskSizeGB := 0
	if agentPool.OsDiskSizeGB != nil {
		osDiskSizeGB = int(*agentPool.OsDiskSizeGB)
	}

	vnetSubnetId := ""
	if agentPool.VnetSubnetID != nil {
		vnetSubnetId = *agentPool.VnetSubnetID
	}

	orchestratorVersion := ""
	if agentPool.OrchestratorVersion != nil {
		orchestratorVersion = *agentPool.OrchestratorVersion
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
			"node_labels":           nodeLabels,
			"node_taints":           []string{},
			"os_disk_size_gb":       osDiskSizeGB,
			"tags":                  tags.Flatten(agentPool.Tags),
			"type":                  string(agentPool.Type),
			"vm_size":               string(agentPool.VMSize),
			"orchestrator_version":  orchestratorVersion,
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
	}

	if agentPool == nil {
		// otherwise we need to fall back to the name of the first agent pool
		for _, v := range *input {
			if v.Name == nil {
				continue
			}
			if v.Mode != containerservice.System {
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
