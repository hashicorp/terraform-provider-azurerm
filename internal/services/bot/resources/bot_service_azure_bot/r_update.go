// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_service_azure_bot

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Bot.BotClient
			id, err := commonids.ParseBotServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config BotServiceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName)
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

			if _, err := client.Update(ctx, id.ResourceGroupName, id.BotServiceName, existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
