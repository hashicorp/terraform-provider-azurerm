// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_source_control_slot

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := webapps.ParseSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			sitePatch := webapps.SitePatchResource{
				Properties: &webapps.SitePatchResourceProperties{
					SiteConfig: &webapps.SiteConfig{
						ScmType: pointer.To(webapps.ScmTypeNone),
					},
				},
			}
			if _, err := client.UpdateSlot(ctx, *id, sitePatch); err != nil {
				return fmt.Errorf("setting App Source Control Type for %s: %v", id, err)
			}

			if _, err := client.DeleteSourceControlSlot(ctx, *id, webapps.DefaultDeleteSourceControlSlotOperationOptions()); err != nil {
				return fmt.Errorf("deleting Source Control for %s: %v", id, err)
			}

			return nil
		},
	}
}
