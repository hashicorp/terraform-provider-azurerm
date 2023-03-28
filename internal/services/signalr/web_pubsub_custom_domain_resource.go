package signalr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CustomDomainWebPubsubModel struct {
	Name                         string `tfschema:"name"`
	WebPubsubId                  string `tfschema:"web_pubsub_id"`
	DomainName                   string `tfschema:"domain_name"`
	WebPubsubCustomCertificateId string `tfschema:"web_pubsub_custom_certificate_id"`
}

type CustomDomainWebPubsubResource struct{}

var _ sdk.Resource = CustomDomainWebPubsubResource{}

func (r CustomDomainWebPubsubResource) Arguments() map[string]*pluginsdk.Schema {
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

		"domain_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"web_pubsub_custom_certificate_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: webpubsub.ValidateCustomCertificateID,
		},
	}
}

func (r CustomDomainWebPubsubResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CustomDomainWebPubsubResource) ModelObject() interface{} {
	return &CustomDomainWebPubsubModel{}
}

func (r CustomDomainWebPubsubResource) ResourceType() string {
	return "azurerm_web_pubsub_custom_domain"
}

func (r CustomDomainWebPubsubResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var customDomainWebPubsubModel CustomDomainWebPubsubModel
			if err := metadata.Decode(&customDomainWebPubsubModel); err != nil {
				return err
			}
			client := metadata.Client.SignalR.WebPubSubClient.WebPubSub

			webPubsubId, err := webpubsub.ParseWebPubSubIDInsensitively(metadata.ResourceData.Get("web_pubsub_id").(string))
			if err != nil {
				return fmt.Errorf("parsing web pubsub id error: %+v", err)
			}

			id := webpubsub.NewCustomDomainID(webPubsubId.SubscriptionId, webPubsubId.ResourceGroupName, webPubsubId.WebPubSubName, metadata.ResourceData.Get("name").(string))
			if _, err := webpubsub.ParseCustomCertificateIDInsensitively(customDomainWebPubsubModel.WebPubsubCustomCertificateId); err != nil {
				return fmt.Errorf("parsing custom certificate for %s: %+v", id, err)
			}

			customDomainObj := webpubsub.CustomDomain{
				Properties: webpubsub.CustomDomainProperties{
					DomainName: customDomainWebPubsubModel.DomainName,
					CustomCertificate: webpubsub.ResourceReference{
						Id: utils.String(customDomainWebPubsubModel.WebPubsubCustomCertificateId),
					},
				},
			}
			if err := client.CustomDomainsCreateOrUpdateThenPoll(ctx, id, customDomainObj); err != nil {
				return fmt.Errorf("creating web pubsub custom domain: %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CustomDomainWebPubsubResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.WebPubSubClient.WebPubSub
			id, err := webpubsub.ParseCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.CustomDomainsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Web Pubsub custom domain %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: got nil model", *id)
			}

			webPubsubId := webpubsub.NewWebPubSubID(id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName).ID()

			state := CustomDomainWebPubsubModel{
				Name:                         id.CustomDomainName,
				WebPubsubId:                  webPubsubId,
				WebPubsubCustomCertificateId: utils.NormalizeNilableString(resp.Model.Properties.CustomCertificate.Id),
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CustomDomainWebPubsubResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.WebPubSubClient.WebPubSub

			id, err := webpubsub.ParseCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.CustomDomainsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CustomDomainWebPubsubResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webpubsub.ValidateCustomDomainID
}
