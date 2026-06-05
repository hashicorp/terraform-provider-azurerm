// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package signalr

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name web_pubsub_custom_domain -service-package-name signalr -properties "name" -compare-values "subscription_id:web_pubsub_id,resource_group_name:web_pubsub_id,web_pubsub_name:web_pubsub_id"

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CustomDomainWebPubsubModel struct {
	Name                         string `tfschema:"name"`
	WebPubsubId                  string `tfschema:"web_pubsub_id"`
	DomainName                   string `tfschema:"domain_name"`
	WebPubsubCustomCertificateId string `tfschema:"web_pubsub_custom_certificate_id"`
}

const webPubsubCustomDomainResourceType = "azurerm_web_pubsub_custom_domain"

type CustomDomainWebPubsubResource struct{}

var (
	_ sdk.Resource             = CustomDomainWebPubsubResource{}
	_ sdk.ResourceWithIdentity = CustomDomainWebPubsubResource{}
)

func (r CustomDomainWebPubsubResource) Identity() resourceids.ResourceId {
	return &webpubsub.CustomDomainId{}
}

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
	return webPubsubCustomDomainResourceType
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

			locks.ByID(webPubsubId.ID())
			defer locks.UnlockByID(webPubsubId.ID())

			if _, err := webpubsub.ParseCustomCertificateIDInsensitively(customDomainWebPubsubModel.WebPubsubCustomCertificateId); err != nil {
				return fmt.Errorf("parsing custom certificate for %s: %+v", id, err)
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
			if err := client.CustomDomainsCreateOrUpdateCallbackThenPoll(ctx, id, customDomainObj, metadata.SetIDCallback(&id)); err != nil {
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
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
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

			webPubsubId := webpubsub.NewWebPubSubID(id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName)

			locks.ByID(webPubsubId.ID())
			defer locks.UnlockByID(webPubsubId.ID())

			if err := client.CustomDomainsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending:                   []string{"Exists"},
				Target:                    []string{"NotFound"},
				Refresh:                   webPubsubCustomDomainDeleteRefreshFunc(ctx, client, *id),
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

func (r CustomDomainWebPubsubResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webpubsub.ValidateCustomDomainID
}

func webPubsubCustomDomainDeleteRefreshFunc(ctx context.Context, client *webpubsub.WebPubSubClient, id webpubsub.CustomDomainId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.CustomDomainsGet(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("checking if %s has been deleted: %+v", id, err)
		}

		return res, "Exists", nil
	}
}

func (r CustomDomainWebPubsubResource) flatten(metadata sdk.ResourceMetaData, id *webpubsub.CustomDomainId, model *webpubsub.CustomDomain) error {
	state := CustomDomainWebPubsubModel{
		Name:        id.CustomDomainName,
		WebPubsubId: webpubsub.NewWebPubSubID(id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName).ID(),
	}

	if model != nil {
		props := model.Properties
		state.DomainName = props.DomainName

		if props.CustomCertificate.Id != nil {
			webPubsubCustomCertificateID, err := webpubsub.ParseCustomCertificateIDInsensitively(*props.CustomCertificate.Id)
			if err != nil {
				return fmt.Errorf("parsing web pubsub custom cert id for %s: %+v", *id, err)
			}
			state.WebPubsubCustomCertificateId = webPubsubCustomCertificateID.ID()
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
