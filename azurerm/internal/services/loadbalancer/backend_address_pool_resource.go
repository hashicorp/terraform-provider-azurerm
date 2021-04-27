package loadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/validate"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var backendAddressPoolResourceName = "azurerm_lb_backend_address_pool"

func resourceArmLoadBalancerBackendAddressPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerBackendAddressPoolCreateUpdate,
		Update: resourceArmLoadBalancerBackendAddressPoolCreateUpdate, // TODO: remove in 3.0 since all fields are ForceNew
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

		Schema: func() map[string]*schema.Schema {
			s := map[string]*schema.Schema{
				"name": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				// TODO 3.0: remove this as it can be inferred from "loadbalancer_id"
				"resource_group_name": azure.SchemaResourceGroupNameDeprecatedComputed(),

				"loadbalancer_id": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validate.LoadBalancerID,
				},

				"backend_ip_configurations": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"load_balancing_rules": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"outbound_rules": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			}

			if !features.ThreePointOh() {
				s["backend_address"] = &schema.Schema{
					Type:       schema.TypeSet,
					Optional:   true,
					Deprecated: "This field is non-functional and will be removed in version 3.0 of the Azure Provider - use the separate `azurerm_lb_backend_address_pool_address` resource instead.",
					MinItems:   1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"virtual_network_id": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: networkValidate.VirtualNetworkID,
							},

							"ip_address": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.IsIPAddress,
							},
						},
					},
				}
			}

			return s
		}(),
	}
}

func resourceArmLoadBalancerBackendAddressPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Load Balancer Name and Group: %+v", err)
	}

	name := d.Get("name").(string)
	id := parse.NewLoadBalancerBackendAddressPoolID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroup, loadBalancerId.Name, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Load Balancer Backend Address Pool %q: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_lb_backend_address_pool", id.ID())
		}
	}

	locks.ByName(name, backendAddressPoolResourceName)
	defer locks.UnlockByName(name, backendAddressPoolResourceName)

	locks.ByID(loadBalancerId.ID())
	defer locks.UnlockByID(loadBalancerId.ID())

	lb, err := lbClient.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(lb.Response) {
			return fmt.Errorf("Load Balancer %q for Backend Address Pool %q was not found", loadBalancerId, id)
		}
		return fmt.Errorf("failed to retrieve Load Balancer %q for Backend Address Pool %q: %+v", loadBalancerId, id, err)
	}

	param := network.BackendAddressPool{
		Name: &id.BackendAddressPoolName,
	}

	// Since API version 2020-05-01, there are two ways to CRUD backend address pool - either via the LB endpoint or via the
	// dedicated BAP endpoint. While based on different sku of the LB, users should insist on interacting one of the two endpoints:
	// - Basic sku: interact with LB endpoint for CUD
	// - Standard sku: interact with BAP endpoint for CUD
	// Particularly, the BAP endpoint can be used for R for bot cases.
	// See: https://github.com/Azure/azure-rest-api-specs/issues/11234 for details.
	sku := lb.Sku
	if sku == nil {
		return fmt.Errorf("nil or empty `sku` for Load Balancer %q for Backend Address Pool %q was not found", loadBalancerId, id)
	}

	if sku.Name == network.LoadBalancerSkuNameBasic {
		// Load balancer backend pool can be configured by either NIC or IP (belongs to a vnet).
		// In case of IP, it can only work for Standard sku LB.
		if len(d.Get("backend_address").(*schema.Set).List()) != 0 {
			return fmt.Errorf("only for Standard (sku) Load Balancer allows IP based Backend Address Pool configuration,"+
				"whilst %q for Backend Address Pool %q is of sku %s", loadBalancerId, id, sku.Name)
		}

		backendAddressPools := append(*lb.LoadBalancerPropertiesFormat.BackendAddressPools, param)
		_, existingPoolIndex, exists := FindLoadBalancerBackEndAddressPoolByName(&lb, id.BackendAddressPoolName)
		if exists {
			// this pool is being updated/reapplied remove the old copy from the slice
			backendAddressPools = append(backendAddressPools[:existingPoolIndex], backendAddressPools[existingPoolIndex+1:]...)
		}

		lb.LoadBalancerPropertiesFormat.BackendAddressPools = &backendAddressPools

		future, err := lbClient.CreateOrUpdate(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, lb)
		if err != nil {
			return fmt.Errorf("updating Load Balancer %q for Backend Address Pool %q: %+v", loadBalancerId, id, err)
		}

		if err = future.WaitForCompletionRef(ctx, lbClient.Client); err != nil {
			return fmt.Errorf("waiting for update of Load Balancer %q for Backend Address Pool %q: %+v", loadBalancerId, id, err)
		}
	} else {
		param.BackendAddressPoolPropertiesFormat = &network.BackendAddressPoolPropertiesFormat{
			// NOTE: Backend Addresses are managed using `azurerm_lb_backend_pool_address`
		}

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, param)
		if err != nil {
			return fmt.Errorf("creating/updating Load Balancer Backend Address Pool %q: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for Creating/Updating of Load Balancer Backend Address Pool %q: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerBackendAddressPoolRead(d, meta)
}

func resourceArmLoadBalancerBackendAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerBackendAddressPoolID(d.Id())
	if err != nil {
		return err
	}

	lbId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Backend Address Pool %q not found - removing from state", id)
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer Backend Address Pool %q: %+v", id, err)
	}

	d.Set("name", id.BackendAddressPoolName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("loadbalancer_id", lbId.ID())

	if props := resp.BackendAddressPoolPropertiesFormat; props != nil {
		// TODO: remove in 3.0
		if !features.ThreePointOh() {
			// @tombuildsstuff: this is a Set so won't be referenced, let's just nil this out for now
			if err := d.Set("backend_address", []interface{}{}); err != nil {
				return fmt.Errorf("setting `backend_address`: %v", err)
			}
		}

		var backendIPConfigurations []string
		if configs := props.BackendIPConfigurations; configs != nil {
			for _, backendConfig := range *configs {
				if backendConfig.ID == nil {
					continue
				}
				backendIPConfigurations = append(backendIPConfigurations, *backendConfig.ID)
			}
		}
		if err := d.Set("backend_ip_configurations", backendIPConfigurations); err != nil {
			return fmt.Errorf("setting `backend_ip_configurations`: %v", err)
		}

		var loadBalancingRules []string
		if rules := props.LoadBalancingRules; rules != nil {
			for _, rule := range *rules {
				if rule.ID == nil {
					continue
				}
				loadBalancingRules = append(loadBalancingRules, *rule.ID)
			}
		}
		if err := d.Set("load_balancing_rules", loadBalancingRules); err != nil {
			return fmt.Errorf("setting `load_balancing_rules`: %v", err)
		}

		var outboundRules []string
		if rules := props.OutboundRules; rules != nil {
			for _, rule := range *rules {
				if rule.ID == nil {
					continue
				}
				outboundRules = append(outboundRules, *rule.ID)
			}
		}
		if err := d.Set("outbound_rules", outboundRules); err != nil {
			return fmt.Errorf("setting `outbound_rules`: %v", err)
		}
	}

	return nil
}

func resourceArmLoadBalancerBackendAddressPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerBackendAddressPoolID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
	loadBalancerID := loadBalancerId.ID()
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
	defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

	lb, err := lbClient.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(lb.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer %q (resource group %q) for Backend Address Pool %q: %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, id.BackendAddressPoolName, err)
	}

	sku := lb.Sku
	if sku == nil {
		return fmt.Errorf("nil or empty `sku` for Load Balancer %q for Backend Address Pool %q was not found", loadBalancerId, id)
	}

	if sku.Name == network.LoadBalancerSkuNameBasic {
		_, index, exists := FindLoadBalancerBackEndAddressPoolByName(&lb, id.BackendAddressPoolName)
		if !exists {
			return nil
		}

		oldBackEndPools := *lb.LoadBalancerPropertiesFormat.BackendAddressPools
		newBackEndPools := append(oldBackEndPools[:index], oldBackEndPools[index+1:]...)
		lb.LoadBalancerPropertiesFormat.BackendAddressPools = &newBackEndPools

		future, err := lbClient.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, lb)
		if err != nil {
			return fmt.Errorf("updating Load Balancer %q (resource group %q) to remove Backend Address Pool %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.BackendAddressPoolName, err)
		}

		if err = future.WaitForCompletionRef(ctx, lbClient.Client); err != nil {
			return fmt.Errorf("waiting for update of Load Balancer %q (resource group %q) for Backend Address Pool %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.BackendAddressPoolName, err)
		}
	} else {
		future, err := client.Delete(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
		if err != nil {
			return fmt.Errorf("deleting Load Balancer Backend Address Pool %q: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for deleting of Load Balancer Backend Address Pool %q: %+v", id, err)
		}
	}

	return nil
}
