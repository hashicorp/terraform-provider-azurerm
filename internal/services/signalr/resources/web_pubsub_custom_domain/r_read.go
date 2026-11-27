// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_custom_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := signalrClient.From(metadata.Client).WebPubSubClient.WebPubSub
			id, err := webpubsub.ParseCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.CustomDomainsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

const webPubsubCustomDomainResourceType = "azurerm_web_pubsub_custom_domain"
