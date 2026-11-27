// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceWebPubSubFlatten(d *pluginsdk.ResourceData, id *webpubsub.WebPubSubId, model *webpubsub.WebPubSubResource, keyModel *webpubsub.WebPubSubKeys) error {
	d.Set("name", id.WebPubSubName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		skuName := ""
		skuCapacity := int64(0)
		if model.Sku != nil {
			skuName = model.Sku.Name
			skuCapacity = *model.Sku.Capacity
		}
		d.Set("sku", skuName)
		d.Set("capacity", skuCapacity)

		if props := model.Properties; props != nil {
			d.Set("external_ip", props.ExternalIP)
			d.Set("hostname", props.HostName)
			d.Set("public_port", props.PublicPort)
			d.Set("server_port", props.ServerPort)
			d.Set("version", props.Version)

			aadAuthEnabled := true
			if props.DisableAadAuth != nil {
				aadAuthEnabled = !(*props.DisableAadAuth)
			}
			d.Set("aad_auth_enabled", aadAuthEnabled)

			disableLocalAuth := false
			if props.DisableLocalAuth != nil {
				disableLocalAuth = !(*props.DisableLocalAuth)
			}
			d.Set("local_auth_enabled", disableLocalAuth)

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

			if err := d.Set("live_trace", flattenLiveTraceConfig(props.LiveTraceConfiguration)); err != nil {
				return fmt.Errorf("setting `live_trace`:%+v", err)
			}

			identityValue, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			if err := d.Set("identity", identityValue); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return fmt.Errorf("setting `tags`: %+v", err)
			}
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

func flattenLiveTraceConfig(input *webpubsub.LiveTraceConfiguration) []any {
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
