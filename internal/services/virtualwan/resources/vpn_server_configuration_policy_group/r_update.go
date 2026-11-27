// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_server_configuration_policy_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNServerConfigurationPolicyGroupUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
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
