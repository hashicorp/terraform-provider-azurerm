package loadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancerBackendAddressPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerBackendAddressPoolCreate,
		Read:   resourceArmLoadBalancerBackendAddressPoolRead,
		Delete: resourceArmLoadBalancerBackendAddressPoolDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*parse.LoadBalancerId, error) {
			id, err := parse.LoadBalancerBackendAddressPoolID(input)
			if err != nil {
				return nil, err
			}

			lbId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
			return &lbId, nil
		}),

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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LoadBalancerID,
			},

			"backend_ip_configurations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Set: schema.HashString,
			},

			"load_balancing_rules": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceArmLoadBalancerBackendAddressPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Load Balancer Name and Group: %+v", err)
	}

	id := parse.NewLoadBalancerBackendAddressPoolID(subscriptionId, loadBalancerId.ResourceGroup, loadBalancerId.Name, d.Get("name").(string))

	loadBalancerID := loadBalancerId.ID("")
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, *loadBalancerId)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.BackendAddressPoolName)
		return nil
	}

	backendAddressPools := append(*loadBalancer.LoadBalancerPropertiesFormat.BackendAddressPools, network.BackendAddressPool{
		Name: utils.String(id.BackendAddressPoolName),
	})
	existingPool, existingPoolIndex, exists := FindLoadBalancerBackEndAddressPoolByName(loadBalancer, id.BackendAddressPoolName)
	if exists {
		if id.BackendAddressPoolName == *existingPool.Name {
			if d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_lb_backend_address_pool", *existingPool.ID)
			}

			// this pool is being updated/reapplied remove old copy from the slice
			backendAddressPools = append(backendAddressPools[:existingPoolIndex], backendAddressPools[existingPoolIndex+1:]...)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.BackendAddressPools = &backendAddressPools

	future, err := client.CreateOrUpdate(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, err)
	}

	d.SetId(id.ID(""))

	return resourceArmLoadBalancerBackendAddressPoolRead(d, meta)
}

func resourceArmLoadBalancerBackendAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerBackendAddressPoolID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, loadBalancerId)
	if err != nil {
		return fmt.Errorf("retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
		return nil
	}

	config, _, exists := FindLoadBalancerBackEndAddressPoolByName(loadBalancer, id.BackendAddressPoolName)
	if !exists {
		log.Printf("[INFO] Load Balancer Backend Address Pool %q not found. Removing from state", id.BackendAddressPoolName)
		d.SetId("")
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
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerBackendAddressPoolID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
	loadBalancerID := loadBalancerId.ID("")
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, loadBalancerId)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		return nil
	}

	_, index, exists := FindLoadBalancerBackEndAddressPoolByName(loadBalancer, d.Get("name").(string))
	if !exists {
		return nil
	}

	oldBackEndPools := *loadBalancer.LoadBalancerPropertiesFormat.BackendAddressPools
	newBackEndPools := append(oldBackEndPools[:index], oldBackEndPools[index+1:]...)
	loadBalancer.LoadBalancerPropertiesFormat.BackendAddressPools = &newBackEndPools

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating LoadBalancer: %+v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion for the LoadBalancer: %+v", err)
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving the Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (resource group %q) ID", id.LoadBalancerName, id.ResourceGroup)
	}

	return nil
}
