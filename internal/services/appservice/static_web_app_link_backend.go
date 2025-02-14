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
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StaticWebAppLinkBackendResource struct{}

var _ sdk.Resource = StaticWebAppLinkBackendResource{}

type StaticWebAppLinkBackendModel struct {
	StaticWebAppID    string `tfschema:"static_web_app_id"`
	BackendResourceId string `tfschema:"backend_resource_id"`
}

type linkedBackend struct {
	id   string
	name string
}

func (r StaticWebAppLinkBackendResource) ResourceType() string {
	return "azurerm_static_web_app_link_backend"
}

func (r StaticWebAppLinkBackendResource) ModelObject() interface{} {
	return &StaticWebAppLinkBackendModel{}
}

func (r StaticWebAppLinkBackendResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return staticsites.ValidateLinkedBackendID
}

func (r StaticWebAppLinkBackendResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"static_web_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: staticsites.ValidateStaticSiteID,
		},

		"backend_resource_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				apimanagementservice.ValidateServiceID,
				commonids.ValidateAppServiceID,
				containerapps.ValidateContainerAppID,
			),
		},
	}
}

func (r StaticWebAppLinkBackendResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r StaticWebAppLinkBackendResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient
			model := StaticWebAppLinkBackendModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			staticAppId, err := staticsites.ParseStaticSiteID(model.StaticWebAppID)
			if err != nil {
				return err
			}

			loc := ""
			linkedBackend := linkedBackend{}

			if apiManagementId, err := apimanagementservice.ParseServiceID(model.BackendResourceId); err == nil {
				linkedBackend.id = apiManagementId.ID()
				linkedBackend.name = apiManagementId.ServiceName

				backendClient := metadata.Client.ApiManagement.ServiceClient

				apiManagement, err := backendClient.Get(ctx, *apiManagementId)
				if err != nil {
					return fmt.Errorf("reading specified %s: %+v", *apiManagementId, err)
				}

				if appModel := apiManagement.Model; appModel != nil {
					loc = location.Normalize(appModel.Location)
				}

			} else if appServiceId, err := commonids.ParseAppServiceID(model.BackendResourceId); err == nil {
				linkedBackend.id = appServiceId.ID()
				linkedBackend.name = appServiceId.SiteName

				backendClient := metadata.Client.AppService.WebAppsClient

				appService, err := backendClient.Get(ctx, *appServiceId)
				if err != nil {
					return fmt.Errorf("reading specified %s: %+v", *appServiceId, err)
				}

				if appModel := appService.Model; appModel != nil {
					loc = location.Normalize(appModel.Location)
				}

			} else if containerAppId, err := containerapps.ParseContainerAppID(model.BackendResourceId); err == nil {
				linkedBackend.id = containerAppId.ID()
				linkedBackend.name = containerAppId.ContainerAppName

				backendClient := metadata.Client.ContainerApps.ContainerAppClient

				containerApp, err := backendClient.Get(ctx, *containerAppId)
				if err != nil {
					return fmt.Errorf("reading specified %s: %+v", *containerAppId, err)
				}

				if appModel := containerApp.Model; appModel != nil {
					loc = location.Normalize(appModel.Location)
				}

			} else {
				return fmt.Errorf("unsupported backend resource type")
			}

			id := staticsites.NewLinkedBackendID(staticAppId.SubscriptionId, staticAppId.ResourceGroupName, staticAppId.StaticSiteName, linkedBackend.name)

			existing, err := client.GetLinkedBackend(ctx, id)
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

			payload := staticsites.StaticSiteLinkedBackendARMResource{
				Properties: &staticsites.StaticSiteLinkedBackendARMResourceProperties{
					Region:            pointer.To(loc),
					BackendResourceId: pointer.To(linkedBackend.id),
				},
			}

			if err = client.LinkBackendThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r StaticWebAppLinkBackendResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient
			state := StaticWebAppLinkBackendModel{}

			id, err := staticsites.ParseLinkedBackendID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			state.StaticWebAppID = staticsites.NewStaticSiteID(id.SubscriptionId, id.ResourceGroupName, id.StaticSiteName).ID()

			result, err := client.GetLinkedBackend(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := result.Model; model != nil {
				if props := model.Properties; props != nil {
					state.BackendResourceId = pointer.From(props.BackendResourceId)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r StaticWebAppLinkBackendResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			id, err := staticsites.ParseLinkedBackendID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			option := staticsites.UnlinkBackendOperationOptions{}
			if _, err := client.UnlinkBackend(ctx, *id, option); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
