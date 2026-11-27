// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_server_configuration_policy_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNServerConfigurationPolicyGroupUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := virtualwans.ParseConfigurationPolicyGroupID(d.Id())
	if err != nil {
		return err
	}

	vpnServerConfigurationId := virtualwans.NewVpnServerConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.VpnServerConfigurationName)

	locks.ByID(vpnServerConfigurationId.ID())
	defer locks.UnlockByID(vpnServerConfigurationId.ID())

	existing, err := client.ConfigurationPolicyGroupsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("policy") {
		payload.Properties.PolicyMembers = expandVPNServerConfigurationPolicyGroupPolicyMembers(d.Get("policy").(*pluginsdk.Set).List())
	}

	if d.HasChange("priority") {
		payload.Properties.Priority = pointer.To(int64(d.Get("priority").(int)))
	}

	if err := client.ConfigurationPolicyGroupsCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVPNServerConfigurationPolicyGroupRead(d, meta)
}
