package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-06-01/containerservice"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKubernetesClusterAgentPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKubernetesClusterAgentPoolCreateUpdate,
		Read:   resourceArmKubernetesClusterAgentPoolRead,
		Update: resourceArmKubernetesClusterAgentPoolCreateUpdate,
		Delete: resourceArmKubernetesClusterAgentPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"agent_pool_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			// AKS name
			// TODO replace by aks ID
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),

			"agent_pool_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(containerservice.AvailabilitySet),
				ValidateFunc: validation.StringInSlice([]string{
					string(containerservice.AvailabilitySet),
					string(containerservice.VirtualMachineScaleSets),
				}, false),
			},

			"orchestrator_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(containerservice.Linux),
				ValidateFunc: validation.StringInSlice([]string{
					string(containerservice.Linux),
					string(containerservice.Windows),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"vm_size": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.NoEmptyStrings,
			},

			"node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"os_disk_size_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"max_pods": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"vnet_subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enable_public_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"node_taints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"autoscale_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},

						"min_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
					},
				},
			},

			"vm_scale_set_profile": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"eviction_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmKubernetesClusterAgentPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containers.AgentPoolsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Managed Kubernetes Cluster Agent Pool create/update.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	agentPoolName := d.Get("agent_pool_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name, agentPoolName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Kubernetes Cluster %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_kubernetes_cluster", *existing.ID)
		}
	}

	agentProfile, err := expandKubernetesClusterAgentPoolProfileProperties(d)
	if err != nil {
		return err
	}

	parameters := containerservice.AgentPool{
		Name:                                     &name,
		ManagedClusterAgentPoolProfileProperties: &agentProfile,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, agentPoolName, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Managed Kubernetes Cluster Agent Pool %q (Resource Group %q): %+v", agentPoolName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Managed Kubernetes Cluster Agent Pool %q (Resource Group %q): %+v", agentPoolName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name, agentPoolName)
	if err != nil {
		return fmt.Errorf("Error retrieving Managed Kubernetes Cluster Agent Pool %q (Resource Group %q): %+v", agentPoolName, resGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Managed Kubernetes Cluster Agent Pool %q (Resource Group %q)", agentPoolName, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmKubernetesClusterAgentPoolRead(d, meta)
}

func resourceArmKubernetesClusterAgentPoolRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient).containers.AgentPoolsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[INFO] %+v", id.Path)
	resGroup := id.ResourceGroup
	agentPoolName := id.Path["agentPools"]
	name := id.Path["managedClusters"]

	resp, err := client.Get(ctx, resGroup, name, agentPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Managed Kubernetes Cluster %q was not found in Resource Group %q - removing from state!", name, resGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("agent_pool_name", agentPoolName)
	d.Set("name", name)
	d.Set("resource_group_name", resGroup)

	profile := resp.ManagedClusterAgentPoolProfileProperties
	if profile.Type != "" {
		d.Set("agent_pool_type", string(profile.Type))
	}

	if profile.Count != nil {
		d.Set("node_count", int(*profile.Count))
	}

	d.Set("availability_zones", utils.FlattenStringSlice(profile.AvailabilityZones))

	if profile.VMSize != "" {
		d.Set("vm_size", string(profile.VMSize))
	}

	if profile.OsDiskSizeGB != nil {
		d.Set("os_disk_size_gb", int(*profile.OsDiskSizeGB))
	}

	if profile.VnetSubnetID != nil {
		d.Set("vnet_subnet_id", *profile.VnetSubnetID)
	}

	if profile.OsType != "" {
		d.Set("os_type", string(profile.OsType))
	}

	if profile.MaxPods != nil {
		d.Set("max_pods", int(*profile.MaxPods))
	}

	if profile.NodeTaints != nil {
		d.Set("node_taints", *profile.NodeTaints)
	}

	if *profile.EnableAutoScaling == true {
		autoscale_configs := make([]interface{}, 0)
		autoscale_config := make(map[string]interface{})
		if profile.MaxCount != nil {
			autoscale_config["max_count"] = int(*profile.MaxCount)
		}
		if profile.MinCount != nil {
			autoscale_config["min_count"] = int(*profile.MinCount)
		}
		autoscale_configs = append(autoscale_configs, autoscale_config)
		d.Set("autoscale_configuration", autoscale_configs)

	}

	if profile.ScaleSetPriority != "" || profile.ScaleSetEvictionPolicy != "" {
		vm_scale_set_profiles := make([]interface{}, 0)
		vm_scale_set_profile := make(map[string]interface{})
		if profile.ScaleSetPriority != "" {
			vm_scale_set_profile["priority"] = profile.ScaleSetPriority
		}
		if profile.ScaleSetEvictionPolicy != "" {
			vm_scale_set_profile["eviction_policy"] = profile.ScaleSetEvictionPolicy
		}
		vm_scale_set_profiles = append(vm_scale_set_profiles, vm_scale_set_profile)
		d.Set("vm_scale_set_profile", vm_scale_set_profiles)
	}

	return nil
}

func resourceArmKubernetesClusterAgentPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containers.AgentPoolsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	agentPoolName := id.Path["agentPools"]
	name := id.Path["managedClusters"]

	future, err := client.Delete(ctx, resGroup, name, agentPoolName)
	if err != nil {
		return fmt.Errorf("Error deleting Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func expandKubernetesClusterAgentPoolProfileProperties(d *schema.ResourceData) (containerservice.ManagedClusterAgentPoolProfileProperties, error) {

	// TODO Default: AvailabilitySet ??
	poolType := d.Get("agent_pool_type").(string)
	// TODO Default Linux ???
	osType := d.Get("os_type").(string)
	vmSize := d.Get("vm_size").(string)
	count := int32(d.Get("node_count").(int))
	osDiskSizeGB := int32(d.Get("os_disk_size_gb").(int))

	profile := containerservice.ManagedClusterAgentPoolProfileProperties{
		Type:         containerservice.AgentPoolType(poolType),
		Count:        utils.Int32(count),
		VMSize:       containerservice.VMSizeTypes(vmSize),
		OsDiskSizeGB: utils.Int32(osDiskSizeGB),
		OsType:       containerservice.OSType(osType),
	}

	// default: get cluster versio
	orchestratorVersion := d.Get("orchestrator_version").(string)
	if orchestratorVersion != "" {
		profile.OrchestratorVersion = utils.String(orchestratorVersion)
	}

	if maxPods := int32(d.Get("max_pods").(int)); maxPods > 0 {
		profile.MaxPods = utils.Int32(maxPods)
	}

	vnetSubnetID := d.Get("vnet_subnet_id").(string)
	if vnetSubnetID != "" {
		profile.VnetSubnetID = utils.String(vnetSubnetID)
	}

	if availavilityZones := utils.ExpandStringSlice(d.Get("availability_zones").([]interface{})); len(*availavilityZones) > 0 {
		profile.AvailabilityZones = availavilityZones
	}

	if nodeTaints := utils.ExpandStringSlice(d.Get("node_taints").([]interface{})); len(*nodeTaints) > 0 {
		profile.NodeTaints = nodeTaints
	}

	if autoscale_configs := d.Get("autoscale_configuration").([]interface{}); len(autoscale_configs) > 0 {
		profile.EnableAutoScaling = utils.Bool(true)
		if len(autoscale_configs) > 0 {
			autoscale_config := autoscale_configs[0].(map[string]interface{})
			profile.MaxCount = utils.Int32(int32(autoscale_config["max_count"].(int)))
			profile.MinCount = utils.Int32(int32(autoscale_config["min_count"].(int)))
		}
	} else {
		profile.EnableAutoScaling = utils.Bool(false)
	}

	if vm_scale_set_profiles := d.Get("vm_scale_set_profile").([]interface{}); len(vm_scale_set_profiles) > 0 {
		if len(vm_scale_set_profiles) > 0 {
			vm_scale_set_profile := vm_scale_set_profiles[0].(map[string]interface{})
			profile.ScaleSetPriority = containerservice.ScaleSetPriority(vm_scale_set_profile["priority"].(string))
			profile.ScaleSetEvictionPolicy = containerservice.ScaleSetEvictionPolicy(vm_scale_set_profile["eviction_policy"].(string))
		}
	}

	return profile, nil
}

func convertKubernetesClusterAgentPoolProfileToKubernetesClusterAgentPoolProfileProperties(agentProfile containerservice.ManagedClusterAgentPoolProfile) containerservice.ManagedClusterAgentPoolProfileProperties {
	agentProfileProperties := containerservice.ManagedClusterAgentPoolProfileProperties{
		Type:                   agentProfile.Type,
		Count:                  agentProfile.Count,
		VMSize:                 agentProfile.VMSize,
		OsDiskSizeGB:           agentProfile.OsDiskSizeGB,
		OsType:                 agentProfile.OsType,
		OrchestratorVersion:    agentProfile.OrchestratorVersion,
		MaxPods:                agentProfile.MaxPods,
		VnetSubnetID:           agentProfile.VnetSubnetID,
		AvailabilityZones:      agentProfile.AvailabilityZones,
		NodeTaints:             agentProfile.NodeTaints,
		EnableAutoScaling:      agentProfile.EnableAutoScaling,
		MaxCount:               agentProfile.MaxCount,
		MinCount:               agentProfile.MinCount,
		ScaleSetPriority:       agentProfile.ScaleSetPriority,
		ScaleSetEvictionPolicy: agentProfile.ScaleSetEvictionPolicy,
	}

	return agentProfileProperties

}
