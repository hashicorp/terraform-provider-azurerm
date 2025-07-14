// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queues"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-10-01-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceServiceBusQueue() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceBusQueueCreateUpdate,
		Read:   resourceServiceBusQueueRead,
		Update: resourceServiceBusQueueCreateUpdate,
		Delete: resourceServiceBusQueueDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := queues.ParseQueueID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceServicebusQueueSchema(),
	}
}

func resourceServicebusQueueSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azValidate.QueueName(),
		},

		// lintignore: S013
		"namespace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},

		"auto_delete_on_idle": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C this gets a default except when using basic sku and can be updated without issues
			Computed:     true,
			ValidateFunc: validate.ISO8601Duration,
		},

		"dead_lettering_on_message_expiration": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"default_message_ttl": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C this gets a default of "P10675199DT2H48M5.4775807S" (Unbounded) and "P14D" in Basic sku and can be updated without issues
			Computed:     true,
			ValidateFunc: validate.ISO8601Duration,
		},

		"duplicate_detection_history_time_window": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT10M", // 10 minutes
			ValidateFunc: validate.ISO8601Duration,
		},

		"batched_operations_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"express_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"partitioning_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		"forward_dead_lettered_messages_to": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azValidate.QueueName(),
		},

		"forward_to": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azValidate.QueueName(),
		},

		"lock_duration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "PT1M", // 1 minute
		},

		"max_delivery_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      10,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"max_message_size_in_kilobytes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C this gets a variable default based on the sku and can be updated without issues
			Computed:     true,
			ValidateFunc: azValidate.ServiceBusMaxMessageSizeInKilobytes(),
		},

		"max_size_in_megabytes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C this gets a variable default based on the sku and can be updated without issues
			Computed:     true,
			ValidateFunc: azValidate.ServiceBusMaxSizeInMegabytes(),
		},

		"requires_duplicate_detection": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		"requires_session": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(queues.EntityStatusActive),
			ValidateFunc: validation.StringInSlice([]string{
				string(queues.EntityStatusActive),
				string(queues.EntityStatusCreating),
				string(queues.EntityStatusDeleting),
				string(queues.EntityStatusDisabled),
				string(queues.EntityStatusReceiveDisabled),
				string(queues.EntityStatusRenaming),
				string(queues.EntityStatusSendDisabled),
				string(queues.EntityStatusUnknown),
			}, false),
		},
	}
}

func resourceServiceBusQueueCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.QueuesClient
	namespaceClient := meta.(*clients.Client).ServiceBus.NamespacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	namespaceId, err := namespaces.ParseNamespaceID(d.Get("namespace_id").(string))
	if err != nil {
		return err
	}

	id := queues.NewQueueID(namespaceId.SubscriptionId, namespaceId.ResourceGroupName, namespaceId.NamespaceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_servicebus_queue", id.ID())
		}
	}

	isPremiumNamespacePartitioned := true
	sbNamespace, err := namespaceClient.Get(ctx, *namespaceId)
	if err != nil {
		return fmt.Errorf("checking the parent namespace %s: %+v", id, err)
	}

	if sbNamespaceModel := sbNamespace.Model; sbNamespaceModel != nil {
		if sbNamespaceModel.Properties != nil &&
			sbNamespaceModel.Properties.PremiumMessagingPartitions != nil && *sbNamespaceModel.Properties.PremiumMessagingPartitions == 1 {
			isPremiumNamespacePartitioned = false
		}
	}

	userConfig := make(map[string]interface{})

	status := queues.EntityStatus(d.Get("status").(string))
	userConfig["status"] = status
	maxDeliveryCount := d.Get("max_delivery_count").(int)
	userConfig["maxDeliveryCount"] = maxDeliveryCount
	deadLetteringOnMesExp := d.Get("dead_lettering_on_message_expiration").(bool)
	userConfig["deadLetteringOnMesExp"] = deadLetteringOnMesExp
	maxSizeInMB := d.Get("max_size_in_megabytes").(int)
	requireDuplicateDetection := d.Get("requires_duplicate_detection").(bool)
	requireSession := d.Get("requires_session").(bool)
	forwardDeadLetteredMessagesTo := d.Get("forward_dead_lettered_messages_to").(string)
	userConfig["forwardDeadLetteredMessagesTo"] = forwardDeadLetteredMessagesTo
	forwardTo := d.Get("forward_to").(string)
	userConfig["forwardTo"] = forwardTo
	lockDuration := d.Get("lock_duration").(string)
	userConfig["lockDuration"] = lockDuration
	defaultMessageTTL := d.Get("default_message_ttl").(string)
	userConfig["defaultMessageTTL"] = defaultMessageTTL
	autoDeleteOnIdle := d.Get("auto_delete_on_idle").(string)
	userConfig["autoDeleteOnIdle"] = autoDeleteOnIdle
	duplicateDetectionHistoryTimeWindow := d.Get("duplicate_detection_history_time_window").(string)

	enableExpress := d.Get("express_enabled").(bool)
	enablePartitioning := d.Get("partitioning_enabled").(bool)
	enableBatchedOperations := d.Get("batched_operations_enabled").(bool)

	userConfig["enableExpress"] = enableExpress
	userConfig["enablePartitioning"] = enablePartitioning
	userConfig["enableBatchOps"] = enableBatchedOperations

	parameters := queues.SBQueue{
		Name: utils.String(id.QueueName),
		Properties: &queues.SBQueueProperties{
			DeadLetteringOnMessageExpiration: utils.Bool(deadLetteringOnMesExp),
			EnableBatchedOperations:          utils.Bool(enableBatchedOperations),
			EnableExpress:                    utils.Bool(enableExpress),
			EnablePartitioning:               utils.Bool(enablePartitioning),
			MaxDeliveryCount:                 utils.Int64(int64(maxDeliveryCount)),
			MaxSizeInMegabytes:               utils.Int64(int64(maxSizeInMB)),
			RequiresDuplicateDetection:       utils.Bool(requireDuplicateDetection),
			RequiresSession:                  utils.Bool(requireSession),
			Status:                           &status,
		},
	}

	if autoDeleteOnIdle != "" {
		parameters.Properties.AutoDeleteOnIdle = &autoDeleteOnIdle
	}

	if defaultMessageTTL != "" {
		parameters.Properties.DefaultMessageTimeToLive = &defaultMessageTTL
	}

	if duplicateDetectionHistoryTimeWindow != "" {
		parameters.Properties.DuplicateDetectionHistoryTimeWindow = &duplicateDetectionHistoryTimeWindow
	}

	if forwardDeadLetteredMessagesTo != "" {
		parameters.Properties.ForwardDeadLetteredMessagesTo = &forwardDeadLetteredMessagesTo
	}

	if forwardTo != "" {
		parameters.Properties.ForwardTo = &forwardTo
	}

	if lockDuration != "" {
		parameters.Properties.LockDuration = &lockDuration
	}

	// We need to retrieve the namespace because Premium namespace works differently from Basic and Standard,
	// so it needs different rules applied to it.
	namespacesClient := meta.(*clients.Client).ServiceBus.NamespacesClient
	namespace, err := namespacesClient.Get(ctx, *namespaceId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *namespaceId, err)
	}

	var sku namespaces.SkuName
	if nsModel := namespace.Model; nsModel != nil {
		sku = nsModel.Sku.Name
	}
	// Enforce Premium namespace to have Express Entities disabled in Terraform since they are not supported for
	// Premium SKU.
	if sku == namespaces.SkuNamePremium && enableExpress {
		return fmt.Errorf("%s does not support Express Entities in Premium SKU and must be disabled", id)
	}

	if sku == namespaces.SkuNamePremium {
		if isPremiumNamespacePartitioned && !enablePartitioning {
			return fmt.Errorf("non-partitioned entities are not allowed in partitioned namespace")
		} else if !isPremiumNamespacePartitioned && enablePartitioning {
			return fmt.Errorf("the parent premium namespace is not partitioned and the partitioning for premium namespace is only available at the namepsace creation")
		}
	}

	// output of `max_message_size_in_kilobytes` is also set in non-Premium namespaces, with a value of 256
	if v, ok := d.GetOk("max_message_size_in_kilobytes"); ok && v.(int) != 256 {
		if sku != namespaces.SkuNamePremium {
			return fmt.Errorf("%s does not support input on `max_message_size_in_kilobytes` in %s SKU and should be removed", id, sku)
		}
		parameters.Properties.MaxMessageSizeInKilobytes = utils.Int64(int64(v.(int)))
	}

	if _, err = client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return err
	}

	if !d.IsNewResource() {
		// wait for property update, api issue is being tracked:https://github.com/Azure/azure-rest-api-specs/issues/21445
		log.Printf("[DEBUG] Waiting for %s status to become ready", id)
		deadline, ok := ctx.Deadline()
		if !ok {
			return fmt.Errorf("internal-error: context had no deadline")
		}
		statusPropertyChangeConf := &pluginsdk.StateChangeConf{
			Pending:                   []string{"Updating"},
			Target:                    []string{"Succeeded"},
			Refresh:                   serviceBusQueueStatusRefreshFunc(ctx, client, id, userConfig),
			ContinuousTargetOccurence: 5,
			Timeout:                   time.Until(deadline),
			MinTimeout:                1 * time.Minute,
		}

		if _, err = statusPropertyChangeConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for status of %s to become ready: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	return resourceServiceBusQueueRead(d, meta)
}

func resourceServiceBusQueueRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.QueuesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)

	d.Set("name", id.QueueName)
	d.Set("namespace_id", namespaceId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("auto_delete_on_idle", props.AutoDeleteOnIdle)
			d.Set("dead_lettering_on_message_expiration", props.DeadLetteringOnMessageExpiration)
			d.Set("default_message_ttl", props.DefaultMessageTimeToLive)
			d.Set("duplicate_detection_history_time_window", props.DuplicateDetectionHistoryTimeWindow)
			d.Set("forward_dead_lettered_messages_to", props.ForwardDeadLetteredMessagesTo)
			d.Set("forward_to", props.ForwardTo)
			d.Set("lock_duration", props.LockDuration)
			d.Set("max_delivery_count", props.MaxDeliveryCount)
			d.Set("max_message_size_in_kilobytes", props.MaxMessageSizeInKilobytes)
			d.Set("requires_duplicate_detection", props.RequiresDuplicateDetection)
			d.Set("requires_session", props.RequiresSession)
			d.Set("status", string(pointer.From(props.Status)))

			d.Set("batched_operations_enabled", props.EnableBatchedOperations)
			d.Set("express_enabled", props.EnableExpress)
			d.Set("partitioning_enabled", props.EnablePartitioning)

			if apiMaxSizeInMegabytes := props.MaxSizeInMegabytes; apiMaxSizeInMegabytes != nil {
				maxSizeInMegabytes := int(*apiMaxSizeInMegabytes)

				// If the queue is NOT in a premium namespace (ie. it is Basic or Standard) and partitioning is enabled
				// then the max size returned by the API will be 16 times greater than the value set.
				if *props.EnablePartitioning {
					namespacesClient := meta.(*clients.Client).ServiceBus.NamespacesClient
					namespace, err := namespacesClient.Get(ctx, namespaceId)
					if err != nil {
						return err
					}

					if nsModel := namespace.Model; nsModel != nil && nsModel.Sku.Name != namespaces.SkuNamePremium {
						const partitionCount = 16
						maxSizeInMegabytes = int(*apiMaxSizeInMegabytes / partitionCount)
					}
				}

				d.Set("max_size_in_megabytes", maxSizeInMegabytes)
			}
		}
	}
	return nil
}

func resourceServiceBusQueueDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.QueuesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := queues.ParseQueueID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}

func serviceBusQueueStatusRefreshFunc(ctx context.Context, client *queues.QueuesClient, id queues.QueueId, userConfig map[string]interface{}) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking servicebus queue %s status...", id)

		resp, err := client.Get(ctx, id)
		if err != nil {
			return nil, "Error", fmt.Errorf("checking servicebus queue %s error: %+v", id, err)
		}

		queueStatus := "Updating"

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				if props.Status != nil && *props.Status == userConfig["status"] &&
					props.MaxDeliveryCount != nil && int(*props.MaxDeliveryCount) == userConfig["maxDeliveryCount"].(int) &&
					props.EnableExpress != nil && *props.EnableExpress == userConfig["enableExpress"].(bool) &&
					props.EnablePartitioning != nil && *props.EnablePartitioning == userConfig["enablePartitioning"].(bool) &&
					props.EnableBatchedOperations != nil && *props.EnableBatchedOperations == userConfig["enableBatchOps"].(bool) {
					queueStatus = "Succeeded"
				}

				if props.DeadLetteringOnMessageExpiration != nil && userConfig["deadLetteringOnMesExp"] != "" {
					if *props.DeadLetteringOnMessageExpiration != userConfig["deadLetteringOnMesExp"].(bool) {
						queueStatus = "Updating"
					}
				}

				if props.LockDuration != nil && userConfig["lockDuration"] != "" {
					if *props.LockDuration != userConfig["lockDuration"].(string) {
						queueStatus = "Updating"
					}
				}

				if props.ForwardTo != nil && userConfig["forwardTo"] != "" {
					if *props.ForwardTo != userConfig["forwardTo"].(string) {
						queueStatus = "Updating"
					}
				}

				if props.ForwardDeadLetteredMessagesTo != nil && userConfig["forwardDeadLetteredMessagesTo"] != nil {
					if *props.ForwardDeadLetteredMessagesTo != userConfig["forwardDeadLetteredMessagesTo"].(string) {
						queueStatus = "Updating"
					}
				}
			}
		}
		return resp, queueStatus, nil
	}
}
