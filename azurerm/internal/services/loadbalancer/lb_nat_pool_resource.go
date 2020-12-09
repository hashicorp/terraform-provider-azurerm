package loadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/state"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancerNatPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerNatPoolCreateUpdate,
		Read:   resourceArmLoadBalancerNatPoolRead,
		Update: resourceArmLoadBalancerNatPoolCreateUpdate,
		Delete: resourceArmLoadBalancerNatPoolDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*parse.LoadBalancerId, error) {
			id, err := parse.LoadBalancerInboundNatPoolID(input)
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
				ValidateFunc: networkValidate.LoadBalancerID,
			},

			"protocol": {
				Type:             schema.TypeString,
				Required:         true,
				StateFunc:        state.IgnoreCase,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.TransportProtocolAll),
					string(network.TransportProtocolTCP),
					string(network.TransportProtocolUDP),
				}, true),
			},

			"frontend_port_start": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"frontend_port_end": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"backend_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"frontend_ip_configuration_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"frontend_ip_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmLoadBalancerNatPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LoadBalancersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Load Balancer Name and Group: %+v", err)
	}

	loadBalancerID := loadBalancerId.ID("")
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, *loadBalancerId)
	if err != nil {
		return fmt.Errorf("retrieving Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", loadBalancerId.Name)
		return nil
	}

	newNatPool, err := expandAzureRmLoadBalancerNatPool(d, loadBalancer)
	if err != nil {
		return fmt.Errorf("expanding NAT Pool: %+v", err)
	}

	natPools := append(*loadBalancer.LoadBalancerPropertiesFormat.InboundNatPools, *newNatPool)

	existingNatPool, existingNatPoolIndex, exists := FindLoadBalancerNatPoolByName(loadBalancer, name)
	if exists {
		if name == *existingNatPool.Name {
			if d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_lb_nat_pool", *existingNatPool.ID)
			}

			// this probe is being updated/reapplied remove old copy from the slice
			natPools = append(natPools[:existingNatPoolIndex], natPools[existingNatPoolIndex+1:]...)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.InboundNatPools = &natPools

	future, err := client.CreateOrUpdate(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of Load Balancer %q (Resource Group %q): %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, err)
	}

	read, err := client.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer %q (Resource Group %q): %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (Resource Group %q) ID", loadBalancerId.Name, loadBalancerId.ResourceGroup)
	}

	var natPoolId string
	for _, InboundNatPool := range *read.LoadBalancerPropertiesFormat.InboundNatPools {
		if *InboundNatPool.Name == name {
			natPoolId = *InboundNatPool.ID
		}
	}

	if natPoolId == "" {
		return fmt.Errorf("Cannot find created Load Balancer NAT Pool ID %q", natPoolId)
	}

	d.SetId(natPoolId)

	return resourceArmLoadBalancerNatPoolRead(d, meta)
}

func resourceArmLoadBalancerNatPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerInboundNatPoolID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, loadBalancerId)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
		return nil
	}

	config, _, exists := FindLoadBalancerNatPoolByName(loadBalancer, id.InboundNatPoolName)
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer Nat Pool %q not found. Removing from state", id.InboundNatPoolName)
		return nil
	}

	d.Set("name", config.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := config.InboundNatPoolPropertiesFormat; props != nil {
		backendPort := 0
		if props.BackendPort != nil {
			backendPort = int(*props.BackendPort)
		}
		d.Set("backend_port", backendPort)

		frontendIPConfigName := ""
		frontendIPConfigID := ""
		if props.FrontendIPConfiguration != nil && props.FrontendIPConfiguration.ID != nil {
			feid, err := parse.LoadBalancerFrontendIpConfigurationID(*props.FrontendIPConfiguration.ID)
			if err != nil {
				return err
			}

			frontendIPConfigName = feid.FrontendIPConfigurationName
			frontendIPConfigID = feid.ID("")
		}
		d.Set("frontend_ip_configuration_id", frontendIPConfigID)
		d.Set("frontend_ip_configuration_name", frontendIPConfigName)

		frontendPortRangeEnd := 0
		if props.FrontendPortRangeEnd != nil {
			frontendPortRangeEnd = int(*props.FrontendPortRangeEnd)
		}
		d.Set("frontend_port_end", frontendPortRangeEnd)

		frontendPortRangeStart := 0
		if props.FrontendPortRangeStart != nil {
			frontendPortRangeStart = int(*props.FrontendPortRangeStart)
		}
		d.Set("frontend_port_start", frontendPortRangeStart)
		d.Set("protocol", string(props.Protocol))
	}

	return nil
}

func resourceArmLoadBalancerNatPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerInboundNatPoolID(d.Id())
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

	_, index, exists := FindLoadBalancerNatPoolByName(loadBalancer, id.InboundNatPoolName)
	if !exists {
		return nil
	}

	oldNatPools := *loadBalancer.LoadBalancerPropertiesFormat.InboundNatPools
	newNatPools := append(oldNatPools[:index], oldNatPools[index+1:]...)
	loadBalancer.LoadBalancerPropertiesFormat.InboundNatPools = &newNatPools

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error creating/updating Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of the Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer: %+v", err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (Resource Group %q) ID", id.LoadBalancerName, id.ResourceGroup)
	}

	return nil
}

func expandAzureRmLoadBalancerNatPool(d *schema.ResourceData, lb *network.LoadBalancer) (*network.InboundNatPool, error) {
	properties := network.InboundNatPoolPropertiesFormat{
		Protocol:               network.TransportProtocol(d.Get("protocol").(string)),
		FrontendPortRangeStart: utils.Int32(int32(d.Get("frontend_port_start").(int))),
		FrontendPortRangeEnd:   utils.Int32(int32(d.Get("frontend_port_end").(int))),
		BackendPort:            utils.Int32(int32(d.Get("backend_port").(int))),
	}

	if v := d.Get("frontend_ip_configuration_name").(string); v != "" {
		rule, exists := FindLoadBalancerFrontEndIpConfigurationByName(lb, v)
		if !exists {
			return nil, fmt.Errorf("[ERROR] Cannot find FrontEnd IP Configuration with the name %s", v)
		}

		properties.FrontendIPConfiguration = &network.SubResource{
			ID: rule.ID,
		}
	}

	return &network.InboundNatPool{
		Name:                           utils.String(d.Get("name").(string)),
		InboundNatPoolPropertiesFormat: &properties,
	}, nil
}
