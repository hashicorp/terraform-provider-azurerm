// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall_policy_rule_collection_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicyrulecollectiongroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFirewallPolicyRuleCollectionGroupDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.FirewallPolicyRuleCollectionGroups
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallpolicyrulecollectiongroups.ParseRuleCollectionGroupID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)
	defer locks.UnlockByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}
