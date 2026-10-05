// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/eventsubscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandEventSubscriptionDestination(d *pluginsdk.ResourceData) eventsubscriptions.EventSubscriptionDestination {
	deliveryMappings := expandEventSubscriptionDeliveryAttributeMappings(d.Get("delivery_property").([]any))

	if val, ok := d.GetOk("azure_alert_monitor"); ok && len(val.([]any)) == 1 {
		return expandEventGridEventSubscriptionAzureAlertMonitor(d.Get("azure_alert_monitor").([]any))
	}

	if val, ok := d.GetOk("azure_function_endpoint"); ok && len(val.([]any)) == 1 {
		return expandEventSubscriptionDestinationAzureFunction(d.Get("azure_function_endpoint").([]any), deliveryMappings)
	}

	eventhubEndpointId, ok := d.GetOk("eventhub_id")
	if ok {
		return expandEventSubscriptionDestinationEventHub(eventhubEndpointId.(string), deliveryMappings)
	}

	hybridConnectionEndpointId, ok := d.GetOk("hybrid_connection_id")
	if ok {
		return expandEventSubscriptionDestinationHybridConnection(hybridConnectionEndpointId.(string), deliveryMappings)
	}

	if val, ok := d.GetOk("service_bus_queue_id"); ok {
		return expandEventSubscriptionDestinationServiceBusQueueEndpoint(val.(string), deliveryMappings)
	}

	if val, ok := d.GetOk("service_bus_topic_id"); ok {
		return expandEventSubscriptionDestinationServiceBusTopicEndpoint(val.(string), deliveryMappings)
	}

	if val, ok := d.GetOk("storage_queue_endpoint"); ok {
		return expandEventSubscriptionStorageQueueEndpoint(val.([]any))
	}

	if val, ok := d.GetOk("webhook_endpoint"); ok {
		return expandEventGridEventSubscriptionWebhookEndpoint(val.([]any), deliveryMappings)
	}

	return nil
}

func expandEventGridEventSubscriptionWebhookEndpoint(input []any, deliveryMappings []eventsubscriptions.DeliveryAttributeMapping) eventsubscriptions.EventSubscriptionDestination {
	props := eventsubscriptions.WebHookEventSubscriptionDestinationProperties{
		DeliveryAttributeMappings: &deliveryMappings,
	}
	webhookDestination := &eventsubscriptions.WebHookEventSubscriptionDestination{
		Properties: &props,
	}

	if len(input) == 0 {
		return webhookDestination
	}

	config := input[0].(map[string]any)

	if v, ok := config["url"]; ok && v != "" {
		props.EndpointURL = pointer.To(v.(string))
	}

	if v, ok := config["max_events_per_batch"]; ok && v != 0 {
		props.MaxEventsPerBatch = pointer.To(int64(v.(int)))
	}

	if v, ok := config["preferred_batch_size_in_kilobytes"]; ok && v != 0 {
		props.PreferredBatchSizeInKilobytes = pointer.To(int64(v.(int)))
	}

	if v, ok := config["active_directory_tenant_id"]; ok && v != "" {
		props.AzureActiveDirectoryTenantId = pointer.To(v.(string))
	}

	if v, ok := config["active_directory_app_id_or_uri"]; ok && v != "" {
		props.AzureActiveDirectoryApplicationIdOrUri = pointer.To(v.(string))
	}

	return webhookDestination
}

func expandEventGridEventSubscriptionAzureAlertMonitor(input []any) eventsubscriptions.MonitorAlertEventSubscriptionDestination {
	item := input[0].(map[string]any)
	props := eventsubscriptions.MonitorAlertEventSubscriptionDestinationProperties{
		Description: pointer.To(item["description"].(string)),
		Severity:    pointer.To(eventsubscriptions.MonitorAlertSeverity(item["severity"].(string))),
	}

	if v, ok := item["action_groups"]; ok && v != nil {
		props.ActionGroups = pluginsdk.ExpandStringSlice(v.([]any))
	}

	return eventsubscriptions.MonitorAlertEventSubscriptionDestination{
		Properties: &props,
	}
}

func flattenEventSubscriptionDestinationAzureAlertMonitor(input any) []any {
	output := make([]any, 0)
	val, ok := input.(eventsubscriptions.MonitorAlertEventSubscriptionDestination)

	if ok && val.Properties != nil {
		output = append(output, map[string]any{
			"severity":      pointer.From(val.Properties.Severity),
			"description":   pointer.From(val.Properties.Description),
			"action_groups": val.Properties.ActionGroups,
		})
	}

	return output
}

func expandEventSubscriptionDestinationAzureFunction(input []any, deliveryMappings []eventsubscriptions.DeliveryAttributeMapping) eventsubscriptions.EventSubscriptionDestination {
	item := input[0].(map[string]any)
	props := eventsubscriptions.AzureFunctionEventSubscriptionDestinationProperties{
		DeliveryAttributeMappings: &deliveryMappings,
	}
	if v, ok := item["function_id"]; ok && v != "" {
		props.ResourceId = pointer.To(v.(string))
	}
	if v, ok := item["max_events_per_batch"]; ok && v != 0 {
		props.MaxEventsPerBatch = pointer.To(int64(v.(int)))
	}
	if v, ok := item["preferred_batch_size_in_kilobytes"]; ok && v != 0 {
		props.PreferredBatchSizeInKilobytes = pointer.To(int64(v.(int)))
	}

	return eventsubscriptions.AzureFunctionEventSubscriptionDestination{
		Properties: &props,
	}
}

func flattenEventSubscriptionDestinationAzureFunction(input eventsubscriptions.EventSubscriptionDestination) []any {
	output := make([]any, 0)

	val, ok := input.(eventsubscriptions.AzureFunctionEventSubscriptionDestination)
	if ok && val.Properties != nil {
		props := *val.Properties
		return append(output, map[string]any{
			"function_id":                       pointer.From(props.ResourceId),
			"max_events_per_batch":              int(pointer.From(props.MaxEventsPerBatch)),
			"preferred_batch_size_in_kilobytes": int(pointer.From(props.PreferredBatchSizeInKilobytes)),
		})
	}

	return output
}

func expandEventSubscriptionDestinationEventHub(eventhubEndpointId string, deliveryMappings []eventsubscriptions.DeliveryAttributeMapping) eventsubscriptions.EventSubscriptionDestination {
	return eventsubscriptions.EventHubEventSubscriptionDestination{
		Properties: &eventsubscriptions.EventHubEventSubscriptionDestinationProperties{
			DeliveryAttributeMappings: pointer.To(deliveryMappings),
			ResourceId:                pointer.To(eventhubEndpointId),
		},
	}
}

func flattenEventSubscriptionDestinationEventHub(input eventsubscriptions.EventSubscriptionDestination) string {
	if val, ok := input.(eventsubscriptions.EventHubEventSubscriptionDestination); ok && val.Properties != nil && val.Properties.ResourceId != nil {
		return *val.Properties.ResourceId
	}

	return ""
}

func expandEventSubscriptionDestinationHybridConnection(hybridConnectionId string, deliveryMappings []eventsubscriptions.DeliveryAttributeMapping) eventsubscriptions.EventSubscriptionDestination {
	return eventsubscriptions.HybridConnectionEventSubscriptionDestination{
		Properties: &eventsubscriptions.HybridConnectionEventSubscriptionDestinationProperties{
			DeliveryAttributeMappings: pointer.To(deliveryMappings),
			ResourceId:                pointer.To(hybridConnectionId),
		},
	}
}

func flattenEventSubscriptionDestinationHybridConnection(input eventsubscriptions.EventSubscriptionDestination) string {
	if val, ok := input.(eventsubscriptions.HybridConnectionEventSubscriptionDestination); ok && val.Properties != nil && val.Properties.ResourceId != nil {
		return *val.Properties.ResourceId
	}

	return ""
}

func expandEventSubscriptionDestinationServiceBusQueueEndpoint(serviceBusQueueEndpointId string, deliveryMappings []eventsubscriptions.DeliveryAttributeMapping) eventsubscriptions.EventSubscriptionDestination {
	return eventsubscriptions.ServiceBusQueueEventSubscriptionDestination{
		Properties: &eventsubscriptions.ServiceBusQueueEventSubscriptionDestinationProperties{
			DeliveryAttributeMappings: pointer.To(deliveryMappings),
			ResourceId:                pointer.To(serviceBusQueueEndpointId),
		},
	}
}

func flattenEventSubscriptionDestinationServiceBusQueueEndpoint(input eventsubscriptions.EventSubscriptionDestination) string {
	if val, ok := input.(eventsubscriptions.ServiceBusQueueEventSubscriptionDestination); ok && val.Properties != nil && val.Properties.ResourceId != nil {
		return *val.Properties.ResourceId
	}

	return ""
}

func expandEventSubscriptionDestinationServiceBusTopicEndpoint(serviceBusTopicEndpointId string, deliveryMappings []eventsubscriptions.DeliveryAttributeMapping) eventsubscriptions.EventSubscriptionDestination {
	return eventsubscriptions.ServiceBusTopicEventSubscriptionDestination{
		Properties: &eventsubscriptions.ServiceBusTopicEventSubscriptionDestinationProperties{
			DeliveryAttributeMappings: pointer.To(deliveryMappings),
			ResourceId:                pointer.To(serviceBusTopicEndpointId),
		},
	}
}

func flattenEventSubscriptionDestinationServiceBusTopicEndpoint(input eventsubscriptions.EventSubscriptionDestination) string {
	if val, ok := input.(eventsubscriptions.ServiceBusTopicEventSubscriptionDestination); ok && val.Properties != nil && val.Properties.ResourceId != nil {
		return *val.Properties.ResourceId
	}

	return ""
}

func expandEventSubscriptionStorageQueueEndpoint(input []any) eventsubscriptions.EventSubscriptionDestination {
	raw := input[0].(map[string]any)
	props := eventsubscriptions.StorageQueueEventSubscriptionDestinationProperties{
		ResourceId: pointer.To(raw["storage_account_id"].(string)),
		QueueName:  pointer.To(raw["queue_name"].(string)),
	}

	if ttlInSeconds := raw["queue_message_time_to_live_in_seconds"]; ttlInSeconds != 0 {
		props.QueueMessageTimeToLiveInSeconds = pointer.To(int64(ttlInSeconds.(int)))
	}

	return eventsubscriptions.StorageQueueEventSubscriptionDestination{
		Properties: &props,
	}
}

func flattenEventSubscriptionDestinationStorageQueueEndpoint(input eventsubscriptions.EventSubscriptionDestination) []any {
	output := make([]any, 0)

	val, ok := input.(eventsubscriptions.StorageQueueEventSubscriptionDestination)
	if ok && val.Properties != nil {
		output = append(output, map[string]any{
			"queue_message_time_to_live_in_seconds": int(pointer.From(val.Properties.QueueMessageTimeToLiveInSeconds)),
			"storage_account_id":                    pointer.From(val.Properties.ResourceId),
			"queue_name":                            pointer.From(val.Properties.QueueName),
		})
	}

	return output
}

func expandEventSubscriptionDeliveryAttributeMappings(input []any) []eventsubscriptions.DeliveryAttributeMapping {
	output := make([]eventsubscriptions.DeliveryAttributeMapping, 0)
	for _, item := range input {
		mappingBlock := item.(map[string]any)

		switch mappingBlock["type"].(string) {
		case "Static":
			output = append(output, eventsubscriptions.StaticDeliveryAttributeMapping{
				Name: pointer.To(mappingBlock["header_name"].(string)),
				Properties: &eventsubscriptions.StaticDeliveryAttributeMappingProperties{
					Value:    pointer.To(mappingBlock["value"].(string)),
					IsSecret: pointer.To(mappingBlock["secret"].(bool)),
				},
			})
		case "Dynamic":
			output = append(output, eventsubscriptions.DynamicDeliveryAttributeMapping{
				Name: pointer.To(mappingBlock["header_name"].(string)),
				Properties: &eventsubscriptions.DynamicDeliveryAttributeMappingProperties{
					SourceField: pointer.To(mappingBlock["source_field"].(string)),
				},
			})
		}
	}

	return output
}

func flattenEventSubscriptionDeliveryAttributeMappings(input eventsubscriptions.EventSubscriptionDestination, mappingsFromState []eventsubscriptions.DeliveryAttributeMapping) []any {
	mappings := make([]eventsubscriptions.DeliveryAttributeMapping, 0)

	if v, ok := input.(eventsubscriptions.AzureFunctionEventSubscriptionDestination); ok && v.Properties != nil && v.Properties.DeliveryAttributeMappings != nil {
		mappings = *v.Properties.DeliveryAttributeMappings
	}
	if v, ok := input.(eventsubscriptions.EventHubEventSubscriptionDestination); ok && v.Properties != nil && v.Properties.DeliveryAttributeMappings != nil {
		mappings = *v.Properties.DeliveryAttributeMappings
	}
	if v, ok := input.(eventsubscriptions.HybridConnectionEventSubscriptionDestination); ok && v.Properties != nil && v.Properties.DeliveryAttributeMappings != nil {
		mappings = *v.Properties.DeliveryAttributeMappings
	}
	// NOTE: `MonitorAlertEventSubscriptionDestination` doesn't contain DeliveryAttributeMappings
	if v, ok := input.(eventsubscriptions.ServiceBusQueueEventSubscriptionDestination); ok && v.Properties != nil && v.Properties.DeliveryAttributeMappings != nil {
		mappings = *v.Properties.DeliveryAttributeMappings
	}
	if v, ok := input.(eventsubscriptions.ServiceBusTopicEventSubscriptionDestination); ok && v.Properties != nil && v.Properties.DeliveryAttributeMappings != nil {
		mappings = *v.Properties.DeliveryAttributeMappings
	}
	// NOTE: `StorageQueueEventSubscriptionDestination` doesn't contain DeliveryAttributeMappings
	if v, ok := input.(eventsubscriptions.WebHookEventSubscriptionDestination); ok && v.Properties != nil && v.Properties.DeliveryAttributeMappings != nil {
		mappings = *v.Properties.DeliveryAttributeMappings
	}

	output := make([]any, 0)
	for _, mapping := range mappings {
		if val, ok := mapping.(eventsubscriptions.StaticDeliveryAttributeMapping); ok {
			secret := false
			value := ""
			if val.Properties != nil {
				if val.Properties.IsSecret != nil {
					secret = *val.Properties.IsSecret
				}
				if val.Properties.Value != nil {
					value = *val.Properties.Value
				}
				if secret {
					// If this is a secret, the Azure API just returns a value of 'Hidden',
					// so we need to lookup the value that was provided from config to return
					for _, v := range mappingsFromState {
						mapping, ok := v.(eventsubscriptions.StaticDeliveryAttributeMapping)
						if ok && mapping.Name != nil && val.Name != nil && *mapping.Name == *val.Name && mapping.Properties != nil && mapping.Properties.Value != nil {
							value = *mapping.Properties.Value
							break
						}
					}
				}
			}
			output = append(output, map[string]any{
				"header_name": pointer.From(val.Name),
				"secret":      secret,
				"type":        "Static",
				"value":       value,
			})
		}

		if val, ok := mapping.(eventsubscriptions.DynamicDeliveryAttributeMapping); ok {
			sourceField := ""
			if val.Properties != nil && val.Properties.SourceField != nil {
				sourceField = *val.Properties.SourceField
			}
			output = append(output, map[string]any{
				"header_name":  pointer.From(val.Name),
				"source_field": sourceField,
				"type":         "Dynamic",
			})
		}
	}

	return output
}

func expandEventSubscriptionIdentity(input []any) (*eventsubscriptions.EventSubscriptionIdentity, error) {
	if len(input) == 0 || input[0] == nil {
		return &eventsubscriptions.EventSubscriptionIdentity{
			Type: pointer.ToEnum[eventsubscriptions.EventSubscriptionIdentityType]("None"),
		}, nil
	}

	identity := input[0].(map[string]any)
	identityType := eventsubscriptions.EventSubscriptionIdentityType(identity["type"].(string))
	eventgridIdentity := eventsubscriptions.EventSubscriptionIdentity{
		Type: pointer.To(identityType),
	}

	userAssignedIdentity := identity["user_assigned_identity"].(string)
	if identityType == eventsubscriptions.EventSubscriptionIdentityTypeUserAssigned {
		eventgridIdentity.UserAssignedIdentity = pointer.To(userAssignedIdentity)
	} else if len(userAssignedIdentity) > 0 {
		return nil, fmt.Errorf("`user_assigned_identity` can only be specified when `type` is `UserAssigned`; but `type` is currently %q", identityType)
	}

	return &eventgridIdentity, nil
}

func flattenEventSubscriptionWebhookEndpoint(input eventsubscriptions.EventSubscriptionDestination, fullUrl *eventsubscriptions.EventSubscriptionFullURL) []any {
	output := make([]any, 0)
	val, ok := input.(eventsubscriptions.WebHookEventSubscriptionDestination)
	if ok {
		webHookUrl := ""
		if fullUrl != nil {
			webHookUrl = *fullUrl.EndpointURL
		}

		azureActiveDirectoryApplicationIdOrUrl := ""
		azureActiveDirectoryTenantId := ""
		maxEventsPerBatch := 0
		preferredBatchSizeInKilobytes := 0
		webhookBaseURL := ""
		if props := val.Properties; props != nil {
			if props.EndpointBaseURL != nil {
				webhookBaseURL = *props.EndpointBaseURL
			}

			if props.MaxEventsPerBatch != nil {
				maxEventsPerBatch = int(*props.MaxEventsPerBatch)
			}

			if props.PreferredBatchSizeInKilobytes != nil {
				preferredBatchSizeInKilobytes = int(*props.PreferredBatchSizeInKilobytes)
			}

			if props.AzureActiveDirectoryTenantId != nil {
				azureActiveDirectoryTenantId = *props.AzureActiveDirectoryTenantId
			}

			if props.AzureActiveDirectoryApplicationIdOrUri != nil {
				azureActiveDirectoryApplicationIdOrUrl = *props.AzureActiveDirectoryApplicationIdOrUri
			}
		}

		output = append(output, map[string]any{
			"url":                               webHookUrl,
			"base_url":                          webhookBaseURL,
			"max_events_per_batch":              maxEventsPerBatch,
			"preferred_batch_size_in_kilobytes": preferredBatchSizeInKilobytes,
			"active_directory_tenant_id":        azureActiveDirectoryTenantId,
			"active_directory_app_id_or_uri":    azureActiveDirectoryApplicationIdOrUrl,
		})
	}

	return output
}

func flattenEventSubscriptionIdentity(input *eventsubscriptions.EventSubscriptionIdentity) []any {
	if input == nil || input.Type == nil || strings.EqualFold(string(*input.Type), "None") {
		return []any{}
	}

	return []any{
		map[string]any{
			"type":                   string(*input.Type),
			"user_assigned_identity": pointer.From(input.UserAssignedIdentity),
		},
	}
}

func flattenEventSubscriptionRetryPolicy(input *eventsubscriptions.RetryPolicy) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"event_time_to_live":    int(pointer.From(input.EventTimeToLiveInMinutes)),
			"max_delivery_attempts": int(pointer.From(input.MaxDeliveryAttempts)),
		},
	}
}

func flattenEventSubscriptionStorageBlobDeadLetterDestination(input eventsubscriptions.DeadLetterDestination) []any {
	if input == nil {
		return []any{}
	}
	val, ok := input.(eventsubscriptions.StorageBlobDeadLetterDestination)
	if !ok || val.Properties == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"storage_account_id":          pointer.From(val.Properties.ResourceId),
			"storage_blob_container_name": pointer.From(val.Properties.BlobContainerName),
		},
	}
}

func flattenEventSubscriptionSubjectFilter(input *eventsubscriptions.EventSubscriptionFilter) []any {
	output := make([]any, 0)
	if input == nil {
		return output
	}
	if (input.SubjectBeginsWith != nil && *input.SubjectBeginsWith == "") && (input.SubjectEndsWith != nil && *input.SubjectEndsWith == "") {
		return output
	}

	output = append(output, map[string]any{
		"subject_begins_with": pointer.From(input.SubjectBeginsWith),
		"subject_ends_with":   pointer.From(input.SubjectEndsWith),
		"case_sensitive":      pointer.From(input.IsSubjectCaseSensitive),
	})

	return output
}

func flattenEventSubscriptionAdvancedFilter(input *eventsubscriptions.EventSubscriptionFilter) []any {
	output := make([]any, 0)
	if input == nil || input.AdvancedFilters == nil {
		return output
	}

	boolEquals := make([]any, 0)
	numberGreaterThan := make([]any, 0)
	numberGreaterThanOrEquals := make([]any, 0)
	numberLessThan := make([]any, 0)
	numberLessThanOrEquals := make([]any, 0)
	numberIn := make([]any, 0)
	numberNotIn := make([]any, 0)
	numberInRange := make([]any, 0)
	numberNotInRange := make([]any, 0)
	stringBeginsWith := make([]any, 0)
	stringNotBeginsWith := make([]any, 0)
	stringEndsWith := make([]any, 0)
	stringNotEndsWith := make([]any, 0)
	stringContains := make([]any, 0)
	stringNotContains := make([]any, 0)
	stringIn := make([]any, 0)
	stringNotIn := make([]any, 0)
	isNotNull := make([]any, 0)
	isNullOrUndefined := make([]any, 0)

	for _, item := range *input.AdvancedFilters {
		switch f := item.(type) {
		case eventsubscriptions.BoolEqualsAdvancedFilter:
			boolEquals = append(boolEquals, flattenValue(f.Key, pointer.To(any(f.Value))))
		case eventsubscriptions.NumberGreaterThanAdvancedFilter:
			numberGreaterThan = append(numberGreaterThan, flattenValue(f.Key, pointer.To(any(f.Value))))
		case eventsubscriptions.NumberGreaterThanOrEqualsAdvancedFilter:
			numberGreaterThanOrEquals = append(numberGreaterThanOrEquals, flattenValue(f.Key, pointer.To(any(f.Value))))
		case eventsubscriptions.NumberLessThanAdvancedFilter:
			numberLessThan = append(numberLessThan, flattenValue(f.Key, pointer.To(any(f.Value))))
		case eventsubscriptions.NumberLessThanOrEqualsAdvancedFilter:
			numberLessThanOrEquals = append(numberLessThanOrEquals, flattenValue(f.Key, pointer.To(any(f.Value))))
		case eventsubscriptions.NumberInAdvancedFilter:
			v := pluginsdk.FlattenSlice(f.Values)
			numberIn = append(numberIn, flattenValues(f.Key, &v))
		case eventsubscriptions.NumberNotInAdvancedFilter:
			v := pluginsdk.FlattenSlice(f.Values)
			numberNotIn = append(numberNotIn, flattenValues(f.Key, &v))
		case eventsubscriptions.StringBeginsWithAdvancedFilter:
			v := pluginsdk.FlattenSlice(f.Values)
			stringBeginsWith = append(stringBeginsWith, flattenValues(f.Key, &v))
		case eventsubscriptions.StringNotBeginsWithAdvancedFilter:
			v := pluginsdk.FlattenSlice(f.Values)
			stringNotBeginsWith = append(stringNotBeginsWith, flattenValues(f.Key, &v))
		case eventsubscriptions.StringEndsWithAdvancedFilter:
			v := pluginsdk.FlattenSlice(f.Values)
			stringEndsWith = append(stringEndsWith, flattenValues(f.Key, &v))
		case eventsubscriptions.StringNotEndsWithAdvancedFilter:
			v := pluginsdk.FlattenSlice(f.Values)
			stringNotEndsWith = append(stringNotEndsWith, flattenValues(f.Key, &v))
		case eventsubscriptions.StringContainsAdvancedFilter:
			v := pluginsdk.FlattenSlice(f.Values)
			stringContains = append(stringContains, flattenValues(f.Key, &v))
		case eventsubscriptions.StringNotContainsAdvancedFilter:
			v := pluginsdk.FlattenSlice(f.Values)
			stringNotContains = append(stringNotContains, flattenValues(f.Key, &v))
		case eventsubscriptions.StringInAdvancedFilter:
			v := pluginsdk.FlattenSlice(f.Values)
			stringIn = append(stringIn, flattenValues(f.Key, &v))
		case eventsubscriptions.StringNotInAdvancedFilter:
			v := pluginsdk.FlattenSlice(f.Values)
			stringNotIn = append(stringNotIn, flattenValues(f.Key, &v))
		case eventsubscriptions.NumberInRangeAdvancedFilter:
			v := pluginsdk.FlattenFloatRangeSlice(f.Values)
			numberInRange = append(numberInRange, flattenRangeValues(f.Key, &v))
		case eventsubscriptions.NumberNotInRangeAdvancedFilter:
			v := pluginsdk.FlattenFloatRangeSlice(f.Values)
			numberNotInRange = append(numberNotInRange, flattenRangeValues(f.Key, &v))
		case eventsubscriptions.IsNotNullAdvancedFilter:
			isNotNull = append(isNotNull, flattenKey(f.Key))
		case eventsubscriptions.IsNullOrUndefinedAdvancedFilter:
			isNullOrUndefined = append(isNullOrUndefined, flattenKey(f.Key))
		}
	}

	return []any{
		map[string][]any{
			"bool_equals":                   boolEquals,
			"number_greater_than":           numberGreaterThan,
			"number_greater_than_or_equals": numberGreaterThanOrEquals,
			"number_less_than":              numberLessThan,
			"number_less_than_or_equals":    numberLessThanOrEquals,
			"number_in":                     numberIn,
			"number_not_in":                 numberNotIn,
			"number_in_range":               numberInRange,
			"number_not_in_range":           numberNotInRange,
			"string_begins_with":            stringBeginsWith,
			"string_not_begins_with":        stringNotBeginsWith,
			"string_ends_with":              stringEndsWith,
			"string_not_ends_with":          stringNotEndsWith,
			"string_contains":               stringContains,
			"string_not_contains":           stringNotContains,
			"string_in":                     stringIn,
			"string_not_in":                 stringNotIn,
			"is_not_null":                   isNotNull,
			"is_null_or_undefined":          isNullOrUndefined,
		},
	}
}

func expandEventSubscriptionRetryPolicy(d *pluginsdk.ResourceData) *eventsubscriptions.RetryPolicy {
	if v, ok := d.GetOk("retry_policy"); ok {
		dest := v.([]any)[0].(map[string]any)
		maxDeliveryAttempts := dest["max_delivery_attempts"].(int)
		eventTimeToLive := dest["event_time_to_live"].(int)
		return &eventsubscriptions.RetryPolicy{
			MaxDeliveryAttempts:      pointer.To(int64(maxDeliveryAttempts)),
			EventTimeToLiveInMinutes: pointer.To(int64(eventTimeToLive)),
		}
	}

	return nil
}

func expandEventSubscriptionFilter(d *pluginsdk.ResourceData) (*eventsubscriptions.EventSubscriptionFilter, error) {
	filter := &eventsubscriptions.EventSubscriptionFilter{}

	if includedEvents, ok := d.GetOk("included_event_types"); ok {
		filter.IncludedEventTypes = pluginsdk.ExpandStringSlice(includedEvents.([]any))
	}

	if v, ok := d.GetOk("subject_filter"); ok {
		if v.([]any)[0] != nil {
			config := v.([]any)[0].(map[string]any)

			filter.SubjectBeginsWith = pointer.To(config["subject_begins_with"].(string))
			filter.SubjectEndsWith = pointer.To(config["subject_ends_with"].(string))
			filter.IsSubjectCaseSensitive = pointer.To(config["case_sensitive"].(bool))
		}
	}

	if advancedFilter, ok := d.GetOk("advanced_filter"); ok {
		advancedFilters := make([]eventsubscriptions.AdvancedFilter, 0)
		for filterKey, filterSchema := range advancedFilter.([]any)[0].(map[string]any) {
			for _, options := range filterSchema.([]any) {
				if filter, err := expandEventSubscriptionAdvancedFilter(filterKey, options.(map[string]any)); err == nil {
					advancedFilters = append(advancedFilters, filter)
				} else {
					return nil, err
				}
			}
		}
		filter.AdvancedFilters = &advancedFilters
	}

	if v, ok := d.GetOk("advanced_filtering_on_arrays_enabled"); ok {
		filter.EnableAdvancedFilteringOnArrays = pointer.To(v.(bool))
	}

	return filter, nil
}

func expandEventSubscriptionAdvancedFilter(operatorType string, config map[string]any) (eventsubscriptions.AdvancedFilter, error) {
	k := config["key"].(string)

	switch operatorType {
	case "bool_equals":
		return eventsubscriptions.BoolEqualsAdvancedFilter{
			Key:   &k,
			Value: pointer.To(config["value"].(bool)),
		}, nil
	case "number_greater_than":
		return eventsubscriptions.NumberGreaterThanAdvancedFilter{
			Key:   &k,
			Value: pointer.To(config["value"].(float64)),
		}, nil
	case "number_greater_than_or_equals":
		return eventsubscriptions.NumberGreaterThanOrEqualsAdvancedFilter{
			Key:   &k,
			Value: pointer.To(config["value"].(float64)),
		}, nil
	case "number_less_than":
		return eventsubscriptions.NumberLessThanAdvancedFilter{
			Key:   &k,
			Value: pointer.To(config["value"].(float64)),
		}, nil
	case "number_less_than_or_equals":
		return eventsubscriptions.NumberLessThanOrEqualsAdvancedFilter{
			Key:   &k,
			Value: pointer.To(config["value"].(float64)),
		}, nil
	case "number_in":
		v := pluginsdk.ExpandFloatSlice(config["values"].([]any))
		return eventsubscriptions.NumberInAdvancedFilter{
			Key:    &k,
			Values: v,
		}, nil
	case "number_not_in":
		v := pluginsdk.ExpandFloatSlice(config["values"].([]any))
		return eventsubscriptions.NumberNotInAdvancedFilter{
			Key:    &k,
			Values: v,
		}, nil
	case "string_begins_with":
		v := pluginsdk.ExpandStringSlice(config["values"].([]any))
		return eventsubscriptions.StringBeginsWithAdvancedFilter{
			Key:    &k,
			Values: v,
		}, nil
	case "string_not_begins_with":
		v := pluginsdk.ExpandStringSlice(config["values"].([]any))
		return eventsubscriptions.StringNotBeginsWithAdvancedFilter{
			Key:    &k,
			Values: v,
		}, nil
	case "string_ends_with":
		v := pluginsdk.ExpandStringSlice(config["values"].([]any))
		return eventsubscriptions.StringEndsWithAdvancedFilter{
			Key:    &k,
			Values: v,
		}, nil
	case "string_not_ends_with":
		v := pluginsdk.ExpandStringSlice(config["values"].([]any))
		return eventsubscriptions.StringNotEndsWithAdvancedFilter{
			Key:    &k,
			Values: v,
		}, nil
	case "string_contains":
		v := pluginsdk.ExpandStringSlice(config["values"].([]any))
		return eventsubscriptions.StringContainsAdvancedFilter{
			Key:    &k,
			Values: v,
		}, nil
	case "string_not_contains":
		v := pluginsdk.ExpandStringSlice(config["values"].([]any))
		return eventsubscriptions.StringNotContainsAdvancedFilter{
			Key:    &k,
			Values: v,
		}, nil
	case "string_in":
		v := pluginsdk.ExpandStringSlice(config["values"].([]any))
		return eventsubscriptions.StringInAdvancedFilter{
			Key:    &k,
			Values: v,
		}, nil
	case "string_not_in":
		v := pluginsdk.ExpandStringSlice(config["values"].([]any))
		return eventsubscriptions.StringNotInAdvancedFilter{
			Key:    &k,
			Values: v,
		}, nil
	case "is_not_null":
		return eventsubscriptions.IsNotNullAdvancedFilter{
			Key: &k,
		}, nil
	case "is_null_or_undefined":
		return eventsubscriptions.IsNullOrUndefinedAdvancedFilter{
			Key: &k,
		}, nil
	case "number_in_range":
		v := pluginsdk.ExpandFloatRangeSlice(config["values"].([]any))
		return eventsubscriptions.NumberInRangeAdvancedFilter{
			Key:    &k,
			Values: v,
		}, nil
	case "number_not_in_range":
		v := pluginsdk.ExpandFloatRangeSlice(config["values"].([]any))
		return eventsubscriptions.NumberNotInRangeAdvancedFilter{
			Key:    &k,
			Values: v,
		}, nil
	default:
		return nil, fmt.Errorf("invalid `advanced_filter` operator_type %q used", operatorType)
	}
}

func expandEventSubscriptionStorageBlobDeadLetterDestination(d *pluginsdk.ResourceData) eventsubscriptions.DeadLetterDestination {
	if v, ok := d.GetOk("storage_blob_dead_letter_destination"); ok {
		dest := v.([]any)[0].(map[string]any)
		return eventsubscriptions.StorageBlobDeadLetterDestination{
			Properties: &eventsubscriptions.StorageBlobDeadLetterDestinationProperties{
				ResourceId:        pointer.To(dest["storage_account_id"].(string)),
				BlobContainerName: pointer.To(dest["storage_blob_container_name"].(string)),
			},
		}
	}

	return nil
}

func flattenValue(inputKey *string, inputValue *any) map[string]any {
	var value any
	if inputValue != nil {
		value = inputValue
	}

	return map[string]any{
		"key":   pointer.From(inputKey),
		"value": value,
	}
}

func flattenValues(inputKey *string, inputValues *[]any) map[string]any {
	values := make([]any, 0)
	if inputValues != nil {
		values = *inputValues
	}

	return map[string]any{
		"key":    pointer.From(inputKey),
		"values": values,
	}
}

func flattenRangeValues(inputKey *string, inputValues *[][]any) map[string]any {
	values := make([]any, 0)
	if inputValues != nil {
		for _, item := range *inputValues {
			values = append(values, item)
		}
	}

	return map[string]any{
		"key":    pointer.From(inputKey),
		"values": values,
	}
}

func flattenKey(inputKey *string) map[string]any {
	return map[string]any{
		"key": pointer.From(inputKey),
	}
}
