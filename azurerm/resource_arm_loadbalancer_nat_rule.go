package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancerNatRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerNatRuleCreateUpdate,
		Read:   resourceArmLoadBalancerNatRuleRead,
		Update: resourceArmLoadBalancerNatRuleCreateUpdate,
		Delete: resourceArmLoadBalancerNatRuleDelete,
		Importer: &schema.ResourceImporter{
			State: loadBalancerSubResourceStateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"location": deprecatedLocationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

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

			"enable_floating_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"frontend_port": {
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
				ValidateFunc: validation.NoZeroValues,
			},

			"frontend_ip_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"backend_ip_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmLoadBalancerNatRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).loadBalancerClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	loadBalancerID := d.Get("loadbalancer_id").(string)

	armMutexKV.Lock(loadBalancerID)
	defer armMutexKV.Unlock(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	props := loadBalancer.LoadBalancerPropertiesFormat
	if props == nil {
		return fmt.Errorf("Error updating Load Inbound NAT Rule: props was nil")
	}

	rules := props.InboundNatRules
	if rules == nil {
		return fmt.Errorf("Error updating Load Balancer Inbound NAT Rule: props.InboundNatRules was nil")
	}

	newNatRule, err := expandAzureRmLoadBalancerNatRule(d, loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Expanding NAT Rule: %+v", err)
	}

	natRules := *rules
	if d.IsNewResource() {
		for _, rule := range natRules {
			if rule.Name != nil && *rule.Name == name {
				return tf.ImportAsExistsError("azurerm_lb_nat_rule", *rule.ID)
			}
		}

		natRules = append(natRules, *newNatRule)
	} else {
		// find and replace
		index := -1
		for i, v := range natRules {
			if v.Name != nil && *v.Name == name {
				index = i
				break
			}
		}

		if index == -1 {
			// note: this should be caught by the read
			return fmt.Errorf("Error: Inbound Nat Rule %q was not found on Load Balancer %q", name, loadBalancerID)
		}

		natRules[index] = *newNatRule
	}

	loadBalancer.LoadBalancerPropertiesFormat.InboundNatRules = &natRules
	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(loadBalancerID)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating / Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(tf.TimeoutForCreateUpdate(d)))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (Resource Group %q) ID", loadBalancerName, resGroup)
	}

	var natRuleId string
	if props := read.LoadBalancerPropertiesFormat; props != nil {
		if rules := props.InboundNatRules; rules != nil {
			for _, rule := range *rules {
				if rule.Name != nil && *rule.Name == name {
					natRuleId = *rule.ID
				}
			}
		}
	}

	if natRuleId == "" {
		return fmt.Errorf("Cannot find created Load Balancer NAT Rule ID %q", natRuleId)
	}

	d.SetId(natRuleId)

	return resourceArmLoadBalancerNatRuleRead(d, meta)
}

func resourceArmLoadBalancerNatRuleRead(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["inboundNatRules"]
	loadBalancerId := d.Get("loadbalancer_id").(string)
	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerId, meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	var rule *network.InboundNatRule
	if props := loadBalancer.LoadBalancerPropertiesFormat; props != nil {
		if rules := props.InboundNatRules; rules != nil {
			for _, r := range *rules {
				if r.Name != nil && *r.Name == name {
					rule = &r
					break
				}
			}
		}
	}

	if rule == nil {
		d.SetId("")
		log.Printf("[INFO] Load Balancer Nat Rule %q not found. Removing from state", name)
		return nil
	}

	d.Set("name", name)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := rule.InboundNatRulePropertiesFormat; props != nil {
		d.Set("protocol", props.Protocol)
		d.Set("frontend_port", props.FrontendPort)
		d.Set("backend_port", props.BackendPort)
		d.Set("enable_floating_ip", props.EnableFloatingIP)

		if ipconfiguration := props.FrontendIPConfiguration; ipconfiguration != nil {
			fipID, err := parseAzureResourceID(*ipconfiguration.ID)
			if err != nil {
				return err
			}

			d.Set("frontend_ip_configuration_name", fipID.Path["frontendIPConfigurations"])
			d.Set("frontend_ip_configuration_id", ipconfiguration.ID)
		}

		if ipconfiguration := props.BackendIPConfiguration; ipconfiguration != nil {
			d.Set("backend_ip_configuration_id", ipconfiguration.ID)
		}
	}

	return nil
}

func resourceArmLoadBalancerNatRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).loadBalancerClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	loadBalancerID := d.Get("loadbalancer_id").(string)

	armMutexKV.Lock(loadBalancerID)
	defer armMutexKV.Unlock(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		return nil
	}

	props := loadBalancer.LoadBalancerPropertiesFormat
	if props == nil {
		return nil
	}

	rules := props.InboundNatRules
	if rules == nil {
		return nil
	}

	natRules := make([]network.InboundNatRule, 0)
	for _, rule := range *rules {
		if rule.Name != nil && *rule.Name != name {
			natRules = append(natRules, rule)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.InboundNatRules = &natRules

	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(loadBalancerID)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Updating Load Balancer %q (Resource Group %q) %+v", loadBalancerName, resGroup, err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the completion of Load Balancer updates for %q (Resource Group %q) %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (resource group %q) ID", loadBalancerName, resGroup)
	}

	return nil
}

func expandAzureRmLoadBalancerNatRule(d *schema.ResourceData, lb *network.LoadBalancer) (*network.InboundNatRule, error) {
	properties := network.InboundNatRulePropertiesFormat{
		Protocol:     network.TransportProtocol(d.Get("protocol").(string)),
		FrontendPort: utils.Int32(int32(d.Get("frontend_port").(int))),
		BackendPort:  utils.Int32(int32(d.Get("backend_port").(int))),
	}

	if v, ok := d.GetOk("enable_floating_ip"); ok {
		properties.EnableFloatingIP = utils.Bool(v.(bool))
	}

	if v := d.Get("frontend_ip_configuration_name").(string); v != "" {
		rule, _, exists := findLoadBalancerFrontEndIpConfigurationByName(lb, v)
		if !exists {
			return nil, fmt.Errorf("[ERROR] Cannot find FrontEnd IP Configuration with the name %s", v)
		}

		properties.FrontendIPConfiguration = &network.SubResource{
			ID: rule.ID,
		}
	}

	natRule := network.InboundNatRule{
		Name: utils.String(d.Get("name").(string)),
		InboundNatRulePropertiesFormat: &properties,
	}

	return &natRule, nil
}
