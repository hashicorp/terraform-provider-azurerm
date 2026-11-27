// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor_firewall_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	frontdoorClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontDoorFirewallPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := frontdoorClient.From(meta).FrontDoorsPolicyClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := webapplicationfirewallpolicies.ParseFrontDoorWebApplicationFirewallPolicyIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	if err := client.PoliciesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
