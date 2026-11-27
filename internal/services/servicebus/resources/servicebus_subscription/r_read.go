// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_subscription

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusSubscriptionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ServiceBus.SubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := subscriptions.ParseSubscriptions2ID(d.Id())
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

	d.Set("topic_id", subscriptions.NewTopicID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicName).ID())

	clientScopedEnabled := false
	userAssignedName := id.SubscriptionName
	clientId := ""

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("auto_delete_on_idle", props.AutoDeleteOnIdle)
			d.Set("default_message_ttl", props.DefaultMessageTimeToLive)
			d.Set("lock_duration", props.LockDuration)
			d.Set("dead_lettering_on_message_expiration", props.DeadLetteringOnMessageExpiration)
			d.Set("dead_lettering_on_filter_evaluation_error", props.DeadLetteringOnFilterEvaluationExceptions)
			d.Set("requires_session", props.RequiresSession)
			d.Set("forward_to", props.ForwardTo)
			d.Set("forward_dead_lettered_messages_to", props.ForwardDeadLetteredMessagesTo)
			d.Set("status", pointer.To(string(*props.Status)))
			d.Set("client_scoped_subscription_enabled", props.IsClientAffine)
			d.Set("batched_operations_enabled", props.EnableBatchedOperations)

			if count := props.MaxDeliveryCount; count != nil {
				d.Set("max_delivery_count", int(*count))
			}
			if props.IsClientAffine != nil && *props.IsClientAffine {
				if props.ClientAffineProperties.ClientId != nil {
					clientId = *props.ClientAffineProperties.ClientId
				}
				clientScopedEnabled = true
			}
			d.Set("client_scoped_subscription", flattenServiceBusNamespaceClientScopedSubscription(props.ClientAffineProperties))
		}
	}

	if clientScopedEnabled {
		suffix := "$" + clientId + "$D"
		userAssignedName = strings.TrimSuffix(userAssignedName, suffix)
	}
	d.Set("name", userAssignedName)

	return nil
}

func flattenServiceBusNamespaceClientScopedSubscription(clientScopedSubsProps *subscriptions.SBClientAffineProperties) []any {
	if clientScopedSubsProps == nil {
		return []any{}
	}

	var clientId string
	var isShareable, isDurable bool

	if clientScopedSubsProps.ClientId != nil {
		clientId = *clientScopedSubsProps.ClientId
	}

	if clientScopedSubsProps.IsShared != nil {
		isShareable = *clientScopedSubsProps.IsShared
	}

	if clientScopedSubsProps.IsDurable != nil {
		isDurable = *clientScopedSubsProps.IsDurable
	}

	return []any{
		map[string]any{
			"client_id": clientId,
			"is_client_scoped_subscription_shareable": isShareable,
			"is_client_scoped_subscription_durable":   isDurable,
		},
	}
}
