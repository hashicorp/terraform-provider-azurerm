package appservice

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ActiveSlotResource struct{}

type ActiveSlotModel struct {
	AppID        string `tfschema:"app_id"`
	SlotName     string `tfschema:"slot_name"`
	PreserveVnet bool   `tfschema:"preserve_vnet"`
	LastSwap     string `tfschema:"last_successful_swap"`
}

func (r ActiveSlotResource) ModelObject() interface{} {
	return ActiveSlotModel{}
}

func (r ActiveSlotResource) ResourceType() string {
	return "" // This is never called
}

var _ sdk.Resource = ActiveSlotResource{}

func (r ActiveSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nil // This is never called
}

func (r ActiveSlotResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"app_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				validate.WebAppID,
				validate.FunctionAppID,
			),
		},

		"slot_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				validate.WebAppName,
				//validate.FunctionAppSlotID, // Uncomment after #14940 is merged
			),
		},

		"preserve_vnet": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},
	}
}

func (r ActiveSlotResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"last_successful_swap": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ActiveSlotResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var activeSlot ActiveSlotModel

			if err := metadata.Decode(&activeSlot); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			appId, err := parse.WebAppID(activeSlot.AppID)
			if err != nil {
				return fmt.Errorf("parsing App ID: %+v", err)
			}

			resp, err := client.Get(ctx, appId.ResourceGroup, appId.SiteName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return fmt.Errorf("%s was not found", appId)
				}
				return fmt.Errorf("reading  %s: %+v", appId, err)
			}

			slotId, err := parse.WebAppSlotID(activeSlot.AppID)
			if err != nil {
				return fmt.Errorf("parsing Slot ID: %+v", err)
			}

			if _, err = client.GetSlot(ctx, slotId.ResourceGroup, slotId.SiteName, slotId.SlotName); err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return fmt.Errorf("%s was not found", slotId)
				}
				return fmt.Errorf("reading %s: %+v", slotId, err)
			}

			csmSlotEntity := web.CsmSlotEntity{
				TargetSlot:   &slotId.SlotName,
				PreserveVnet: &activeSlot.PreserveVnet,
			}
			locks.ByID(appId.ID())
			defer locks.UnlockByID(appId.ID())

			future, err := client.SwapSlotWithProduction(ctx, appId.ResourceGroup, appId.SiteName, csmSlotEntity)
			if err != nil {
				return fmt.Errorf("making %s the active slot: %+v", slotId.SlotName, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for slot swap to complete: %+v", err)
			}

			metadata.SetID(appId)

			return nil
		},
	}
}

func (r ActiveSlotResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := parse.WebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			app, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading active slot for %s: %+v", id.SiteName, err)
			}

			if app.SiteProperties == nil || app.SiteProperties.SlotSwapStatus == nil {
				return fmt.Errorf("reading site properties to determine active slot status: %+v", err)
			}

			activeSlot := ActiveSlotModel{
				AppID:        utils.NormalizeNilableString(app.ID),
				SlotName:     utils.NormalizeNilableString(app.SiteProperties.SlotSwapStatus.SourceSlotName),
				PreserveVnet: metadata.ResourceData.Get("preserve_vnet").(bool), // This cannot be read, so we'll grab it from state
				LastSwap:     app.SiteProperties.SlotSwapStatus.TimestampUtc.String(),
			}

			return metadata.Encode(&activeSlot)
		},
	}
}

func (r ActiveSlotResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// Nothing to do here - there's no actual resource to delete
			// Note: deleting does not change the active slot, nor revert to any previous state.
			return nil
		},
	}
}
