// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func flattenIoTHubSku(input *devices.IotHubSkuInfo) []any {
	output := make(map[string]any)

	output["name"] = string(input.Name)
	if capacity := input.Capacity; capacity != nil {
		output["capacity"] = int(*capacity)
	}

	return []any{output}
}

func flattenIoTHubSharedAccessPolicy(input *[]devices.SharedAccessSignatureAuthorizationRule) []any {
	results := make([]any, 0)

	if keys := input; keys != nil {
		for _, key := range *keys {
			keyMap := make(map[string]any)

			if keyName := key.KeyName; keyName != nil {
				keyMap["key_name"] = *keyName
			}

			if primaryKey := key.PrimaryKey; primaryKey != nil {
				keyMap["primary_key"] = *primaryKey
			}

			if secondaryKey := key.SecondaryKey; secondaryKey != nil {
				keyMap["secondary_key"] = *secondaryKey
			}

			keyMap["permissions"] = string(key.Rights)
			results = append(results, keyMap)
		}
	}

	return results
}

func flattenIoTHubFileUpload(storageEndpoints map[string]*devices.StorageEndpointProperties, messagingEndpoints map[string]*devices.MessagingEndpointProperties, enableFileUploadNotifications *bool) []any {
	results := make([]any, 0)
	output := make(map[string]any)

	if storageEndpointProperties, ok := storageEndpoints["$default"]; ok {
		if connString := storageEndpointProperties.ConnectionString; connString != nil {
			output["connection_string"] = *connString
		}
		if containerName := storageEndpointProperties.ContainerName; containerName != nil {
			output["container_name"] = *containerName
		}
		if sasTTLAsIso8601 := storageEndpointProperties.SasTTLAsIso8601; sasTTLAsIso8601 != nil {
			output["sas_ttl"] = *sasTTLAsIso8601
		}

		authenticationType := string(devices.AuthenticationTypeKeyBased)
		if v := string(storageEndpointProperties.AuthenticationType); v != "" {
			authenticationType = v
		}
		output["authentication_type"] = authenticationType

		identityId := ""
		if storageEndpointProperties.Identity != nil && storageEndpointProperties.Identity.UserAssignedIdentity != nil {
			identityId = *storageEndpointProperties.Identity.UserAssignedIdentity
		}
		output["identity_id"] = identityId

		if messagingEndpointProperties, ok := messagingEndpoints["fileNotifications"]; ok {
			if lockDurationAsIso8601 := messagingEndpointProperties.LockDurationAsIso8601; lockDurationAsIso8601 != nil {
				output["lock_duration"] = *lockDurationAsIso8601
			}
			if ttlAsIso8601 := messagingEndpointProperties.TTLAsIso8601; ttlAsIso8601 != nil {
				output["default_ttl"] = *ttlAsIso8601
			}
			if maxDeliveryCount := messagingEndpointProperties.MaxDeliveryCount; maxDeliveryCount != nil {
				output["max_delivery_count"] = *maxDeliveryCount
			}
		}

		if enableFileUploadNotifications != nil {
			output["notifications"] = *enableFileUploadNotifications
		}

		results = append(results, output)
	}

	return results
}

func flattenIoTHubEndpoint(input *devices.RoutingProperties) []any {
	results := make([]any, 0)

	if input != nil && input.Endpoints != nil {
		if containers := input.Endpoints.StorageContainers; containers != nil {
			for _, container := range *containers {
				output := make(map[string]any)

				authenticationType := string(devices.AuthenticationTypeKeyBased)
				if string(container.AuthenticationType) != "" {
					authenticationType = string(container.AuthenticationType)
				}
				output["authentication_type"] = authenticationType

				output["connection_string"] = pointer.From(container.ConnectionString)

				output["endpoint_uri"] = pointer.From(container.EndpointURI)

				identityId := ""
				if container.Identity != nil && container.Identity.UserAssignedIdentity != nil {
					identityId = *container.Identity.UserAssignedIdentity
				}
				output["identity_id"] = identityId

				if name := container.Name; name != nil {
					output["name"] = *name
				}
				if containerName := container.ContainerName; containerName != nil {
					output["container_name"] = *containerName
				}
				if fileNameFmt := container.FileNameFormat; fileNameFmt != nil {
					output["file_name_format"] = *fileNameFmt
				}
				if batchFreq := container.BatchFrequencyInSeconds; batchFreq != nil {
					output["batch_frequency_in_seconds"] = *batchFreq
				}
				if chunkSize := container.MaxChunkSizeInBytes; chunkSize != nil {
					output["max_chunk_size_in_bytes"] = *chunkSize
				}
				if resourceGroup := container.ResourceGroup; resourceGroup != nil {
					output["resource_group_name"] = *resourceGroup
				}

				output["encoding"] = string(container.Encoding)
				output["type"] = "AzureIotHub.StorageContainer"
				output["subscription_id"] = pointer.From(container.SubscriptionID)

				results = append(results, output)
			}
		}

		if queues := input.Endpoints.ServiceBusQueues; queues != nil {
			for _, queue := range *queues {
				output := make(map[string]any)

				authenticationType := string(devices.AuthenticationTypeKeyBased)
				if string(queue.AuthenticationType) != "" {
					authenticationType = string(queue.AuthenticationType)
				}
				output["authentication_type"] = authenticationType

				output["connection_string"] = pointer.From(queue.ConnectionString)

				output["endpoint_uri"] = pointer.From(queue.EndpointURI)

				output["entity_path"] = pointer.From(queue.EntityPath)

				identityId := ""
				if queue.Identity != nil && queue.Identity.UserAssignedIdentity != nil {
					identityId = *queue.Identity.UserAssignedIdentity
				}
				output["identity_id"] = identityId

				if name := queue.Name; name != nil {
					output["name"] = *name
				}
				if resourceGroup := queue.ResourceGroup; resourceGroup != nil {
					output["resource_group_name"] = *resourceGroup
				}

				output["type"] = "AzureIotHub.ServiceBusQueue"
				output["subscription_id"] = pointer.From(queue.SubscriptionID)

				results = append(results, output)
			}
		}

		if topics := input.Endpoints.ServiceBusTopics; topics != nil {
			for _, topic := range *topics {
				output := make(map[string]any)

				authenticationType := string(devices.AuthenticationTypeKeyBased)
				if string(topic.AuthenticationType) != "" {
					authenticationType = string(topic.AuthenticationType)
				}
				output["authentication_type"] = authenticationType

				output["connection_string"] = pointer.From(topic.ConnectionString)

				output["endpoint_uri"] = pointer.From(topic.EndpointURI)

				output["entity_path"] = pointer.From(topic.EntityPath)

				identityId := ""
				if topic.Identity != nil && topic.Identity.UserAssignedIdentity != nil {
					identityId = *topic.Identity.UserAssignedIdentity
				}
				output["identity_id"] = identityId

				if name := topic.Name; name != nil {
					output["name"] = *name
				}
				if resourceGroup := topic.ResourceGroup; resourceGroup != nil {
					output["resource_group_name"] = *resourceGroup
				}

				output["type"] = "AzureIotHub.ServiceBusTopic"
				output["subscription_id"] = pointer.From(topic.SubscriptionID)

				results = append(results, output)
			}
		}

		if eventHubs := input.Endpoints.EventHubs; eventHubs != nil {
			for _, eventHub := range *eventHubs {
				output := make(map[string]any)

				authenticationType := string(devices.AuthenticationTypeKeyBased)
				if string(eventHub.AuthenticationType) != "" {
					authenticationType = string(eventHub.AuthenticationType)
				}
				output["authentication_type"] = authenticationType

				output["connection_string"] = pointer.From(eventHub.ConnectionString)

				output["endpoint_uri"] = pointer.From(eventHub.EndpointURI)

				output["entity_path"] = pointer.From(eventHub.EntityPath)

				identityId := ""
				if eventHub.Identity != nil && eventHub.Identity.UserAssignedIdentity != nil {
					identityId = *eventHub.Identity.UserAssignedIdentity
				}
				output["identity_id"] = identityId

				if name := eventHub.Name; name != nil {
					output["name"] = *name
				}
				if resourceGroup := eventHub.ResourceGroup; resourceGroup != nil {
					output["resource_group_name"] = *resourceGroup
				}

				output["type"] = "AzureIotHub.EventHub"
				output["subscription_id"] = pointer.From(eventHub.SubscriptionID)

				results = append(results, output)
			}
		}
	}

	return results
}

func flattenIoTHubRoute(input *devices.RoutingProperties) []any {
	results := make([]any, 0)

	if input != nil && input.Routes != nil {
		for _, route := range *input.Routes {
			output := make(map[string]any)

			if name := route.Name; name != nil {
				output["name"] = *name
			}
			if condition := route.Condition; condition != nil {
				output["condition"] = *condition
			}
			if endpointNames := route.EndpointNames; endpointNames != nil {
				output["endpoint_names"] = *endpointNames
			}
			if isEnabled := route.IsEnabled; isEnabled != nil {
				output["enabled"] = *isEnabled
			}
			output["source"] = route.Source

			results = append(results, output)
		}
	}

	return results
}

func flattenIoTHubEnrichment(input *devices.RoutingProperties) []any {
	results := make([]any, 0)

	if input != nil && input.Enrichments != nil {
		for _, enrichment := range *input.Enrichments {
			output := make(map[string]any)

			if key := enrichment.Key; key != nil {
				output["key"] = *key
			}
			if value := enrichment.Value; value != nil {
				output["value"] = *value
			}
			if endpointNames := enrichment.EndpointNames; endpointNames != nil {
				output["endpoint_names"] = *endpointNames
			}

			results = append(results, output)
		}
	}

	return results
}

func flattenIoTHubFallbackRoute(input *devices.RoutingProperties) []any {
	if input.FallbackRoute == nil {
		return []any{}
	}

	output := make(map[string]any)
	route := input.FallbackRoute

	if condition := route.Condition; condition != nil {
		output["condition"] = *condition
	}
	if isEnabled := route.IsEnabled; isEnabled != nil {
		output["enabled"] = *isEnabled
	}
	if source := route.Source; source != nil {
		output["source"] = *source
	}

	output["endpoint_names"] = pluginsdk.FlattenSlice(route.EndpointNames)

	return []any{output}
}

func flattenIoTHubCloudToDevice(input *devices.CloudToDeviceProperties) []any {
	if input == nil {
		return []any{}
	}

	output := make(map[string]any)

	if maxDeliveryCount := input.MaxDeliveryCount; maxDeliveryCount != nil {
		output["max_delivery_count"] = *maxDeliveryCount
	}
	if defaultTimeToLive := input.DefaultTTLAsIso8601; defaultTimeToLive != nil {
		output["default_ttl"] = *defaultTimeToLive
	}

	output["feedback"] = flattenIoTHubCloudToDeviceFeedback(input.Feedback)

	return []any{output}
}

func flattenIoTHubCloudToDeviceFeedback(input *devices.FeedbackProperties) []any {
	if input == nil {
		return []any{}
	}

	feedback := make(map[string]any)
	if feedbackMaxDeliveryCount := input.MaxDeliveryCount; feedbackMaxDeliveryCount != nil {
		feedback["max_delivery_count"] = *feedbackMaxDeliveryCount
	}
	if feedbackTimeToLive := input.TTLAsIso8601; feedbackTimeToLive != nil {
		feedback["time_to_live"] = *feedbackTimeToLive
	}
	if feedbackLockDuration := input.LockDurationAsIso8601; feedbackLockDuration != nil {
		feedback["lock_duration"] = *feedbackLockDuration
	}

	return []any{feedback}
}

func flattenNetworkRuleSetProperties(input *devices.NetworkRuleSetProperties) []any {
	if input == nil {
		return []any{}
	}

	output := make(map[string]any)
	output["default_action"] = input.DefaultAction
	output["apply_to_builtin_eventhub_endpoint"] = input.ApplyToBuiltInEventHubEndpoint
	rules := make([]any, 0)

	for _, r := range *input.IPRules {
		rawRule := make(map[string]any)

		if r.FilterName != nil {
			rawRule["name"] = *r.FilterName
		}

		rawRule["action"] = string(r.Action)

		if r.IPMask != nil {
			rawRule["ip_mask"] = *r.IPMask
		}
		rules = append(rules, rawRule)
	}

	output["ip_rule"] = rules
	return []any{output}
}

func flattenIotHubIdentity(input *devices.ArmIdentity) (*[]any, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(input.Type),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}
