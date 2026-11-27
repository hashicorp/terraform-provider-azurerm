// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_server_configuration_policy_group

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNServerConfigurationPolicyGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := virtualwans.ParseConfigurationPolicyGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ConfigurationPolicyGroupsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ConfigurationPolicyGroupName)

	vpnServerConfigurationId := virtualwans.NewVpnServerConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.VpnServerConfigurationName)
	d.Set("vpn_server_configuration_id", vpnServerConfigurationId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("is_default", props.IsDefault)
			d.Set("priority", props.Priority)

			if err := d.Set("policy", flattenVPNServerConfigurationPolicyGroupPolicyMembers(props.PolicyMembers)); err != nil {
				return fmt.Errorf("setting `policy`: %+v", err)
			}
		}
	}

	return nil
}

func flattenVPNServerConfigurationPolicyGroupPolicyMembers(input *[]virtualwans.VpnServerConfigurationPolicyGroupMember) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		attributeType := ""
		if item.AttributeType != nil {
			attributeType = string(*item.AttributeType)
		}

		results = append(results, map[string]any{
			"name":  pointer.From(item.Name),
			"type":  attributeType,
			"value": pointer.From(item.AttributeValue),
		})
	}

	return results
}
