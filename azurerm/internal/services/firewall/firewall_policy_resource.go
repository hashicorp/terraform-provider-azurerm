package firewall

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-07-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/firewall/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/firewall/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const azureFirewallPolicyResourceName = "azurerm_firewall_policy"

func resourceFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallPolicyCreateUpdate,
		Read:   resourceFirewallPolicyRead,
		Update: resourceFirewallPolicyCreateUpdate,
		Delete: resourceFirewallPolicyDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.FirewallPolicyID(id)
			return err
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
				ValidateFunc: validate.FirewallPolicyName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.FirewallPolicySkuTierPremium),
					string(network.FirewallPolicySkuTierStandard),
				}, false),
			},

			"location": location.Schema(),

			"base_policy_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.FirewallPolicyID,
			},

			"dns": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"servers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsIPv4Address,
							},
						},
						"proxy_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						// TODO 3.0 - remove this property
						"network_rule_fqdn_enabled": {
							Type:       schema.TypeBool,
							Optional:   true,
							Computed:   true,
							Deprecated: "This property has been deprecated as the service team has removed it from all API versions and is no longer supported by Azure. It will be removed in v3.0 of the provider.",
						},
					},
				},
			},

			"threat_intelligence_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(network.AzureFirewallThreatIntelModeAlert),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.AzureFirewallThreatIntelModeAlert),
					string(network.AzureFirewallThreatIntelModeDeny),
					string(network.AzureFirewallThreatIntelModeOff),
				}, false),
			},

			"threat_intelligence_allowlist": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_addresses": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.Any(validation.IsIPv4Range, validation.IsIPv4Address),
							},
							AtLeastOneOf: []string{"threat_intelligence_allowlist.0.ip_addresses", "threat_intelligence_allowlist.0.fqdns"},
						},
						"fqdns": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							AtLeastOneOf: []string{"threat_intelligence_allowlist.0.ip_addresses", "threat_intelligence_allowlist.0.fqdns"},
						},
					},
				},
			},

			"child_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"firewalls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"rule_collection_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.SchemaEnforceLowerCaseKeys(),
		},
	}
}

func resourceFirewallPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if resp.ID != nil && *resp.ID != "" {
			return tf.ImportAsExistsError("azurerm_firewall_policy", *resp.ID)
		}
	}

	props := network.FirewallPolicy{
		FirewallPolicyPropertiesFormat: &network.FirewallPolicyPropertiesFormat{
			ThreatIntelMode:      network.AzureFirewallThreatIntelMode(d.Get("threat_intelligence_mode").(string)),
			ThreatIntelWhitelist: expandFirewallPolicyThreatIntelWhitelist(d.Get("threat_intelligence_allowlist").([]interface{})),
			DNSSettings:          expandFirewallPolicyDNSSetting(d.Get("dns").([]interface{})),
		},
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if id, ok := d.GetOk("base_policy_id"); ok {
		props.FirewallPolicyPropertiesFormat.BasePolicy = &network.SubResource{ID: utils.String(id.(string))}
	}

	if v, ok := d.GetOk("sku"); ok {
		props.FirewallPolicyPropertiesFormat.Sku = &network.FirewallPolicySku{
			Tier: network.FirewallPolicySkuTier(v.(string)),
		}
	}

	locks.ByName(name, azureFirewallPolicyResourceName)
	defer locks.UnlockByName(name, azureFirewallPolicyResourceName)

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, props); err != nil {
		return fmt.Errorf("creating Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("retrieving Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Firewall Policy %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceFirewallPolicyRead(d, meta)
}

func resourceFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Firewall Policy %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Firewall Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if prop := resp.FirewallPolicyPropertiesFormat; prop != nil {
		basePolicyID := ""
		if resp.BasePolicy != nil && resp.BasePolicy.ID != nil {
			basePolicyID = *resp.BasePolicy.ID
		}
		d.Set("base_policy_id", basePolicyID)

		d.Set("threat_intelligence_mode", string(prop.ThreatIntelMode))

		if sku := prop.Sku; sku != nil {
			d.Set("sku", string(sku.Tier))
		}

		if err := d.Set("threat_intelligence_allowlist", flattenFirewallPolicyThreatIntelWhitelist(resp.ThreatIntelWhitelist)); err != nil {
			return fmt.Errorf(`setting "threat_intelligence_allowlist": %+v`, err)
		}

		if err := d.Set("dns", flattenFirewallPolicyDNSSetting(prop.DNSSettings)); err != nil {
			return fmt.Errorf(`setting "dns": %+v`, err)
		}

		if err := d.Set("child_policies", flattenNetworkSubResourceID(prop.ChildPolicies)); err != nil {
			return fmt.Errorf(`setting "child_policies": %+v`, err)
		}

		if err := d.Set("firewalls", flattenNetworkSubResourceID(prop.Firewalls)); err != nil {
			return fmt.Errorf(`setting "firewalls": %+v`, err)
		}

		if err := d.Set("rule_collection_groups", flattenNetworkSubResourceID(prop.RuleCollectionGroups)); err != nil {
			return fmt.Errorf(`setting "rule_collection_groups": %+v`, err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceFirewallPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, azureFirewallPolicyResourceName)
	defer locks.UnlockByName(id.Name, azureFirewallPolicyResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Firewall Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting Firewall Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandFirewallPolicyThreatIntelWhitelist(input []interface{}) *network.FirewallPolicyThreatIntelWhitelist {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &network.FirewallPolicyThreatIntelWhitelist{
		IPAddresses: utils.ExpandStringSlice(raw["ip_addresses"].(*schema.Set).List()),
		Fqdns:       utils.ExpandStringSlice(raw["fqdns"].(*schema.Set).List()),
	}

	return output
}

func expandFirewallPolicyDNSSetting(input []interface{}) *network.DNSSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &network.DNSSettings{
		Servers:     utils.ExpandStringSlice(raw["servers"].(*schema.Set).List()),
		EnableProxy: utils.Bool(raw["proxy_enabled"].(bool)),
	}

	return output
}

func flattenFirewallPolicyThreatIntelWhitelist(input *network.FirewallPolicyThreatIntelWhitelist) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"ip_addresses": utils.FlattenStringSlice(input.IPAddresses),
			"fqdns":        utils.FlattenStringSlice(input.Fqdns),
		},
	}
}

func flattenFirewallPolicyDNSSetting(input *network.DNSSettings) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	proxyEnabled := false
	if input.EnableProxy != nil {
		proxyEnabled = *input.EnableProxy
	}

	return []interface{}{
		map[string]interface{}{
			"servers":       utils.FlattenStringSlice(input.Servers),
			"proxy_enabled": proxyEnabled,
			// TODO 3.0: remove the setting zero value for property below.
			"network_rule_fqdn_enabled": false,
		},
	}
}
