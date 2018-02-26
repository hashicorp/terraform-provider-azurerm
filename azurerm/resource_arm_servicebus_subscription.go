package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceBusSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusSubscriptionCreate,
		Read:   resourceArmServiceBusSubscriptionRead,
		Update: resourceArmServiceBusSubscriptionCreate,
		Delete: resourceArmServiceBusSubscriptionDelete,
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

			"topic_name": {
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

			"lock_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dead_lettering_on_message_expiration": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"enable_batched_operations": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"max_delivery_count": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"requires_session": {
				Type:     schema.TypeBool,
				Optional: true,
				// cannot be modified
				ForceNew: true,
			},

			"forward_to": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// TODO: remove in the next major version
			"dead_lettering_on_filter_evaluation_exceptions": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "This field has been deprecated by Azure",
			},
		},
	}
}

func resourceArmServiceBusSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusSubscriptionsClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM ServiceBus Subscription creation.")

	name := d.Get("name").(string)
	topicName := d.Get("topic_name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	deadLetteringExpiration := d.Get("dead_lettering_on_message_expiration").(bool)
	enableBatchedOps := d.Get("enable_batched_operations").(bool)
	maxDeliveryCount := int32(d.Get("max_delivery_count").(int))
	requiresSession := d.Get("requires_session").(bool)

	parameters := servicebus.SBSubscription{
		SBSubscriptionProperties: &servicebus.SBSubscriptionProperties{
			DeadLetteringOnMessageExpiration: &deadLetteringExpiration,
			EnableBatchedOperations:          &enableBatchedOps,
			MaxDeliveryCount:                 &maxDeliveryCount,
			RequiresSession:                  &requiresSession,
		},
	}

	if autoDeleteOnIdle := d.Get("auto_delete_on_idle").(string); autoDeleteOnIdle != "" {
		parameters.SBSubscriptionProperties.AutoDeleteOnIdle = &autoDeleteOnIdle
	}

	if lockDuration := d.Get("lock_duration").(string); lockDuration != "" {
		parameters.SBSubscriptionProperties.LockDuration = &lockDuration
	}

	if forwardTo := d.Get("forward_to").(string); forwardTo != "" {
		parameters.SBSubscriptionProperties.ForwardTo = &forwardTo
	}

	_, err := client.CreateOrUpdate(ctx, resourceGroup, namespaceName, topicName, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, namespaceName, topicName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ServiceBus Subscription %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceBusSubscriptionRead(d, meta)
}

func resourceArmServiceBusSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusSubscriptionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	topicName := id.Path["topics"]
	name := id.Path["subscriptions"]

	resp, err := client.Get(ctx, resourceGroup, namespaceName, topicName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure ServiceBus Subscription %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("namespace_name", namespaceName)
	d.Set("topic_name", topicName)

	if props := resp.SBSubscriptionProperties; props != nil {
		d.Set("auto_delete_on_idle", props.AutoDeleteOnIdle)
		d.Set("default_message_ttl", props.DefaultMessageTimeToLive)
		d.Set("lock_duration", props.LockDuration)
		d.Set("dead_lettering_on_message_expiration", props.DeadLetteringOnMessageExpiration)
		d.Set("enable_batched_operations", props.EnableBatchedOperations)
		d.Set("requires_session", props.RequiresSession)
		d.Set("forward_to", props.ForwardTo)

		if count := props.MaxDeliveryCount; count != nil {
			d.Set("max_delivery_count", int(*count))
		}
	}

	return nil
}

func resourceArmServiceBusSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusSubscriptionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	topicName := id.Path["topics"]
	name := id.Path["subscriptions"]

	_, err = client.Delete(ctx, resourceGroup, namespaceName, topicName, name)

	return err
}
