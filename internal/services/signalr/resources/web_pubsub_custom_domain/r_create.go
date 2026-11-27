// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_custom_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var customDomainWebPubsubModel CustomDomainWebPubsubModel
			if err := metadata.Decode(&customDomainWebPubsubModel); err != nil {
				return err
			}
			client := signalrClient.From(metadata.Client).WebPubSubClient.WebPubSub

			webPubsubId, err := webpubsub.ParseWebPubSubIDInsensitively(metadata.ResourceData.Get("web_pubsub_id").(string))
			if err != nil {
				return fmt.Errorf("parsing web pubsub id error: %+v", err)
			}

			id := webpubsub.NewCustomDomainID(webPubsubId.SubscriptionId, webPubsubId.ResourceGroupName, webPubsubId.WebPubSubName, metadata.ResourceData.Get("name").(string))

			locks.ByID(webPubsubId.ID())
			defer locks.UnlockByID(webPubsubId.ID())

			if _, err := webpubsub.ParseCustomCertificateIDInsensitively(customDomainWebPubsubModel.WebPubsubCustomCertificateId); err != nil {
				return fmt.Errorf("parsing custom certificate for %s: %+v", id, err)
			}

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.CustomDomainsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			customDomainObj := webpubsub.CustomDomain{
				Properties: webpubsub.CustomDomainProperties{
					DomainName: customDomainWebPubsubModel.DomainName,
					CustomCertificate: webpubsub.ResourceReference{
						Id: pointer.To(customDomainWebPubsubModel.WebPubsubCustomCertificateId),
					},
				},
			}
			if err := client.CustomDomainsCreateOrUpdateCallbackThenPoll(ctx, id, customDomainObj, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}
