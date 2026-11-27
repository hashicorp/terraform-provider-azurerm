// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_topic

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/topics"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	servicebusClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusTopicCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := servicebusClient.From(meta).TopicsClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	var id topics.TopicId
	if namespaceIdLit := d.Get("namespace_id").(string); namespaceIdLit != "" {
		namespaceId, err := topics.ParseNamespaceID(namespaceIdLit)
		if err != nil {
			return err
		}
		id = topics.NewTopicID(namespaceId.SubscriptionId, namespaceId.ResourceGroupName, namespaceId.NamespaceName, d.Get("name").(string))
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
				return tf.ImportAsExistsError("azurerm_servicebus_topic", id.ID())
			}
		}
	}

	enableBatchedOperations := d.Get("batched_operations_enabled").(bool)
	enableExpress := d.Get("express_enabled").(bool)
	enablePartitioning := d.Get("partitioning_enabled").(bool)

	status := topics.EntityStatus(d.Get("status").(string))
	parameters := topics.SBTopic{
		Name: pointer.To(id.TopicName),
		Properties: &topics.SBTopicProperties{
			Status:                     &status,
			EnableBatchedOperations:    pointer.To(enableBatchedOperations),
			EnableExpress:              pointer.To(enableExpress),
			EnablePartitioning:         pointer.To(enablePartitioning),
			MaxSizeInMegabytes:         pointer.To(int64(d.Get("max_size_in_megabytes").(int))),
			RequiresDuplicateDetection: pointer.To(d.Get("requires_duplicate_detection").(bool)),
			SupportOrdering:            pointer.To(d.Get("support_ordering").(bool)),
		},
	}

	if autoDeleteOnIdle := d.Get("auto_delete_on_idle").(string); autoDeleteOnIdle != "" {
		parameters.Properties.AutoDeleteOnIdle = pointer.To(autoDeleteOnIdle)
	}

	if defaultTTL := d.Get("default_message_ttl").(string); defaultTTL != "" {
		parameters.Properties.DefaultMessageTimeToLive = pointer.To(defaultTTL)
	}

	if duplicateWindow := d.Get("duplicate_detection_history_time_window").(string); duplicateWindow != "" {
		parameters.Properties.DuplicateDetectionHistoryTimeWindow = pointer.To(duplicateWindow)
	}

	// We need to retrieve the namespace because Premium namespace works differently from Basic and Standard
	namespacesClient := servicebusClient.From(meta).NamespacesClient
	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	resp, err := namespacesClient.Get(ctx, namespaceId)
	if err != nil {
		return fmt.Errorf("retrieving ServiceBus Namespace %q (Resource Group %q): %+v", id.NamespaceName, id.ResourceGroupName, err)
	}

	isPremiumNamespacePartitioned := true
	var sku namespaces.SkuName
	if nsModel := resp.Model; nsModel != nil {
		sku = nsModel.Sku.Name
		if props := nsModel.Properties; props != nil && props.PremiumMessagingPartitions != nil && *props.PremiumMessagingPartitions == 1 {
			isPremiumNamespacePartitioned = false
		}
	}

	if sku == namespaces.SkuNamePremium {
		if isPremiumNamespacePartitioned && !enablePartitioning {
			return errors.New("topic must have `partitioning_enabled` set to `true` when the parent namespace is partitioned")
		} else if !isPremiumNamespacePartitioned && enablePartitioning {
			return errors.New("topic partitioning is only available if the parent namespace is partitioned")
		}
	}

	// output of `max_message_size_in_kilobytes` is also set in non-Premium namespaces, with a value of 256
	if v, ok := d.GetOk("max_message_size_in_kilobytes"); ok && v.(int) != 256 {
		if model := resp.Model; model != nil {
			if model.Sku.Name != namespaces.SkuNamePremium {
				return fmt.Errorf("%s does not support input on `max_message_size_in_kilobytes` in %s SKU and should be removed", id, model.Sku.Name)
			}
			parameters.Properties.MaxMessageSizeInKilobytes = pointer.To(int64(v.(int)))
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceServiceBusTopicRead(d, meta)
}
