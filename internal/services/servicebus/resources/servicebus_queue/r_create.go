// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_queue

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusQueueCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(serviceBusQueueResourceName, id.ID())
			}
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

	userConfig := make(map[string]any)

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
		Name: pointer.To(id.QueueName),
		Properties: &queues.SBQueueProperties{
			DeadLetteringOnMessageExpiration: pointer.To(deadLetteringOnMesExp),
			EnableBatchedOperations:          pointer.To(enableBatchedOperations),
			EnableExpress:                    pointer.To(enableExpress),
			EnablePartitioning:               pointer.To(enablePartitioning),
			MaxDeliveryCount:                 pointer.To(int64(maxDeliveryCount)),
			MaxSizeInMegabytes:               pointer.To(int64(maxSizeInMB)),
			RequiresDuplicateDetection:       pointer.To(requireDuplicateDetection),
			RequiresSession:                  pointer.To(requireSession),
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
			return fmt.Errorf("the parent premium namespace is not partitioned and the partitioning for premium namespace is only available at the namespace creation")
		}
	}

	// output of `max_message_size_in_kilobytes` is also set in non-Premium namespaces, with a value of 256
	if v, ok := d.GetOk("max_message_size_in_kilobytes"); ok && v.(int) != 256 {
		if sku != namespaces.SkuNamePremium {
			return fmt.Errorf("%s does not support input on `max_message_size_in_kilobytes` in %s SKU and should be removed", id, sku)
		}
		parameters.Properties.MaxMessageSizeInKilobytes = pointer.To(int64(v.(int)))
	}

	if _, err = client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return err
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
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

	return resourceServiceBusQueueRead(d, meta)
}
