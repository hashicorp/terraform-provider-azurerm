// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app_function_app_registration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
