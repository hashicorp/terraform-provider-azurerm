// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_app_active_slot

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := commonids.ParseWebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			app, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(app.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading active slot for %s: %+v", id.SiteName, err)
			}

			if app.Model == nil || app.Model.Properties == nil || app.Model.Properties.SlotSwapStatus == nil {
				return fmt.Errorf("reading site properties to determine active slot status: %+v", err)
			}

			activeSlot := WebAppActiveSlotModel{
				LastSwap: pointer.From(app.Model.Properties.SlotSwapStatus.TimestampUtc),
			}

			if slotName := app.Model.Properties.SlotSwapStatus.SourceSlotName; slotName != nil {
				activeSlot.SlotID = webapps.NewSlotID(id.SubscriptionId, id.ResourceGroupName, id.SiteName, *slotName).ID()
			}

			// Default value here for imports as this cannot be read from service as it's part of the swap request only and not stored
			overwriteNetworking := true
			if p, ok := metadata.ResourceData.GetOk("overwrite_network_config"); ok {
				overwriteNetworking = p.(bool)
			}
			activeSlot.OverwriteNetworking = overwriteNetworking

			return metadata.Encode(&activeSlot)
		},
	}
}
