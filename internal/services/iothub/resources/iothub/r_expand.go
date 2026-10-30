// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func expandIoTHubRoutes(d *pluginsdk.ResourceData) *[]devices.RouteProperties {
	routeList := d.Get("route").([]any)

	routeProperties := make([]devices.RouteProperties, 0)

	for _, routeRaw := range routeList {
		route := routeRaw.(map[string]any)

		source := devices.RoutingSource(route["source"].(string))

		endpointNamesRaw := route["endpoint_names"].([]any)

		routeProperties = append(routeProperties, devices.RouteProperties{
			Name:          pointer.To(route["name"].(string)),
			Source:        source,
			Condition:     pointer.To(route["condition"].(string)),
			EndpointNames: pluginsdk.ExpandStringSlice(endpointNamesRaw),
			IsEnabled:     pointer.To(route["enabled"].(bool)),
		})
	}

	return &routeProperties
}

func expandIoTHubEnrichments(d *pluginsdk.ResourceData) *[]devices.EnrichmentProperties {
	enrichmentList := d.Get("enrichment").([]any)

	enrichmentProperties := make([]devices.EnrichmentProperties, 0)

	for _, enrichmentRaw := range enrichmentList {
		enrichment := enrichmentRaw.(map[string]any)

		endpointNamesRaw := enrichment["endpoint_names"].([]any)

		enrichmentProperties = append(enrichmentProperties, devices.EnrichmentProperties{
			Key:           pointer.To(enrichment["key"].(string)),
			Value:         pointer.To(enrichment["value"].(string)),
			EndpointNames: pluginsdk.ExpandStringSlice(endpointNamesRaw),
		})
	}

	return &enrichmentProperties
}

func expandIoTHubFileUpload(d *pluginsdk.ResourceData) (map[string]*devices.StorageEndpointProperties, map[string]*devices.MessagingEndpointProperties, bool, error) {
	fileUploadList := d.Get("file_upload").([]any)

	storageEndpointProperties := make(map[string]*devices.StorageEndpointProperties)
	messagingEndpointProperties := make(map[string]*devices.MessagingEndpointProperties)
	notifications := false

	if len(fileUploadList) > 0 {
		fileUploadMap := fileUploadList[0].(map[string]any)

		authenticationType := devices.AuthenticationType(fileUploadMap["authentication_type"].(string))
		identityId := fileUploadMap["identity_id"].(string)
		notifications = fileUploadMap["notifications"].(bool)
		sasTTL := fileUploadMap["sas_ttl"].(string)

		storageEndpointProperties["$default"] = &devices.StorageEndpointProperties{
			AuthenticationType: authenticationType,
			ConnectionString:   pointer.To(fileUploadMap["connection_string"].(string)),
			ContainerName:      pointer.To(fileUploadMap["container_name"].(string)),
		}

		if sasTTL != "" {
			storageEndpointProperties["$default"].SasTTLAsIso8601 = &sasTTL
		}

		if identityId != "" {
			if authenticationType != devices.AuthenticationTypeIdentityBased {
				return nil, nil, false, fmt.Errorf("`identity_id` can only be specified when `authentication_type` is `identityBased`")
			}
			storageEndpointProperties["$default"].Identity = &devices.ManagedIdentity{
				UserAssignedIdentity: &identityId,
			}
		}

		messagingEndpointProperties["fileNotifications"] = &devices.MessagingEndpointProperties{
			LockDurationAsIso8601: pointer.To(fileUploadMap["lock_duration"].(string)),
			TTLAsIso8601:          pointer.To(fileUploadMap["default_ttl"].(string)),
			MaxDeliveryCount:      pointer.To(int32(fileUploadMap["max_delivery_count"].(int))),
		}
	}

	return storageEndpointProperties, messagingEndpointProperties, notifications, nil
}

func expandIoTHubEndpoints(d *pluginsdk.ResourceData, subscriptionId string) (*devices.RoutingEndpoints, error) {
	routeEndpointList := d.Get("endpoint").([]any)

	serviceBusQueueEndpointProperties := make([]devices.RoutingServiceBusQueueEndpointProperties, 0)
	serviceBusTopicEndpointProperties := make([]devices.RoutingServiceBusTopicEndpointProperties, 0)
	eventHubProperties := make([]devices.RoutingEventHubProperties, 0)
	storageContainerProperties := make([]devices.RoutingStorageContainerProperties, 0)

	for k, endpointRaw := range routeEndpointList {
		endpoint := endpointRaw.(map[string]any)

		t := endpoint["type"]
		name := endpoint["name"].(string)
		resourceGroup := endpoint["resource_group_name"].(string)
		authenticationType := devices.AuthenticationType(endpoint["authentication_type"].(string))

		var subscriptionID string
		// To align with the previous TF behaviour, `subscription_id` needs to be set with the provider's subscription Id when it isn't specified in the tf config, otherwise TF behaviour is different than before and it may block the existing users
		// From the business perspective, the raw config handling is only meant for the case that the user has an EventHub whose Endpoint's subscription is not the provider's one. Then the user wants to reset it to the provider's one by unset the subscription_id
		// From the TF code perspective, given `Computed: true` is enabled, TF would always get the value from the last apply when this property isn't set in the tf config. So `d.GetRawConfig()` is required to determine if it's set in the tf config
		if v := d.GetRawConfig().AsValueMap()["endpoint"].AsValueSlice()[k].AsValueMap()["subscription_id"]; v.IsNull() {
			subscriptionID = subscriptionId
		} else {
			subscriptionID = endpoint["subscription_id"].(string)
		}

		var identity *devices.ManagedIdentity
		var endpointUri *string
		var entityPath *string
		var connectionStr *string
		if v := endpoint["identity_id"].(string); v != "" {
			identity = &devices.ManagedIdentity{
				UserAssignedIdentity: pointer.To(v),
			}
		}
		if v := endpoint["endpoint_uri"].(string); v != "" {
			endpointUri = pointer.To(v)
		}
		if v := endpoint["entity_path"].(string); v != "" {
			entityPath = pointer.To(v)
		}
		if v := endpoint["connection_string"].(string); v != "" {
			connectionStr = pointer.To(v)
		}

		if authenticationType == devices.AuthenticationTypeKeyBased {
			if connectionStr == nil {
				return nil, fmt.Errorf("`connection_string` must be specified when `authentication_type` is `keyBased`")
			}
			if identity != nil || endpointUri != nil || entityPath != nil {
				return nil, fmt.Errorf("`identity_id`, `endpoint_uri` or `entity_path` cannot be specified when `authentication_type` is `keyBased`")
			}
		} else {
			if endpointUri == nil {
				return nil, fmt.Errorf("`endpoint_uri` must be specified when `authentication_type` is `identityBased`")
			}

			if entityPath == nil && t != "AzureIotHub.StorageContainer" {
				return nil, fmt.Errorf("`entity_path` must be specified when `authentication_type` is `identityBased` and `type` is `%s`", t)
			}

			if connectionStr != nil {
				return nil, fmt.Errorf("`connection_string` cannot be specified when `authentication_type` is `identityBased`")
			}
		}

		switch t {
		case "AzureIotHub.StorageContainer":
			containerName := endpoint["container_name"].(string)
			if containerName == "" {
				return nil, fmt.Errorf("`container_name` must be specified when `type` is `AzureIotHub.StorageContainer`")
			}

			encoding := endpoint["encoding"].(string)

			storageContainer := devices.RoutingStorageContainerProperties{
				AuthenticationType:      authenticationType,
				Identity:                identity,
				EndpointURI:             endpointUri,
				ConnectionString:        connectionStr,
				Name:                    &name,
				SubscriptionID:          &subscriptionID,
				ResourceGroup:           &resourceGroup,
				ContainerName:           &containerName,
				FileNameFormat:          pointer.To(endpoint["file_name_format"].(string)),
				BatchFrequencyInSeconds: pointer.To(int32(endpoint["batch_frequency_in_seconds"].(int))),
				MaxChunkSizeInBytes:     pointer.To(int32(endpoint["max_chunk_size_in_bytes"].(int))),
				Encoding:                devices.Encoding(encoding),
			}
			storageContainerProperties = append(storageContainerProperties, storageContainer)

		case "AzureIotHub.ServiceBusQueue":
			sbQueue := devices.RoutingServiceBusQueueEndpointProperties{
				AuthenticationType: authenticationType,
				Identity:           identity,
				EndpointURI:        endpointUri,
				EntityPath:         entityPath,
				ConnectionString:   connectionStr,
				Name:               &name,
				SubscriptionID:     &subscriptionID,
				ResourceGroup:      &resourceGroup,
			}
			serviceBusQueueEndpointProperties = append(serviceBusQueueEndpointProperties, sbQueue)

		case "AzureIotHub.ServiceBusTopic":
			sbTopic := devices.RoutingServiceBusTopicEndpointProperties{
				AuthenticationType: authenticationType,
				Identity:           identity,
				EndpointURI:        endpointUri,
				EntityPath:         entityPath,
				ConnectionString:   connectionStr,
				Name:               &name,
				SubscriptionID:     &subscriptionID,
				ResourceGroup:      &resourceGroup,
			}
			serviceBusTopicEndpointProperties = append(serviceBusTopicEndpointProperties, sbTopic)

		case "AzureIotHub.EventHub":
			eventHub := devices.RoutingEventHubProperties{
				AuthenticationType: authenticationType,
				Identity:           identity,
				EndpointURI:        endpointUri,
				EntityPath:         entityPath,
				ConnectionString:   connectionStr,
				Name:               &name,
				SubscriptionID:     &subscriptionID,
				ResourceGroup:      &resourceGroup,
			}
			eventHubProperties = append(eventHubProperties, eventHub)
		}
	}

	return &devices.RoutingEndpoints{
		ServiceBusQueues:  &serviceBusQueueEndpointProperties,
		ServiceBusTopics:  &serviceBusTopicEndpointProperties,
		EventHubs:         &eventHubProperties,
		StorageContainers: &storageContainerProperties,
	}, nil
}

func expandIoTHubFallbackRoute(d *pluginsdk.ResourceData) *devices.FallbackRouteProperties {
	fallbackRouteList := d.Get("fallback_route").([]any)
	if len(fallbackRouteList) == 0 {
		return nil
	}

	fallbackRouteMap := fallbackRouteList[0].(map[string]any)

	return &devices.FallbackRouteProperties{
		Source:        pointer.To(fallbackRouteMap["source"].(string)),
		Condition:     pointer.To(fallbackRouteMap["condition"].(string)),
		EndpointNames: pluginsdk.ExpandStringSlice(fallbackRouteMap["endpoint_names"].([]any)),
		IsEnabled:     pointer.To(fallbackRouteMap["enabled"].(bool)),
	}
}

func expandIoTHubSku(d *pluginsdk.ResourceData) *devices.IotHubSkuInfo {
	skuList := d.Get("sku").([]any)
	skuMap := skuList[0].(map[string]any)

	return &devices.IotHubSkuInfo{
		Name:     devices.IotHubSku(skuMap["name"].(string)),
		Capacity: pointer.To(int64(skuMap["capacity"].(int))),
	}
}

func expandIoTHubCloudToDevice(d *pluginsdk.ResourceData) *devices.CloudToDeviceProperties {
	ctdList := d.Get("cloud_to_device").([]any)
	if len(ctdList) == 0 {
		return nil
	}
	cloudToDevice := devices.CloudToDeviceProperties{}
	ctdMap := ctdList[0].(map[string]any)

	cloudToDevice.DefaultTTLAsIso8601 = pointer.To(ctdMap["default_ttl"].(string))
	cloudToDevice.MaxDeliveryCount = pointer.To(int32(ctdMap["max_delivery_count"].(int)))
	feedback := ctdMap["feedback"].([]any)

	cloudToDeviceFeedback := devices.FeedbackProperties{}
	if len(feedback) > 0 {
		feedbackMap := feedback[0].(map[string]any)

		cloudToDeviceFeedback.TTLAsIso8601 = pointer.To(feedbackMap["time_to_live"].(string))
		cloudToDeviceFeedback.LockDurationAsIso8601 = pointer.To(feedbackMap["lock_duration"].(string))
		cloudToDeviceFeedback.MaxDeliveryCount = pointer.To(int32(feedbackMap["max_delivery_count"].(int)))
	}

	cloudToDevice.Feedback = &cloudToDeviceFeedback

	return &cloudToDevice
}

func expandNetworkRuleSetProperties(d *pluginsdk.ResourceData) *devices.NetworkRuleSetProperties {
	networkRuleSet := d.Get("network_rule_set").([]any)
	networkRuleSetProps := devices.NetworkRuleSetProperties{}
	nrsMap := networkRuleSet[0].(map[string]any)

	networkRuleSetProps.DefaultAction = devices.DefaultAction(nrsMap["default_action"].(string))
	networkRuleSetProps.ApplyToBuiltInEventHubEndpoint = pointer.To(nrsMap["apply_to_builtin_eventhub_endpoint"].(bool))
	ipRules := nrsMap["ip_rule"].([]any)

	if len(ipRules) != 0 {
		rules := make([]devices.NetworkRuleSetIPRule, 0)

		for _, r := range ipRules {
			rawRule := r.(map[string]any)
			rule := &devices.NetworkRuleSetIPRule{
				FilterName: pointer.To(rawRule["name"].(string)),
				Action:     devices.NetworkRuleIPAction(rawRule["action"].(string)),
				IPMask:     pointer.To(rawRule["ip_mask"].(string)),
			}
			rules = append(rules, *rule)
		}
		networkRuleSetProps.IPRules = &rules
	}
	return &networkRuleSetProps
}

func expandIotHubIdentity(input []any) (*devices.ArmIdentity, error) {
	config, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	identity := devices.ArmIdentity{
		Type: devices.ResourceIdentityType(config.Type),
	}

	if len(config.IdentityIds) != 0 {
		identityIds := make(map[string]*devices.ArmUserIdentity, len(config.IdentityIds))
		for id := range config.IdentityIds {
			identityIds[id] = &devices.ArmUserIdentity{}
		}
		identity.UserAssignedIdentities = identityIds
	}

	return &identity, nil
}
