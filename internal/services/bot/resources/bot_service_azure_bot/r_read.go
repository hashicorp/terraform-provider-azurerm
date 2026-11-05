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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Bot.BotClient

			id, err := commonids.ParseBotServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName)
			if err != nil {
				if response.WasNotFound(resp.Response.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := BotServiceModel{
				Name:              id.BotServiceName,
				ResourceGroupName: id.ResourceGroupName,
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

				state.StreamingEndpointEnabled = pointer.From(props.IsStreamingSupported)
			}

			return metadata.Encode(&state)
		},
	}
}
