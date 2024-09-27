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
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/sdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StaticWebAppResource struct{}

var _ sdk.ResourceWithUpdate = StaticWebAppResource{}

var _ sdk.ResourceWithCustomizeDiff = StaticWebAppResource{}

type StaticWebAppResourceModel struct {
	Name                string                                     `tfschema:"name"`
	ResourceGroupName   string                                     `tfschema:"resource_group_name"`
	Location            string                                     `tfschema:"location"`
	AppSettings         map[string]string                          `tfschema:"app_settings"`
	BasicAuth           []helpers.BasicAuth                        `tfschema:"basic_auth"`
	ConfigFileChanges   bool                                       `tfschema:"configuration_file_changes_enabled"`
	Identity            []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PreviewEnvironments bool                                       `tfschema:"preview_environments_enabled"`
	PublicNetworkAccess bool                                       `tfschema:"public_network_access_enabled"`
	SkuTier             string                                     `tfschema:"sku_tier"`
	SkuSize             string                                     `tfschema:"sku_size"`
	Tags                map[string]string                          `tfschema:"tags"`

	ApiKey          string `tfschema:"api_key"`
	DefaultHostName string `tfschema:"default_host_name"`
}

func (r StaticWebAppResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StaticWebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"configuration_file_changes_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"preview_environments_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"sku_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(resourceproviders.SkuNameFree),
			ValidateFunc: validation.StringInSlice([]string{
				string(resourceproviders.SkuNameStandard),
				string(resourceproviders.SkuNameFree),
			}, false),
		},

		"sku_size": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(resourceproviders.SkuNameFree),
			ValidateFunc: validation.StringInSlice([]string{
				string(resourceproviders.SkuNameStandard),
				string(resourceproviders.SkuNameFree),
			}, false),
		},

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"basic_auth": helpers.BasicAuthSchema(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": tags.Schema(),
	}
}

func (r StaticWebAppResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r StaticWebAppResource) ModelObject() interface{} {
	return &StaticWebAppResourceModel{}
}

func (r StaticWebAppResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return staticsites.ValidateStaticSiteID
}

func (r StaticWebAppResource) ResourceType() string {
	return "azurerm_static_web_app"
}

func (r StaticWebAppResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			model := StaticWebAppResourceModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := staticsites.NewStaticSiteID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.GetStaticSite(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			envelope := staticsites.StaticSiteARMResource{
				Location:   location.Normalize(model.Location),
				Properties: nil,
				Sku: &staticsites.SkuDescription{
					Name: pointer.To(model.SkuSize),
					Tier: pointer.To(model.SkuTier),
				},
				Tags: pointer.To(model.Tags),
			}

			ident, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding identity for %s: %+v", id, err)
			}
			if ident.Type != identity.TypeNone {
				envelope.Identity = ident
			}

			props := &staticsites.StaticSite{
				AllowConfigFileUpdates:   pointer.To(model.ConfigFileChanges),
				StagingEnvironmentPolicy: pointer.To(staticsites.StagingEnvironmentPolicyEnabled),
				PublicNetworkAccess:      pointer.To(helpers.PublicNetworkAccessEnabled),
			}

			if !model.PreviewEnvironments {
				props.StagingEnvironmentPolicy = pointer.To(staticsites.StagingEnvironmentPolicyDisabled)
			}

			if !model.PublicNetworkAccess {
				props.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessDisabled)
			}

			envelope.Properties = props

			if err := client.CreateOrUpdateStaticSiteThenPoll(ctx, id, envelope); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if len(model.AppSettings) > 0 {
				appSettings := staticsites.StringDictionary{
					Properties: pointer.To(model.AppSettings),
				}

				if _, err = client.CreateOrUpdateStaticSiteAppSettings(ctx, id, appSettings); err != nil {
					return fmt.Errorf("updating app settings for %s: %+v", id, err)
				}
			}

			if len(model.BasicAuth) > 0 {
				sdkHackClient := sdkhacks.NewStaticWebAppClient(client)

				auth := model.BasicAuth[0]

				authProps := staticsites.StaticSiteBasicAuthPropertiesARMResource{
					Properties: &staticsites.StaticSiteBasicAuthPropertiesARMResourceProperties{
						ApplicableEnvironmentsMode: auth.Environments,
						Password:                   pointer.To(auth.Password),
						SecretState:                pointer.To("Password"),
					},
				}

				if _, err := sdkHackClient.CreateOrUpdateBasicAuth(ctx, id, authProps); err != nil {
					return fmt.Errorf("setting basic auth on %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r StaticWebAppResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			id, err := staticsites.ParseStaticSiteID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			staticSite, err := client.GetStaticSite(ctx, *id)
			if err != nil {
				if response.WasNotFound(staticSite.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := StaticWebAppResourceModel{
				Name:              id.StaticSiteName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := staticSite.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				ident, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening identity for %s: %+v", *id, err)
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

				sec, err := client.ListStaticSiteSecrets(ctx, *id)
				if err != nil || sec.Model == nil {
					return fmt.Errorf("retrieving secrets for %s: %+v", *id, err)
				}

				if secProps := sec.Model.Properties; secProps != nil {
					propsMap := pointer.From(secProps)
					apiKey := ""
					apiKey = propsMap["apiKey"]
					state.ApiKey = apiKey
				}
			}

			appSettings, err := client.ListStaticSiteAppSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving app_settings for %s: %+v", *id, err)
			}
			if appSettingsModel := appSettings.Model; appSettingsModel != nil {
				state.AppSettings = pointer.From(appSettingsModel.Properties)
			}

			sdkHackClient := sdkhacks.NewStaticWebAppClient(client)
			auth, err := sdkHackClient.GetBasicAuth(ctx, *id)
			if err != nil && !response.WasNotFound(auth.HttpResponse) { // If basic auth is not configured then this 404's
				return fmt.Errorf("retrieving auth config for %s: %+v", *id, err)
			}
			if !response.WasNotFound(auth.HttpResponse) {
				if authModel := auth.Model; authModel != nil && authModel.Properties != nil && !strings.EqualFold(authModel.Properties.ApplicableEnvironmentsMode, helpers.EnvironmentsTypeSpecifiedEnvironments) {
					state.BasicAuth = []helpers.BasicAuth{
						{
							Password:     metadata.ResourceData.Get("basic_auth.0.password").(string), // Grab this from the config, if we can.
							Environments: authModel.Properties.ApplicableEnvironmentsMode,
						},
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r StaticWebAppResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			id, err := staticsites.ParseStaticSiteID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteStaticSiteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r StaticWebAppResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			config := StaticWebAppResourceModel{}

			if err := metadata.Decode(&config); err != nil {
				return err
			}

			id, err := staticsites.ParseStaticSiteID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.GetStaticSite(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("retrieving %s for update: %+v", *id, err)
			}

			model := *existing.Model

			if metadata.ResourceData.HasChange("identity") {
				ident, err := identity.ExpandSystemAndUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return err
				}
				model.Identity = ident
				// If we're changing to `Free` we must remove the Identity first
				if strings.EqualFold(string(resourceproviders.SkuNameFree), config.SkuTier) {
					if err := client.CreateOrUpdateStaticSiteThenPoll(ctx, *id, model); err != nil {
						return fmt.Errorf("creating %s: %+v", id, err)
					}
					// Once removed, the identity payload needs to be nilled or the API validation for `Free` will reject the request
					model.Identity = nil
				}
			}

			if metadata.ResourceData.HasChanges("sku_tier", "sku_size") {
				model.Sku = &staticsites.SkuDescription{
					Name: pointer.To(config.SkuSize),
					Tier: pointer.To(config.SkuTier),
				}
			}

			if metadata.ResourceData.HasChange("configuration_file_changes_enabled") {
				model.Properties.AllowConfigFileUpdates = pointer.To(config.ConfigFileChanges)
			}

			if metadata.ResourceData.HasChange("preview_environments_enabled") {
				if !config.PreviewEnvironments {
					model.Properties.StagingEnvironmentPolicy = pointer.To(staticsites.StagingEnvironmentPolicyDisabled)
				} else {
					model.Properties.StagingEnvironmentPolicy = pointer.To(staticsites.StagingEnvironmentPolicyEnabled)
				}
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				if !config.PublicNetworkAccess {
					model.Properties.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessDisabled)
				} else {
					model.Properties.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessEnabled)
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateStaticSiteThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("app_settings") {
				appSettings := staticsites.StringDictionary{
					Properties: pointer.To(config.AppSettings),
				}

				if _, err = client.CreateOrUpdateStaticSiteAppSettings(ctx, *id, appSettings); err != nil {
					return fmt.Errorf("updating app settings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("basic_auth") {
				sdkHackClient := sdkhacks.NewStaticWebAppClient(client)
				authProps := staticsites.StaticSiteBasicAuthPropertiesARMResource{}
				if len(config.BasicAuth) > 0 {
					auth := config.BasicAuth[0]
					authProps.Properties = &staticsites.StaticSiteBasicAuthPropertiesARMResourceProperties{
						ApplicableEnvironmentsMode: auth.Environments,
						Password:                   pointer.To(auth.Password),
						SecretState:                pointer.To("Password"),
					}
				} else {
					authProps.Properties = &staticsites.StaticSiteBasicAuthPropertiesARMResourceProperties{
						ApplicableEnvironmentsMode: "SpecifiedEnvironments",
						Password:                   nil,
						SecretState:                nil,
					}
				}

				if _, err := sdkHackClient.CreateOrUpdateBasicAuth(ctx, *id, authProps); err != nil {
					return fmt.Errorf("setting basic auth on %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r StaticWebAppResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			skuTier := rd.Get("sku_tier").(string)
			skuSize := rd.Get("sku_size").(string)

			if strings.EqualFold(skuTier, string(resourceproviders.SkuNameFree)) && strings.EqualFold(skuSize, string(resourceproviders.SkuNameFree)) {
				basicAuth, authOk := rd.GetOk("basic_auth")
				if authOk && len(basicAuth.([]interface{})) > 0 {
					return fmt.Errorf("basic_auth cannot be used with the Free tier of Static Web Apps")
				}
				ident, identOk := rd.GetOk("identity")
				if identOk && len(ident.([]interface{})) > 0 {
					return fmt.Errorf("identities cannot be used with the Free tier of Static Web Apps")
				}
			}

			return nil
		},
	}
}
