package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancerOutboundRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerOutboundRuleCreateUpdate,
		Read:   resourceArmLoadBalancerOutboundRuleRead,
		Update: resourceArmLoadBalancerOutboundRuleCreateUpdate,
		Delete: resourceArmLoadBalancerOutboundRuleDelete,

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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"frontend_ip_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"backend_address_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.TransportProtocolAll),
					string(network.TransportProtocolTCP),
					string(network.TransportProtocolUDP),
				}, false),
			},

			"enable_tcp_reset": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"allocated_outbound_ports": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1024,
			},

			"idle_timeout_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
			},
		},
	}
}

func resourceArmLoadBalancerOutboundRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.LoadBalancersClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	loadBalancerID := d.Get("loadbalancer_id").(string)
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

	newOutboundRule, err := expandAzureRmLoadBalancerOutboundRule(d, loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Exanding Load Balancer Rule: %+v", err)
	}

	outboundRules := make([]network.OutboundRule, 0)

	if loadBalancer.LoadBalancerPropertiesFormat.OutboundRules != nil {
		outboundRules = *loadBalancer.LoadBalancerPropertiesFormat.OutboundRules
	}

	existingOutboundRule, existingOutboundRuleIndex, exists := findLoadBalancerOutboundRuleByName(loadBalancer, name)
	if exists {
		if name == *existingOutboundRule.Name {
			if features.ShouldResourcesBeImported() && d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_lb_outbound_rule", *existingOutboundRule.ID)
			}

			// this outbound rule is being updated/reapplied remove old copy from the slice
			outboundRules = append(outboundRules[:existingOutboundRuleIndex], outboundRules[existingOutboundRuleIndex+1:]...)
		}
	}

	outboundRules = append(outboundRules, *newOutboundRule)

	loadBalancer.LoadBalancerPropertiesFormat.OutboundRules = &outboundRules
	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(loadBalancerID)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating LoadBalancer: %+v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion for Load Balancer updates: %+v", err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error Getting LoadBalancer: %+v", err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %s (resource group %s) ID", loadBalancerName, resGroup)
	}

	var outboundRuleId string
	for _, OutboundRule := range *(*read.LoadBalancerPropertiesFormat).OutboundRules {
		if *OutboundRule.Name == name {
			outboundRuleId = *OutboundRule.ID
		}
	}

	if outboundRuleId == "" {
		return fmt.Errorf("Cannot find created Load Balancer Outbound Rule ID %q", outboundRuleId)
	}

	d.SetId(outboundRuleId)

	return resourceArmLoadBalancerOutboundRuleRead(d, meta)
}

func resourceArmLoadBalancerOutboundRuleRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["outboundRules"]

	loadBalancer, exists, err := retrieveLoadBalancerById(d.Get("loadbalancer_id").(string), meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	config, _, exists := findLoadBalancerOutboundRuleByName(loadBalancer, name)
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer Outbound Rule %q not found. Removing from state", name)
		return nil
	}

	d.Set("name", config.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if properties := config.OutboundRulePropertiesFormat; properties != nil {
		d.Set("protocol", properties.Protocol)
		d.Set("backend_address_pool_id", properties.BackendAddressPool.ID)

		frontendIpConfigurations := make([]interface{}, 0)
		for _, feConfig := range *properties.FrontendIPConfigurations {
			if feConfig.ID == nil {
				continue
			}

			feConfigId, err := azure.ParseAzureResourceID(*feConfig.ID)
			if err != nil {
				return nil
			}

			name := feConfigId.Path["frontendIPConfigurations"]
			frontendConfiguration := map[string]interface{}{
				"id":   *feConfig.ID,
				"name": name,
			}
			frontendIpConfigurations = append(frontendIpConfigurations, frontendConfiguration)
		}
		d.Set("frontend_ip_configuration", frontendIpConfigurations)

		if properties.EnableTCPReset != nil {
			d.Set("enable_tcp_reset", properties.EnableTCPReset)
		}

		if properties.IdleTimeoutInMinutes != nil {
			d.Set("idle_timeout_in_minutes", properties.IdleTimeoutInMinutes)
		}

		if properties.AllocatedOutboundPorts != nil {
			d.Set("allocated_outbound_ports", properties.AllocatedOutboundPorts)
		}

	}

	return nil
}

func resourceArmLoadBalancerOutboundRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.LoadBalancersClient
	ctx := meta.(*ArmClient).StopContext

	loadBalancerID := d.Get("loadbalancer_id").(string)
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		return nil
	}

	_, index, exists := findLoadBalancerOutboundRuleByName(loadBalancer, d.Get("name").(string))
	if !exists {
		return nil
	}

	oldOutboundRules := *loadBalancer.LoadBalancerPropertiesFormat.OutboundRules
	newOutboundRules := append(oldOutboundRules[:index], oldOutboundRules[index+1:]...)
	loadBalancer.LoadBalancerPropertiesFormat.OutboundRules = &newOutboundRules

	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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

func expandAzureRmLoadBalancerOutboundRule(d *schema.ResourceData, lb *network.LoadBalancer) (*network.OutboundRule, error) {

	properties := network.OutboundRulePropertiesFormat{
		Protocol: network.LoadBalancerOutboundRuleProtocol(d.Get("protocol").(string)),
	}

	feConfigs := d.Get("frontend_ip_configuration").([]interface{})
	feConfigSubResources := make([]network.SubResource, 0)

	for _, raw := range feConfigs {
		v := raw.(map[string]interface{})
		rule, exists := findLoadBalancerFrontEndIpConfigurationByName(lb, v["name"].(string))
		if !exists {
			return nil, fmt.Errorf("[ERROR] Cannot find FrontEnd IP Configuration with the name %s", v["name"])
		}

		feConfigSubResource := network.SubResource{
			ID: rule.ID,
		}

		feConfigSubResources = append(feConfigSubResources, feConfigSubResource)
	}

	properties.FrontendIPConfigurations = &feConfigSubResources

	if v := d.Get("backend_address_pool_id").(string); v != "" {
		properties.BackendAddressPool = &network.SubResource{
			ID: &v,
		}
	}

	if v, ok := d.GetOk("idle_timeout_in_minutes"); ok {
		properties.IdleTimeoutInMinutes = utils.Int32(int32(v.(int)))
	}

	if v, ok := d.GetOk("enable_tcp_reset"); ok {
		properties.EnableTCPReset = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("allocated_outbound_ports"); ok {
		properties.AllocatedOutboundPorts = utils.Int32(int32(v.(int)))
	}

	return &network.OutboundRule{
		Name:                         utils.String(d.Get("name").(string)),
		OutboundRulePropertiesFormat: &properties,
	}, nil
}
