// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/firewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func FirewallDataSourcePolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPolicies
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := firewallpolicies.NewFirewallPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id, firewallpolicies.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.FirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			basePolicyID := ""
			if props.BasePolicy != nil && props.BasePolicy.Id != nil {
				basePolicyID = *props.BasePolicy.Id
			}
			d.Set("base_policy_id", basePolicyID)
			if err := d.Set("child_policies", flattenNetworkSubResourceID(props.ChildPolicies)); err != nil {
				return fmt.Errorf(`setting "child_policies": %+v`, err)
			}
			if err := d.Set("dns", flattenFirewallPolicyDNSSetting(props.DnsSettings)); err != nil {
				return fmt.Errorf(`setting "dns": %+v`, err)
			}
			if err := d.Set("firewalls", flattenNetworkSubResourceID(props.Firewalls)); err != nil {
				return fmt.Errorf(`setting "firewalls": %+v`, err)
			}
			if err := d.Set("rule_collection_groups", flattenNetworkSubResourceID(props.RuleCollectionGroups)); err != nil {
				return fmt.Errorf(`setting "rule_collection_groups": %+v`, err)
			}
			d.Set("threat_intelligence_mode", string(pointer.From(props.ThreatIntelMode)))
			if err := d.Set("threat_intelligence_allowlist", flattenFirewallPolicyThreatIntelWhitelist(props.ThreatIntelWhitelist)); err != nil {
				return fmt.Errorf(`setting "threat_intelligence_allowlist": %+v`, err)
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}
