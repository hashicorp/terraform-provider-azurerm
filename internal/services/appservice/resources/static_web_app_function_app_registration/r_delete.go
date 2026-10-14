// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app_function_app_registration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			id, err := staticsites.ParseUserProvidedFunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.DetachUserProvidedFunctionAppFromStaticSite(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
