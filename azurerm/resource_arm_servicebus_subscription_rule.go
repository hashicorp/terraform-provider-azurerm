package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceBusSubscriptionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusSubscriptionRuleCreate,
		Read:   resourceArmServiceBusSubscriptionRuleRead,
		Update: resourceArmServiceBusSubscriptionRuleCreate,
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

			"subscription_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"filter_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(servicebus.FilterTypeSQLFilter),
					string(servicebus.FilterTypeCorrelationFilter),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"resource_group_name": resourceGroupNameSchema(),

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

			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceArmServiceBusSubscriptionRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusSubscriptionRulesClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM Service Bus Subscription Rule creation.")

	name := d.Get("name").(string)
	topicName := d.Get("topic_name").(string)
	subscriptionName := d.Get("subscription_name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	rule, err := buildAzureRmServiceBusSubscriptionRule(name, d)
	if err != nil {
		return err
	}

	err = validateArmServiceBusSubscriptionRule(name, rule)
	if err != nil {
		return err
	}

	_, err = client.CreateOrUpdate(ctx, resourceGroup, namespaceName, topicName, subscriptionName, name, rule)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, namespaceName, topicName, subscriptionName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Service Bus Subscription Rule %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceBusSubscriptionRuleRead(d, meta)
}

func resourceArmServiceBusSubscriptionRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusSubscriptionRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
	client := meta.(*ArmClient).serviceBusSubscriptionRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	topicName := id.Path["topics"]
	subscriptionName := id.Path["subscriptions"]
	name := id.Path["rules"]

	_, err = client.Delete(ctx, resourceGroup, namespaceName, topicName, subscriptionName, name)

	return err
}

func buildAzureRmServiceBusSubscriptionRule(name string, d *schema.ResourceData) (rule servicebus.Rule, err error) {
	correlationfilter, err := getAzureRmServiceBusCorrelationFilter(name, d)
	if err != nil {
		return
	}

	filterType := d.Get("filter_type").(string)
	rule = servicebus.Rule{
		Ruleproperties: &servicebus.Ruleproperties{
			FilterType:        servicebus.FilterType(filterType),
			Action:            getAzureRmServiceBusAction(d),
			SQLFilter:         getAzureRmServiceBusSQLFilter(d),
			CorrelationFilter: correlationfilter,
		},
	}

	return
}

func getAzureRmServiceBusAction(d *schema.ResourceData) *servicebus.Action {
	if action := d.Get("action").(string); action != "" {
		serviceBusAction := servicebus.Action{
			SQLExpression: &action,
		}
		return &serviceBusAction
	}
	return nil
}

func getAzureRmServiceBusSQLFilter(d *schema.ResourceData) *servicebus.SQLFilter {
	if sqlFilter := d.Get("sql_filter").(string); sqlFilter != "" {
		serviceBusSQLFilter := servicebus.SQLFilter{
			SQLExpression: &sqlFilter,
		}
		return &serviceBusSQLFilter
	}
	return nil
}

func getAzureRmServiceBusCorrelationFilter(name string, d *schema.ResourceData) (*servicebus.CorrelationFilter, error) {
	if config := d.Get("correlation_filter").([]interface{}); len(config) > 0 {

		if config[0] == nil {
			err := fmt.Errorf("Cannot create Service Bus Subscription Rule (%s). At least one property must be set in the `correlation_filter` block", name)
			return nil, err
		}

		filter := config[0].(map[string]interface{})
		correlationFilter := servicebus.CorrelationFilter{}

		if correlationID := filter["correlation_id"].(string); correlationID != "" {
			correlationFilter.CorrelationID = &correlationID
		}

		if messageID := filter["message_id"].(string); messageID != "" {
			correlationFilter.MessageID = &messageID
		}

		if to := filter["to"].(string); to != "" {
			correlationFilter.To = &to
		}

		if replyTo := filter["reply_to"].(string); replyTo != "" {
			correlationFilter.ReplyTo = &replyTo
		}

		if label := filter["label"].(string); label != "" {
			correlationFilter.Label = &label
		}

		if sessionID := filter["session_id"].(string); sessionID != "" {
			correlationFilter.SessionID = &sessionID
		}

		if replyToSessionID := filter["reply_to_session_id"].(string); replyToSessionID != "" {
			correlationFilter.ReplyToSessionID = &replyToSessionID
		}

		if contentType := filter["content_type"].(string); contentType != "" {
			correlationFilter.ContentType = &contentType
		}

		return &correlationFilter, nil
	}

	return nil, nil
}

func flattenAzureRmServiceBusCorrelationFilter(f *servicebus.CorrelationFilter) []interface{} {
	filters := make([]interface{}, 0, 1)
	if f == nil {
		return filters
	}

	filter := make(map[string]interface{})

	if f.CorrelationID != nil {
		filter["correlation_id"] = *f.CorrelationID
	}

	if f.MessageID != nil {
		filter["message_id"] = *f.MessageID
	}

	if f.To != nil {
		filter["to"] = *f.To
	}

	if f.ReplyTo != nil {
		filter["reply_to"] = *f.ReplyTo
	}

	if f.Label != nil {
		filter["label"] = *f.Label
	}

	if f.SessionID != nil {
		filter["session_id"] = *f.SessionID
	}

	if f.ReplyToSessionID != nil {
		filter["reply_to_session_id"] = *f.ReplyToSessionID
	}

	if f.ContentType != nil {
		filter["content_type"] = *f.ContentType
	}

	filters = append(filters, filter)
	return filters
}

func validateArmServiceBusSubscriptionRule(name string, rule servicebus.Rule) error {
	if rule.Ruleproperties.FilterType == servicebus.FilterTypeSQLFilter {
		if rule.Ruleproperties.SQLFilter == nil {
			return fmt.Errorf("Cannot create Service Bus Subscription Rule (%s). 'sql_filter' must be specified when 'filter_type' is set to 'SqlFilter'", name)
		}
		if rule.Ruleproperties.CorrelationFilter != nil {
			return fmt.Errorf("Service Bus Subscription Rule (%s) does not support `correlation_filter` when 'filter_type' is set to 'SqlFilter'", name)
		}
	}

	if rule.Ruleproperties.FilterType == servicebus.FilterTypeCorrelationFilter {
		if rule.Ruleproperties.CorrelationFilter == nil {
			return fmt.Errorf("Cannot create Service Bus Subscription Rule (%s). 'correlation_filter' must be specified when 'filter_type' is set to 'CorrelationFilter'", name)
		}
		if rule.Ruleproperties.SQLFilter != nil {
			return fmt.Errorf("Service Bus Subscription Rule (%s) does not support `sql_filter` when 'filter_type' is set to 'CorrelationFilter'", name)
		}
	}

	return nil
}
