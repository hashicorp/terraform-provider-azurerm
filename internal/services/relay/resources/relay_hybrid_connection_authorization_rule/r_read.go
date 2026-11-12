// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_hybrid_connection_authorization_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRelayHybridConnectionAuthorizationRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Relay.HybridConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hybridconnections.ParseHybridConnectionAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetAuthorizationRule(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	keysResp, err := client.ListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	d.Set("name", id.AuthorizationRuleName)
	d.Set("hybrid_connection_name", id.HybridConnectionName)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		listen, send, manage := flattenHybridConnectionAuthorizationRuleRights(model.Properties.Rights)
		d.Set("manage", manage)
		d.Set("listen", listen)
		d.Set("send", send)
	}

	d.Set("primary_key", keysResp.Model.PrimaryKey)
	d.Set("primary_connection_string", keysResp.Model.PrimaryConnectionString)
	d.Set("secondary_key", keysResp.Model.SecondaryKey)
	d.Set("secondary_connection_string", keysResp.Model.SecondaryConnectionString)

	return nil
}
