// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_environment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiCenter.EnvironmentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ApiCenterEnvironmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			if model.Title == "" {
				model.Title = model.Name
			}

			service, err := services.ParseServiceID(model.ServiceId)
			if err != nil {
				return err
			}

			// @favoretti: can't find any workspace creation buttons in the portal, assume everything defaults to "default" for now?
			id := environments.NewEnvironmentID(subscriptionId, service.ResourceGroupName, service.ServiceName, "default", model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing ApiCenter Environment %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			customProperties, err := helpers.ExpandCustomProperties(model.CustomProperties)
			if err != nil {
				return err
			}

			apiCenterEnvironment := environments.Environment{
				Name: pointer.To(model.Name),
				Properties: &environments.EnvironmentProperties{
					CustomProperties: customProperties,
					Description:      new(model.Description),
					Kind:             environments.EnvironmentKind(model.Type),
					Onboarding: &environments.Onboarding{
						DeveloperPortalUri: pointer.To([]string{model.DevPortalUri}),
						Instructions:       pointer.To(model.Instructions),
					},
					Server: &environments.EnvironmentServer{
						Type:                pointer.ToEnum[environments.EnvironmentServerType](model.ServerType),
						ManagementPortalUri: pointer.To([]string{model.MgmtPortalUri}),
					},
					Title: model.Title,
				},
				Type: pointer.To(model.Type),
			}

			if _, err = client.CreateOrUpdate(ctx, id, apiCenterEnvironment); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
