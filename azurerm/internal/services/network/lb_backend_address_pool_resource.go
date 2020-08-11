package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancerBackendAddressPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerBackendAddressPoolCreateOrUpdate,
		Read:   resourceArmLoadBalancerBackendAddressPoolRead,
		Update: resourceArmLoadBalancerBackendAddressPoolCreateOrUpdate,
		Delete: resourceArmLoadBalancerBackendAddressPoolDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*parse.LoadBalancerId, error) {
			id, err := parse.LoadBalancerBackendAddressPoolID(input)
			if err != nil {
				return nil, err
			}

			lbId := parse.NewLoadBalancerID(id.ResourceGroup, id.LoadBalancerName)
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

			"ip_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"virtual_network_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"ip_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
					},
				},
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

func resourceArmLoadBalancerBackendAddressPoolCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	loadBalancerID := d.Get("loadbalancer_id").(string)
	name := d.Get("name").(string)
	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Load Balancer Name and Group: %+v", err)
	}

	loadBalancerID := loadBalancerId.ID(subscriptionId)
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	lb, exists, err := retrieveLoadBalancerById(d, loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Load Balancer Name and Group: %+v", err)
	}

	if strings.EqualFold(string(lb.Sku.Name), string(network.LoadBalancerSkuNameBasic)) {
		if _, ok := d.GetOk("ip_address"); ok {
			return fmt.Errorf("LoadBalancer uses Basic Sku and cannot currently configured for a backend address pool with pre-allocated IP addresses")
		}

		client := meta.(*clients.Client).Network.LoadBalancersClient
		ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
		defer cancel()

		backendAddressPools := append(*lb.LoadBalancerPropertiesFormat.BackendAddressPools, expandAzureRmLoadBalancerBackendAddressPools(d))
		existingPool, existingPoolIndex, exists := FindLoadBalancerBackEndAddressPoolByName(lb, name)
		if exists {
			if name == *existingPool.Name {
				if features.ShouldResourcesBeImported() && d.IsNewResource() {
					return tf.ImportAsExistsError("azurerm_lb_backend_address_pool", *existingPool.ID)
				}

				// this pool is being updated/reapplied remove old copy from the slice
				backendAddressPools = append(backendAddressPools[:existingPoolIndex], backendAddressPools[existingPoolIndex+1:]...)
			}
		}
		lb.LoadBalancerPropertiesFormat.BackendAddressPools = &backendAddressPools

		future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *lb)
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
		for _, BackendAddressPool := range *read.LoadBalancerPropertiesFormat.BackendAddressPools {
			if *BackendAddressPool.Name == name {
				poolId = *BackendAddressPool.ID
			}
		}

		if poolId == "" {
			return fmt.Errorf("Cannot find created Load Balancer Backend Address Pool ID %q", poolId)
		}

		d.SetId(poolId)
	} else {
		client := meta.(*clients.Client).Network.LoadBalancerBackendAddressPoolsClient
		ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
		defer cancel()

		backendPoolProperties := network.BackendAddressPoolPropertiesFormat{}
		if v, ok := d.GetOk("ip_address"); ok {
			backendIPAddresses := expandArmLoadBalancerBackendIPAddressPool(v.([]interface{}))
			backendPoolProperties.LoadBalancerBackendAddresses = backendIPAddresses
		}
		future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, name, network.BackendAddressPool{
			BackendAddressPoolPropertiesFormat: &backendPoolProperties,
		})
		if err != nil {
			return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
		}

		read, err := client.Get(ctx, resGroup, loadBalancerName, name)
		if err != nil {
			return fmt.Errorf("Retrieving Load Balancer backend address pool %s on %q (Resource Group %q): %+v", name, loadBalancerName, resGroup, err)
		}
		if read.ID == nil {
			return fmt.Errorf("New backend address pool %s on Load Balancer %q (Resource Group %q) does not provide any ID", name, loadBalancerName, resGroup)
		}

		d.SetId(*read.ID)
	}

	return resourceArmLoadBalancerBackendAddressPoolRead(d, meta)
}

func resourceArmLoadBalancerBackendAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	return readArmLoadBalancerBackendIPAddressPool(d, meta, false)
}

func resourceArmLoadBalancerBackendAddressPoolDelete(d *schema.ResourceData, meta interface{}) error {
	loadBalancerID := d.Get("loadbalancer_id").(string)
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	lb, exists, err := retrieveLoadBalancerById(d, loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		return nil
	}

	name := d.Get("name").(string)
	_, index, exists := FindLoadBalancerBackEndAddressPoolByName(lb, name)
	if !exists {
		log.Printf("[INFO] Load Balancer Backend Address Pool %q not found. Removing from state", id.Name)
		d.SetId("")
		return nil
	}

	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
	}

	if strings.EqualFold(string(lb.Sku.Name), string(network.LoadBalancerSkuNameBasic)) {
		client := meta.(*clients.Client).Network.LoadBalancersClient
		ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
		defer cancel()

		oldBackEndPools := *lb.LoadBalancerPropertiesFormat.BackendAddressPools
		newBackEndPools := append(oldBackEndPools[:index], oldBackEndPools[index+1:]...)
		lb.LoadBalancerPropertiesFormat.BackendAddressPools = &newBackEndPools

		resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
		if err != nil {
			return fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
		}

		future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *lb)
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
	} else {
		client := meta.(*clients.Client).Network.LoadBalancerBackendAddressPoolsClient
		ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
		defer cancel()

		future, err := client.Delete(ctx, resGroup, loadBalancerName, name)
		if err != nil {
			return fmt.Errorf("Deleting LoadBalancer address pool [%s]: %+v", name, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Waiting for the completion for the LoadBalancer: %+v", err)
		}
	}

	d.SetId("")
	return nil
}

func readArmLoadBalancerBackendIPAddressPool(d *schema.ResourceData, meta interface{}, dataSourceMode bool) error {
	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Load Balancer Name and Group: %+v", err)
	}

	var name string
	if dataSourceMode {
		name = d.Get("name").(string)
	} else {
		id, err := azure.ParseAzureResourceID(d.Id())
		if err != nil {
			return err
		}
		name = id.Path["backendAddressPools"]
		resGroup = id.ResourceGroup
	}

	lb, exists, err := retrieveLoadBalancerById(d, d.Get("loadbalancer_id").(string), meta)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	var backendIPConfigurations []string
	var loadBalancingRules []string
	var backendIPAddressPool []interface{}

	if strings.EqualFold(string(lb.Sku.Name), string(network.LoadBalancerSkuNameBasic)) {
		config, _, exists := FindLoadBalancerBackEndAddressPoolByName(lb, name)
		if !exists {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Backend Address Pool %q not found. Removing from state", name)
			return nil
		}

		if dataSourceMode {
			d.SetId(*config.ID)
		}
		d.Set("name", config.Name)

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
	} else {
		client := meta.(*clients.Client).Network.LoadBalancerBackendAddressPoolsClient
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		read, err := client.Get(ctx, resGroup, loadBalancerName, name)
		if err != nil {
			return fmt.Errorf("Retrieving Load Balancer backend address pool %s on %q (Resource Group %q): %+v", name, loadBalancerName, resGroup, err)
		}

		if dataSourceMode {
			d.SetId(*read.ID)
		}
		d.Set("name", read.Name)

		if props := read.BackendAddressPoolPropertiesFormat; props != nil {
			backendIPConfigurations = flattenArmLoadBalancerBackendIPConfigurations(props.BackendIPConfigurations)
			loadBalancingRules = flattenArmLoadBalancerLoadBalancingRules(props.LoadBalancingRules)
			backendIPAddressPool = flattenArmLoadBalancerBackendIPAddressPool(props.LoadBalancerBackendAddresses)
		}
	}

	if !dataSourceMode {
		d.Set("resource_group_name", resGroup)
	}
	d.Set("backend_ip_configurations", backendIPConfigurations)
	d.Set("ip_address", backendIPAddressPool)
	d.Set("load_balancing_rules", loadBalancingRules)
	return nil
}

func expandAzureRmLoadBalancerBackendAddressPools(d *schema.ResourceData) network.BackendAddressPool {
	return network.BackendAddressPool{
		Name: utils.String(d.Get("name").(string)),
	}
}

func flattenArmLoadBalancerBackendIPConfigurations(configs *[]network.InterfaceIPConfiguration) []string {
	backendIPConfigurations := make([]string, 0)
	if configs != nil {
		for _, backendConfig := range *configs {
			backendIPConfigurations = append(backendIPConfigurations, *backendConfig.ID)
		}
	}
	return backendIPConfigurations
}

func flattenArmLoadBalancerLoadBalancingRules(rules *[]network.SubResource) []string {
	loadBalancingRules := make([]string, 0)
	if rules != nil {
		for _, rule := range *rules {
			loadBalancingRules = append(loadBalancingRules, *rule.ID)
		}
	}
	return loadBalancingRules
}

func flattenArmLoadBalancerBackendIPAddressPool(loadBalancerBackendAddresses *[]network.LoadBalancerBackendAddress) []interface{} {
	backendIPAddresses := make([]interface{}, 0)
	if loadBalancerBackendAddresses != nil {
		for _, lbba := range *loadBalancerBackendAddresses {
			ipAddress := make(map[string]interface{})
			if name := lbba.Name; name != nil {
				ipAddress["name"] = name
			}
			if properties := lbba.LoadBalancerBackendAddressPropertiesFormat; properties != nil {
				if vnet := lbba.LoadBalancerBackendAddressPropertiesFormat.VirtualNetwork; vnet != nil {
					ipAddress["virtual_network_id"] = vnet.ID
				}
				if addr := lbba.LoadBalancerBackendAddressPropertiesFormat.IPAddress; addr != nil {
					ipAddress["ip_address"] = addr
				}
			}
			backendIPAddresses = append(backendIPAddresses, ipAddress)
		}
	}
	return backendIPAddresses
}

func expandArmLoadBalancerBackendIPAddressPool(input []interface{}) *[]network.LoadBalancerBackendAddress {
	output := make([]network.LoadBalancerBackendAddress, 0)

	for _, item := range input {
		vals := item.(map[string]interface{})

		var name *string
		if v, ok := vals["name"]; ok {
			name = utils.String(v.(string))
		}
		vnetID := vals["virtual_network_id"].(string)
		ipAddress := vals["ip_address"].(string)

		output = append(output, network.LoadBalancerBackendAddress{
			&network.LoadBalancerBackendAddressPropertiesFormat{
				VirtualNetwork: &network.SubResource{
					ID: &vnetID,
				},
				IPAddress: &ipAddress,
			},
			name,
		})
	}

	return &output
}
