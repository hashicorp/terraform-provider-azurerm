package servicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceServiceBusSubscription() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceServiceBusSubscriptionRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"topic_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.TopicID,
				AtLeastOneOf: []string{"topic_id", "resource_group_name", "namespace_name", "topic_name"},
			},

			"namespace_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.NamespaceName,
				AtLeastOneOf: []string{"topic_id", "resource_group_name", "namespace_name", "topic_name"},
			},

			"resource_group_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: resourcegroups.ValidateName,
				AtLeastOneOf: []string{"topic_id", "resource_group_name", "namespace_name", "topic_name"},
			},

			"topic_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.TopicName(),
				AtLeastOneOf: []string{"topic_id", "resource_group_name", "namespace_name", "topic_name"},
			},

			"auto_delete_on_idle": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"default_message_ttl": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"lock_duration": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"dead_lettering_on_message_expiration": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"dead_lettering_on_filter_evaluation_error": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_batched_operations": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"max_delivery_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"requires_session": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"forward_to": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"forward_dead_lettered_messages_to": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServiceBusSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.SubscriptionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var rgName string
	var nsName string
	var topicName string
	if v, ok := d.Get("topic_id").(string); ok && v != "" {
		topicId, err := parse.TopicID(v)
		if err != nil {
			return fmt.Errorf("parsing topic ID %q: %+v", v, err)
		}
		rgName = topicId.ResourceGroup
		nsName = topicId.NamespaceName
		topicName = topicId.Name
	} else {
		rgName = d.Get("resource_group_name").(string)
		nsName = d.Get("namespace_name").(string)
		topicName = d.Get("topic_name").(string)
	}

	id := parse.NewSubscriptionID(subscriptionId, rgName, nsName, topicName, d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.NamespaceName, id.TopicName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if props := existing.SBSubscriptionProperties; props != nil {
		d.Set("auto_delete_on_idle", props.AutoDeleteOnIdle)
		d.Set("default_message_ttl", props.DefaultMessageTimeToLive)
		d.Set("lock_duration", props.LockDuration)
		d.Set("dead_lettering_on_message_expiration", props.DeadLetteringOnMessageExpiration)
		d.Set("dead_lettering_on_filter_evaluation_error", props.DeadLetteringOnFilterEvaluationExceptions)
		d.Set("enable_batched_operations", props.EnableBatchedOperations)
		d.Set("requires_session", props.RequiresSession)
		d.Set("forward_dead_lettered_messages_to", props.ForwardDeadLetteredMessagesTo)
		d.Set("forward_to", props.ForwardTo)

		maxDeliveryCount := 0
		if props.MaxDeliveryCount != nil {
			maxDeliveryCount = int(*props.MaxDeliveryCount)
		}

		d.Set("max_delivery_count", maxDeliveryCount)
	}

	return nil
}
