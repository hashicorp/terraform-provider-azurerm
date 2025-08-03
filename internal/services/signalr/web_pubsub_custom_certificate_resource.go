// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CustomCertWebPubsubModel struct {
	Name               string `tfschema:"name"`
	WebPubsubId        string `tfschema:"web_pubsub_id"`
	CustomCertId       string `tfschema:"custom_certificate_id"`
	CertificateVersion string `tfschema:"certificate_version"`
}

type CustomCertWebPubsubResource struct{}

var _ sdk.Resource = CustomCertWebPubsubResource{}

func (r CustomCertWebPubsubResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"web_pubsub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: webpubsub.ValidateWebPubSubID,
		},

		"custom_certificate_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				keyVaultValidate.NestedItemId,
				keyVaultValidate.NestedItemIdWithOptionalVersion,
			),
		},
	}
}

func (r CustomCertWebPubsubResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"certificate_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r CustomCertWebPubsubResource) ModelObject() interface{} {
	return &CustomCertWebPubsubModel{}
}

func (r CustomCertWebPubsubResource) ResourceType() string {
	return "azurerm_web_pubsub_custom_certificate"
}

func (r CustomCertWebPubsubResource) Create() sdk.ResourceFunc {
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

			keyVaultCertificateId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(metadata.ResourceData.Get("custom_certificate_id").(string))
			if err != nil {
				return fmt.Errorf("parsing custom certificate id error: %+v", err)
			}

			keyVaultUri := keyVaultCertificateId.KeyVaultBaseUrl
			keyVaultSecretName := keyVaultCertificateId.Name

			id := webpubsub.NewCustomCertificateID(webPubsubId.SubscriptionId, webPubsubId.ResourceGroupName, webPubsubId.WebPubSubName, customCertWebPubsub.Name)

			locks.ByID(webPubsubId.ID())
			defer locks.UnlockByID(webPubsubId.ID())

			existing, err := client.CustomCertificatesGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			customCertObj := webpubsub.CustomCertificate{
				Properties: webpubsub.CustomCertificateProperties{
					KeyVaultBaseUri:    keyVaultUri,
					KeyVaultSecretName: keyVaultSecretName,
				},
			}
			if keyVaultCertificateId.Version != "" {
				customCertObj.Properties.KeyVaultSecretVersion = utils.String(keyVaultCertificateId.Version)
			}

			if err := client.CustomCertificatesCreateOrUpdateThenPoll(ctx, id, customCertObj); err != nil {
				return fmt.Errorf("creating web pubsub custom certificate: %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CustomCertWebPubsubResource) Read() sdk.ResourceFunc {
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
			certName := resp.Model.Properties.KeyVaultSecretName

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
			certVersion := ""
			if resp.Model.Properties.KeyVaultSecretVersion != nil {
				certVersion = *resp.Model.Properties.KeyVaultSecretVersion
			}
			nestedItem, err := keyVaultParse.NewNestedItemID(vaultBasedUri, keyVaultParse.NestedItemTypeCertificate, certName, certVersion)
			if err != nil {
				return err
			}

			certId := nestedItem.ID()

			state := CustomCertWebPubsubModel{
				Name:               id.CustomCertificateName,
				CustomCertId:       certId,
				WebPubsubId:        webpubsub.NewWebPubSubID(id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName).ID(),
				CertificateVersion: pointer.From(resp.Model.Properties.KeyVaultSecretVersion),
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CustomCertWebPubsubResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.WebPubSubClient.WebPubSub

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

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending:                   []string{"Exists"},
				Target:                    []string{"NotFound"},
				Refresh:                   webPubsubCustomCertificateDeleteRefreshFunc(ctx, client, *id),
				Timeout:                   time.Until(deadline),
				PollInterval:              10 * time.Second,
				ContinuousTargetOccurence: 20,
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be fully deleted: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r CustomCertWebPubsubResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webpubsub.ValidateCustomCertificateID
}

func webPubsubCustomCertificateDeleteRefreshFunc(ctx context.Context, client *webpubsub.WebPubSubClient, id webpubsub.CustomCertificateId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.CustomCertificatesGet(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("checking if %s has been deleted: %+v", id, err)
		}

		return res, "Exists", nil
	}
}
