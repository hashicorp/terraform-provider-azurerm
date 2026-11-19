// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_custom_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.WebPubSubClient.WebPubSub
			keyVaultClient := metadata.Client.KeyVault
			id, err := webpubsub.ParseCustomCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.CustomCertificatesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: got nil model", *id)
			}

			vaultBasedUri := resp.Model.Properties.KeyVaultBaseUri
			subscriptionResourceId := commonids.NewSubscriptionID(id.SubscriptionId)
			keyVaultIdRaw, err := keyVaultClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, vaultBasedUri)
			if err != nil {
				return fmt.Errorf("getting key vault base uri from %s: %+v", id, err)
			}
			if keyVaultIdRaw != nil {
				vaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
				if err != nil {
					return fmt.Errorf("parsing key vault %s: %+v", vaultId, err)
				}
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

const webPubsubCustomCertificateResourceType = "azurerm_web_pubsub_custom_certificate"
