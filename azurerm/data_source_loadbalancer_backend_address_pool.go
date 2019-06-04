package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func dataSourceArmLoadBalancerBackendAddressPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmLoadBalancerBackendAddressPoolRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func dataSourceArmLoadBalancerBackendAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	loadBalancerId := d.Get("loadbalancer_id").(string)
	name := d.Get("name").(string)

	loadBalancer, exists, err := retrieveLoadBalancerById(d.Get("loadbalancer_id").(string), meta)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		return fmt.Errorf("Unable to retrieve Backend Address Pool %q since Load Balancer %q was not found", name, loadBalancerId)
	}

	config, _, exists := findLoadBalancerBackEndAddressPoolByName(loadBalancer, name)
	if !exists {
		return fmt.Errorf("Backend Address Pool %q was not found in Load Balancer %q", name, loadBalancerId)
	}

	d.SetId(*config.ID)

	return nil
}
