// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_custom_certificate

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := signalrClient.From(metadata.Client).SignalRClient

			id, err := signalr.ParseCustomCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			signalrId := signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName)

			locks.ByID(signalrId.ID())
			defer locks.UnlockByID(signalrId.ID())

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
