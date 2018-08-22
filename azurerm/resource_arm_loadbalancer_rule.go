package azurerm

import (
	"context"
	"fmt"
	"log"
	"regexp"
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

func resourceArmLoadBalancerRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerRuleCreateUpdate,
		Read:   resourceArmLoadBalancerRuleRead,
		Update: resourceArmLoadBalancerRuleCreateUpdate,
		Delete: resourceArmLoadBalancerRuleDelete,
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
				ValidateFunc: validateArmLoadBalancerRuleName,
			},

			"location": deprecatedLocationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
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

			"backend_address_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

			"probe_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enable_floating_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"idle_timeout_in_minutes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(4, 30),
			},

			"load_distribution": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArmLoadBalancerRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error updating Load Balancer Rules: props was nil")
	}

	rules := props.LoadBalancingRules
	if rules == nil {
		return fmt.Errorf("Error updating Load Balancer Rules: props.LoadBalancingRules was nil")
	}
	lbRules := *rules

	newLbRule, err := expandAzureRmLoadBalancerRule(d, loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Exanding Load Balancer Rule: %+v", err)
	}

	if d.IsNewResource() {
		// check if it requires import
		for _, rule := range lbRules {
			if rule.Name != nil && *rule.Name == name {
				return tf.ImportAsExistsError("azurerm_lb_rule", *rule.ID)
			}
		}

		// otherwise just append it
		lbRules = append(lbRules, *newLbRule)
	} else {
		index := -1
		for i, v := range lbRules {
			if v.Name != nil && *v.Name == name {
				index = i
				break
			}
		}

		if index == -1 {
			// should be caught by the Read
		}

		lbRules[index] = *newLbRule
	}

	loadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules = &lbRules
	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(loadBalancerID)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating LoadBalancer: %+v", err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(tf.TimeoutForCreateUpdate(d)))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion for Load Balancer updates: %+v", err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error Getting LoadBalancer: %+v", err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %s (resource group %s) ID", loadBalancerName, resGroup)
	}

	var ruleId string
	if props := read.LoadBalancerPropertiesFormat; props != nil {
		if rules := props.LoadBalancingRules; rules != nil {
			for _, rule := range *rules {
				if rule.Name != nil && *rule.Name == name {
					ruleId = *rule.ID
				}
			}
		}
	}

	if ruleId == "" {
		return fmt.Errorf("Cannot find created Load Balancer Rule ID %q", ruleId)
	}

	d.SetId(ruleId)

	return resourceArmLoadBalancerRuleRead(d, meta)
}

func resourceArmLoadBalancerRuleRead(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["loadBalancingRules"]
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

	var rule *network.LoadBalancingRule
	if props := loadBalancer.LoadBalancerPropertiesFormat; props != nil {
		if rules := props.LoadBalancingRules; rules != nil {
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
		log.Printf("[INFO] Load Balancer Rule %q not found. Removing from state", name)
		return nil
	}

	d.Set("name", name)
	d.Set("resource_group_name", id.ResourceGroup)

	if properties := rule.LoadBalancingRulePropertiesFormat; properties != nil {
		d.Set("protocol", properties.Protocol)
		d.Set("frontend_port", properties.FrontendPort)
		d.Set("backend_port", properties.BackendPort)

		if properties.EnableFloatingIP != nil {
			d.Set("enable_floating_ip", properties.EnableFloatingIP)
		}

		if properties.IdleTimeoutInMinutes != nil {
			d.Set("idle_timeout_in_minutes", properties.IdleTimeoutInMinutes)
		}

		if properties.FrontendIPConfiguration != nil {
			fipID, err := parseAzureResourceID(*properties.FrontendIPConfiguration.ID)
			if err != nil {
				return err
			}

			d.Set("frontend_ip_configuration_name", fipID.Path["frontendIPConfigurations"])
			d.Set("frontend_ip_configuration_id", properties.FrontendIPConfiguration.ID)
		}

		if properties.BackendAddressPool != nil {
			d.Set("backend_address_pool_id", properties.BackendAddressPool.ID)
		}

		if properties.Probe != nil {
			d.Set("probe_id", properties.Probe.ID)
		}

		if properties.LoadDistribution != "" {
			d.Set("load_distribution", properties.LoadDistribution)
		}
	}

	return nil
}

func resourceArmLoadBalancerRuleDelete(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error updating Load Balancer Rules: props was nil")
	}

	rules := props.LoadBalancingRules
	if rules == nil {
		return fmt.Errorf("Error updating Load Balancer Rules: props.LoadBalancingRules was nil")
	}
	lbRules := *rules

	newRules := make([]network.LoadBalancingRule, 0)
	for _, rule := range lbRules {
		if rule.Name != nil && *rule.Name != name {
			newRules = append(newRules, rule)
		}
	}
	loadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules = &newRules

	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(loadBalancerID)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error Getting LoadBalancer: %+v", err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ID of Load Balancer %q (resource group %s)", loadBalancerName, resGroup)
	}

	return nil
}

func expandAzureRmLoadBalancerRule(d *schema.ResourceData, lb *network.LoadBalancer) (*network.LoadBalancingRule, error) {
	properties := network.LoadBalancingRulePropertiesFormat{
		Protocol:         network.TransportProtocol(d.Get("protocol").(string)),
		FrontendPort:     utils.Int32(int32(d.Get("frontend_port").(int))),
		BackendPort:      utils.Int32(int32(d.Get("backend_port").(int))),
		EnableFloatingIP: utils.Bool(d.Get("enable_floating_ip").(bool)),
	}

	if v, ok := d.GetOk("idle_timeout_in_minutes"); ok {
		properties.IdleTimeoutInMinutes = utils.Int32(int32(v.(int)))
	}

	if v := d.Get("load_distribution").(string); v != "" {
		properties.LoadDistribution = network.LoadDistribution(v)
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

	if v := d.Get("backend_address_pool_id").(string); v != "" {
		properties.BackendAddressPool = &network.SubResource{
			ID: &v,
		}
	}

	if v := d.Get("probe_id").(string); v != "" {
		properties.Probe = &network.SubResource{
			ID: &v,
		}
	}

	return &network.LoadBalancingRule{
		Name: utils.String(d.Get("name").(string)),
		LoadBalancingRulePropertiesFormat: &properties,
	}, nil
}

func validateArmLoadBalancerRuleName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z_0-9.-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only word characters, numbers, underscores, periods, and hyphens allowed in %q: %q",
			k, value))
	}

	if len(value) > 80 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 80 characters: %q", k, value))
	}

	if len(value) == 0 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be an empty string: %q", k, value))
	}
	if !regexp.MustCompile(`[a-zA-Z0-9_]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must end with a word character, number, or underscore: %q", k, value))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9]`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must start with a word character or number: %q", k, value))
	}

	return
}
