// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func FirewallDataSourcePolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: FirewallDataSourcePolicyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FirewallPolicyName(),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"base_policy_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"child_policies": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"dns": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"servers": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"proxy_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"network_rule_fqdn_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"firewalls": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"rule_collection_groups": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"threat_intelligence_mode": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"threat_intelligence_allowlist": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ip_addresses": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"fqdns": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func FirewallDataSourcePolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewFirewallPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
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
		if err := d.Set("dns", flattenFirewallPolicyDNSSetting(prop.DNSSettings)); err != nil {
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
