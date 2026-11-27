// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_subscription

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	servicebusClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusSubscriptionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := servicebusClient.From(meta).SubscriptionsClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	var id subscriptions.Subscriptions2Id
	name := d.Get("name").(string)

	clientId := ""
	var isSubShared bool
	isClintScopedEnabled := d.Get("client_scoped_subscription_enabled").(bool)
	if isClintScopedEnabled {
		clientScopedSubsRawData := d.Get("client_scoped_subscription").([]any)
		if len(clientScopedSubsRawData) > 0 {
			clientScopedSubsProps := clientScopedSubsRawData[0].(map[string]any)
			if clientScopedSubsProps["client_id"] != "" {
				clientId = clientScopedSubsProps["client_id"].(string)
			}
			name = name + "$" + clientId + "$D"

			isSubShared = clientScopedSubsProps["is_client_scoped_subscription_shareable"].(bool)
		}
	}

	if topicIdLit := d.Get("topic_id").(string); topicIdLit != "" {
		topicId, _ := subscriptions.ParseTopicID(topicIdLit)
		id = subscriptions.NewSubscriptions2ID(topicId.SubscriptionId, topicId.ResourceGroupName, topicId.NamespaceName, topicId.TopicName, name)
	}

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_servicebus_subscription", id.ID())
			}
		}
	}

	enableBatchedOperations := d.Get("batched_operations_enabled").(bool)

	status := subscriptions.EntityStatus(d.Get("status").(string))
	parameters := subscriptions.SBSubscription{
		Properties: &subscriptions.SBSubscriptionProperties{
			DeadLetteringOnMessageExpiration:          pointer.To(d.Get("dead_lettering_on_message_expiration").(bool)),
			DeadLetteringOnFilterEvaluationExceptions: pointer.To(d.Get("dead_lettering_on_filter_evaluation_error").(bool)),
			EnableBatchedOperations:                   pointer.To(enableBatchedOperations),
			MaxDeliveryCount:                          pointer.To(int64(d.Get("max_delivery_count").(int))),
			RequiresSession:                           pointer.To(d.Get("requires_session").(bool)),
			Status:                                    &status,
			IsClientAffine:                            pointer.To(isClintScopedEnabled),
			ClientAffineProperties:                    &subscriptions.SBClientAffineProperties{},
		},
	}

	if isClintScopedEnabled {
		parameters.Properties.ClientAffineProperties.IsShared = &isSubShared
		parameters.Properties.ClientAffineProperties.ClientId = &clientId
	}

	if autoDeleteOnIdle := d.Get("auto_delete_on_idle").(string); autoDeleteOnIdle != "" {
		parameters.Properties.AutoDeleteOnIdle = &autoDeleteOnIdle
	}

	if lockDuration := d.Get("lock_duration").(string); lockDuration != "" {
		parameters.Properties.LockDuration = &lockDuration
	}

	if forwardTo := d.Get("forward_to").(string); forwardTo != "" {
		parameters.Properties.ForwardTo = &forwardTo
	}

	if forwardDeadLetteredMessagesTo := d.Get("forward_dead_lettered_messages_to").(string); forwardDeadLetteredMessagesTo != "" {
		parameters.Properties.ForwardDeadLetteredMessagesTo = &forwardDeadLetteredMessagesTo
	}

	if defaultMessageTtl := d.Get("default_message_ttl").(string); defaultMessageTtl != "" {
		parameters.Properties.DefaultMessageTimeToLive = &defaultMessageTtl
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	if d.IsNewResource() && features.From(meta).ServiceBus.AutoDeleteSubscriptionDefaultRule {
		rulesClient := servicebusClient.From(meta).SubscriptionRulesClient
		defaultRuleId := rules.NewRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicName, id.SubscriptionName, "$Default")

		if resp, err := rulesClient.Delete(ctx, defaultRuleId); err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("deleting default rule for %s: %+v", id, err)
			}
		}
	}

	return resourceServiceBusSubscriptionRead(d, meta)
}
