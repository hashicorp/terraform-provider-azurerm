// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_hub

import (
	"errors"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandEventHandler(input []any) *[]webpubsub.EventHandler {
	if len(input) == 0 {
		return nil
	}

	results := make([]webpubsub.EventHandler, 0)

	for _, eventHandlerItem := range input {
		block := eventHandlerItem.(map[string]any)
		eventHandlerSettings := webpubsub.EventHandler{
			UrlTemplate: block["url_template"].(string),
		}

		if v, ok := block["user_event_pattern"]; ok {
			eventHandlerSettings.UserEventPattern = pointer.To(v.(string))
		}

		if v, ok := block["system_events"]; ok {
			systemEvents := make([]string, 0)
			for _, item := range v.(*pluginsdk.Set).List() {
				systemEvents = append(systemEvents, item.(string))
			}
			eventHandlerSettings.SystemEvents = &systemEvents
		}

		if v, ok := block["auth"].([]any); ok {
			if len(v) > 0 {
				eventHandlerSettings.Auth = expandAuth(v)
			}
		}

		results = append(results, eventHandlerSettings)
	}
	return &results
}

func expandEventListener(input []any) (*[]webpubsub.EventListener, error) {
	result := make([]webpubsub.EventListener, 0)
	if len(input) == 0 {
		return &result, nil
	}

	for _, eventListenerItem := range input {
		block := eventListenerItem.(map[string]any)
		systemEvents := make([]string, 0)
		userEventPattern := ""
		if v, ok := block["user_event_name_filter"]; ok && len(v.([]any)) > 0 {
			userEventPatternList := pluginsdk.ExpandStringSlice(v.([]any))
			userEventPattern = strings.Join(*userEventPatternList, ",")
		}

		if v, ok := block["system_event_name_filter"]; ok {
			for _, item := range v.([]any) {
				systemEvents = append(systemEvents, item.(string))
			}
		}
		filter := webpubsub.EventNameFilter{
			SystemEvents:     &systemEvents,
			UserEventPattern: pointer.To(userEventPattern),
		}

		endpointName := block["eventhub_namespace_name"].(string)
		fullQualifiedName := endpointName + ".servicebus.windows.net"
		if _, ok := block["eventhub_name"]; !ok {
			return nil, errors.New("no event hub is specified")
		}
		ehName := block["eventhub_name"].(string)
		endpoint := webpubsub.EventHubEndpoint{
			FullyQualifiedNamespace: fullQualifiedName,
			EventHubName:            ehName,
		}

		result = append(result, webpubsub.EventListener{
			Filter:   filter,
			Endpoint: endpoint,
		})
	}
	return &result, nil
}

func expandAuth(input []any) *webpubsub.UpstreamAuthSettings {
	if len(input) == 0 || input[0] == nil {
		return &webpubsub.UpstreamAuthSettings{
			Type: pointer.To(webpubsub.UpstreamAuthTypeNone),
		}
	}

	authRaw := input[0].(map[string]any)
	return &webpubsub.UpstreamAuthSettings{
		Type: pointer.To(webpubsub.UpstreamAuthTypeManagedIdentity),
		ManagedIdentity: &webpubsub.ManagedIdentitySettings{
			Resource: pointer.To(authRaw["managed_identity_id"].(string)),
		},
	}
}
