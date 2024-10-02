// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/sdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StaticWebAppDataSource struct{}

var _ sdk.DataSource = StaticWebAppDataSource{}

type StaticWebAppDataSourceModel struct {
	Name                string                                     `tfschema:"name"`
	ResourceGroupName   string                                     `tfschema:"resource_group_name"`
	Location            string                                     `tfschema:"location"`
	ApiKey              string                                     `tfschema:"api_key"`
	AppSettings         map[string]string                          `tfschema:"app_settings"`
	BasicAuth           []helpers.BasicAuthComputed                `tfschema:"basic_auth"`
	ConfigFileChanges   bool                                       `tfschema:"configuration_file_changes_enabled"`
	DefaultHostName     string                                     `tfschema:"default_host_name"`
	Identity            []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PreviewEnvironments bool                                       `tfschema:"preview_environments_enabled"`
	PublicNetworkAccess bool                                       `tfschema:"public_network_access_enabled"`
	SkuTier             string                                     `tfschema:"sku_tier"`
	SkuSize             string                                     `tfschema:"sku_size"`
	Tags                map[string]string                          `tfschema:"tags"`
}

func (s StaticWebAppDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StaticWebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (s StaticWebAppDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"configuration_file_changes_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"preview_environments_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"sku_tier": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"sku_size": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"basic_auth": helpers.BasicAuthSchemaComputed(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"api_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (s StaticWebAppDataSource) ModelObject() interface{} {
	return &StaticWebAppDataSourceModel{}
}

func (s StaticWebAppDataSource) ResourceType() string {
	return "azurerm_static_web_app"
}

func (s StaticWebAppDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state StaticWebAppDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id := staticsites.NewStaticSiteID(subscriptionId, state.ResourceGroupName, state.Name)

			staticSite, err := client.GetStaticSite(ctx, id)
			if err != nil {
				if response.WasNotFound(staticSite.HttpResponse) {
					return fmt.Errorf("%s not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := staticSite.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				ident, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening identity for %s: %+v", id, err)
				}
				state.Identity = pointer.From(ident)

				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					state.ConfigFileChanges = pointer.From(props.AllowConfigFileUpdates)
					state.DefaultHostName = pointer.From(props.DefaultHostname)
					state.PreviewEnvironments = pointer.From(props.StagingEnvironmentPolicy) == staticsites.StagingEnvironmentPolicyEnabled
					state.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)
				}

				if sku := model.Sku; sku != nil {
					state.SkuSize = pointer.From(sku.Name)
					state.SkuTier = pointer.From(sku.Tier)
				}

				sec, err := client.ListStaticSiteSecrets(ctx, id)
				if err != nil || sec.Model == nil {
					return fmt.Errorf("retrieving secrets for %s: %+v", id, err)
				}

				if secProps := sec.Model.Properties; secProps != nil {
					propsMap := pointer.From(secProps)
					apiKey := ""
					apiKey = propsMap["apiKey"]
					state.ApiKey = apiKey
				}
			}

			appSettings, err := client.ListStaticSiteAppSettings(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving app_settings for %s: %+v", id, err)
			}
			if appSettingsModel := appSettings.Model; appSettingsModel != nil {
				state.AppSettings = pointer.From(appSettingsModel.Properties)
			}

			sdkHackClient := sdkhacks.NewStaticWebAppClient(client)
			auth, err := sdkHackClient.GetBasicAuth(ctx, id)
			if err != nil && !response.WasNotFound(auth.HttpResponse) { // If basic auth is not configured then this 404's
				return fmt.Errorf("retrieving auth config for %s: %+v", id, err)
			}
			if !response.WasNotFound(auth.HttpResponse) {
				if authModel := auth.Model; authModel != nil && authModel.Properties != nil && !strings.EqualFold(authModel.Properties.ApplicableEnvironmentsMode, helpers.EnvironmentsTypeSpecifiedEnvironments) {
					state.BasicAuth = []helpers.BasicAuthComputed{
						{
							Environments: authModel.Properties.ApplicableEnvironmentsMode,
						},
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
