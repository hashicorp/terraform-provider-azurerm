// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceArmSignalRServiceFlatten(d *pluginsdk.ResourceData, id *signalr.SignalRId, model *signalr.SignalRResource, keyModel *signalr.SignalRKeys) error {
	d.Set("name", id.SignalRName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		if err := d.Set("sku", flattenSignalRServiceSku(model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("hostname", props.HostName)
			d.Set("ip_address", props.ExternalIP)
			d.Set("public_port", props.PublicPort)
			d.Set("server_port", props.ServerPort)

			connectivityLogsEnabled := false
			messagingLogsEnabled := false
			httpLogsEnabled := false
			liveTraceEnabled := false
			serviceMode := "Default"
			if props.Features != nil {
				for _, feature := range *props.Features {
					if feature.Flag == "EnableLiveTrace" {
						liveTraceEnabled = strings.EqualFold(feature.Value, "True")
					}
					if feature.Flag == signalr.FeatureFlagsServiceMode {
						serviceMode = feature.Value
					}
				}
			}

			d.Set("live_trace_enabled", liveTraceEnabled)
			d.Set("service_mode", serviceMode)

			aadAuthEnabled := true
			if props.DisableAadAuth != nil {
				aadAuthEnabled = !(*props.DisableAadAuth)
			}
			d.Set("aad_auth_enabled", aadAuthEnabled)

			localAuthEnabled := true
			if props.DisableLocalAuth != nil {
				localAuthEnabled = !(*props.DisableLocalAuth)
			}
			d.Set("local_auth_enabled", localAuthEnabled)

			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil {
				publicNetworkAccessEnabled = strings.EqualFold(*props.PublicNetworkAccess, "Enabled")
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

			tlsClientCertEnabled := false
			if props.Tls != nil && props.Tls.ClientCertEnabled != nil {
				tlsClientCertEnabled = *props.Tls.ClientCertEnabled
			}
			d.Set("tls_client_cert_enabled", tlsClientCertEnabled)

			if props.Serverless != nil && props.Serverless.ConnectionTimeoutInSeconds != nil {
				d.Set("serverless_connection_timeout_in_seconds", int(*props.Serverless.ConnectionTimeoutInSeconds))
			}

			if err := d.Set("cors", flattenSignalRCors(props.Cors)); err != nil {
				return fmt.Errorf("setting `cors`: %+v", err)
			}

			if err := d.Set("upstream_endpoint", flattenUpstreamSettings(props.Upstream)); err != nil {
				return fmt.Errorf("setting `upstream_endpoint`: %+v", err)
			}

			if err := d.Set("live_trace", flattenSignalRLiveTraceConfig(props.LiveTraceConfiguration)); err != nil {
				return fmt.Errorf("setting `live_trace`:%+v", err)
			}

			if props.ResourceLogConfiguration != nil && props.ResourceLogConfiguration.Categories != nil {
				for _, item := range *props.ResourceLogConfiguration.Categories {
					name := pointer.From(item.Name)

					cateEnabled := pointer.From(item.Enabled)

					switch name {
					case "MessagingLogs":
						messagingLogsEnabled = strings.EqualFold(cateEnabled, "true")
					case "ConnectivityLogs":
						connectivityLogsEnabled = strings.EqualFold(cateEnabled, "true")
					case "HttpRequestLogs":
						httpLogsEnabled = strings.EqualFold(cateEnabled, "true")
					default:
						continue
					}
				}
			}
			d.Set("connectivity_logs_enabled", connectivityLogsEnabled)
			d.Set("messaging_logs_enabled", messagingLogsEnabled)
			d.Set("http_request_logs_enabled", httpLogsEnabled)
		}

		identityValue, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identityValue); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	if keyModel != nil {
		d.Set("primary_access_key", keyModel.PrimaryKey)
		d.Set("primary_connection_string", keyModel.PrimaryConnectionString)
		d.Set("secondary_access_key", keyModel.SecondaryKey)
		d.Set("secondary_connection_string", keyModel.SecondaryConnectionString)
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenUpstreamSettings(upstreamSettings *signalr.ServerlessUpstreamSettings) []any {
	result := make([]any, 0)
	if upstreamSettings == nil || upstreamSettings.Templates == nil {
		return result
	}

	for _, settings := range *upstreamSettings.Templates {
		categoryPattern := make([]any, 0)
		if settings.CategoryPattern != nil {
			categoryPatterns := strings.Split(*settings.CategoryPattern, ",")
			categoryPattern = pluginsdk.FlattenSlice(&categoryPatterns)
		}

		eventPattern := make([]any, 0)
		if settings.EventPattern != nil {
			eventPatterns := strings.Split(*settings.EventPattern, ",")
			eventPattern = pluginsdk.FlattenSlice(&eventPatterns)
		}

		hubPattern := make([]any, 0)
		if settings.HubPattern != nil {
			hubPatterns := strings.Split(*settings.HubPattern, ",")
			hubPattern = pluginsdk.FlattenSlice(&hubPatterns)
		}

		var managedIdentityId string
		if upstreamAuth := settings.Auth; upstreamAuth != nil && upstreamAuth.Type != nil && *upstreamAuth.Type != signalr.UpstreamAuthTypeNone {
			if upstreamAuth.ManagedIdentity != nil && upstreamAuth.ManagedIdentity.Resource != nil {
				managedIdentityId = *upstreamAuth.ManagedIdentity.Resource
			}
		}

		result = append(result, map[string]any{
			"url_template":              settings.UrlTemplate,
			"hub_pattern":               hubPattern,
			"event_pattern":             eventPattern,
			"category_pattern":          categoryPattern,
			"user_assigned_identity_id": managedIdentityId,
		})
	}
	return result
}

func flattenSignalRCors(input *signalr.SignalRCorsSettings) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	result := make(map[string]any)

	allowedOrigins := make([]any, 0)
	if s := input.AllowedOrigins; s != nil {
		for _, v := range *s {
			allowedOrigins = append(allowedOrigins, v)
		}
	}
	result["allowed_origins"] = pluginsdk.NewSet(pluginsdk.HashString, allowedOrigins)

	return append(results, result)
}

func flattenSignalRServiceSku(input *signalr.ResourceSku) []any {
	if input == nil {
		return []any{}
	}

	capacity := 0
	if input.Capacity != nil {
		capacity = int(*input.Capacity)
	}

	return []any{
		map[string]any{
			"capacity": capacity,
			"name":     input.Name,
		},
	}
}

func flattenSignalRLiveTraceConfig(input *signalr.LiveTraceConfiguration) []any {
	result := make([]any, 0)
	if input == nil {
		return result
	}

	var enabled bool
	if input.Enabled != nil {
		enabled = strings.EqualFold(*input.Enabled, "true")
	}

	var (
		messagingLogEnabled    bool
		connectivityLogEnabled bool
		httpLogsEnabled        bool
	)

	if input.Categories != nil {
		for _, item := range *input.Categories {
			name := pointer.From(item.Name)

			cateEnabled := pointer.From(item.Enabled)

			switch name {
			case "MessagingLogs":
				messagingLogEnabled = strings.EqualFold(cateEnabled, "true")
			case "ConnectivityLogs":
				connectivityLogEnabled = strings.EqualFold(cateEnabled, "true")
			case "HttpRequestLogs":
				httpLogsEnabled = strings.EqualFold(cateEnabled, "true")
			default:
				continue
			}
		}
	}
	return []any{map[string]any{
		"enabled":                   enabled,
		"messaging_logs_enabled":    messagingLogEnabled,
		"connectivity_logs_enabled": connectivityLogEnabled,
		"http_request_logs_enabled": httpLogsEnabled,
	}}
}
