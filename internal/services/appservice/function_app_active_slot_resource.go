// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

type FunctionAppActiveSlotResource struct{}

type FunctionAppActiveSlotModel struct {
	SlotID              string `tfschema:"slot_id"`
	OverwriteNetworking bool   `tfschema:"overwrite_network_config"` // Note: This setting controls the ambiguously named `PreserveVnet`
	LastSwap            string `tfschema:"last_successful_swap"`
}

var _ sdk.ResourceWithUpdate = FunctionAppActiveSlotResource{}

func (r FunctionAppActiveSlotResource) ModelObject() interface{} {
	return &FunctionAppActiveSlotModel{}
}

func (r FunctionAppActiveSlotResource) ResourceType() string {
	return "azurerm_function_app_active_slot"
}

func (r FunctionAppActiveSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.FunctionAppID
}

func (r FunctionAppActiveSlotResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"slot_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Description:  "The ID of the Slot to swap with `Production`.",
			ValidateFunc: validate.FunctionAppSlotID,
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

func (r FunctionAppActiveSlotResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"last_successful_swap": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The timestamp of the last successful swap with `Production`",
		},
	}
}

func (r FunctionAppActiveSlotResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var activeSlot FunctionAppActiveSlotModel

			if err := metadata.Decode(&activeSlot); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.FunctionAppSlotID(activeSlot.SlotID)
			appId := parse.NewWebAppID(id.SubscriptionId, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("parsing App ID: %+v", err)
			}

			app, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if utils.ResponseWasNotFound(app.Response) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			csmSlotEntity := web.CsmSlotEntity{
				TargetSlot:   &id.SlotName,
				PreserveVnet: &activeSlot.OverwriteNetworking,
			}

			locks.ByID(appId.ID())
			defer locks.UnlockByID(appId.ID())

			future, err := client.SwapSlotWithProduction(ctx, id.ResourceGroup, id.SiteName, csmSlotEntity)
			if err != nil {
				return fmt.Errorf("making %s the active slot: %+v", id.SlotName, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for slot swap to complete: %+v", err)
			}

			metadata.SetID(appId)

			return nil
		},
	}
}

func (r FunctionAppActiveSlotResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := parse.FunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			app, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if utils.ResponseWasNotFound(app.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading active slot for %s: %+v", id.SiteName, err)
			}

			if app.SiteProperties == nil || app.SiteProperties.SlotSwapStatus == nil {
				return fmt.Errorf("reading site properties to determine active slot status: %+v", err)
			}

			activeSlot := FunctionAppActiveSlotModel{
				LastSwap: app.SiteProperties.SlotSwapStatus.TimestampUtc.String(),
			}

			if slotName := app.SiteProperties.SlotSwapStatus.SourceSlotName; slotName != nil {
				activeSlot.SlotID = parse.NewWebAppSlotID(id.SubscriptionId, id.ResourceGroup, id.SiteName, *slotName).ID()
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

func (r FunctionAppActiveSlotResource) Delete() sdk.ResourceFunc {
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

func (r FunctionAppActiveSlotResource) Update() sdk.ResourceFunc {
	return r.Create()
}
