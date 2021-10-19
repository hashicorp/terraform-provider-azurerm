package loadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	loadBalancerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/state"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmLoadBalancerNatRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"loadbalancer_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: loadBalancerValidate.LoadBalancerID,
			},

			"protocol": {
				Type:             pluginsdk.TypeString,
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
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"backend_port": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"frontend_ip_configuration_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"enable_floating_ip": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"enable_tcp_reset": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"idle_timeout_in_minutes": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(4, 30),
			},

			"frontend_ip_configuration_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"backend_ip_configuration_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmLoadBalancerNatRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerInboundNatRulesClient
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("retrieving Load Balancer Name and Group: %+v", err)
	}
	id := parse.NewLoadBalancerInboundNatRuleID(subscriptionId, loadBalancerId.ResourceGroup, loadBalancerId.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.InboundNatRuleName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Load Balancer Inbound Nat Rule %q: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_lb_nat_rule", id.ID())
		}
	}

	loadBalancer, err := lbClient.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(loadBalancer.Response) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer %q (resource group %q) for Nat Rule %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.InboundNatRuleName, err)
	}

	newNatRule, err := expandAzureRmLoadBalancerNatRule(d, &loadBalancer, *loadBalancerId)
	if err != nil {
		return fmt.Errorf("expanding NAT Rule: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.InboundNatRuleName, *newNatRule)
	if err != nil {
		return fmt.Errorf("updating Load Balancer %q (Resource Group %q) for Nat Rule %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.InboundNatRuleName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Load Balancer %q (Resource Group %q) for Nat Rule %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.InboundNatRuleName, err)
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerNatRuleRead(d, meta)
}

func resourceArmLoadBalancerNatRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerInboundNatRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerInboundNatRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.InboundNatRuleName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Inbound Nat Rule %q not found. Removing from state", id)
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer %q (resource group %q) for Nat Rule %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.InboundNatRuleName, err)
	}

	d.Set("name", id.InboundNatRuleName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.InboundNatRulePropertiesFormat; props != nil {
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
			frontendIPConfigID = feid.ID()
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

func resourceArmLoadBalancerNatRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerInboundNatRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerInboundNatRuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.LoadBalancerName, id.InboundNatRuleName)
	if err != nil {
		return fmt.Errorf("deleting Load Balancer Nat Rule %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of deleting Load Balancer Nat Rule %q: %+v", id, err)
	}

	return nil
}

func expandAzureRmLoadBalancerNatRule(d *pluginsdk.ResourceData, lb *network.LoadBalancer, loadBalancerId parse.LoadBalancerId) (*network.InboundNatRule, error) {
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

		id := parse.NewLoadBalancerFrontendIpConfigurationID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroup, loadBalancerId.Name, v).ID()
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
