package loadbalancer

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
				ValidateFunc: validate.LoadBalancerID,
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
	client := meta.(*clients.Client).Network.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}

	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, *loadBalancerId)
	if err != nil {
		return fmt.Errorf("retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		return fmt.Errorf("Load Balancer %q (Resource Group %q) was not found", loadBalancerId.Name, loadBalancerId.ResourceGroup)
	}

	bap, _, exists := FindLoadBalancerBackEndAddressPoolByName(loadBalancer, name)
	if !exists {
		return fmt.Errorf("Backend Address Pool %q was not found in Load Balancer %q (Resource Group %q)", name, loadBalancerId.Name, loadBalancerId.ResourceGroup)
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
