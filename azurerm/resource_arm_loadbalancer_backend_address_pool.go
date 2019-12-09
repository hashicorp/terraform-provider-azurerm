package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancerBackendAddressPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerBackendAddressPoolCreate,
		Read:   resourceArmLoadBalancerBackendAddressPoolRead,
		Delete: resourceArmLoadBalancerBackendAddressPoolDelete,
		Importer: &schema.ResourceImporter{
			State: loadBalancerSubResourceStateImporter,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationDeprecated(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"backend_ip_configurations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.NoEmptyStrings,
				},
				Set: schema.HashString,
			},

			"load_balancing_rules": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.NoEmptyStrings,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceArmLoadBalancerBackendAddressPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.LoadBalancersClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	loadBalancerID := d.Get("loadbalancer_id").(string)
	name := d.Get("name").(string)
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(d, loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	backendAddressPools := append(*loadBalancer.LoadBalancerPropertiesFormat.BackendAddressPools, expandAzureRmLoadBalancerBackendAddressPools(d))
	existingPool, existingPoolIndex, exists := findLoadBalancerBackEndAddressPoolByName(loadBalancer, name)
	if exists {
		if name == *existingPool.Name {
			if features.ShouldResourcesBeImported() && d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_lb_backend_address_pool", *existingPool.ID)
			}

			// this pool is being updated/reapplied remove old copy from the slice
			backendAddressPools = append(backendAddressPools[:existingPoolIndex], backendAddressPools[existingPoolIndex+1:]...)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.BackendAddressPools = &backendAddressPools
	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Load Balancer Name and Group: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (Resource Group %q) ID", loadBalancerName, resGroup)
	}

	var poolId string
	for _, BackendAddressPool := range *(*read.LoadBalancerPropertiesFormat).BackendAddressPools {
		if *BackendAddressPool.Name == name {
			poolId = *BackendAddressPool.ID
		}
	}

	if poolId == "" {
		return fmt.Errorf("Cannot find created Load Balancer Backend Address Pool ID %q", poolId)
	}

	d.SetId(poolId)

	return resourceArmLoadBalancerBackendAddressPoolRead(d, meta)
}

func resourceArmLoadBalancerBackendAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["backendAddressPools"]

	loadBalancer, exists, err := retrieveLoadBalancerById(d, d.Get("loadbalancer_id").(string), meta)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	config, _, exists := findLoadBalancerBackEndAddressPoolByName(loadBalancer, name)
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer Backend Address Pool %q not found. Removing from state", name)
		return nil
	}

	d.Set("name", config.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	var backendIpConfigurations []string
	var loadBalancingRules []string

	if props := config.BackendAddressPoolPropertiesFormat; props != nil {
		if configs := props.BackendIPConfigurations; configs != nil {
			for _, backendConfig := range *configs {
				backendIpConfigurations = append(backendIpConfigurations, *backendConfig.ID)
			}
		}

		if rules := props.LoadBalancingRules; rules != nil {
			for _, rule := range *rules {
				loadBalancingRules = append(loadBalancingRules, *rule.ID)
			}
		}
	}

	d.Set("backend_ip_configurations", backendIpConfigurations)
	d.Set("load_balancing_rules", loadBalancingRules)

	return nil
}

func resourceArmLoadBalancerBackendAddressPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	loadBalancerID := d.Get("loadbalancer_id").(string)
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(d, loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		return nil
	}

	_, index, exists := findLoadBalancerBackEndAddressPoolByName(loadBalancer, d.Get("name").(string))
	if !exists {
		return nil
	}

	oldBackEndPools := *loadBalancer.LoadBalancerPropertiesFormat.BackendAddressPools
	newBackEndPools := append(oldBackEndPools[:index], oldBackEndPools[index+1:]...)
	loadBalancer.LoadBalancerPropertiesFormat.BackendAddressPools = &newBackEndPools

	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating LoadBalancer: %+v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion for the LoadBalancer: %+v", err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving the Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (resource group %q) ID", loadBalancerName, resGroup)
	}

	return nil
}

func expandAzureRmLoadBalancerBackendAddressPools(d *schema.ResourceData) network.BackendAddressPool {
	return network.BackendAddressPool{
		Name: utils.String(d.Get("name").(string)),
	}
}
