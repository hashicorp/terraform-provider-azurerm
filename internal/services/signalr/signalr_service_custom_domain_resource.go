package signalr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2023-02-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CustomDomainSignalrServiceModel struct {
	Name                       string `tfschema:"name"`
	SignalRServiceId           string `tfschema:"signalr_service_id"`
	DomainName                 string `tfschema:"domain_name"`
	SignalrCustomCertificateId string `tfschema:"signalr_custom_certificate_id"`
}

type CustomDomainSignalrServiceResource struct{}

var _ sdk.Resource = CustomDomainSignalrServiceResource{}

func (r CustomDomainSignalrServiceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"signalr_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: signalr.ValidateSignalRID,
		},

		"domain_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"signalr_custom_certificate_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: signalr.ValidateCustomCertificateID,
		},
	}
}

func (r CustomDomainSignalrServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CustomDomainSignalrServiceResource) ModelObject() interface{} {
	return &CustomDomainSignalrServiceModel{}
}

func (r CustomDomainSignalrServiceResource) ResourceType() string {
	return "azurerm_signalr_custom_domain"
}

func (r CustomDomainSignalrServiceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var customDomainSignalrServiceModel CustomDomainSignalrServiceModel
			if err := metadata.Decode(&customDomainSignalrServiceModel); err != nil {
				return err
			}
			client := metadata.Client.SignalR.SignalRClient

			signalRServiceId, err := signalr.ParseSignalRID(metadata.ResourceData.Get("signalr_service_id").(string))
			if err != nil {
				return fmt.Errorf("parsing signalr service id error: %+v", err)
			}

			id := signalr.NewCustomDomainID(signalRServiceId.SubscriptionId, signalRServiceId.ResourceGroupName, signalRServiceId.SignalRName, metadata.ResourceData.Get("name").(string))
			if _, err := signalr.ParseCustomCertificateIDInsensitively(customDomainSignalrServiceModel.SignalrCustomCertificateId); err != nil {
				return fmt.Errorf("parsing custom certificate for %s: %+v", id, err)
			}

			customDomainObj := signalr.CustomDomain{
				Properties: signalr.CustomDomainProperties{
					DomainName: customDomainSignalrServiceModel.DomainName,
					CustomCertificate: signalr.ResourceReference{
						Id: utils.String(customDomainSignalrServiceModel.SignalrCustomCertificateId),
					},
				},
			}
			if err := client.CustomDomainsCreateOrUpdateThenPoll(ctx, id, customDomainObj); err != nil {
				return fmt.Errorf("creating signalR custom domain: %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CustomDomainSignalrServiceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.SignalRClient
			id, err := signalr.ParseCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.CustomDomainsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading signalR custom domain %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: got nil model", *id)
			}

			signalrServiceId := signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName).ID()

			state := CustomDomainSignalrServiceModel{
				Name:                       id.CustomDomainName,
				SignalRServiceId:           signalrServiceId,
				SignalrCustomCertificateId: utils.NormalizeNilableString(resp.Model.Properties.CustomCertificate.Id),
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CustomDomainSignalrServiceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.SignalRClient

			id, err := signalr.ParseCustomDomainID(metadata.ResourceData.Id())
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

func (r CustomDomainSignalrServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return signalr.ValidateCustomDomainID
}
