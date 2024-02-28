package appservice

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StaticWebAppResource struct{}

var _ sdk.ResourceWithUpdate = StaticWebAppResource{}

var _ sdk.ResourceWithCustomizeDiff = StaticWebAppResource{}

type StaticWebAppResourceModel struct {
	Name              string                                     `tfschema:"name"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	Location          string                                     `tfschema:"location"`
	SkuTier           string                                     `tfschema:"sku_tier"`
	SkuSize           string                                     `tfschema:"sku_size"`
	AppSettings       map[string]string                          `tfschema:"app_settings"`
	Identity          []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Tags              map[string]string                          `tfschema:"tags"`

	ApiKey          string `tfschema:"api_key"`
	DefaultHostName string `tfschema:"default_host_name"`
}

func (r StaticWebAppResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StaticWebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

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

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": tags.Schema(),
	}
}

func (r StaticWebAppResource) Attributes() map[string]*schema.Schema {
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
					return fmt.Errorf("failed checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			ident, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)

			envelope := staticsites.StaticSiteARMResource{
				Identity:   ident,
				Location:   location.Normalize(model.Location),
				Properties: nil,
				Sku: &staticsites.SkuDescription{
					Name: pointer.To(model.SkuSize),
					Tier: pointer.To(model.SkuTier),
				},
				Tags: pointer.To(model.Tags),
			}

			props := &staticsites.StaticSite{ // TODO - Can we support this now?
				AllowConfigFileUpdates:      nil,
				Branch:                      nil,
				BuildProperties:             nil,
				ContentDistributionEndpoint: nil,
				CustomDomains:               nil,
				DatabaseConnections:         nil,
				DefaultHostname:             nil,
				EnterpriseGradeCdnStatus:    nil,
				KeyVaultReferenceIdentity:   nil,
				LinkedBackends:              nil,
				PrivateEndpointConnections:  nil,
				Provider:                    nil,
				PublicNetworkAccess:         nil,
				RepositoryToken:             nil,
				RepositoryUrl:               nil,
				StagingEnvironmentPolicy:    nil,
				TemplateProperties:          nil,
				UserProvidedFunctionApps:    nil,
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
					state.DefaultHostName = pointer.From(props.DefaultHostname)
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

func (r StaticWebAppResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return staticsites.ValidateStaticSiteID
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
			}

			if metadata.ResourceData.HasChanges("sku_tier", "sku_size") {
				model.Sku = &staticsites.SkuDescription{
					Name: pointer.To(config.SkuSize),
					Tier: pointer.To(config.SkuTier),
				}
			}

			props := &staticsites.StaticSite{ // TODO - Can we support this now?
				AllowConfigFileUpdates:      nil,
				Branch:                      nil,
				BuildProperties:             nil,
				ContentDistributionEndpoint: nil,
				CustomDomains:               nil,
				DatabaseConnections:         nil,
				DefaultHostname:             nil,
				EnterpriseGradeCdnStatus:    nil,
				KeyVaultReferenceIdentity:   nil,
				LinkedBackends:              nil,
				PrivateEndpointConnections:  nil,
				Provider:                    nil,
				PublicNetworkAccess:         nil,
				RepositoryToken:             nil,
				RepositoryUrl:               nil,
				StagingEnvironmentPolicy:    nil,
				TemplateProperties:          nil,
				UserProvidedFunctionApps:    nil,
			}

			model.Properties = props

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

			return nil
		},
	}
}

func (r StaticWebAppResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Free SKU cannot be used with Identity
			return nil
		},

		Timeout: 5 * time.Minute,
	}
}
