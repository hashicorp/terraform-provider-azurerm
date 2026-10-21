// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_domain

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/domains"
)

func flattenDomainInputMapping(input domains.InputSchemaMapping) []any {
	output := make([]any, 0)
	val, ok := input.(domains.JsonInputSchemaMapping)
	if ok {
		if props := val.Properties; props != nil {
			eventTime := ""
			if props.EventTime != nil && props.EventTime.SourceField != nil {
				eventTime = *props.EventTime.SourceField
			}

			id := ""
			if props.Id != nil && props.Id.SourceField != nil {
				id = *props.Id.SourceField
			}

			topic := ""
			if props.Topic != nil && props.Topic.SourceField != nil {
				topic = *props.Topic.SourceField
			}

			dataVersion := ""
			if props.DataVersion != nil && props.DataVersion.SourceField != nil {
				dataVersion = *props.DataVersion.SourceField
			}

			eventType := ""
			if props.EventType != nil && props.EventType.SourceField != nil {
				eventType = *props.EventType.SourceField
			}

			subject := ""
			if props.Subject != nil && props.Subject.SourceField != nil {
				subject = *props.Subject.SourceField
			}

			output = append(output, map[string]any{
				"data_version": dataVersion,
				"event_type":   eventType,
				"event_time":   eventTime,
				"id":           id,
				"topic":        topic,
				"subject":      subject,
			})
		}
	}
	return output
}

func flattenDomainInputMappingDefaultValues(input domains.InputSchemaMapping) []any {
	output := make([]any, 0)
	val, ok := input.(domains.JsonInputSchemaMapping)
	if ok {
		if props := val.Properties; props != nil {
			dataVersion := ""
			if props.DataVersion != nil && props.DataVersion.DefaultValue != nil {
				dataVersion = *props.DataVersion.DefaultValue
			}

			eventType := ""
			if props.EventType != nil && props.EventType.DefaultValue != nil {
				eventType = *props.EventType.DefaultValue
			}

			subject := ""
			if props.Subject != nil && props.Subject.DefaultValue != nil {
				subject = *props.Subject.DefaultValue
			}

			output = append(output, map[string]any{
				"data_version": dataVersion,
				"event_type":   eventType,
				"subject":      subject,
			})
		}
	}

	return output
}

func flattenDomainInboundIPRules(input *[]domains.InboundIPRule) []any {
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
