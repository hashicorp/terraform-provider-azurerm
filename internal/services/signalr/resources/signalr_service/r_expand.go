// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandUpstreamSettings(input []any) *signalr.ServerlessUpstreamSettings {
	upstreamTemplates := make([]signalr.UpstreamTemplate, 0)

	for _, upstreamSetting := range input {
		setting := upstreamSetting.(map[string]any)
		auth := signalr.UpstreamAuthSettings{
			Type: pointer.To(signalr.UpstreamAuthTypeNone),
		}
		upstreamTemplate := signalr.UpstreamTemplate{
			HubPattern:      pointer.To(strings.Join(*pluginsdk.ExpandStringSlice(setting["hub_pattern"].([]any)), ",")),
			EventPattern:    pointer.To(strings.Join(*pluginsdk.ExpandStringSlice(setting["event_pattern"].([]any)), ",")),
			CategoryPattern: pointer.To(strings.Join(*pluginsdk.ExpandStringSlice(setting["category_pattern"].([]any)), ",")),
			UrlTemplate:     setting["url_template"].(string),
			Auth:            &auth,
		}

		if setting["user_assigned_identity_id"].(string) != "" {
			upstreamTemplate.Auth = &signalr.UpstreamAuthSettings{
				Type: pointer.To(signalr.UpstreamAuthTypeManagedIdentity),
				ManagedIdentity: &signalr.ManagedIdentitySettings{
					Resource: pointer.To(setting["user_assigned_identity_id"].(string)),
				},
			}
		}

		upstreamTemplates = append(upstreamTemplates, upstreamTemplate)
	}

	return &signalr.ServerlessUpstreamSettings{
		Templates: &upstreamTemplates,
	}
}

func expandSignalRCors(input []any) *signalr.SignalRCorsSettings {
	corsSettings := signalr.SignalRCorsSettings{}

	if len(input) == 0 || input[0] == nil {
		return &corsSettings
	}

	setting := input[0].(map[string]any)
	origins := setting["allowed_origins"].(*pluginsdk.Set).List()

	allowedOrigins := make([]string, 0)
	for _, param := range origins {
		allowedOrigins = append(allowedOrigins, param.(string))
	}

	corsSettings.AllowedOrigins = &allowedOrigins

	return &corsSettings
}

func expandSignalRServiceSku(input []any) *signalr.ResourceSku {
	v := input[0].(map[string]any)
	return &signalr.ResourceSku{
		Name:     v["name"].(string),
		Capacity: pointer.To(int64(v["capacity"].(int))),
	}
}

func expandSignalRLiveTraceConfig(input []any) *signalr.LiveTraceConfiguration {
	resourceCategories := make([]signalr.LiveTraceCategory, 0)
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)

	enabled := "false"
	if v["enabled"].(bool) {
		enabled = "true"
	}

	messageLogEnabled := "false"
	if v["messaging_logs_enabled"].(bool) {
		messageLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, signalr.LiveTraceCategory{
		Name:    pointer.To("MessagingLogs"),
		Enabled: pointer.To(messageLogEnabled),
	})

	connectivityLogEnabled := "false"
	if v["connectivity_logs_enabled"].(bool) {
		connectivityLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, signalr.LiveTraceCategory{
		Name:    pointer.To("ConnectivityLogs"),
		Enabled: pointer.To(connectivityLogEnabled),
	})

	httpLogEnabled := "false"
	if v["http_request_logs_enabled"].(bool) {
		httpLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, signalr.LiveTraceCategory{
		Name:    pointer.To("HttpRequestLogs"),
		Enabled: pointer.To(httpLogEnabled),
	})

	return &signalr.LiveTraceConfiguration{
		Enabled:    &enabled,
		Categories: &resourceCategories,
	}
}

func expandSignalRResourceLogConfig(connectivityLogEnabled bool, messagingLogEnabled bool, httpLogEnabled bool) *signalr.ResourceLogConfiguration {
	resourceLogCategories := make([]signalr.ResourceLogCategory, 0)

	messagingLog := "false"
	if messagingLogEnabled {
		messagingLog = "true"
	}
	resourceLogCategories = append(resourceLogCategories, signalr.ResourceLogCategory{
		Name:    pointer.To("MessagingLogs"),
		Enabled: pointer.To(messagingLog),
	})

	connectivityLog := "false"
	if connectivityLogEnabled {
		connectivityLog = "true"
	}
	resourceLogCategories = append(resourceLogCategories, signalr.ResourceLogCategory{
		Name:    pointer.To("ConnectivityLogs"),
		Enabled: pointer.To(connectivityLog),
	})

	httpLog := "false"
	if httpLogEnabled {
		httpLog = "true"
	}
	resourceLogCategories = append(resourceLogCategories, signalr.ResourceLogCategory{
		Name:    pointer.To("HttpRequestLogs"),
		Enabled: pointer.To(httpLog),
	})

	return &signalr.ResourceLogConfiguration{
		Categories: &resourceLogCategories,
	}
}
