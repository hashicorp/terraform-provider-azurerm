// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

var (
	_ sdk.ResourceWithUpdate         = AzureBotServiceResource{}
	_ sdk.ResourceWithCustomImporter = AzureBotServiceResource{}
)

type AzureBotServiceResource struct{}

type BotServiceModel struct {
	Name                              string                 `tfschema:"name"`
	ResourceGroupName                 string                 `tfschema:"resource_group_name"`
	Location                          string                 `tfschema:"location"`
	Sku                               string                 `tfschema:"sku"`
	MicrosoftAppId                    string                 `tfschema:"microsoft_app_id"`
	DisplayName                       string                 `tfschema:"display_name"`
	Endpoint                          string                 `tfschema:"endpoint"`
	DeveloperAppInsightsKey           string                 `tfschema:"developer_app_insights_key"`
	DeveloperAppInsightsApiKey        string                 `tfschema:"developer_app_insights_api_key"`
	DeveloperAppInsightsApplicationId string                 `tfschema:"developer_app_insights_application_id"`
	CmkKeyVaultKeyUrl                 string                 `tfschema:"cmk_key_vault_key_url"`
	MicrosoftAppMsiId                 string                 `tfschema:"microsoft_app_msi_id"`
	MicrosoftAppTenantId              string                 `tfschema:"microsoft_app_tenant_id"`
	MicrosoftAppType                  string                 `tfschema:"microsoft_app_type"`
	LocalAuthenticationEnabled        bool                   `tfschema:"local_authentication_enabled"`
	LuisAppIds                        []string               `tfschema:"luis_app_ids"`
	LuisKey                           string                 `tfschema:"luis_key"`
	PublicNetworkAccessEnabled        bool                   `tfschema:"public_network_access_enabled"`
	StreamingEndpointEnabled          bool                   `tfschema:"streaming_endpoint_enabled"`
	IconUrl                           string                 `tfschema:"icon_url"`
	Tags                              map[string]interface{} `tfschema:"tags"`
}

func (r AzureBotServiceResource) ModelObject() interface{} {
	return &BotServiceModel{}
}

func (r AzureBotServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.BotServiceID
}

func (r AzureBotServiceResource) ResourceType() string {
	return "azurerm_bot_service_azure_bot"
}

func (r AzureBotServiceResource) Arguments() map[string]*pluginsdk.Schema {
	output := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(botservice.SkuNameF0),
				string(botservice.SkuNameS1),
			}, false),
		},

		"microsoft_app_id": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"endpoint": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"developer_app_insights_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"developer_app_insights_api_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"developer_app_insights_application_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cmk_key_vault_key_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
		},

		"microsoft_app_msi_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"microsoft_app_tenant_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"microsoft_app_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(botservice.MsaAppTypeMultiTenant),
				string(botservice.MsaAppTypeSingleTenant),
				string(botservice.MsaAppTypeUserAssignedMSI),
			}, false),
		},

		"local_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"luis_app_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"luis_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"streaming_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"icon_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "https://docs.botframework.com/static/devportal/client/images/bot-framework-default.png",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}

	if !features.FivePointOh() {
		output["cmk_key_vault_key_url"].ValidateFunc = keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeAny)

		output["microsoft_app_type"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because Azure sets a value for this if omitted
			Computed: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(botservice.MsaAppTypeMultiTenant),
				string(botservice.MsaAppTypeSingleTenant),
				string(botservice.MsaAppTypeUserAssignedMSI),
			}, false),
		}
	}

	return output
}

func (r AzureBotServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AzureBotServiceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Bot.BotClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config BotServiceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := parse.NewBotServiceID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			displayName := config.DisplayName
			if displayName == "" {
				displayName = id.Name
			}

			publicNetworkEnabled := botservice.PublicNetworkAccessEnabled
			if !config.PublicNetworkAccessEnabled {
				publicNetworkEnabled = botservice.PublicNetworkAccessDisabled
			}

			props := botservice.Bot{
				Location: pointer.To(config.Location),
				Sku: &botservice.Sku{
					Name: botservice.SkuName(config.Sku),
				},
				Kind: botservice.KindAzurebot,
				Properties: &botservice.BotProperties{
					DisplayName:                       pointer.To(displayName),
					Endpoint:                          pointer.To(config.Endpoint),
					MsaAppID:                          pointer.To(config.MicrosoftAppId),
					DeveloperAppInsightKey:            pointer.To(config.DeveloperAppInsightsKey),
					DeveloperAppInsightsAPIKey:        pointer.To(config.DeveloperAppInsightsApiKey),
					DeveloperAppInsightsApplicationID: pointer.To(config.DeveloperAppInsightsApplicationId),
					DisableLocalAuth:                  pointer.To(!config.LocalAuthenticationEnabled),
					IsCmekEnabled:                     pointer.To(false),
					CmekKeyVaultURL:                   pointer.To(config.CmkKeyVaultKeyUrl),
					LuisAppIds:                        &config.LuisAppIds,
					LuisKey:                           pointer.To(config.LuisKey),
					PublicNetworkAccess:               publicNetworkEnabled,
					IsStreamingSupported:              pointer.To(config.StreamingEndpointEnabled),
					IconURL:                           pointer.To(config.IconUrl),
				},
				Tags: tags.Expand(config.Tags),
			}

			if config.CmkKeyVaultKeyUrl != "" {
				props.Properties.IsCmekEnabled = pointer.To(true)
			}

			if config.MicrosoftAppType != "" {
				props.Properties.MsaAppType = botservice.MsaAppType(config.MicrosoftAppType)
			}

			if config.MicrosoftAppTenantId != "" {
				props.Properties.MsaAppTenantID = pointer.To(config.MicrosoftAppTenantId)
			}

			if config.MicrosoftAppMsiId != "" {
				props.Properties.MsaAppMSIResourceID = pointer.To(config.MicrosoftAppMsiId)
			}

			if _, err := client.Create(ctx, id.ResourceGroup, id.Name, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AzureBotServiceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Bot.BotClient

			id, err := parse.BotServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := BotServiceModel{
				Name:              id.Name,
				ResourceGroupName: id.ResourceGroup,
				Location:          location.NormalizeNilable(resp.Location),
			}

			if v := resp.Sku; v != nil {
				state.Sku = string(v.Name)
			}

			state.Tags = tags.Flatten(resp.Tags)

			// The API doesn't return this property, so we need to preserve the value from config/state
			if apiKey, ok := metadata.ResourceData.GetOk("developer_app_insights_api_key"); ok && apiKey.(string) != "" {
				state.DeveloperAppInsightsApiKey = apiKey.(string)
			}

			if props := resp.Properties; props != nil {
				state.MicrosoftAppId = pointer.From(props.MsaAppID)
				state.DisplayName = pointer.From(props.DisplayName)
				state.Endpoint = pointer.From(props.Endpoint)
				state.DeveloperAppInsightsKey = pointer.From(props.DeveloperAppInsightKey)
				state.DeveloperAppInsightsApplicationId = pointer.From(props.DeveloperAppInsightsApplicationID)
				state.MicrosoftAppType = string(props.MsaAppType)
				state.MicrosoftAppTenantId = pointer.From(props.MsaAppTenantID)
				state.MicrosoftAppMsiId = pointer.From(props.MsaAppMSIResourceID)
				state.CmkKeyVaultKeyUrl = pointer.From(props.CmekKeyVaultURL)
				state.IconUrl = pointer.From(props.IconURL)

				localAuthEnabled := true
				if v := props.DisableLocalAuth; v != nil {
					localAuthEnabled = !*v
				}
				state.LocalAuthenticationEnabled = localAuthEnabled

				publicNetworkAccessEnabled := true
				if v := props.PublicNetworkAccess; v != botservice.PublicNetworkAccessEnabled {
					publicNetworkAccessEnabled = false
				}
				state.PublicNetworkAccessEnabled = publicNetworkAccessEnabled

				if v := props.LuisAppIds; v != nil {
					state.LuisAppIds = *v
				}

				streamingEndpointEnabled := false
				if v := props.IsStreamingSupported; v != nil {
					streamingEndpointEnabled = *v
				}
				state.StreamingEndpointEnabled = streamingEndpointEnabled
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AzureBotServiceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Bot.BotClient
			id, err := parse.BotServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r AzureBotServiceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Bot.BotClient
			id, err := parse.BotServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config BotServiceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("display_name") {
				existing.Properties.DisplayName = pointer.To(config.DisplayName)
			}

			if metadata.ResourceData.HasChange("endpoint") {
				existing.Properties.Endpoint = pointer.To(config.Endpoint)
			}

			if metadata.ResourceData.HasChange("developer_app_insights_key") {
				existing.Properties.DeveloperAppInsightKey = pointer.To(config.DeveloperAppInsightsKey)
			}

			if metadata.ResourceData.HasChange("developer_app_insights_api_key") {
				existing.Properties.DeveloperAppInsightsAPIKey = pointer.To(config.DeveloperAppInsightsApiKey)
			}

			if metadata.ResourceData.HasChange("developer_app_insights_application_id") {
				existing.Properties.DeveloperAppInsightsApplicationID = pointer.To(config.DeveloperAppInsightsApplicationId)
			}

			if metadata.ResourceData.HasChange("local_authentication_enabled") {
				existing.Properties.DisableLocalAuth = pointer.To(!config.LocalAuthenticationEnabled)
			}

			if metadata.ResourceData.HasChange("luis_app_ids") {
				existing.Properties.LuisAppIds = &config.LuisAppIds
			}

			if metadata.ResourceData.HasChange("luis_key") {
				existing.Properties.LuisKey = pointer.To(config.LuisKey)
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				if config.PublicNetworkAccessEnabled {
					existing.Properties.PublicNetworkAccess = botservice.PublicNetworkAccessEnabled
				} else {
					existing.Properties.PublicNetworkAccess = botservice.PublicNetworkAccessDisabled
				}
			}

			if metadata.ResourceData.HasChange("streaming_endpoint_enabled") {
				existing.Properties.IsStreamingSupported = pointer.To(config.StreamingEndpointEnabled)
			}

			if metadata.ResourceData.HasChange("icon_url") {
				existing.Properties.IconURL = pointer.To(config.IconUrl)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.Expand(config.Tags)
			}

			if _, err := client.Update(ctx, id.ResourceGroup, id.Name, existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r AzureBotServiceResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Bot.BotClient

		id, err := parse.BotServiceID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if actualKind := string(resp.Kind); actualKind != string(botservice.KindAzurebot) {
			return fmt.Errorf("bot has mismatched type, expected: %q, got %q", string(botservice.KindAzurebot), actualKind)
		}

		return nil
	}
}
