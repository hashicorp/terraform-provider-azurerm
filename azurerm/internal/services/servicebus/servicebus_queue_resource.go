package servicebus

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceBusQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusQueueCreateUpdate,
		Read:   resourceArmServiceBusQueueRead,
		Update: resourceArmServiceBusQueueCreateUpdate,
		Delete: resourceArmServiceBusQueueDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.QueueID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azValidate.QueueName(),
			},

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azValidate.NamespaceName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			// Optional
			"auto_delete_on_idle": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ISO8601Duration,
			},

			"dead_lettering_on_message_expiration": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"default_message_ttl": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ISO8601Duration,
			},

			"duplicate_detection_history_time_window": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ISO8601Duration,
			},

			"enable_batched_operations": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"enable_express": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enable_partitioning": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"forward_dead_lettered_messages_to": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azValidate.QueueName(),
			},

			"forward_to": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azValidate.QueueName(),
			},

			"lock_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_delivery_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"max_size_in_megabytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ServiceBusMaxSizeInMegabytes(),
			},

			"requires_duplicate_detection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"requires_session": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(servicebus.Active),
				ValidateFunc: validation.StringInSlice([]string{
					string(servicebus.Active),
					string(servicebus.Creating),
					string(servicebus.Deleting),
					string(servicebus.Disabled),
					string(servicebus.ReceiveDisabled),
					string(servicebus.Renaming),
					string(servicebus.SendDisabled),
					string(servicebus.Unknown),
				}, false),
			},
		},
	}
}

func resourceArmServiceBusQueueCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.QueuesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for ServiceBus Queue creation/update.")

	deadLetteringOnMessageExpiration := d.Get("dead_lettering_on_message_expiration").(bool)
	enableBatchedOperations := d.Get("enable_batched_operations").(bool)
	enableExpress := d.Get("enable_express").(bool)
	enablePartitioning := d.Get("enable_partitioning").(bool)
	maxDeliveryCount := int32(d.Get("max_delivery_count").(int))
	maxSizeInMegabytes := int32(d.Get("max_size_in_megabytes").(int))
	requiresDuplicateDetection := d.Get("requires_duplicate_detection").(bool)
	requiresSession := d.Get("requires_session").(bool)
	status := servicebus.EntityStatus(d.Get("status").(string))

	resourceId := parse.NewQueueID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.NamespaceName, resourceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", resourceId, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_servicebus_queue", resourceId.ID(""))
		}
	}

	parameters := servicebus.SBQueue{
		Name: utils.String(resourceId.Name),
		SBQueueProperties: &servicebus.SBQueueProperties{
			DeadLetteringOnMessageExpiration: &deadLetteringOnMessageExpiration,
			EnableBatchedOperations:          &enableBatchedOperations,
			EnableExpress:                    &enableExpress,
			EnablePartitioning:               &enablePartitioning,
			MaxDeliveryCount:                 &maxDeliveryCount,
			MaxSizeInMegabytes:               &maxSizeInMegabytes,
			RequiresDuplicateDetection:       &requiresDuplicateDetection,
			RequiresSession:                  &requiresSession,
			Status:                           status,
		},
	}

	if autoDeleteOnIdle := d.Get("auto_delete_on_idle").(string); autoDeleteOnIdle != "" {
		parameters.SBQueueProperties.AutoDeleteOnIdle = &autoDeleteOnIdle
	}

	if defaultMessageTTL := d.Get("default_message_ttl").(string); defaultMessageTTL != "" {
		parameters.SBQueueProperties.DefaultMessageTimeToLive = &defaultMessageTTL
	}

	if duplicateDetectionHistoryTimeWindow := d.Get("duplicate_detection_history_time_window").(string); duplicateDetectionHistoryTimeWindow != "" {
		parameters.SBQueueProperties.DuplicateDetectionHistoryTimeWindow = &duplicateDetectionHistoryTimeWindow
	}

	if forwardDeadLetteredMessagesTo := d.Get("forward_dead_lettered_messages_to").(string); forwardDeadLetteredMessagesTo != "" {
		parameters.SBQueueProperties.ForwardDeadLetteredMessagesTo = &forwardDeadLetteredMessagesTo
	}

	if forwardTo := d.Get("forward_to").(string); forwardTo != "" {
		parameters.SBQueueProperties.ForwardTo = &forwardTo
	}

	if lockDuration := d.Get("lock_duration").(string); lockDuration != "" {
		parameters.SBQueueProperties.LockDuration = &lockDuration
	}

	// We need to retrieve the namespace because Premium namespace works differently from Basic and Standard,
	// so it needs different rules applied to it.
	namespacesClient := meta.(*clients.Client).ServiceBus.NamespacesClient
	namespace, err := namespacesClient.Get(ctx, resourceId.ResourceGroup, resourceId.NamespaceName)
	if err != nil {
		return fmt.Errorf("retrieving ServiceBus Namespace %q (Resource Group %q): %+v", resourceId.NamespaceName, resourceId.ResourceGroup, err)
	}

	// Enforce Premium namespace to have Express Entities disabled in Terraform since they are not supported for
	// Premium SKU.
	if namespace.Sku.Name == servicebus.Premium && d.Get("enable_express").(bool) {
		return fmt.Errorf("ServiceBus Queue %q does not support Express Entities in Premium SKU and must be disabled", resourceId.Name)
	}

	if _, err = client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.NamespaceName, resourceId.Name, parameters); err != nil {
		return err
	}

	d.SetId(resourceId.ID(""))
	return resourceArmServiceBusQueueRead(d, meta)
}

func resourceArmServiceBusQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.QueuesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.QueueID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NamespaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.SBQueueProperties; props != nil {
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
		d.Set("requires_duplicate_detection", props.RequiresDuplicateDetection)
		d.Set("requires_session", props.RequiresSession)
		d.Set("status", props.Status)

		if apiMaxSizeInMegabytes := props.MaxSizeInMegabytes; apiMaxSizeInMegabytes != nil {
			maxSizeInMegabytes := int(*apiMaxSizeInMegabytes)

			// If the queue is NOT in a premium namespace (ie. it is Basic or Standard) and partitioning is enabled
			// then the max size returned by the API will be 16 times greater than the value set.
			if *props.EnablePartitioning {
				namespacesClient := meta.(*clients.Client).ServiceBus.NamespacesClient
				namespace, err := namespacesClient.Get(ctx, id.ResourceGroup, id.NamespaceName)
				if err != nil {
					return err
				}

				if namespace.Sku.Name != servicebus.Premium {
					const partitionCount = 16
					maxSizeInMegabytes = int(*apiMaxSizeInMegabytes / partitionCount)
				}
			}

			d.Set("max_size_in_megabytes", maxSizeInMegabytes)
		}
	}

	return nil
}

func resourceArmServiceBusQueueDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.QueuesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.QueueID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.NamespaceName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}
