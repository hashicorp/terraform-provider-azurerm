// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_queue_authorization_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusQueueAuthorizationRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ServiceBus.QueuesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := queues.ParseQueueAuthorizationRuleID(d.Id())
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
	d.Set("queue_id", queues.NewQueueID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.QueueName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			listen, send, manage := flattenQueueAuthorizationRuleRights(&props.Rights)
			d.Set("manage", manage)
			d.Set("listen", listen)
			d.Set("send", send)
		}
	}

	keysResp, err := client.ListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	if keysModel := keysResp.Model; keysModel != nil {
		d.Set("primary_key", keysModel.PrimaryKey)
		d.Set("primary_connection_string", keysModel.PrimaryConnectionString)
		d.Set("secondary_key", keysModel.SecondaryKey)
		d.Set("secondary_connection_string", keysModel.SecondaryConnectionString)
		d.Set("primary_connection_string_alias", keysModel.AliasPrimaryConnectionString)
		d.Set("secondary_connection_string_alias", keysModel.AliasSecondaryConnectionString)
	}

	return nil
}

func flattenQueueAuthorizationRuleRights(rights *[]queues.AccessRights) (listen, send, manage bool) {
	// zero (initial) value for a bool in go is false

	if rights != nil {
		for _, right := range *rights {
			switch right {
			case queues.AccessRightsListen:
				listen = true
			case queues.AccessRightsSend:
				send = true
			case queues.AccessRightsManage:
				manage = true
			default:
				log.Printf("[DEBUG] Unknown Authorization Rule Right '%s'", right)
			}
		}
	}

	return listen, send, manage
}
