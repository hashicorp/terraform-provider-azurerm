// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_socketio

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
)

func (w Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := signalrClient.From(metadata.Client).WebPubSubClient.WebPubSub

			id, err := webpubsub.ParseWebPubSubID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config WebPubSubSocketIOResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Sku == nil {
				return fmt.Errorf("retrieving %s: `model.Sku` was nil", *id)
			}

			payload := webpubsub.WebPubSubResource{
				Location:   location.Normalize(config.Location),
				Properties: &webpubsub.WebPubSubProperties{},
				Sku:        existing.Model.Sku,
			}
			props := payload.Properties

			rd := metadata.ResourceData
			if rd.HasChange("sku") {
				payload.Sku = expandWebPubSubSocketIOSkuFromModel(config.Sku)
			}

			if rd.HasChange("aad_auth_enabled") {
				props.DisableAadAuth = pointer.To(!config.AADAuthEnabled)
			}

			if rd.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemOrUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if rd.HasChanges(
				"live_trace_enabled",
				"live_trace_connectivity_logs_enabled",
				"live_trace_http_request_logs_enabled",
				"live_trace_messaging_logs_enabled",
			) {
				props.LiveTraceConfiguration = expandLiveTraceConfigFromModel(config)
			}

			if rd.HasChange("local_auth_enabled") {
				props.DisableLocalAuth = pointer.To(!config.LocalAuthEnabled)
			}

			if rd.HasChange("public_network_access") {
				props.PublicNetworkAccess = pointer.To(config.PublicNetworkAccess)
			}

			if rd.HasChange("service_mode") {
				props.SocketIO = &webpubsub.WebPubSubSocketIOSettings{
					ServiceMode: pointer.To(config.ServiceMode),
				}
			}

			if rd.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if rd.HasChange("tls_client_cert_enabled") {
				props.Tls = &webpubsub.WebPubSubTlsSettings{
					ClientCertEnabled: pointer.To(config.TlsClientCertEnabled),
				}
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
