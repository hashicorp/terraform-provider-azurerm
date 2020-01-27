package servicebus

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/servicebus/mgmt/2018-01-01-preview/servicebus"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var serviceBusNamespaceDefaultAuthorizationRule = "RootManageSharedAccessKey"

func resourceArmServiceBusNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusNamespaceCreateUpdate,
		Read:   resourceArmServiceBusNamespaceRead,
		Update: resourceArmServiceBusNamespaceCreateUpdate,
		Delete: resourceArmServiceBusNamespaceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		MigrateState:  ResourceAzureRMServiceBusNamespaceMigrateState,
		SchemaVersion: 1,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{0,100}[a-zA-Z0-9]$"),
					"The namespace can contain only letters, numbers, and hyphens. The namespace must start with a letter, and it must end with a letter or number.",
				),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(servicebus.Basic),
					string(servicebus.Standard),
					string(servicebus.Premium),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"network_rulesets": {
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

			"capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validate.IntInSlice([]int{0, 1, 2, 4, 8}),
			},

			"default_primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"zone_redundant": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmServiceBusNamespaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClientPreview
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM ServiceBus Namespace creation.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	sku := d.Get("sku").(string)
	t := d.Get("tags").(map[string]interface{})

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing ServiceBus Namespace %q (resource group %q) ID", name, resourceGroup)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_servicebus_namespace", *existing.ID)
		}
	}

	parameters := servicebus.SBNamespace{
		Location: &location,
		Sku: &servicebus.SBSku{
			Name: servicebus.SkuName(sku),
			Tier: servicebus.SkuTier(sku),
		},
		Tags: tags.Expand(t),
	}

	if capacity := d.Get("capacity"); capacity != nil {
		if !strings.EqualFold(sku, string(servicebus.Premium)) && capacity.(int) > 0 {
			return fmt.Errorf("Service Bus SKU %q only supports `capacity` of 0", sku)
		}
		if strings.EqualFold(sku, string(servicebus.Premium)) && capacity.(int) == 0 {
			return fmt.Errorf("Service Bus SKU %q only supports `capacity` of 1, 2, 4 or 8", sku)
		}
		parameters.Sku.Capacity = utils.Int32(int32(capacity.(int)))
	}

	if zoneRedundant, ok := d.GetOkExists("zone_redundant"); ok {
		properties := servicebus.SBNamespaceProperties{
			ZoneRedundant: utils.Bool(zoneRedundant.(bool)),
		}
		parameters.SBNamespaceProperties = &properties
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ServiceBus Namespace %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	ruleSets, hasRuleSets := d.GetOk("network_rulesets")
	if hasRuleSets {
		rulesets := servicebus.NetworkRuleSet{
			NetworkRuleSetProperties: expandServiceBusNetworkRuleset(ruleSets.([]interface{})),
		}

		// cannot use network rulesets with the basic SKU
		if parameters.Sku.Name != servicebus.Basic && parameters.Sku.Name != servicebus.Standard {
			if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, resourceGroup, name, rulesets); err != nil {
				return fmt.Errorf("Error setting network ruleset properties for Service Bus %q (resource group %q): %v", name, resourceGroup, err)
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

	return resourceArmServiceBusNamespaceRead(d, meta)
}

func resourceArmServiceBusNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClientPreview
	clientStable := meta.(*clients.Client).ServiceBus.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure ServiceBus Namespace %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", strings.ToLower(string(sku.Name)))
		d.Set("capacity", sku.Capacity)
	}

	if properties := resp.SBNamespaceProperties; properties != nil {
		d.Set("zone_redundant", properties.ZoneRedundant)
	}

	ruleset, err := client.GetNetworkRuleSet(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Service Bus %q Network Ruleset: %+v", name, err)
	}

	if err := d.Set("network_rulesets", flattenServiceBusNetworkRuleset(ruleset)); err != nil {
		return fmt.Errorf("Error setting `network_ruleset` for Service Bus %s: %v", name, err)
	}

	keys, err := clientStable.ListKeys(ctx, resourceGroup, name, serviceBusNamespaceDefaultAuthorizationRule)
	if err != nil {
		log.Printf("[WARN] Unable to List default keys for Namespace %q (Resource Group %q): %+v", name, resourceGroup, err)
	} else {
		d.Set("default_primary_connection_string", keys.PrimaryConnectionString)
		d.Set("default_secondary_connection_string", keys.SecondaryConnectionString)
		d.Set("default_primary_key", keys.PrimaryKey)
		d.Set("default_secondary_key", keys.SecondaryKey)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmServiceBusNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClientPreview
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Service Bus %q: %+v", name, err)
	}

	return nil
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
