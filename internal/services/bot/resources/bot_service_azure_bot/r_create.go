// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_service_azure_bot

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Bot.BotClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config BotServiceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := commonids.NewBotServiceID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName)
				if err != nil {
					if !response.WasNotFound(existing.Response.Response) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.Response.Response) {
					return tf.ImportAsExistsError(r.ResourceType(), id.ID())
				}
			}

			displayName := config.DisplayName
			if displayName == "" {
				displayName = id.BotServiceName
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

			if _, err := client.Create(ctx, id.ResourceGroupName, id.BotServiceName, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
