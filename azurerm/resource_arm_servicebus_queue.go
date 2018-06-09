package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceBusQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusQueueCreateUpdate,
		Read:   resourceArmServiceBusQueueRead,
		Update: resourceArmServiceBusQueueCreateUpdate,
		Delete: resourceArmServiceBusQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": deprecatedLocationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"auto_delete_on_idle": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"default_message_ttl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"duplicate_detection_history_time_window": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enable_express": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"enable_partitioning": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
				ForceNew: true,
			},

			"lock_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_size_in_megabytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"requires_duplicate_detection": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
				ForceNew: true,
			},

			"requires_session": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"dead_lettering_on_message_expiration": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			// TODO: remove these in the next major release
			"enable_batched_operations": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "This field has been removed by Azure.",
			},
			"support_ordering": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "This field has been removed by Azure.",
			},
		},
	}
}

func resourceArmServiceBusQueueCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusQueuesClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM ServiceBus Queue creation/update.")

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	enableExpress := d.Get("enable_express").(bool)
	enablePartitioning := d.Get("enable_partitioning").(bool)
	maxSize := int32(d.Get("max_size_in_megabytes").(int))
	requiresDuplicateDetection := d.Get("requires_duplicate_detection").(bool)
	requiresSession := d.Get("requires_session").(bool)
	deadLetteringOnMessageExpiration := d.Get("dead_lettering_on_message_expiration").(bool)

	parameters := servicebus.SBQueue{
		Name: &name,
		SBQueueProperties: &servicebus.SBQueueProperties{
			EnableExpress:                    &enableExpress,
			EnablePartitioning:               &enablePartitioning,
			MaxSizeInMegabytes:               &maxSize,
			RequiresDuplicateDetection:       &requiresDuplicateDetection,
			RequiresSession:                  &requiresSession,
			DeadLetteringOnMessageExpiration: &deadLetteringOnMessageExpiration,
		},
	}

	if autoDeleteOnIdle := d.Get("auto_delete_on_idle").(string); autoDeleteOnIdle != "" {
		parameters.SBQueueProperties.AutoDeleteOnIdle = &autoDeleteOnIdle
	}

	if defaultTTL := d.Get("default_message_ttl").(string); defaultTTL != "" {
		parameters.SBQueueProperties.DefaultMessageTimeToLive = &defaultTTL
	}

	if duplicateWindow := d.Get("duplicate_detection_history_time_window").(string); duplicateWindow != "" {
		parameters.SBQueueProperties.DuplicateDetectionHistoryTimeWindow = &duplicateWindow
	}

	if lockDuration := d.Get("lock_duration").(string); lockDuration != "" {
		parameters.SBQueueProperties.LockDuration = &lockDuration
	}

	// We need to retrieve the namespace because Premium namespace works differently from Basic and Standard,
	// so it needs different rules applied to it.
	namespacesClient := meta.(*ArmClient).serviceBusNamespacesClient
	namespace, nsErr := namespacesClient.Get(ctx, resourceGroup, namespaceName)
	if nsErr != nil {
		return nsErr
	}

	// Enforce Premium namespace to have partitioning enabled in Terraform. It is always enabled in Azure for
	// Premium SKU.
	if namespace.Sku.Name == servicebus.Premium && !d.Get("enable_partitioning").(bool) {
		return fmt.Errorf("ServiceBus Queue (%s) must have Partitioning enabled for Premium SKU", name)
	}

	// Enforce Premium namespace to have Express Entities disabled in Terraform since they are not supported for
	// Premium SKU.
	if namespace.Sku.Name == servicebus.Premium && d.Get("enable_express").(bool) {
		return fmt.Errorf("ServiceBus Queue (%s) does not support Express Entities in Premium SKU and must be disabled", name)
	}

	_, err := client.CreateOrUpdate(ctx, resourceGroup, namespaceName, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ServiceBus Queue %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceBusQueueRead(d, meta)
}

func resourceArmServiceBusQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusQueuesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["queues"]

	resp, err := client.Get(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure ServiceBus Queue %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("namespace_name", namespaceName)

	if props := resp.SBQueueProperties; props != nil {
		d.Set("auto_delete_on_idle", props.AutoDeleteOnIdle)
		d.Set("default_message_ttl", props.DefaultMessageTimeToLive)
		d.Set("duplicate_detection_history_time_window", props.DuplicateDetectionHistoryTimeWindow)
		d.Set("lock_duration", props.LockDuration)

		d.Set("enable_express", props.EnableExpress)
		d.Set("enable_partitioning", props.EnablePartitioning)
		d.Set("requires_duplicate_detection", props.RequiresDuplicateDetection)
		d.Set("requires_session", props.RequiresSession)
		d.Set("dead_lettering_on_message_expiration", props.DeadLetteringOnMessageExpiration)

		if maxSizeMB := props.MaxSizeInMegabytes; maxSizeMB != nil {
			maxSize := int(*maxSizeMB)

			// If the queue is NOT in a premium namespace (ie. it is Basic or Standard) and partitioning is enabled
			// then the max size returned by the API will be 16 times greater than the value set.
			if *props.EnablePartitioning {
				namespacesClient := meta.(*ArmClient).serviceBusNamespacesClient
				namespace, err := namespacesClient.Get(ctx, resourceGroup, namespaceName)
				if err != nil {
					return err
				}

				if namespace.Sku.Name != servicebus.Premium {
					const partitionCount = 16
					maxSize = int(*props.MaxSizeInMegabytes / partitionCount)
				}
			}

			d.Set("max_size_in_megabytes", maxSize)
		}
	}

	return nil
}

func resourceArmServiceBusQueueDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusQueuesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["queues"]

	resp, err := client.Delete(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}
