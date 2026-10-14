// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_app_active_slot

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/custompollers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var activeSlot WebAppActiveSlotModel

			if err := metadata.Decode(&activeSlot); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			id, err := webapps.ParseSlotID(activeSlot.SlotID)
			appId := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName)
			if err != nil {
				return fmt.Errorf("parsing App ID: %+v", err)
			}

			app, err := client.Get(ctx, appId)
			if err != nil {
				if response.WasNotFound(app.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			csmSlotEntity := webapps.CsmSlotEntity{
				TargetSlot:   id.SlotName,
				PreserveVnet: activeSlot.OverwriteNetworking,
			}

			locks.ByID(appId.ID())
			defer locks.UnlockByID(appId.ID())

			if _, err := client.SwapSlotWithProduction(ctx, appId, csmSlotEntity); err != nil {
				return fmt.Errorf("making %s the active slot: %+v", id.SlotName, err)
			}

			pollerType := custompollers.NewAppServiceActiveSlotPoller(client, appId, *id)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return err
			}

			metadata.SetID(appId)

			return nil
		},
	}
}
