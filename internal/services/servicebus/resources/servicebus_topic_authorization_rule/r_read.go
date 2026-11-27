// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_topic_authorization_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusTopicAuthorizationRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ServiceBus.TopicsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := topics.ParseTopicAuthorizationRuleID(d.Id())
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

	d.Set("name", id.AuthorizationRuleName)
	d.Set("topic_id", topics.NewTopicID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			listen, send, manage := flattenTopicAuthorizationRuleRights(&props.Rights)
			d.Set("listen", listen)
			d.Set("send", send)
			d.Set("manage", manage)
		}
	}

	keysResp, err := client.ListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	if model := keysResp.Model; model != nil {
		d.Set("primary_key", model.PrimaryKey)
		d.Set("primary_connection_string", model.PrimaryConnectionString)
		d.Set("secondary_key", model.SecondaryKey)
		d.Set("secondary_connection_string", model.SecondaryConnectionString)
		d.Set("primary_connection_string_alias", model.AliasPrimaryConnectionString)
		d.Set("secondary_connection_string_alias", model.AliasSecondaryConnectionString)
	}

	return nil
}

func flattenTopicAuthorizationRuleRights(rights *[]topics.AccessRights) (listen, send, manage bool) {
	if rights != nil {
		for _, right := range *rights {
			switch right {
			case topics.AccessRightsListen:
				listen = true
			case topics.AccessRightsSend:
				send = true
			case topics.AccessRightsManage:
				manage = true
			default:
				log.Printf("[DEBUG] Unknown Authorization Rule Right '%s'", right)
			}
		}
	}

	return listen, send, manage
}
