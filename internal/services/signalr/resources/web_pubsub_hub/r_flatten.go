// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_hub

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenEventHandler(input *[]webpubsub.EventHandler) []any {
	eventHandlerBlock := make([]any, 0)
	if input == nil {
		return eventHandlerBlock
	}

	for _, item := range *input {
		sysEvents := make([]any, 0)
		if item.SystemEvents != nil {
			sysEvents = pluginsdk.FlattenSlice(item.SystemEvents)
		}

		authBlock := make([]any, 0)
		if item.Auth != nil {
			authBlock = flattenAuth(item.Auth)
		}

		eventHandlerBlock = append(eventHandlerBlock, map[string]any{
			"url_template":       item.UrlTemplate,
			"user_event_pattern": pointer.From(item.UserEventPattern),
			"system_events":      sysEvents,
			"auth":               authBlock,
		})
	}
	return eventHandlerBlock
}

func flattenEventListener(listener *[]webpubsub.EventListener) []any {
	eventListenerBlocks := make([]any, 0)
	if listener == nil {
		return eventListenerBlocks
	}

	for _, item := range *listener {
		listenerBlock := make(map[string]any, 0)
		// todo use the type Assertion or Type field in sdk to get the different sub-class
		if eventFilter := item.Filter; eventFilter != nil {
			eventNameFilter := item.Filter.(webpubsub.EventNameFilter)
			userNameFilterList := make([]any, 0)
			if eventNameFilter.SystemEvents != nil {
				listenerBlock["system_event_name_filter"] = pluginsdk.FlattenSlice(eventNameFilter.SystemEvents)
			}
			if eventNameFilter.UserEventPattern != nil && *eventNameFilter.UserEventPattern != "" {
				v := strings.SplitSeq(*eventNameFilter.UserEventPattern, ",")
				for s := range v {
					userNameFilterList = append(userNameFilterList, s)
				}
				listenerBlock["user_event_name_filter"] = userNameFilterList
			}
		}

		if eventEndpoint := item.Endpoint; eventEndpoint != nil {
			eventhubEndpoint := item.Endpoint.(webpubsub.EventHubEndpoint)
			listenerBlock["eventhub_namespace_name"] = strings.TrimSuffix(eventhubEndpoint.FullyQualifiedNamespace, ".servicebus.windows.net")
			listenerBlock["eventhub_name"] = eventhubEndpoint.EventHubName
		}
		eventListenerBlocks = append(eventListenerBlocks, listenerBlock)
	}

	return eventListenerBlocks
}

func flattenAuth(input *webpubsub.UpstreamAuthSettings) []any {
	if input == nil || input.Type == nil || *input.Type == webpubsub.UpstreamAuthTypeNone || input.ManagedIdentity == nil || input.ManagedIdentity.Resource == nil {
		return make([]any, 0)
	}

	authId := *input.ManagedIdentity.Resource

	return []any{
		map[string]any{
			"managed_identity_id": authId,
		},
	}
}
