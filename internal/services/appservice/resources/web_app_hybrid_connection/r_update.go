// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_app_hybrid_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseRelayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var appHybridConn WebAppHybridConnectionModel
			if err := metadata.Decode(&appHybridConn); err != nil {
				return err
			}

			existing, err := client.GetHybridConnection(ctx, *id)
			if err != nil || existing.Model == nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", id, err)
			}

			model := *existing.Model

			if metadata.ResourceData.HasChange("hostname") {
				model.Properties.Hostname = pointer.To(appHybridConn.HostName)
			}

			if metadata.ResourceData.HasChange("port") {
				model.Properties.Port = pointer.To(appHybridConn.HostPort)
			}

			if metadata.ResourceData.HasChange("send_key_name") {
				relayId, err := hybridconnections.ParseHybridConnectionID(appHybridConn.RelayId)
				if err != nil {
					return err
				}

				sendKeyValue, err := helpers.GetSendKeyValue(ctx, metadata, *relayId, appHybridConn.SendKeyName)
				if err != nil {
					return err
				}
				model.Properties.SendKeyValue = sendKeyValue
			}

			if _, err = client.CreateOrUpdateHybridConnection(ctx, *id, model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
