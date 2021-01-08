package firewall

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/firewall/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func FirewallDataSourcePolicy() *schema.Resource {
	return &schema.Resource{
		Read: FirewallDataSourcePolicyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.FirewallPolicyName(),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"base_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"child_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"dns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"proxy_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"network_rule_fqdn_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
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

			"threat_intelligence_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"threat_intelligence_allowlist": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fqdns": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func FirewallDataSourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Firewall Policy %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("retrieving Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Firewall Policy %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if prop := resp.FirewallPolicyPropertiesFormat; prop != nil {
		basePolicyID := ""
		if resp.BasePolicy != nil && resp.BasePolicy.ID != nil {
			basePolicyID = *resp.BasePolicy.ID
		}
		d.Set("base_policy_id", basePolicyID)
		if err := d.Set("child_policies", flattenNetworkSubResourceID(prop.ChildPolicies)); err != nil {
			return fmt.Errorf(`setting "child_policies": %+v`, err)
		}
		if err := d.Set("dns", flattenFirewallPolicyDNSSetting(resp.DNSSettings)); err != nil {
			return fmt.Errorf(`setting "dns": %+v`, err)
		}
		if err := d.Set("firewalls", flattenNetworkSubResourceID(prop.Firewalls)); err != nil {
			return fmt.Errorf(`setting "firewalls": %+v`, err)
		}
		if err := d.Set("rule_collection_groups", flattenNetworkSubResourceID(prop.RuleCollectionGroups)); err != nil {
			return fmt.Errorf(`setting "rule_collection_groups": %+v`, err)
		}
		d.Set("threat_intelligence_mode", string(prop.ThreatIntelMode))
		if err := d.Set("threat_intelligence_allowlist", flattenFirewallPolicyThreatIntelWhitelist(resp.ThreatIntelWhitelist)); err != nil {
			return fmt.Errorf(`setting "threat_intelligence_allowlist": %+v`, err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
