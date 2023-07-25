// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/subscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topics"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceServiceBusSubscription() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceBusSubscriptionCreateUpdate,
		Read:   resourceServiceBusSubscriptionRead,
		Update: resourceServiceBusSubscriptionCreateUpdate,
		Delete: resourceServiceBusSubscriptionDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ServiceBusSubscriptionV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := subscriptions.ParseSubscriptions2ID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceServicebusSubscriptionSchema(),
	}
}

func resourceServicebusSubscriptionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SubscriptionName(),
		},

		//lintignore: S013
		"topic_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: topics.ValidateTopicID,
		},

		"auto_delete_on_idle": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"default_message_ttl": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"lock_duration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"dead_lettering_on_message_expiration": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"dead_lettering_on_filter_evaluation_error": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		// TODO 4.0: change this from enable_* to *_enabled
		"enable_batched_operations": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"max_delivery_count": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"requires_session": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			// cannot be modified
			ForceNew: true,
		},

		"forward_to": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"forward_dead_lettered_messages_to": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(subscriptions.EntityStatusActive),
			ValidateFunc: validation.StringInSlice([]string{
				string(subscriptions.EntityStatusActive),
				string(subscriptions.EntityStatusDisabled),
				string(subscriptions.EntityStatusReceiveDisabled),
			}, false),
		},

		"client_scoped_subscription_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"client_scoped_subscription": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"is_client_scoped_subscription_shareable": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
						ForceNew: true,
					},

					"is_client_scoped_subscription_durable": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},
	}
}

func resourceServiceBusSubscriptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.SubscriptionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for ServiceBus Subscription creation.")

	var id subscriptions.Subscriptions2Id
	name := d.Get("name").(string)

	clientId := ""
	var isSubShared bool
	isClintScopedEnabled := d.Get("client_scoped_subscription_enabled").(bool)
	if isClintScopedEnabled {
		clientScopedSubsRawData := d.Get("client_scoped_subscription").([]interface{})
		if len(clientScopedSubsRawData) > 0 {
			clientScopedSubsProps := clientScopedSubsRawData[0].(map[string]interface{})
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

	status := subscriptions.EntityStatus(d.Get("status").(string))
	parameters := subscriptions.SBSubscription{
		Properties: &subscriptions.SBSubscriptionProperties{
			DeadLetteringOnMessageExpiration:          utils.Bool(d.Get("dead_lettering_on_message_expiration").(bool)),
			DeadLetteringOnFilterEvaluationExceptions: utils.Bool(d.Get("dead_lettering_on_filter_evaluation_error").(bool)),
			EnableBatchedOperations:                   utils.Bool(d.Get("enable_batched_operations").(bool)),
			MaxDeliveryCount:                          utils.Int64(int64(d.Get("max_delivery_count").(int))),
			RequiresSession:                           utils.Bool(d.Get("requires_session").(bool)),
			Status:                                    &status,
			IsClientAffine:                            utils.Bool(isClintScopedEnabled),
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

	d.SetId(id.ID())
	return resourceServiceBusSubscriptionRead(d, meta)
}

func resourceServiceBusSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			d.Set("enable_batched_operations", props.EnableBatchedOperations)
			d.Set("requires_session", props.RequiresSession)
			d.Set("forward_to", props.ForwardTo)
			d.Set("forward_dead_lettered_messages_to", props.ForwardDeadLetteredMessagesTo)
			d.Set("status", utils.String(string(*props.Status)))
			d.Set("client_scoped_subscription_enabled", props.IsClientAffine)

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

func resourceServiceBusSubscriptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.SubscriptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := subscriptions.ParseSubscriptions2ID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func flattenServiceBusNamespaceClientScopedSubscription(clientScopedSubsProps *subscriptions.SBClientAffineProperties) []interface{} {
	if clientScopedSubsProps == nil {
		return []interface{}{}
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

	return []interface{}{
		map[string]interface{}{
			"client_id": clientId,
			"is_client_scoped_subscription_shareable": isShareable,
			"is_client_scoped_subscription_durable":   isDurable,
		},
	}
}
