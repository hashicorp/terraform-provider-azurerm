package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceBusSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusSubscriptionCreateUpdate,
		Read:   resourceArmServiceBusSubscriptionRead,
		Update: resourceArmServiceBusSubscriptionCreateUpdate,
		Delete: resourceArmServiceBusSubscriptionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceBusSubscriptionName(),
			},

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceBusNamespaceName(),
			},

			"topic_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceBusTopicName(),
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

func resourceArmServiceBusSubscriptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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

	if d.IsNewResource() {
		// first check if there's one in this subscription requiring import
		resp, err := client.Get(ctx, resourceGroup, namespaceName, topicName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for the existence of Service Bus Subscription %q (Resource Group %q / Namespace %q / Topic %q): %+v", name, resourceGroup, namespaceName, topicName, err)
			}
		}
		if resp.ID != nil {
			return tf.ImportAsExistsError("azurerm_servicebus_subscription", *resp.ID)
		}
	}

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

	if defaultMessageTtl := d.Get("default_message_ttl").(string); defaultMessageTtl != "" {
		parameters.DefaultMessageTimeToLive = &defaultMessageTtl
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(tf.TimeoutForCreateUpdate(d)))
	defer cancel()
	if _, err := client.CreateOrUpdate(waitCtx, resourceGroup, namespaceName, topicName, name, parameters); err != nil {
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

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	if resp, err := client.Delete(waitCtx, resourceGroup, namespaceName, topicName, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting ServiceBus Subscription %q: %+v", name, err)
		}
	}

	return err
}
