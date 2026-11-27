// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_socketio

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (w Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.WebPubSubClient.WebPubSub

			id, err := webpubsub.ParseWebPubSubID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			keys, err := client.ListKeys(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", *id, err)
			}

			return w.flatten(metadata, id, resp.Model, keys.Model)
		},
	}
}
