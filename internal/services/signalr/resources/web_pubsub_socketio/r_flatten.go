// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_socketio

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (w Resource) flatten(metadata sdk.ResourceMetaData, id *webpubsub.WebPubSubId, model *webpubsub.WebPubSubResource, keyModel *webpubsub.WebPubSubKeys) error {
	state := WebPubSubSocketIOResourceModel{
		Name:              id.WebPubSubName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMapToModel(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		state.Identity = pointer.From(flattenedIdentity)

		if sku := model.Sku; sku != nil {
			state.Sku = flattenWebPubSubSocketIOSkuToModel(sku)
		}

		if props := model.Properties; props != nil {
			liveTrace := flattenLiveTraceConfigToMap(props.LiveTraceConfiguration)

			state.AADAuthEnabled = !pointer.From(props.DisableAadAuth)
			state.ExternalIP = pointer.From(props.ExternalIP)
			state.HostName = pointer.From(props.HostName)
			state.PublicPort = pointer.From(props.PublicPort)
			state.ServerPort = pointer.From(props.ServerPort)

			state.LiveTraceEnabled = liveTrace["enabled"]
			state.LiveTraceConnectivityLogsEnabled = liveTrace["connectivityLogsEnabled"]
			state.LiveTraceHttpRequestLogsEnabled = liveTrace["httpLogsEnabled"]
			state.LiveTraceMessagingLogsEnabled = liveTrace["messagingLogsEnabled"]

			state.LocalAuthEnabled = !pointer.From(props.DisableLocalAuth)
			state.PublicNetworkAccess = pointer.From(props.PublicNetworkAccess)

			if socketio := props.SocketIO; socketio != nil {
				state.ServiceMode = pointer.From(socketio.ServiceMode)
			}

			if tls := props.Tls; tls != nil {
				state.TlsClientCertEnabled = pointer.From(tls.ClientCertEnabled)
			}
		}
	}

	if keyModel != nil {
		state.PrimaryAccessKey = pointer.From(keyModel.PrimaryKey)
		state.PrimaryConnectionString = pointer.From(keyModel.PrimaryConnectionString)
		state.SecondaryAccessKey = pointer.From(keyModel.SecondaryKey)
		state.SecondaryConnectionString = pointer.From(keyModel.SecondaryConnectionString)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenLiveTraceConfigToMap(input *webpubsub.LiveTraceConfiguration) map[string]bool {
	result := make(map[string]bool)
	if input == nil {
		return result
	}

	result["enabled"] = strings.EqualFold(pointer.From(input.Enabled), "true")

	if input.Categories != nil {
		for _, item := range *input.Categories {
			name := pointer.From(item.Name)
			categoryEnabled := pointer.From(item.Enabled)

			switch name {
			case "ConnectivityLogs":
				result["connectivityLogsEnabled"] = strings.EqualFold(categoryEnabled, "true")
			case "HttpRequestLogs":
				result["httpLogsEnabled"] = strings.EqualFold(categoryEnabled, "true")
			case "MessagingLogs":
				result["messagingLogsEnabled"] = strings.EqualFold(categoryEnabled, "true")
			default:
				continue
			}
		}
	}

	return result
}

func flattenWebPubSubSocketIOSkuToModel(input *webpubsub.ResourceSku) []WebPubSubSocketIOSkuModel {
	result := make([]WebPubSubSocketIOSkuModel, 0, 1)
	if input == nil {
		return result
	}

	result = append(result, WebPubSubSocketIOSkuModel{
		Name:     input.Name,
		Capacity: pointer.From(input.Capacity),
	})

	return result
}
