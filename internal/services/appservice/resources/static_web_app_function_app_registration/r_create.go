// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app_function_app_registration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient
			appClient := metadata.Client.AppService.WebAppsClient

			model := StaticWebAppFunctionAppRegistrationModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			staticAppId, err := staticsites.ParseStaticSiteID(model.StaticWebAppID)
			if err != nil {
				return err
			}

			functionAppId, err := commonids.ParseAppServiceID(model.FunctionAppID)
			if err != nil {
				return err
			}

			app, err := appClient.Get(ctx, *functionAppId)
			if err != nil {
				return fmt.Errorf("reading specified %s: %+v", *functionAppId, err)
			}

			loc := ""
			if appModel := app.Model; appModel != nil {
				loc = location.Normalize(appModel.Location)
			}

			id := staticsites.NewUserProvidedFunctionAppID(staticAppId.SubscriptionId, staticAppId.ResourceGroupName, staticAppId.StaticSiteName, functionAppId.SiteName)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetUserProvidedFunctionAppForStaticSite(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			backends, err := client.GetLinkedBackends(ctx, *staticAppId)
			if err != nil {
				return fmt.Errorf("checking for existing Static Site backends for %s: %+v", id, err)
			}

			if backendList := backends.Model; backendList != nil {
				if len(*backendList) != 0 {
					return fmt.Errorf("%s already has a backend and cannot have another", id)
				}
			}

			payload := staticsites.StaticSiteUserProvidedFunctionAppARMResource{
				Properties: &staticsites.StaticSiteUserProvidedFunctionAppARMResourceProperties{
					FunctionAppRegion:     pointer.To(loc),
					FunctionAppResourceId: pointer.To(functionAppId.ID()),
				},
			}

			if err := client.RegisterUserProvidedFunctionAppWithStaticSiteCallbackThenPoll(ctx, id, payload, staticsites.DefaultRegisterUserProvidedFunctionAppWithStaticSiteOperationOptions(), metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
