package servicebus

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/servicebus/mgmt/2018-01-01-preview/servicebus"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceBusNamespaceNetworkRuleSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusNamespaceNetworkRuleSetCreateUpdate,
		Read:   resourceArmServiceBusNamespaceNetworkRuleSetRead,
		Update: resourceArmServiceBusNamespaceNetworkRuleSetCreateUpdate,
		Delete: resourceArmServiceBusNamespaceNetworkRuleSetDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"resource_group_name": azure.SchemaResourceGroupName(),

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceBusNamespaceName(),
			},

			"properties": {
				Type:       schema.TypeList,
				Optional:   true,
				MaxItems:   1,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(servicebus.Allow),
								string(servicebus.Deny),
							}, false),
						},

						"virtual_network_rule": {
							Type:       schema.TypeList,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// the API returns the subnet ID's resource group name in lowercase
									// https://github.com/Azure/azure-sdk-for-go/issues/5855
									"subnet_id": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateFunc:     azure.ValidateResourceID,
										DiffSuppressFunc: suppress.CaseDifference,
									},

									"ignore_missing_virtual_network_service_endpoint": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},

						"ip_rule": {
							Type:       schema.TypeList,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_mask": {
										Type:     schema.TypeString,
										Required: true,
									},

									"action": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  string(servicebus.NetworkRuleIPActionAllow),
										ValidateFunc: validation.StringInSlice([]string{
											string(servicebus.NetworkRuleIPActionAllow),
										}, false),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmServiceBusNamespaceNetworkRuleSetCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClientPreview
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM ServiceBus Namespace Authorization Rule creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	namespaceName := d.Get("namespace_name").(string)
	ruleSets, hasRuleSets := d.GetOk("properties")

	resp, err := client.Get(ctx, resourceGroup, namespaceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure ServiceBus Namespace %q: %+v", namespaceName, err)
	}

	sku := resp.Sku

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetNetworkRuleSet(ctx, resourceGroup, namespaceName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing ServiceBus Namespace Network Rule (Resource Group %q): %+v", resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_servicebus_namespace_network_rule_set", *existing.ID)
		}
	}

	if hasRuleSets {
		rulesets := servicebus.NetworkRuleSet{
			NetworkRuleSetProperties: expandServiceBusNetworkRuleset(ruleSets.([]interface{})),
		}

		// cannot use network rulesets with the basic SKU
		if sku.Name != servicebus.Basic && sku.Name != servicebus.Standard {
			if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, resourceGroup, namespaceName, rulesets); err != nil {
				return fmt.Errorf("Error setting network ruleset properties for Service Bus %q (resource group %q): %v", namespaceName, resourceGroup, err)
			}
		} else {
			// so if the user has specified the non default rule sets throw a validation error
			if rulesets.DefaultAction != servicebus.Deny ||
				(rulesets.IPRules != nil && len(*rulesets.IPRules) > 0) ||
				(rulesets.VirtualNetworkRules != nil && len(*rulesets.VirtualNetworkRules) > 0) {
				return fmt.Errorf("network_rulesets cannot be used when the SKU is basic or Standard")
			}
		}
	}

	read, err := client.GetNetworkRuleSet(ctx, resourceGroup, namespaceName)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ServiceBus Namespace Network Rule Set (resource group %s)", resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceBusNamespaceNetworkRuleSetRead(d, meta)
}

func resourceArmServiceBusNamespaceNetworkRuleSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClientPreview
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]

	resp, err := client.GetNetworkRuleSet(ctx, resGroup, namespaceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure ServiceBus Namespace Network Rule Set: %+v", err)
	}

	d.Set("namespace_name", namespaceName)
	d.Set("resource_group_name", resGroup)

	if err := d.Set("properties", flattenServiceBusNetworkRuleset(resp)); err != nil {
		return fmt.Errorf("Error setting `properties` for Service Bus Network Rule Set: %v", err)
	}

	return nil
}

func resourceArmServiceBusNamespaceNetworkRuleSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClientPreview
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM ServiceBus Namespace Authorization Rule creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	namespaceName := d.Get("namespace_name").(string)

	resp, err := client.Get(ctx, resourceGroup, namespaceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure ServiceBus Namespace %q: %+v", namespaceName, err)
	}

	sku := resp.Sku

	rulesetproperties := servicebus.NetworkRuleSetProperties{
		DefaultAction: servicebus.DefaultAction("Allow"),
	}
	rulesets := servicebus.NetworkRuleSet{
		NetworkRuleSetProperties: &rulesetproperties,
	}

	// cannot use network rulesets with the basic SKU
	if sku.Name != servicebus.Basic && sku.Name != servicebus.Standard {
		if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, resourceGroup, namespaceName, rulesets); err != nil {
			return fmt.Errorf("Error setting network ruleset properties for Service Bus %q (resource group %q): %v", namespaceName, resourceGroup, err)
		}
	} else {
		// so if the user has specified the non default rule sets throw a validation error
		if rulesets.DefaultAction != servicebus.Deny ||
			(rulesets.IPRules != nil && len(*rulesets.IPRules) > 0) ||
			(rulesets.VirtualNetworkRules != nil && len(*rulesets.VirtualNetworkRules) > 0) {
			return fmt.Errorf("network_rulesets cannot be used when the SKU is basic or Standard")
		}
	}

	read, err := client.GetNetworkRuleSet(ctx, resourceGroup, namespaceName)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ServiceBus Namespace Network Rule Set (resource group %s)", resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceBusNamespaceNetworkRuleSetRead(d, meta)
}

func expandServiceBusNetworkRuleset(input []interface{}) *servicebus.NetworkRuleSetProperties {
	if len(input) == 0 {
		return nil
	}

	block := input[0].(map[string]interface{})

	ruleset := servicebus.NetworkRuleSetProperties{
		DefaultAction: servicebus.DefaultAction(block["default_action"].(string)),
	}

	if v, ok := block["virtual_network_rule"].([]interface{}); ok {
		if len(v) > 0 {
			var rules []servicebus.NWRuleSetVirtualNetworkRules
			for _, r := range v {
				rblock := r.(map[string]interface{})
				rules = append(rules, servicebus.NWRuleSetVirtualNetworkRules{
					Subnet: &servicebus.Subnet{
						ID: utils.String(rblock["subnet_id"].(string)),
					},
					IgnoreMissingVnetServiceEndpoint: utils.Bool(rblock["ignore_missing_virtual_network_service_endpoint"].(bool)),
				})
			}

			ruleset.VirtualNetworkRules = &rules
		}
	}

	if v, ok := block["ip_rule"].([]interface{}); ok {
		if len(v) > 0 {
			var rules []servicebus.NWRuleSetIPRules
			for _, r := range v {
				rblock := r.(map[string]interface{})
				rules = append(rules, servicebus.NWRuleSetIPRules{
					IPMask: utils.String(rblock["ip_mask"].(string)),
					Action: servicebus.NetworkRuleIPAction(rblock["action"].(string)),
				})
			}

			ruleset.IPRules = &rules
		}
	}

	return &ruleset
}

func flattenServiceBusNetworkRuleset(ruleset servicebus.NetworkRuleSet) []interface{} {
	if ruleset.NetworkRuleSetProperties == nil {
		return nil
	}

	vnetBlocks := make([]interface{}, 0)
	if vnetRules := ruleset.NetworkRuleSetProperties.VirtualNetworkRules; vnetRules != nil {
		for _, vnetRule := range *vnetRules {
			block := make(map[string]interface{})

			if s := vnetRule.Subnet; s != nil {
				if v := s.ID; v != nil {
					block["subnet_id"] = *v
				}
			}

			if v := vnetRule.IgnoreMissingVnetServiceEndpoint; v != nil {
				block["ignore_missing_virtual_network_service_endpoint"] = *v
			}

			vnetBlocks = append(vnetBlocks, block)
		}
	}
	ipBlocks := make([]interface{}, 0)
	if ipRules := ruleset.NetworkRuleSetProperties.IPRules; ipRules != nil {
		for _, ipRule := range *ipRules {
			block := make(map[string]interface{})

			block["action"] = string(ipRule.Action)

			if v := ipRule.IPMask; v != nil {
				block["ip_mask"] = *v
			}

			ipBlocks = append(ipBlocks, block)
		}
	}

	return []interface{}{map[string]interface{}{
		"default_action":       string(ruleset.DefaultAction),
		"virtual_network_rule": vnetBlocks,
		"ip_rule":              ipBlocks,
	}}
}
