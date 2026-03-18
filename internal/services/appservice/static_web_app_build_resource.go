// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StaticWebAppBuildResource struct{}

var _ sdk.ResourceWithUpdate = StaticWebAppBuildResource{}

type StaticWebAppBuildResourceModel struct {
	Name                 string            `tfschema:"name"`
	ResourceGroupName    string            `tfschema:"resource_group_name"`
	EnvironmentVariables map[string]string `tfschema:"environment_variables"`
	Build                string            `tfschema:"build"`
}

func (r StaticWebAppBuildResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StaticWebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"environment_variables": {
			Type:     pluginsdk.TypeMap,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"build": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r StaticWebAppBuildResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StaticWebAppBuildResource) ModelObject() interface{} {
	return &StaticWebAppBuildResourceModel{}
}

func (r StaticWebAppBuildResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return staticsites.ValidateStaticSiteID
}

func (r StaticWebAppBuildResource) ResourceType() string {
	return "azurerm_static_web_app_build"
}

func (r StaticWebAppBuildResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			model := StaticWebAppBuildResourceModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := staticsites.NewBuildID(subscriptionId, model.ResourceGroupName, model.Name, model.Build)

			metadata.SetID(id)

			appSettings := staticsites.StringDictionary{
				Properties: pointer.To(model.EnvironmentVariables),
			}

			if _, err := client.CreateOrUpdateStaticSiteBuildAppSettings(ctx, id, appSettings); err != nil {
				return fmt.Errorf("updating app settings for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r StaticWebAppBuildResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			id, err := staticsites.ParseBuildID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			appSettings, err := client.ListStaticSiteBuildAppSettings(ctx, *id)
			if err != nil || appSettings.Model == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := StaticWebAppBuildResourceModel{
				Name:              id.StaticSiteName,
				ResourceGroupName: id.ResourceGroupName,
				Build:             id.BuildName,
			}

			if appSettingsModel := appSettings.Model; appSettingsModel != nil {
				state.EnvironmentVariables = pointer.From(appSettingsModel.Properties)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r StaticWebAppBuildResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			config := StaticWebAppBuildResourceModel{}

			if err := metadata.Decode(&config); err != nil {
				return err
			}

			id, err := staticsites.ParseBuildID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			emptyAppSettings := staticsites.StringDictionary{
				Properties: pointer.To(map[string]string{}),
			}

			if _, err := client.CreateOrUpdateStaticSiteBuildAppSettings(ctx, *id, emptyAppSettings); err != nil {
				return fmt.Errorf("deleting app settings for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r StaticWebAppBuildResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			config := StaticWebAppBuildResourceModel{}

			if err := metadata.Decode(&config); err != nil {
				return err
			}

			id, err := staticsites.ParseBuildID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			appSettings := staticsites.StringDictionary{
				Properties: pointer.To(config.EnvironmentVariables),
			}

			if _, err := client.CreateOrUpdateStaticSiteBuildAppSettings(ctx, *id, appSettings); err != nil {
				return fmt.Errorf("updating app settings for %s: %+v", id, err)
			}

			return nil
		},
	}
}
