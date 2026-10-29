// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall_policy_rule_collection_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicyrulecollectiongroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFirewallPolicyRuleCollectionGroupCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.FirewallPolicyRuleCollectionGroups
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	policyId, err := firewallpolicies.ParseFirewallPolicyID(d.Get("firewall_policy_id").(string))
	if err != nil {
		return err
	}

	id := firewallpolicyrulecollectiongroups.NewRuleCollectionGroupID(policyId.SubscriptionId, policyId.ResourceGroupName, policyId.FirewallPolicyName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if resp.Model != nil {
				return tf.ImportAsExistsError("azurerm_firewall_policy_rule_collection_group", id.ID())
			}
		}
	}

	locks.ByName(policyId.FirewallPolicyName, AzureFirewallPolicyResourceName)
	defer locks.UnlockByName(policyId.FirewallPolicyName, AzureFirewallPolicyResourceName)

	param := firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollectionGroup{
		Properties: &firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollectionGroupProperties{
			Priority: pointer.To(int64(d.Get("priority").(int))),
		},
	}
	var rulesCollections []firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection
	rulesCollections = append(rulesCollections, expandFirewallPolicyRuleCollectionApplication(d.Get("application_rule_collection").([]any))...)
	rulesCollections = append(rulesCollections, expandFirewallPolicyRuleCollectionNetwork(d.Get("network_rule_collection").([]any))...)

	natRules, err := expandFirewallPolicyRuleCollectionNat(d.Get("nat_rule_collection").([]any))
	if err != nil {
		return fmt.Errorf("expanding NAT rule collection: %w", err)
	}
	rulesCollections = append(rulesCollections, natRules...)

	param.Properties.RuleCollections = &rulesCollections

	if d.IsNewResource() {
		if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, param, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err = client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceFirewallPolicyRuleCollectionGroupRead(d, meta)
}
