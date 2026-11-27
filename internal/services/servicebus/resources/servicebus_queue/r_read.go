// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_queue

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	servicebusClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusQueueRead(d *pluginsdk.ResourceData, meta any) error {
	client := servicebusClient.From(meta).QueuesClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := queues.ParseQueueID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	namespacesClient := servicebusClient.From(meta).NamespacesClient

	return resourceServiceBusQueueFlatten(ctx, namespacesClient, d, id, resp.Model, true)
}

func serviceBusQueueStatusRefreshFunc(ctx context.Context, client *queues.QueuesClient, id queues.QueueId, userConfig map[string]any) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		log.Printf("[DEBUG] Checking servicebus queue %s status...", id)

		resp, err := client.Get(ctx, id)
		if err != nil {
			return nil, "Error", fmt.Errorf("checking servicebus queue %s error: %+v", id, err)
		}

		queueStatus := "Updating"

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				if pointer.From(props.Status) == userConfig["status"] &&
					int(pointer.From(props.MaxDeliveryCount)) == userConfig["maxDeliveryCount"].(int) &&
					pointer.From(props.EnableExpress) == userConfig["enableExpress"].(bool) &&
					pointer.From(props.EnablePartitioning) == userConfig["enablePartitioning"].(bool) &&
					pointer.From(props.EnableBatchedOperations) == userConfig["enableBatchOps"].(bool) {
					queueStatus = "Succeeded"
				}

				if userConfig["deadLetteringOnMesExp"] != "" {
					if pointer.From(props.DeadLetteringOnMessageExpiration) != userConfig["deadLetteringOnMesExp"].(bool) {
						queueStatus = "Updating"
					}
				}

				if userConfig["lockDuration"] != "" {
					if pointer.From(props.LockDuration) != userConfig["lockDuration"].(string) {
						queueStatus = "Updating"
					}
				}

				if userConfig["forwardTo"] != "" {
					if pointer.From(props.ForwardTo) != userConfig["forwardTo"].(string) {
						queueStatus = "Updating"
					}
				}

				if userConfig["forwardDeadLetteredMessagesTo"] != nil {
					if pointer.From(props.ForwardDeadLetteredMessagesTo) != userConfig["forwardDeadLetteredMessagesTo"].(string) {
						queueStatus = "Updating"
					}
				}
			}
		}
		return resp, queueStatus, nil
	}
}
