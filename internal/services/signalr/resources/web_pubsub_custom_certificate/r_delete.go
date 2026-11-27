// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_custom_certificate

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := signalrClient.From(metadata.Client).WebPubSubClient.WebPubSub

			id, err := webpubsub.ParseCustomCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			webPubsubId := webpubsub.NewWebPubSubID(id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName)

			locks.ByID(webPubsubId.ID())
			defer locks.UnlockByID(webPubsubId.ID())

			if _, err := client.CustomCertificatesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			poller := custompollers.NewEventualConsistencyPoller(20, func(pollerCtx context.Context) (*http.Response, error) {
				resp, err := client.CustomCertificatesGet(pollerCtx, *id)
				return resp.HttpResponse, err
			}, custompollers.DefaultDeletionEventualConsistencyPollerOptions())
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be fully deleted: %+v", *id, err)
			}
			return nil
		},
	}
}
