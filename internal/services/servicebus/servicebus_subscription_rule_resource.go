package servicebus

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceServiceBusSubscriptionRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceBusSubscriptionRuleCreateUpdate,
		Read:   resourceServiceBusSubscriptionRuleRead,
		Update: resourceServiceBusSubscriptionRuleCreateUpdate,
		Delete: resourceServiceBusSubscriptionRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := rules.ParseRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceServicebusSubscriptionRuleSchema(),
	}
}

func resourceServicebusSubscriptionRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		//lintignore: S013
		"subscription_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     subscriptions.ValidateSubscriptions2ID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"filter_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(subscriptions.FilterTypeSqlFilter),
				string(subscriptions.FilterTypeCorrelationFilter),
			}, false),
		},

		"action": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"sql_filter": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.SqlFilter,
		},

		"correlation_filter": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"sql_filter"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"correlation_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"message_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"to": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"reply_to": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"label": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"session_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"reply_to_session_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"content_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"properties": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
				},
			},
		},
	}
}

func resourceServiceBusSubscriptionRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.SubscriptionRulesClient
	subscriptionClient := meta.(*clients.Client).ServiceBus.SubscriptionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Azure Service Bus Subscription Rule creation.")

	filterType := d.Get("filter_type").(string)

	var id subscriptions.RuleId
	if subscriptionIdLit := d.Get("subscription_id").(string); subscriptionIdLit != "" {
		subscriptionId, _ := rules.ParseSubscriptions2ID(subscriptionIdLit)
		id = subscriptions.NewRuleID(subscriptionId.SubscriptionId,
			subscriptionId.ResourceGroupName,
			subscriptionId.NamespaceName,
			subscriptionId.TopicName,
			subscriptionId.SubscriptionName,
			d.Get("name").(string),
		)
	}

	if d.IsNewResource() {
		existing, err := subscriptionClient.RulesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_servicebus_subscription_rule", id.ID())
		}
	}

	filter := rules.FilterType(filterType)
	rule := rules.Rule{
		Properties: &rules.Ruleproperties{
			FilterType: &filter,
		},
	}

	if action := d.Get("action").(string); action != "" {
		rule.Properties.Action = &rules.Action{
			SqlExpression: &action,
		}
	}

	if *rule.Properties.FilterType == rules.FilterTypeCorrelationFilter {
		correlationFilter, err := expandAzureRmServiceBusCorrelationFilter(d)
		if err != nil {
			return fmt.Errorf("expanding `correlation_filter`: %+v", err)
		}

		rule.Properties.CorrelationFilter = correlationFilter
	}

	if *rule.Properties.FilterType == rules.FilterTypeSqlFilter {
		sqlFilter := d.Get("sql_filter").(string)
		rule.Properties.SqlFilter = &rules.SqlFilter{
			SqlExpression: &sqlFilter,
		}
	}

	// switch this around so that the rule id from subscription is generated for the get method
	ruleId, err := rules.ParseRuleID(id.ID())
	if err != nil {
		return err
	}
	if _, err := client.CreateOrUpdate(ctx, *ruleId, rule); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceServiceBusSubscriptionRuleRead(d, meta)
}

func resourceServiceBusSubscriptionRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.SubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := subscriptions.ParseRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.RulesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("subscription_id", subscriptions.NewSubscriptions2ID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicName, id.SubscriptionName).ID())
	d.Set("name", id.RuleName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("filter_type", props.FilterType)

			if props.Action != nil {
				d.Set("action", props.Action.SqlExpression)
			}

			if props.SqlFilter != nil {
				d.Set("sql_filter", props.SqlFilter.SqlExpression)
			}

			if err := d.Set("correlation_filter", flattenAzureRmServiceBusCorrelationFilter(props.CorrelationFilter)); err != nil {
				return fmt.Errorf("setting `correlation_filter`: %+v", err)
			}
		}
	}

	return nil
}

func resourceServiceBusSubscriptionRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.SubscriptionRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rules.ParseRuleID(d.Id())
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

func expandAzureRmServiceBusCorrelationFilter(d *pluginsdk.ResourceData) (*rules.CorrelationFilter, error) {
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

	properties := expandProperties(config["properties"].(map[string]interface{}))

	if contentType == "" && correlationID == "" && label == "" && messageID == "" && replyTo == "" && replyToSessionID == "" && sessionID == "" && to == "" && len(*properties) == 0 {
		return nil, fmt.Errorf("at least one property must be set in the `correlation_filter` block")
	}

	correlationFilter := rules.CorrelationFilter{}

	if correlationID != "" {
		correlationFilter.CorrelationId = utils.String(correlationID)
	}

	if messageID != "" {
		correlationFilter.MessageId = utils.String(messageID)
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
		correlationFilter.SessionId = utils.String(sessionID)
	}

	if replyToSessionID != "" {
		correlationFilter.ReplyToSessionId = utils.String(replyToSessionID)
	}

	if contentType != "" {
		correlationFilter.ContentType = utils.String(contentType)
	}

	if len(*properties) > 0 {
		correlationFilter.Properties = properties
	}

	return &correlationFilter, nil
}

func flattenAzureRmServiceBusCorrelationFilter(input *subscriptions.CorrelationFilter) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	filter := make(map[string]interface{})

	if input.CorrelationId != nil {
		filter["correlation_id"] = *input.CorrelationId
	}

	if input.MessageId != nil {
		filter["message_id"] = *input.MessageId
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

	if input.SessionId != nil {
		filter["session_id"] = *input.SessionId
	}

	if input.ReplyToSessionId != nil {
		filter["reply_to_session_id"] = *input.ReplyToSessionId
	}

	if input.ContentType != nil {
		filter["content_type"] = *input.ContentType
	}

	if input.Properties != nil {
		filter["properties"] = flattenProperties(input.Properties)
	}

	return []interface{}{filter}
}

func expandProperties(input map[string]interface{}) *map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = v.(string)
	}
	return &output
}

func flattenProperties(input *map[string]string) map[string]*string {
	output := make(map[string]*string)

	if input != nil {
		for k, v := range *input {
			output[k] = utils.String(v)
		}
	}

	return output
}
