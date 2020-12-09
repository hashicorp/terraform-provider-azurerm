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

func resourceArmLoadBalancerNatRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerNatRuleCreateUpdate,
		Read:   resourceArmLoadBalancerNatRuleRead,
		Update: resourceArmLoadBalancerNatRuleCreateUpdate,
		Delete: resourceArmLoadBalancerNatRuleDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*parse.LoadBalancerId, error) {
			id, err := parse.LoadBalancerInboundNatRuleID(input)
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"enable_floating_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"enable_tcp_reset": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"idle_timeout_in_minutes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(4, 30),
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
	client := meta.(*clients.Client).Network.LoadBalancersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("retrieving Load Balancer Name and Group: %+v", err)
	}
	loadBalancerIdRaw := loadBalancerId.ID("")
	locks.ByID(loadBalancerIdRaw)
	defer locks.UnlockByID(loadBalancerIdRaw)

	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, *loadBalancerId)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	newNatRule, err := expandAzureRmLoadBalancerNatRule(d, loadBalancer, *loadBalancerId)
	if err != nil {
		return fmt.Errorf("Error Expanding NAT Rule: %+v", err)
	}

	natRules := append(*loadBalancer.LoadBalancerPropertiesFormat.InboundNatRules, *newNatRule)

	existingNatRule, existingNatRuleIndex, exists := FindLoadBalancerNatRuleByName(loadBalancer, name)
	if exists {
		if name == *existingNatRule.Name {
			if d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_lb_nat_rule", *existingNatRule.ID)
			}

			// this nat rule is being updated/reapplied remove old copy from the slice
			natRules = append(natRules[:existingNatRuleIndex], natRules[existingNatRuleIndex+1:]...)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.InboundNatRules = &natRules

	future, err := client.CreateOrUpdate(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating / Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Load Balancer %q (Resource Group %q): %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, err)
	}

	read, err := client.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer %q (Resource Group %q): %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (Resource Group %q) ID", loadBalancerId.Name, loadBalancerId.ResourceGroup)
	}

	var natRuleId string
	for _, InboundNatRule := range *read.LoadBalancerPropertiesFormat.InboundNatRules {
		if *InboundNatRule.Name == name {
			natRuleId = *InboundNatRule.ID
		}
	}

	if natRuleId != "" {
		d.SetId(natRuleId)
	} else {
		return fmt.Errorf("Cannot find created Load Balancer NAT Rule ID %q", natRuleId)
	}

	return resourceArmLoadBalancerNatRuleRead(d, meta)
}

func resourceArmLoadBalancerNatRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerInboundNatRuleID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, loadBalancerId)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
		return nil
	}

	config, _, exists := FindLoadBalancerNatRuleByName(loadBalancer, id.InboundNatRuleName)
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer Nat Rule %q not found. Removing from state", id.InboundNatRuleName)
		return nil
	}

	d.Set("name", config.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := config.InboundNatRulePropertiesFormat; props != nil {
		backendIPConfigId := ""
		if props.BackendIPConfiguration != nil && props.BackendIPConfiguration.ID != nil {
			backendIPConfigId = *props.BackendIPConfiguration.ID
		}
		d.Set("backend_ip_configuration_id", backendIPConfigId)

		backendPort := 0
		if props.BackendPort != nil {
			backendPort = int(*props.BackendPort)
		}
		d.Set("backend_port", backendPort)
		d.Set("enable_floating_ip", props.EnableFloatingIP)
		d.Set("enable_tcp_reset", props.EnableTCPReset)

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
		d.Set("frontend_ip_configuration_name", frontendIPConfigName)
		d.Set("frontend_ip_configuration_id", frontendIPConfigID)

		frontendPort := 0
		if props.FrontendPort != nil {
			frontendPort = int(*props.FrontendPort)
		}
		d.Set("frontend_port", frontendPort)

		idleTimeoutInMinutes := 0
		if props.IdleTimeoutInMinutes != nil {
			idleTimeoutInMinutes = int(*props.IdleTimeoutInMinutes)
		}
		d.Set("idle_timeout_in_minutes", idleTimeoutInMinutes)
		d.Set("protocol", string(props.Protocol))
	}

	return nil
}

func resourceArmLoadBalancerNatRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerInboundNatRuleID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
	loadBalancerID := loadBalancerId.ID("")
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, loadBalancerId)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		return nil
	}

	_, index, exists := FindLoadBalancerNatRuleByName(loadBalancer, id.InboundNatRuleName)
	if !exists {
		return nil
	}

	oldNatRules := *loadBalancer.LoadBalancerPropertiesFormat.InboundNatRules
	newNatRules := append(oldNatRules[:index], oldNatRules[index+1:]...)
	loadBalancer.LoadBalancerPropertiesFormat.InboundNatRules = &newNatRules

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q) %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of Load Balancer updates for %q (Resource Group %q) %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (resource group %q) ID", id.LoadBalancerName, id.ResourceGroup)
	}

	return nil
}

func expandAzureRmLoadBalancerNatRule(d *schema.ResourceData, lb *network.LoadBalancer, loadBalancerId parse.LoadBalancerId) (*network.InboundNatRule, error) {
	properties := network.InboundNatRulePropertiesFormat{
		Protocol:       network.TransportProtocol(d.Get("protocol").(string)),
		FrontendPort:   utils.Int32(int32(d.Get("frontend_port").(int))),
		BackendPort:    utils.Int32(int32(d.Get("backend_port").(int))),
		EnableTCPReset: utils.Bool(d.Get("enable_tcp_reset").(bool)),
	}

	if v, ok := d.GetOk("enable_floating_ip"); ok {
		properties.EnableFloatingIP = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("idle_timeout_in_minutes"); ok {
		properties.IdleTimeoutInMinutes = utils.Int32(int32(v.(int)))
	}

	if v := d.Get("frontend_ip_configuration_name").(string); v != "" {
		if _, exists := FindLoadBalancerFrontEndIpConfigurationByName(lb, v); !exists {
			return nil, fmt.Errorf("[ERROR] Cannot find FrontEnd IP Configuration with the name %s", v)
		}

		id := parse.NewLoadBalancerFrontendIpConfigurationID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroup, loadBalancerId.Name, v).ID("")
		properties.FrontendIPConfiguration = &network.SubResource{
			ID: utils.String(id),
		}
	}

	natRule := network.InboundNatRule{
		Name:                           utils.String(d.Get("name").(string)),
		InboundNatRulePropertiesFormat: &properties,
	}

	return &natRule, nil
}
