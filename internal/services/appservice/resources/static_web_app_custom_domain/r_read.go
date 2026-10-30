// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app_custom_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			id, err := staticsites.ParseCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Some values are not retrievable from the API so we try and load the config.
			state := StaticWebAppCustomDomainResourceModel{}
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.GetStaticSiteCustomDomain(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state.DomainName = id.CustomDomainName
			state.StaticSiteId = staticsites.NewStaticSiteID(id.SubscriptionId, id.ResourceGroupName, id.StaticSiteName).ID()

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					state.ValidationToken = pointer.From(props.ValidationToken)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
