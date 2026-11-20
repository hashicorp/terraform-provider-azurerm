// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFirewallPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.FirewallPolicies
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := firewallpolicies.NewFirewallPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			resp, err := client.Get(ctx, id, firewallpolicies.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if resp.Model != nil {
				return tf.ImportAsExistsError("azurerm_firewall_policy", id.ID())
			}
		}
	}

	props := firewallpolicies.FirewallPolicy{
		Properties: &firewallpolicies.FirewallPolicyPropertiesFormat{
			ThreatIntelMode:      pointer.ToEnum[firewallpolicies.AzureFirewallThreatIntelMode](d.Get("threat_intelligence_mode").(string)),
			ThreatIntelWhitelist: expandFirewallPolicyThreatIntelWhitelist(d.Get("threat_intelligence_allowlist").([]any)),
			DnsSettings:          expandFirewallPolicyDNSSetting(d.Get("dns").([]any)),
			IntrusionDetection:   expandFirewallPolicyIntrusionDetection(d.Get("intrusion_detection").([]any)),
			TransportSecurity:    expandFirewallPolicyTransportSecurity(d.Get("tls_certificate").([]any)),
			Insights:             expandFirewallPolicyInsights(d.Get("insights").([]any)),
			ExplicitProxy:        expandFirewallPolicyExplicitProxy(d.Get("explicit_proxy").([]any)),
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}
	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	props.Identity = expandedIdentity

	if id, ok := d.GetOk("base_policy_id"); ok {
		props.Properties.BasePolicy = &firewallpolicies.SubResource{Id: pointer.To(id.(string))}
	}

	if v, ok := d.GetOk("sku"); ok {
		props.Properties.Sku = &firewallpolicies.FirewallPolicySku{
			Tier: pointer.ToEnum[firewallpolicies.FirewallPolicySkuTier](v.(string)),
		}
	}

	if v, ok := d.GetOk("sql_redirect_allowed"); ok {
		props.Properties.Sql = &firewallpolicies.FirewallPolicySQL{
			AllowSqlRedirect: pointer.To(v.(bool)),
		}
	}

	if v, ok := d.GetOk("private_ip_ranges"); ok {
		privateIPRanges := pluginsdk.ExpandStringSlice(v.([]any))
		props.Properties.Snat = &firewallpolicies.FirewallPolicySNAT{
			PrivateRanges: privateIPRanges,
		}
	}

	if v, ok := d.GetOk("auto_learn_private_ranges_enabled"); ok {
		if props.Properties.Snat == nil {
			props.Properties.Snat = &firewallpolicies.FirewallPolicySNAT{}
		}
		if v.(bool) {
			props.Properties.Snat.AutoLearnPrivateRanges = pointer.To(firewallpolicies.AutoLearnPrivateRangesModeEnabled)
		}
	}

	locks.ByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)
	defer locks.UnlockByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, props, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, props); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceFirewallPolicyRead(d, meta)
}
