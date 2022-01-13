package loadbalancer

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"log"
	"time"
)

func resourceArmCrossRegionLoadBalancerBackendAddressPoolAddress() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressCreateUpdate,
		Update: resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressCreateUpdate,
		Read:   resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressRead,
		Delete: resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*parse.LoadBalancerId, error) {
			id, err := parse.LoadBalancerBackendAddressPoolID(input)
			if err != nil {
				return nil, err
			}

			lbId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
			return &lbId, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: func() map[string]*schema.Schema {
			s := map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"backend_address_pool_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				//my-loadbalancer-R1 & my-loadbalancer-R2 -> order mattered
				"backend_addresses": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*schema.Schema{
							"load_balancer": {
								Type:     pluginsdk.TypeList,
								Required: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*schema.Schema{
										"lb_name": {
											Type:     pluginsdk.TypeString,
											Required: true,
											//todo: validate function: Append the name of the load balancers, virtual machines, and other resources in each region with a -R1 and -R2.
											ValidateFunc: validation.StringIsNotEmpty,
										},

										"lb_ip_address": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.IsIPAddress,
										},

										"lb_frontend_ip_configuration_id": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: azure.ValidateResourceID,
										},
									},
								},
							},
						},
					},
				},
			}
			return s
		}(),
	}
}

func resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	//todo: forCreate or forCreatUpdate?
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	poolId, err := parse.LoadBalancerBackendAddressPoolID(d.Get("backend_address_pool_id").(string))
	if err != nil {
		return err
	}

	//Outer load balancer which contains the backend address pool
	lb, err := meta.(*clients.Client).LoadBalancers.LoadBalancersClient.Get(ctx, poolId.ResourceGroup, poolId.LoadBalancerName, "")
	if err != nil {
		return fmt.Errorf("retrieving Load Balancer %q (Resource Group %q): %+v", poolId.LoadBalancerName, poolId.ResourceGroup, err)
	}

	if lb.Sku != nil && lb.Sku.Tier != network.LoadBalancerSkuTierGlobal {
		return fmt.Errorf("Different Regional Backend Address Pool Addresses can only be set under the Global SKU tier")
	}

	addressName := d.Get("name").(string)
	id := parse.NewBackendAddressPoolAddressID(subscriptionId, poolId.ResourceGroup, poolId.LoadBalancerName, poolId.BackendAddressPoolName, addressName)
	pool, err := client.Get(ctx, poolId.ResourceGroup, poolId.LoadBalancerName, poolId.BackendAddressPoolName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *poolId, err)
	}
	if pool.BackendAddressPoolPropertiesFormat == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *poolId)
	}

	addresses := make([]network.LoadBalancerBackendAddress, 0)
	if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
		addresses = *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses
	}

	log.Printf("checking for existing %s..", id)
	for _, address := range addresses {
		if address.Name == nil {
			continue
		}

		if *address.Name == id.AddressName {
			return tf.ImportAsExistsError("azurerm_crlb_backend_address_pool_address", id.ID())
		}
	}
	addresses = expandBackendAddressPoolAddresses(d.Get("backend_addresses").([]interface{}))

	pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = &addresses

	log.Printf("adding %s..", id)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}
	log.Printf("waiting for update %s..", id)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerRead(d, meta)
}

func resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	//backend address pool address id
	id, err := parse.BackendAddressPoolAddressID(d.Id())
	if err != nil {
		return err
	}

	pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if pool.BackendAddressPoolPropertiesFormat == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	d.Set("name", id.AddressName)
	d.Set("backend_address_pool_id", (parse.NewLoadBalancerBackendAddressPoolID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)).ID())
	if err := d.Set("backend_addresses", flattenBackendAddressPoolAddresses(pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses)); err != nil {
		return fmt.Errorf("setting `backend_addresses` for backend address pool %s: %+v", id.BackendAddressPoolName, err)
	}

	return nil
}

func resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackendAddressPoolAddressID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
	defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

	pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if pool.BackendAddressPoolPropertiesFormat == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = nil
	log.Printf("removing %s..", id)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
	if err != nil {
		return fmt.Errorf("removing %s: %+v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for removal of %s: %+v", *id, err)
	}
	return nil
}

func expandBackendAddressPoolAddresses(input []interface{}) []network.LoadBalancerBackendAddress {
	if len(input) == 0 {
		return nil
	}

	backendAddresses := input[0].(map[string]interface{})
	var lbs []network.LoadBalancerBackendAddress

	if v, ok := backendAddresses["load_balancer"].([]interface{}); ok {
		if len(v) > 0 {
			for _, r := range v {
				lbBlocks := r.(map[string]interface{})
				lbs = append(lbs, network.LoadBalancerBackendAddress{
					Name: utils.String(lbBlocks["lb_name"].(string)),
					LoadBalancerBackendAddressPropertiesFormat: &network.LoadBalancerBackendAddressPropertiesFormat{
						LoadBalancerFrontendIPConfiguration: &network.SubResource{
							ID: utils.String(lbBlocks["lb_frontend_ip_configuration_id"].(string)),
						},
						IPAddress: utils.String(lbBlocks["lb_ip_address"].(string)),
					},
				})
			}
		}
	}

	return lbs
}

func flattenBackendAddressPoolAddresses(backendAddresses *[]network.LoadBalancerBackendAddress) []interface{} {
	if backendAddresses == nil {
		return nil
	}

	lbBlock := make([]interface{}, 0)
	for _, backendAddress := range *backendAddresses {
		block := make(map[string]interface{})

		if backendAddress.Name != nil {
			block["lb_name"] = *backendAddress.Name
		}

		if backendAddress.LoadBalancerBackendAddressPropertiesFormat.IPAddress != nil {
			block["lb_ip_address"] = *backendAddress.LoadBalancerBackendAddressPropertiesFormat.IPAddress
		}

		if backendAddress.LoadBalancerBackendAddressPropertiesFormat.LoadBalancerFrontendIPConfiguration.ID != nil {
			block["lb_frontend_ip_configuration_id"] = *backendAddress.LoadBalancerBackendAddressPropertiesFormat.LoadBalancerFrontendIPConfiguration.ID
		}

		lbBlock = append(lbBlock, block)
	}

	return []interface{}{map[string]interface{}{
		"load_balancer": lbBlock,
	}}
}
