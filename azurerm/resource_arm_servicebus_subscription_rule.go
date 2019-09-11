package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceBusSubscriptionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusSubscriptionRuleCreateUpdate,
		Read:   resourceArmServiceBusSubscriptionRuleRead,
		Update: resourceArmServiceBusSubscriptionRuleCreateUpdate,
		Delete: resourceArmServiceBusSubscriptionRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

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

			"subscription_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceBusSubscriptionName(),
			},

			"filter_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(servicebus.FilterTypeSQLFilter),
					string(servicebus.FilterTypeCorrelationFilter),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"sql_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"correlation_filter": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"sql_filter"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"correlation_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"message_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"to": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"reply_to": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"label": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"session_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"reply_to_session_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmServiceBusSubscriptionRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicebus.SubscriptionRulesClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure Service Bus Subscription Rule creation.")

	name := d.Get("name").(string)
	topicName := d.Get("topic_name").(string)
	subscriptionName := d.Get("subscription_name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	filterType := d.Get("filter_type").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, namespaceName, topicName, subscriptionName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Service Bus Subscription %q (Resource Group %q, namespace %q): %+v", name, resourceGroup, namespaceName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_servicebus_subscription_rule", *existing.ID)
		}
	}

	rule := servicebus.Rule{
		Ruleproperties: &servicebus.Ruleproperties{
			FilterType: servicebus.FilterType(filterType),
		},
	}

	if action := d.Get("action").(string); action != "" {
		rule.Ruleproperties.Action = &servicebus.Action{
			SQLExpression: &action,
		}
	}

	if rule.Ruleproperties.FilterType == servicebus.FilterTypeCorrelationFilter {
		correlationFilter, err := expandAzureRmServiceBusCorrelationFilter(d)
		if err != nil {
			return fmt.Errorf("Cannot create Service Bus Subscription Rule %q: %+v", name, err)
		}

		rule.Ruleproperties.CorrelationFilter = correlationFilter
	}

	if rule.Ruleproperties.FilterType == servicebus.FilterTypeSQLFilter {
		sqlFilter := d.Get("sql_filter").(string)
		rule.Ruleproperties.SQLFilter = &servicebus.SQLFilter{
			SQLExpression: &sqlFilter,
		}
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, namespaceName, topicName, subscriptionName, name, rule); err != nil {
		return fmt.Errorf("Error issuing create/update request for Service Bus Subscription %q (Resource Group %q, namespace %q): %+v", name, resourceGroup, namespaceName, err)
	}

	read, err := client.Get(ctx, resourceGroup, namespaceName, topicName, subscriptionName, name)
	if err != nil {
		return fmt.Errorf("Error issuing get request for Service Bus Subscription %q (Resource Group %q, namespace %q): %+v", name, resourceGroup, namespaceName, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Service Bus Subscription Rule %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceBusSubscriptionRuleRead(d, meta)
}

func resourceArmServiceBusSubscriptionRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicebus.SubscriptionRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	topicName := id.Path["topics"]
	subscriptionName := id.Path["subscriptions"]
	name := id.Path["rules"]

	resp, err := client.Get(ctx, resourceGroup, namespaceName, topicName, subscriptionName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Service Bus Subscription Rule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("namespace_name", namespaceName)
	d.Set("topic_name", topicName)
	d.Set("subscription_name", subscriptionName)

	if properties := resp.Ruleproperties; properties != nil {
		d.Set("filter_type", properties.FilterType)

		if properties.Action != nil {
			d.Set("action", properties.Action.SQLExpression)
		}

		if properties.SQLFilter != nil {
			d.Set("sql_filter", properties.SQLFilter.SQLExpression)
		}

		if err := d.Set("correlation_filter", flattenAzureRmServiceBusCorrelationFilter(properties.CorrelationFilter)); err != nil {
			return fmt.Errorf("Error setting `correlation_filter` on Azure Service Bus Subscription Rule (%q): %+v", name, err)
		}
	}

	return nil
}

func resourceArmServiceBusSubscriptionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicebus.SubscriptionRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	topicName := id.Path["topics"]
	subscriptionName := id.Path["subscriptions"]
	name := id.Path["rules"]

	resp, err := client.Delete(ctx, resourceGroup, namespaceName, topicName, subscriptionName, name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting ServiceBus Subscription Rule %q: %+v", name, err)
		}
	}

	return nil
}

func expandAzureRmServiceBusCorrelationFilter(d *schema.ResourceData) (*servicebus.CorrelationFilter, error) {
	configs := d.Get("correlation_filter").([]interface{})
	if len(configs) == 0 {
		return nil, fmt.Errorf("`correlation_filter` is required when `filter_type` is set to `CorrelationFilter`")
	}

	config := configs[0].(map[string]interface{})

	contentType := config["content_type"].(string)
	correlationID := config["correlation_id"].(string)
	label := config["label"].(string)
	messageID := config["message_id"].(string)
	replyTo := config["reply_to"].(string)
	replyToSessionID := config["reply_to_session_id"].(string)
	sessionID := config["session_id"].(string)
	to := config["to"].(string)

	if contentType == "" && correlationID == "" && label == "" && messageID == "" && replyTo == "" && replyToSessionID == "" && sessionID == "" && to == "" {
		return nil, fmt.Errorf("At least one property must be set in the `correlation_filter` block")
	}

	correlationFilter := servicebus.CorrelationFilter{}

	if correlationID != "" {
		correlationFilter.CorrelationID = utils.String(correlationID)
	}

	if messageID != "" {
		correlationFilter.MessageID = utils.String(messageID)
	}

	if to != "" {
		correlationFilter.To = utils.String(to)
	}

	if replyTo != "" {
		correlationFilter.ReplyTo = utils.String(replyTo)
	}

	if label != "" {
		correlationFilter.Label = utils.String(label)
	}

	if sessionID != "" {
		correlationFilter.SessionID = utils.String(sessionID)
	}

	if replyToSessionID != "" {
		correlationFilter.ReplyToSessionID = utils.String(replyToSessionID)
	}

	if contentType != "" {
		correlationFilter.ContentType = utils.String(contentType)
	}

	return &correlationFilter, nil
}

func flattenAzureRmServiceBusCorrelationFilter(input *servicebus.CorrelationFilter) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	filter := make(map[string]interface{})

	if input.CorrelationID != nil {
		filter["correlation_id"] = *input.CorrelationID
	}

	if input.MessageID != nil {
		filter["message_id"] = *input.MessageID
	}

	if input.To != nil {
		filter["to"] = *input.To
	}

	if input.ReplyTo != nil {
		filter["reply_to"] = *input.ReplyTo
	}

	if input.Label != nil {
		filter["label"] = *input.Label
	}

	if input.SessionID != nil {
		filter["session_id"] = *input.SessionID
	}

	if input.ReplyToSessionID != nil {
		filter["reply_to_session_id"] = *input.ReplyToSessionID
	}

	if input.ContentType != nil {
		filter["content_type"] = *input.ContentType
	}

	return []interface{}{filter}
}
