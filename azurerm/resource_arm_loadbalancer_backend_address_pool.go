package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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
			Create: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"location": deprecatedLocationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

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
					ValidateFunc: validation.NoZeroValues,
				},
				Set: schema.HashString,
			},

			"load_balancing_rules": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceArmLoadBalancerBackendAddressPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).loadBalancerClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	loadBalancerID := d.Get("loadbalancer_id").(string)
	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(loadBalancerID)
	if err != nil {
		return fmt.Errorf("Error parsing Load Balancer Name and Group: %+v", err)
	}

	armMutexKV.Lock(loadBalancerID)
	defer armMutexKV.Unlock(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	props := loadBalancer.LoadBalancerPropertiesFormat
	if props == nil {
		return fmt.Errorf("[ERROR] Load Balancer Properties was nil - unable to configure Backend Address Pool")
	}

	pools := props.BackendAddressPools
	if pools == nil {
		return fmt.Errorf("[ERROR] Load Balancer Backend Address Pools was nil - unable to configure Backend Address Pool")
	}

	// firstly does this already exist
	for _, pool := range *pools {
		if pool.Name != nil && *pool.Name == name {
			return tf.ImportAsExistsError("azurerm_lb_backend_address_pool", *pool.ID)
		}
	}

	// as such we should just be able to append and send
	newPool := network.BackendAddressPool{
		Name: utils.String(name),
	}
	backendAddressPools := append(*pools, newPool)
	loadBalancer.LoadBalancerPropertiesFormat.BackendAddressPools = &backendAddressPools

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		return fmt.Errorf("Error Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (Resource Group %q) ID", loadBalancerName, resGroup)
	}

	var poolId string
	if props := read.LoadBalancerPropertiesFormat; props != nil {
		if pools := props.BackendAddressPools; pools != nil {
			for _, pool := range *pools {
				if pool.Name != nil && *pool.Name == name {
					poolId = *pool.ID
				}
			}
		}
	}

	if poolId == "" {
		return fmt.Errorf("Cannot find created Load Balancer Backend Address Pool ID %q", poolId)
	}

	d.SetId(poolId)

	return resourceArmLoadBalancerBackendAddressPoolRead(d, meta)
}

func resourceArmLoadBalancerBackendAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["backendAddressPools"]
	loadBalancerId := d.Get("loadbalancer_id").(string)

	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerId, meta)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	var pool *network.BackendAddressPool
	if props := loadBalancer.LoadBalancerPropertiesFormat; props != nil {
		if pools := props.BackendAddressPools; pools != nil {
			for _, p := range *pools {
				if p.Name != nil && *p.Name == name {
					pool = &p
					break
				}
			}
		}
	}

	if pool == nil {
		d.SetId("")
		log.Printf("[INFO] Backend Address Pool %q was not found on Load Balancer %q. Removing from state", name, *loadBalancer.ID)
		return nil
	}

	d.Set("name", name)
	d.Set("resource_group_name", id.ResourceGroup)

	var backendIpConfigurations []string
	var loadBalancingRules []string

	if props := pool.BackendAddressPoolPropertiesFormat; props != nil {
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
	client := meta.(*ArmClient).loadBalancerClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	loadBalancerID := d.Get("loadbalancer_id").(string)
	armMutexKV.Lock(loadBalancerID)
	defer armMutexKV.Unlock(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		return nil
	}

	props := loadBalancer.LoadBalancerPropertiesFormat
	if props == nil {
		return nil
	}

	pools := props.BackendAddressPools
	if pools == nil {
		return nil
	}

	newBackEndPools := make([]network.BackendAddressPool, 0)
	for _, pool := range *pools {
		if pool.Name != nil && *pool.Name != name {
			newBackEndPools = append(newBackEndPools, pool)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.BackendAddressPools = &newBackEndPools

	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(loadBalancerID)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating LoadBalancer: %+v", err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
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
