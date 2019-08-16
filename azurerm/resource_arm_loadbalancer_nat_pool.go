package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancerNatPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerNatPoolCreateUpdate,
		Read:   resourceArmLoadBalancerNatPoolRead,
		Update: resourceArmLoadBalancerNatPoolCreateUpdate,
		Delete: resourceArmLoadBalancerNatPoolDelete,
		Importer: &schema.ResourceImporter{
			State: loadBalancerSubResourceStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationDeprecated(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"protocol": {
				Type:             schema.TypeString,
				Required:         true,
				StateFunc:        ignoreCaseStateFunc,
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
				ValidateFunc: validate.NoEmptyStrings,
			},

			"frontend_ip_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmLoadBalancerNatPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.LoadBalancersClient
	ctx := meta.(*ArmClient).StopContext

	loadBalancerID := d.Get("loadbalancer_id").(string)
	name := d.Get("name").(string)
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	newNatPool, err := expandAzureRmLoadBalancerNatPool(d, loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Expanding NAT Pool: %+v", err)
	}

	natPools := append(*loadBalancer.LoadBalancerPropertiesFormat.InboundNatPools, *newNatPool)

	existingNatPool, existingNatPoolIndex, exists := findLoadBalancerNatPoolByName(loadBalancer, name)
	if exists {
		if name == *existingNatPool.Name {
			if requireResourcesToBeImported && d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_lb_nat_pool", *existingNatPool.ID)
			}

			// this probe is being updated/reapplied remove old copy from the slice
			natPools = append(natPools[:existingNatPoolIndex], natPools[existingNatPoolIndex+1:]...)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.InboundNatPools = &natPools
	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(loadBalancerID)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (Resource Group %q) ID", loadBalancerName, resGroup)
	}

	var natPoolId string
	for _, InboundNatPool := range *(*read.LoadBalancerPropertiesFormat).InboundNatPools {
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
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["inboundNatPools"]

	loadBalancer, exists, err := retrieveLoadBalancerById(d.Get("loadbalancer_id").(string), meta)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	config, _, exists := findLoadBalancerNatPoolByName(loadBalancer, name)
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer Nat Pool %q not found. Removing from state", name)
		return nil
	}

	d.Set("name", config.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := config.InboundNatPoolPropertiesFormat; props != nil {
		d.Set("protocol", props.Protocol)
		d.Set("frontend_port_start", props.FrontendPortRangeStart)
		d.Set("frontend_port_end", props.FrontendPortRangeEnd)
		d.Set("backend_port", props.BackendPort)

		if feipConfig := props.FrontendIPConfiguration; feipConfig != nil {
			fipID, err := azure.ParseAzureResourceID(*feipConfig.ID)
			if err != nil {
				return err
			}

			d.Set("frontend_ip_configuration_name", fipID.Path["frontendIPConfigurations"])
			d.Set("frontend_ip_configuration_id", feipConfig.ID)
		}
	}

	return nil
}

func resourceArmLoadBalancerNatPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.LoadBalancersClient
	ctx := meta.(*ArmClient).StopContext

	loadBalancerID := d.Get("loadbalancer_id").(string)
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		return nil
	}

	_, index, exists := findLoadBalancerNatPoolByName(loadBalancer, d.Get("name").(string))
	if !exists {
		return nil
	}

	oldNatPools := *loadBalancer.LoadBalancerPropertiesFormat.InboundNatPools
	newNatPools := append(oldNatPools[:index], oldNatPools[index+1:]...)
	loadBalancer.LoadBalancerPropertiesFormat.InboundNatPools = &newNatPools

	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error creating/updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of the Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer: %+v", err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (Resource Group %q) ID", loadBalancerName, resGroup)
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
		rule, exists := findLoadBalancerFrontEndIpConfigurationByName(lb, v)
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
