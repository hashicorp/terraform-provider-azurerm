package servicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queues"
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

		//lintignore: S013
		"namespace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},

		// Optional
		"auto_delete_on_idle": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validate.ISO8601Duration,
		},

		"dead_lettering_on_message_expiration": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"default_message_ttl": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validate.ISO8601Duration,
		},

		"duplicate_detection_history_time_window": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validate.ISO8601Duration,
		},

		// TODO 4.0: change this from enable_* to *_enabled
		"enable_batched_operations": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		// TODO 4.0: change this from enable_* to *_enabled
		"enable_express": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// TODO 4.0: change this from enable_* to *_enabled
		"enable_partitioning": {
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
			Computed: true,
		},

		"max_delivery_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      10,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"max_message_size_in_kilobytes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: azValidate.ServiceBusMaxMessageSizeInKilobytes(),
		},

		"max_size_in_megabytes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
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
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	namespaceId, err := namespaces.ParseNamespaceID(d.Get("namespace_id").(string))
	if err != nil {
		return err
	}

	id := queues.NewQueueID(namespaceId.SubscriptionId, namespaceId.ResourceGroupName, namespaceId.NamespaceName, d.Get("name").(string))

	isPartitioningEnabled := false
	if d.HasChange("enable_partitioning") {
		existingQueue, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existingQueue.HttpResponse) {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
		}

		if model := existingQueue.Model; model != nil {
			if props := model.Properties; props != nil {
				if model.Id != nil && props.EnablePartitioning != nil && *props.EnablePartitioning {
					isPartitioningEnabled = true
				}
			}
		}
	}

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

	status := queues.EntityStatus(d.Get("status").(string))
	parameters := queues.SBQueue{
		Name: utils.String(id.QueueName),
		Properties: &queues.SBQueueProperties{
			DeadLetteringOnMessageExpiration: utils.Bool(d.Get("dead_lettering_on_message_expiration").(bool)),
			EnableBatchedOperations:          utils.Bool(d.Get("enable_batched_operations").(bool)),
			EnableExpress:                    utils.Bool(d.Get("enable_express").(bool)),
			EnablePartitioning:               utils.Bool(d.Get("enable_partitioning").(bool)),
			MaxDeliveryCount:                 utils.Int64(int64(d.Get("max_delivery_count").(int))),
			MaxSizeInMegabytes:               utils.Int64(int64(d.Get("max_size_in_megabytes").(int))),
			RequiresDuplicateDetection:       utils.Bool(d.Get("requires_duplicate_detection").(bool)),
			RequiresSession:                  utils.Bool(d.Get("requires_session").(bool)),
			Status:                           &status,
		},
	}

	if autoDeleteOnIdle := d.Get("auto_delete_on_idle").(string); autoDeleteOnIdle != "" {
		parameters.Properties.AutoDeleteOnIdle = &autoDeleteOnIdle
	}

	if defaultMessageTTL := d.Get("default_message_ttl").(string); defaultMessageTTL != "" {
		parameters.Properties.DefaultMessageTimeToLive = &defaultMessageTTL
	}

	if duplicateDetectionHistoryTimeWindow := d.Get("duplicate_detection_history_time_window").(string); duplicateDetectionHistoryTimeWindow != "" {
		parameters.Properties.DuplicateDetectionHistoryTimeWindow = &duplicateDetectionHistoryTimeWindow
	}

	if forwardDeadLetteredMessagesTo := d.Get("forward_dead_lettered_messages_to").(string); forwardDeadLetteredMessagesTo != "" {
		parameters.Properties.ForwardDeadLetteredMessagesTo = &forwardDeadLetteredMessagesTo
	}

	if forwardTo := d.Get("forward_to").(string); forwardTo != "" {
		parameters.Properties.ForwardTo = &forwardTo
	}

	if lockDuration := d.Get("lock_duration").(string); lockDuration != "" {
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
	if sku == namespaces.SkuNamePremium && d.Get("enable_express").(bool) {
		return fmt.Errorf("%s does not support Express Entities in Premium SKU and must be disabled", id)
	}

	if sku == namespaces.SkuNamePremium && d.Get("enable_partitioning").(bool) && !isPartitioningEnabled {
		return fmt.Errorf("partitioning Entities is not supported in Premium SKU and must be disabled")
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
			d.Set("enable_batched_operations", props.EnableBatchedOperations)
			d.Set("enable_express", props.EnableExpress)
			d.Set("enable_partitioning", props.EnablePartitioning)
			d.Set("forward_dead_lettered_messages_to", props.ForwardDeadLetteredMessagesTo)
			d.Set("forward_to", props.ForwardTo)
			d.Set("lock_duration", props.LockDuration)
			d.Set("max_delivery_count", props.MaxDeliveryCount)
			d.Set("max_message_size_in_kilobytes", props.MaxMessageSizeInKilobytes)
			d.Set("requires_duplicate_detection", props.RequiresDuplicateDetection)
			d.Set("requires_session", props.RequiresSession)
			d.Set("status", props.Status)

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
