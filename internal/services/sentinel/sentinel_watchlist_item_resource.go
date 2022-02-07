package sentinel

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WatchlistItemResource struct{}

var _ sdk.ResourceWithUpdate = WatchlistItemResource{}

type WatchlistItemModel struct {
	Name        string                 `tfschema:"name"`
	WatchlistID string                 `tfschema:"watchlist_id"`
	Properties  map[string]interface{} `tfschema:"properties"`
}

func (r WatchlistItemResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
		"watchlist_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WatchlistID,
		},
		"properties": {
			Type:     pluginsdk.TypeMap,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r WatchlistItemResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WatchlistItemResource) ResourceType() string {
	return "azurerm_sentinel_watchlist_item"
}

func (r WatchlistItemResource) ModelObject() interface{} {
	return &WatchlistItemModel{}
}

func (r WatchlistItemResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.WatchlistItemID
}

func (r WatchlistItemResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistItemsClient

			var model WatchlistItemModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			// Generate a random UUID as the resource name if the user doesn't specify it.
			if model.Name == "" {
				model.Name = uuid.New().String()
			}

			watchlistId, err := parse.WatchlistID(model.WatchlistID)
			if err != nil {
				return err
			}

			id := parse.NewWatchlistItemID(watchlistId.SubscriptionId, watchlistId.ResourceGroup, watchlistId.WorkspaceName, watchlistId.Name, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.WatchlistName, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			params := securityinsight.WatchlistItem{
				WatchlistItemProperties: &securityinsight.WatchlistItemProperties{
					ItemsKeyValue: model.Properties,
				},
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.WatchlistName, id.Name, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WatchlistItemResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistItemsClient
			id, err := parse.WatchlistItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.WatchlistName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			watchlistId := parse.NewWatchlistID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.WatchlistName)

			var properties map[string]interface{}
			if props := resp.WatchlistItemProperties; props != nil {
				if itemsKV := props.ItemsKeyValue; itemsKV != nil {
					properties = itemsKV.(map[string]interface{})
				}
			}
			model := WatchlistItemModel{
				WatchlistID: watchlistId.ID(),
				Name:        id.Name,
				Properties:  properties,
			}

			return metadata.Encode(&model)
		},
	}
}

func (r WatchlistItemResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistItemsClient

			id, err := parse.WatchlistItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.WatchlistName, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r WatchlistItemResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistItemsClient

			var model WatchlistItemModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			watchlistId, err := parse.WatchlistID(model.WatchlistID)
			if err != nil {
				return err
			}
			id := parse.NewWatchlistItemID(watchlistId.SubscriptionId, watchlistId.ResourceGroup, watchlistId.WorkspaceName, watchlistId.Name, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.WatchlistName, id.Name)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			update := securityinsight.WatchlistItem{
				WatchlistItemProperties: existing.WatchlistItemProperties,
			}

			if metadata.ResourceData.HasChange("properties") {
				if update.WatchlistItemProperties == nil {
					update.WatchlistItemProperties = &securityinsight.WatchlistItemProperties{}
				}
				update.WatchlistItemProperties.ItemsKeyValue = model.Properties
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.WatchlistName, id.Name, update); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
