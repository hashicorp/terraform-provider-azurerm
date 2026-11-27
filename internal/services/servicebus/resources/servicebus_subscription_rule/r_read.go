// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_subscription_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	servicebusClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusSubscriptionRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := servicebusClient.From(meta).SubscriptionRulesClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := rules.ParseRuleID(d.Id())
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

	d.Set("subscription_id", subscriptions.NewSubscriptions2ID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicName, id.SubscriptionName).ID())
	d.Set("name", id.RuleName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("filter_type", pointer.FromEnum(props.FilterType))

			if props.Action != nil {
				d.Set("action", props.Action.SqlExpression)
			}

			if props.SqlFilter != nil {
				d.Set("sql_filter", props.SqlFilter.SqlExpression)
				if props.SqlFilter.CompatibilityLevel != nil {
					d.Set("sql_filter_compatibility_level", props.SqlFilter.CompatibilityLevel)
				}
			}

			if err := d.Set("correlation_filter", flattenAzureRmServiceBusCorrelationFilter(props.CorrelationFilter)); err != nil {
				return fmt.Errorf("setting `correlation_filter`: %+v", err)
			}
		}
	}

	return nil
}

func flattenAzureRmServiceBusCorrelationFilter(input *rules.CorrelationFilter) []any {
	if input == nil {
		return []any{}
	}

	filter := make(map[string]any)

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

	return []any{filter}
}

func flattenProperties(input *map[string]string) map[string]*string {
	output := make(map[string]*string)

	if input != nil {
		for k, v := range *input {
			output[k] = pointer.To(v)
		}
	}

	return output
}
