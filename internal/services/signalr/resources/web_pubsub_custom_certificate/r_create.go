// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_custom_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var customCertWebPubsub CustomCertWebPubsubModel
			if err := metadata.Decode(&customCertWebPubsub); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			client := metadata.Client.SignalR.WebPubSubClient.WebPubSub

			webPubsubId, err := webpubsub.ParseWebPubSubID(metadata.ResourceData.Get("web_pubsub_id").(string))
			if err != nil {
				return fmt.Errorf("parsing web pubsub service id error: %+v", err)
			}

			keyVaultCertificateId, err := keyvault.ParseNestedItemID(metadata.ResourceData.Get("custom_certificate_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeCertificate)
			if err != nil {
				return fmt.Errorf("parsing custom certificate id error: %+v", err)
			}

			keyVaultUri := keyVaultCertificateId.KeyVaultBaseURL
			keyVaultSecretName := keyVaultCertificateId.Name

			id := webpubsub.NewCustomCertificateID(webPubsubId.SubscriptionId, webPubsubId.ResourceGroupName, webPubsubId.WebPubSubName, customCertWebPubsub.Name)

			locks.ByID(webPubsubId.ID())
			defer locks.UnlockByID(webPubsubId.ID())

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.CustomCertificatesGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			customCertObj := webpubsub.CustomCertificate{
				Properties: webpubsub.CustomCertificateProperties{
					KeyVaultBaseUri:    keyVaultUri,
					KeyVaultSecretName: keyVaultSecretName,
				},
			}
			if keyVaultCertificateId.Version != "" {
				customCertObj.Properties.KeyVaultSecretVersion = pointer.To(keyVaultCertificateId.Version)
			}

			if err := client.CustomCertificatesCreateOrUpdateCallbackThenPoll(ctx, id, customCertObj, metadata.SetIDAndIdentityCallback(&id)); err != nil {
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
