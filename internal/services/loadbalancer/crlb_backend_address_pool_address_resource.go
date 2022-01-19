package loadbalancer

import (
	"fmt"
	"log"
	"time"

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
)

func resourceArmCrossRegionLoadBalancerBackendAddressPoolAddress() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressCreate,
		Update: resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressUpdate,
		Read:   resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressRead,
		Delete: resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.BackendAddressPoolAddressID(id)
			return err
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
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"backend_address_pool_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				//"backend_address_ip_address": {
				//	Type:         pluginsdk.TypeString,
				//	Optional:     true,
				//	ValidateFunc: validation.IsIPAddress,
				//},

				"backend_address_ip_configuration_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
			}
			return s
		}(),
	}
}

func resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := client.SubscriptionID
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	poolId, err := parse.LoadBalancerBackendAddressPoolID(d.Get("backend_address_pool_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(poolId.BackendAddressPoolName, backendAddressPoolResourceName)
	defer locks.UnlockByName(poolId.BackendAddressPoolName, backendAddressPoolResourceName)

	lb, err := lbClient.Get(ctx, poolId.ResourceGroup, poolId.LoadBalancerName, "")
	if err != nil {
		return fmt.Errorf("retrieving Load Balancer %q (Resource Group %q): %+v", poolId.LoadBalancerName, poolId.ResourceGroup, err)
	}

	if lb.Sku != nil && lb.Sku.Tier != network.LoadBalancerSkuTierGlobal {
		return fmt.Errorf("Regional Backend Address Pool Addresses can only be set under the Global SKU tier")
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

	backendAddress := make([]network.LoadBalancerBackendAddress, 0)
	if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
		backendAddress = *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses
	}

	log.Printf("checking for existing %s..", id)
	for _, address := range backendAddress {
		if address.Name == nil {
			continue
		}

		if *address.Name == id.AddressName {
			return tf.ImportAsExistsError("azurerm_crlb_backend_address_pool_address", id.ID())
		}
	}

	//var ipAddress string
	//if v := d.Get("backend_address_ip_address"); v != "" {
	//	ipAddress = v.(string)
	//}

	backendAddress = append(backendAddress, network.LoadBalancerBackendAddress{
		Name: utils.String(addressName),
		LoadBalancerBackendAddressPropertiesFormat: &network.LoadBalancerBackendAddressPropertiesFormat{
			LoadBalancerFrontendIPConfiguration: &network.SubResource{
				ID: utils.String(d.Get("backend_address_ip_configuration_id").(string)),
			},
		},
	})

	pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = &backendAddress

	log.Printf("adding backend addresses %s for backend address pool..", id)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
	if err != nil {
		return fmt.Errorf("updating backend addresses %s for backend address pool: %+v", id, err)
	}
	log.Printf("waiting for updating backend addresses %s for backend address pool ..", id)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of backend addresses %s for backend address pool: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressRead(d, meta)
}

func resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	id, err := parse.BackendAddressPoolAddressID(d.Id())
	if err != nil {
		return nil
	}
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if pool.BackendAddressPoolPropertiesFormat == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	var backendAddress *network.LoadBalancerBackendAddress
	if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
		for _, address := range *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses {
			if address.Name == nil {
				continue
			}

			if *address.Name == id.AddressName {
				backendAddress = &address
				break
			}
		}
	}
	if backendAddress == nil {
		log.Printf("[DEBUG] %s was not found - removing from state", *id)
		d.SetId("")
		return nil
	}

	backendAddressPoolId := parse.NewLoadBalancerBackendAddressPoolID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	d.Set("backend_address_pool_id", backendAddressPoolId.ID())

	d.Set("name", id.AddressName)
	if props := backendAddress.LoadBalancerBackendAddressPropertiesFormat; props != nil {
		//if props.IPAddress != nil {
		//	d.Set("backend_address_ip_address", *props.IPAddress)
		//}
		if props.LoadBalancerFrontendIPConfiguration != nil && props.LoadBalancerFrontendIPConfiguration.ID != nil {
			d.Set("backend_address_ip_configuration_id", *props.LoadBalancerFrontendIPConfiguration.ID)
		}
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

	addresses := make([]network.LoadBalancerBackendAddress, 0)
	if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
		addresses = *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses
	}

	newAddresses := make([]network.LoadBalancerBackendAddress, 0)
	for _, address := range addresses {
		if address.Name == nil {
			continue
		}

		if *address.Name != id.AddressName {
			newAddresses = append(newAddresses, address)
		}
	}
	pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = &newAddresses

	pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = nil
	log.Printf("removing backend address %s for backend address pool..", id)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
	if err != nil {
		return fmt.Errorf("removing backend address %s for backend address pool: %+v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for removal of backend address %s for backend address pool: %+v", *id, err)
	}
	return nil
}

func resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
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

	backendAddress := make([]network.LoadBalancerBackendAddress, 0)
	if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
		backendAddress = *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses
	}
	index := -1
	for i, address := range backendAddress {
		if address.Name == nil {
			continue
		}

		if *address.Name == id.AddressName {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("%s was not found", *id)
	}

	//var ipAddress string
	//if v := d.Get("backend_address_ip_address"); v != "" {
	//	ipAddress = v.(string)
	//}
	backendAddress[index] = network.LoadBalancerBackendAddress{
		Name: utils.String(d.Get("name").(string)),
		LoadBalancerBackendAddressPropertiesFormat: &network.LoadBalancerBackendAddressPropertiesFormat{
			//IPAddress: &ipAddress,
			LoadBalancerFrontendIPConfiguration: &network.SubResource{
				ID: utils.String(d.Get("backend_address_ip_configuration_id").(string)),
			},
		},
	}

	pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = &backendAddress

	log.Printf("updating backend addresses %s for backend address pool..", id)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
	if err != nil {
		return fmt.Errorf("updating backend addresses %s for backend address pool: %+v", id, err)
	}
	log.Printf("waiting for updating backend addresses %s for backend address pool ..", id)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of backend addresses %s for backend address pool: %+v", id, err)
	}

	return resourceArmCrossRegionLoadBalancerBackendAddressPoolAddressRead(d, meta)
}
