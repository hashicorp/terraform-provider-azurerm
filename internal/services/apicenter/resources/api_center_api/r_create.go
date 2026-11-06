package api_center_api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiCenter.ApisClient

			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ApiCenterApiModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			if config.Title == "" {
				config.Title = config.Name
			}

			service, err := services.ParseServiceID(config.ServiceId)
			if err != nil {
				return err
			}

			id := apis.NewApiID(subscriptionId, service.ResourceGroupName, service.ServiceName, "default", config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := apis.Api{
				Name: pointer.To(config.Name),
				Properties: &apis.ApiProperties{
					Contacts:              expandContactList(config.Contacts),
					Description:           pointer.To(config.Description),
					ExternalDocumentation: expandExternalDocumentationList(config.ExternalDocumentation),
					Kind:                  apis.ApiKind(config.Kind),
					License:               expandLicense(config.License),
					Summary:               pointer.To(config.Summary),
					Title:                 config.Title,
				},
			}

			if config.TermsOfServiceUrl != "" {
				payload.Properties.TermsOfService = &apis.TermsOfService{
					Url: config.TermsOfServiceUrl,
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
