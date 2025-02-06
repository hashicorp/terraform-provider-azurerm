// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StaticWebAppFunctionAppRegistrationResource struct{}

var _ sdk.Resource = StaticWebAppFunctionAppRegistrationResource{}

type StaticWebAppFunctionAppRegistrationModel struct {
	StaticWebAppID string `tfschema:"static_web_app_id"`
	FunctionAppID  string `tfschema:"function_app_id"`
}

func (r StaticWebAppFunctionAppRegistrationResource) ResourceType() string {
	return "azurerm_static_web_app_function_app_registration"
}

func (r StaticWebAppFunctionAppRegistrationResource) ModelObject() interface{} {
	return &StaticWebAppFunctionAppRegistrationModel{}
}

func (r StaticWebAppFunctionAppRegistrationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return staticsites.ValidateUserProvidedFunctionAppID
}

func (r StaticWebAppFunctionAppRegistrationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"static_web_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: staticsites.ValidateStaticSiteID,
		},

		"function_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateFunctionAppID,
		},
	}
}

func (r StaticWebAppFunctionAppRegistrationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StaticWebAppFunctionAppRegistrationResource) Create() sdk.ResourceFunc {
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

			existing, err := client.GetUserProvidedFunctionAppForStaticSite(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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

			if err = client.RegisterUserProvidedFunctionAppWithStaticSiteThenPoll(ctx, id, payload, staticsites.DefaultRegisterUserProvidedFunctionAppWithStaticSiteOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r StaticWebAppFunctionAppRegistrationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			state := StaticWebAppFunctionAppRegistrationModel{}

			id, err := staticsites.ParseUserProvidedFunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			state.StaticWebAppID = staticsites.NewStaticSiteID(id.SubscriptionId, id.ResourceGroupName, id.StaticSiteName).ID()

			result, err := client.GetUserProvidedFunctionAppForStaticSite(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := result.Model; model != nil {
				if props := model.Properties; props != nil {
					functionAppId, err := commonids.ParseAppServiceIDInsensitively(pointer.From(props.FunctionAppResourceId))
					if err != nil {
						return err
					}

					state.FunctionAppID = functionAppId.ID()
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r StaticWebAppFunctionAppRegistrationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			id, err := staticsites.ParseUserProvidedFunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.DetachUserProvidedFunctionAppFromStaticSite(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
