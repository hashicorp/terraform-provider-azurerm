package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func dataSourceArmLoadBalancerBackendAddressPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmLoadBalancerBackendAddressPoolRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

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

			"backend_ip_configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmLoadBalancerBackendAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	loadBalancerID := d.Get("loadbalancer_id").(string)
	name := d.Get("name").(string)

	loadBalancer, exists, err := retrieveLoadBalancerById(d, d.Get("loadbalancer_id").(string), meta)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		return fmt.Errorf("Unable to retrieve Backend Address Pool %q since Load Balancer %q was not found", name, loadBalancerID)
	}

	bap, _, exists := FindLoadBalancerBackEndAddressPoolByName(loadBalancer, name)
	if !exists {
		return fmt.Errorf("Backend Address Pool %q was not found in Load Balancer %q", name, loadBalancerID)
	}

	d.SetId(*bap.ID)

	backendIPConfigurations := make([]interface{}, 0)
	if props := bap.BackendAddressPoolPropertiesFormat; props != nil {
		if beipConfigs := props.BackendIPConfigurations; beipConfigs != nil {
			for _, config := range *beipConfigs {
				ipConfig := make(map[string]interface{})
				if id := config.ID; id != nil {
					ipConfig["id"] = *id
					backendIPConfigurations = append(backendIPConfigurations, ipConfig)
				}
			}
		}
	}

	d.Set("backend_ip_configurations", backendIPConfigurations)

	return nil
}
