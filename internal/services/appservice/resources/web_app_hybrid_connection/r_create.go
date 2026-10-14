// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_app_hybrid_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var appHybridConn WebAppHybridConnectionModel

			client := metadata.Client.AppService.WebAppsClient

			if err := metadata.Decode(&appHybridConn); err != nil {
				return err
			}
			appId, err := commonids.ParseWebAppID(appHybridConn.WebAppId)
			if err != nil {
				return err
			}
			relayId, err := hybridconnections.ParseHybridConnectionID(appHybridConn.RelayId)
			if err != nil {
				return err
			}

			id := webapps.NewRelayID(appId.SubscriptionId, appId.ResourceGroupName, appId.SiteName, relayId.NamespaceName, relayId.HybridConnectionName)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetHybridConnection(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %s", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError(r.ResourceType(), id.ID())
				}
			}

			sendKeyValue, err := helpers.GetSendKeyValue(ctx, metadata, *relayId, appHybridConn.SendKeyName)
			if err != nil {
				return err
			}

			envelope := webapps.HybridConnection{
				Properties: &webapps.HybridConnectionProperties{
					RelayArmUri:  pointer.To(relayId.ID()),
					Hostname:     pointer.To(appHybridConn.HostName),
					Port:         pointer.To(appHybridConn.HostPort),
					SendKeyName:  pointer.To(appHybridConn.SendKeyName),
					SendKeyValue: sendKeyValue,
				},
			}

			if _, err = client.CreateOrUpdateHybridConnection(ctx, id, envelope); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
