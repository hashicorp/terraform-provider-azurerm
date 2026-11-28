// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_subscription_rule

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusSubscriptionRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ServiceBus.SubscriptionRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	filterType := d.Get("filter_type").(string)

	var id rules.RuleId
	if subscriptionIdLit := d.Get("subscription_id").(string); subscriptionIdLit != "" {
		subscriptionId, _ := rules.ParseSubscriptions2ID(subscriptionIdLit)
		id = rules.NewRuleID(
			subscriptionId.SubscriptionId,
			subscriptionId.ResourceGroupName,
			subscriptionId.NamespaceName,
			subscriptionId.TopicName,
			subscriptionId.SubscriptionName,
			d.Get("name").(string),
		)
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_servicebus_subscription_rule", id.ID())
			}
		}
	}

	rule := rules.Rule{
		Properties: &rules.Ruleproperties{
			FilterType: pointer.ToEnum[rules.FilterType](filterType),
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

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceServiceBusSubscriptionRuleRead(d, meta)
}

func expandAzureRmServiceBusCorrelationFilter(d *pluginsdk.ResourceData) (*rules.CorrelationFilter, error) {
	configs := d.Get("correlation_filter").([]any)
	if len(configs) == 0 {
		return nil, errors.New("`correlation_filter` is required when `filter_type` is set to `CorrelationFilter`")
	}

	if configs[0] == nil {
		return nil, errors.New("at least one property must not be empty in the `correlation_filter` block")
	}
	config := configs[0].(map[string]any)

	contentType := config["content_type"].(string)
	correlationID := config["correlation_id"].(string)
	label := config["label"].(string)
	messageID := config["message_id"].(string)
	replyTo := config["reply_to"].(string)
	replyToSessionID := config["reply_to_session_id"].(string)
	sessionID := config["session_id"].(string)
	to := config["to"].(string)

	properties := expandProperties(config["properties"].(map[string]any))

	if contentType == "" && correlationID == "" && label == "" && messageID == "" && replyTo == "" && replyToSessionID == "" && sessionID == "" && to == "" && len(*properties) == 0 {
		return nil, errors.New("at least one property must be set in the `correlation_filter` block")
	}

	correlationFilter := rules.CorrelationFilter{}

	if correlationID != "" {
		correlationFilter.CorrelationId = pointer.To(correlationID)
	}

	if messageID != "" {
		correlationFilter.MessageId = pointer.To(messageID)
	}

	if to != "" {
		correlationFilter.To = pointer.To(to)
	}

	if replyTo != "" {
		correlationFilter.ReplyTo = pointer.To(replyTo)
	}

	if label != "" {
		correlationFilter.Label = pointer.To(label)
	}

	if sessionID != "" {
		correlationFilter.SessionId = pointer.To(sessionID)
	}

	if replyToSessionID != "" {
		correlationFilter.ReplyToSessionId = pointer.To(replyToSessionID)
	}

	if contentType != "" {
		correlationFilter.ContentType = pointer.To(contentType)
	}

	if len(*properties) > 0 {
		correlationFilter.Properties = properties
	}

	return &correlationFilter, nil
}

func expandProperties(input map[string]any) *map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = v.(string)
	}
	return &output
}
