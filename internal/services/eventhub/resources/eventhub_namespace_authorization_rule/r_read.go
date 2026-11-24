// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_authorization_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationrulesnamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceAuthorizationRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.NamespaceAuthorizationRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := authorizationrulesnamespaces.ParseAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.NamespacesGetAuthorizationRule(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.AuthorizationRuleName)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			listen, send, manage := flattenEventHubNamespaceAuthorizationRuleRights(props.Rights)
			d.Set("manage", manage)
			d.Set("listen", listen)
			d.Set("send", send)
		}
	}

	keysResp, err := client.NamespacesListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	if model := keysResp.Model; model != nil {
		d.Set("primary_key", model.PrimaryKey)
		d.Set("secondary_key", model.SecondaryKey)
		d.Set("primary_connection_string", model.PrimaryConnectionString)
		d.Set("secondary_connection_string", model.SecondaryConnectionString)
		d.Set("primary_connection_string_alias", model.AliasPrimaryConnectionString)
		d.Set("secondary_connection_string_alias", model.AliasSecondaryConnectionString)
	}

	return nil
}

func flattenEventHubNamespaceAuthorizationRuleRights(rights []authorizationrulesnamespaces.AccessRights) (listen, send, manage bool) {
	// zero (initial) value for a bool in go is false

	for _, right := range rights {
		switch right {
		case authorizationrulesnamespaces.AccessRightsListen:
			listen = true
		case authorizationrulesnamespaces.AccessRightsSend:
			send = true
		case authorizationrulesnamespaces.AccessRightsManage:
			manage = true
		default:
			log.Printf("[DEBUG] Unknown Authorization Rule Right '%s'", right)
		}
	}

	return listen, send, manage
}
