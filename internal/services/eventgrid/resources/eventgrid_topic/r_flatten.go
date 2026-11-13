// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_topic

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topics"
)

func flattenTopicInputMapping(input topics.InputSchemaMapping) []any {
	val, ok := input.(topics.JsonInputSchemaMapping)
	if !ok {
		return []any{}
	}

	dataVersion := ""
	eventTime := ""
	eventType := ""
	id := ""
	subject := ""
	topic := ""
	if props := val.Properties; props != nil {
		if props.EventTime != nil && props.EventTime.SourceField != nil {
			eventTime = *props.EventTime.SourceField
		}

		if props.Id != nil && props.Id.SourceField != nil {
			id = *props.Id.SourceField
		}

		if props.Topic != nil && props.Topic.SourceField != nil {
			topic = *props.Topic.SourceField
		}

		if props.DataVersion != nil && props.DataVersion.SourceField != nil {
			dataVersion = *props.DataVersion.SourceField
		}

		if props.EventType != nil && props.EventType.SourceField != nil {
			eventType = *props.EventType.SourceField
		}

		if props.Subject != nil && props.Subject.SourceField != nil {
			subject = *props.Subject.SourceField
		}
	}

	return []any{
		map[string]any{
			"data_version": dataVersion,
			"event_time":   eventTime,
			"event_type":   eventType,
			"id":           id,
			"subject":      subject,
			"topic":        topic,
		},
	}
}

func flattenTopicInputMappingDefaultValues(input topics.InputSchemaMapping) []any {
	val, ok := input.(topics.JsonInputSchemaMapping)
	if !ok || val.Properties == nil {
		return []any{}
	}

	dataVersion := ""
	eventType := ""
	subject := ""
	if val.Properties != nil {
		if val.Properties.DataVersion != nil && val.Properties.DataVersion.DefaultValue != nil {
			dataVersion = *val.Properties.DataVersion.DefaultValue
		}
		if val.Properties.EventType != nil && val.Properties.EventType.DefaultValue != nil {
			eventType = *val.Properties.EventType.DefaultValue
		}
		if val.Properties.Subject != nil && val.Properties.Subject.DefaultValue != nil {
			subject = *val.Properties.Subject.DefaultValue
		}
	}

	return []any{
		map[string]any{
			"data_version": dataVersion,
			"event_type":   eventType,
			"subject":      subject,
		},
	}
}

func flattenTopicInboundIPRules(input *[]topics.InboundIPRule) []any {
	rules := make([]any, 0)
	if input == nil {
		return rules
	}

	for _, r := range *input {
		action := ""
		if r.Action != nil {
			action = string(*r.Action)
		}

		rules = append(rules, map[string]any{
			"action":  action,
			"ip_mask": pointer.From(r.IPMask),
		})
	}
	return rules
}
