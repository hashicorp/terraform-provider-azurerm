// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_server_configuration_policy_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNServerConfigurationPolicyGroupCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vpnServerConfigurationId, err := virtualwans.ParseVpnServerConfigurationID(d.Get("vpn_server_configuration_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(vpnServerConfigurationId.ID())
	defer locks.UnlockByID(vpnServerConfigurationId.ID())

	id := virtualwans.NewConfigurationPolicyGroupID(subscriptionId, vpnServerConfigurationId.ResourceGroupName, vpnServerConfigurationId.VpnServerConfigurationName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.ConfigurationPolicyGroupsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_vpn_server_configuration_policy_group", id.ID())
		}
	}

	payload := virtualwans.VpnServerConfigurationPolicyGroup{
		Properties: &virtualwans.VpnServerConfigurationPolicyGroupProperties{
			IsDefault:     pointer.To(d.Get("is_default").(bool)),
			PolicyMembers: expandVPNServerConfigurationPolicyGroupPolicyMembers(d.Get("policy").(*pluginsdk.Set).List()),
			Priority:      pointer.To(int64(d.Get("priority").(int))),
		},
	}

	if err := client.ConfigurationPolicyGroupsCreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVPNServerConfigurationPolicyGroupRead(d, meta)
}

func expandVPNServerConfigurationPolicyGroupPolicyMembers(input []any) *[]virtualwans.VpnServerConfigurationPolicyGroupMember {
	results := make([]virtualwans.VpnServerConfigurationPolicyGroupMember, 0)

	for _, item := range input {
		v := item.(map[string]any)

		results = append(results, virtualwans.VpnServerConfigurationPolicyGroupMember{
			Name:           pointer.To(v["name"].(string)),
			AttributeType:  pointer.ToEnum[virtualwans.VpnPolicyMemberAttributeType](v["type"].(string)),
			AttributeValue: pointer.To(v["value"].(string)),
		})
	}

	return &results
}
