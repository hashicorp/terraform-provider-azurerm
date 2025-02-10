// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WebAppActiveSlotResource struct{}

type WebAppActiveSlotModel struct {
	SlotID              string `tfschema:"slot_id"`
	OverwriteNetworking bool   `tfschema:"overwrite_network_config"` // Note: This setting controls the ambiguously named `PreserveVnet`
	LastSwap            string `tfschema:"last_successful_swap"`
}

var _ sdk.ResourceWithUpdate = WebAppActiveSlotResource{}

func (r WebAppActiveSlotResource) ModelObject() interface{} {
	return &WebAppActiveSlotModel{}
}

func (r WebAppActiveSlotResource) ResourceType() string {
	return "azurerm_web_app_active_slot"
}

func (r WebAppActiveSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateAppServiceID
}

func (r WebAppActiveSlotResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"slot_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Description:  "The ID of the Slot to swap with `Production`.",
			ValidateFunc: webapps.ValidateSlotID,
		},

		"overwrite_network_config": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "The swap action should overwrite the Production slot's network configuration with the configuration from this slot. Defaults to `true`.",
			ForceNew:    true,
		},
	}
}

func (r WebAppActiveSlotResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"last_successful_swap": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The timestamp of the last successful swap with `Production`",
		},
	}
}

func (r WebAppActiveSlotResource) Create() sdk.ResourceFunc {
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

func (r WebAppActiveSlotResource) Read() sdk.ResourceFunc {
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

func (r WebAppActiveSlotResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// Nothing to do here - there's no actual resource to delete
			// Note: deleting does not change the active slot, nor revert to any previous state.
			return nil
		},
	}
}

// Note: `Update` re-uses `Create` as there is no actual resource being managed, this meta-resource simply triggers a
// swap operations between the named slot and `Production`. Without this changing which slot is `Active` would result in
// Terraform deleting and recreating this resource, which may cause concern that the operation is somehow destructive.

func (r WebAppActiveSlotResource) Update() sdk.ResourceFunc {
	return r.Create()
}
