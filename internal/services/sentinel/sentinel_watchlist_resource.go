package sentinel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type WatchlistResource struct{}

var _ sdk.Resource = WatchlistResource{}

type WatchlistModel struct {
	Name                    string   `tfschema:"name"`
	LogAnalyticsWorkspaceId string   `tfschema:"log_analytics_workspace_id"`
	DisplayName             string   `tfschema:"display_name"`
	Description             string   `tfschema:"description"`
	Labels                  []string `tfschema:"labels"`
	DefaultDuration         string   `tfschema:"default_duration"`
	ItemSearchKey           string   `tfschema:"item_search_key"`
}

func (r WatchlistResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},
		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"item_search_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"labels": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"default_duration": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonValidate.ISO8601Duration,
		},
	}
}

func (r WatchlistResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WatchlistResource) ResourceType() string {
	return "azurerm_sentinel_watchlist"
}

func (r WatchlistResource) ModelObject() interface{} {
	return &WatchlistModel{}
}

func (r WatchlistResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.WatchlistID
}

func (r WatchlistResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistsClient

			var model WatchlistModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			workspaceId, err := workspaces.ParseWorkspaceID(model.LogAnalyticsWorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing Log Analytics Workspace ID: %w", err)
			}

			id := parse.NewWatchlistID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := securityinsight.Watchlist{
				WatchlistProperties: &securityinsight.WatchlistProperties{
					DisplayName: &model.DisplayName,
					// The only supported provider for now is "Microsoft"
					Provider: utils.String("Microsoft"),

					ItemsSearchKey: utils.String(model.ItemSearchKey),
				},
			}

			if model.Description != "" {
				param.WatchlistProperties.Description = &model.Description
			}
			if len(model.Labels) != 0 {
				param.WatchlistProperties.Labels = &model.Labels
			}
			if model.DefaultDuration != "" {
				param.WatchlistProperties.DefaultDuration = &model.DefaultDuration
			}

			_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, param)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WatchlistResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistsClient
			id, err := parse.WatchlistID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := WatchlistModel{
				Name:                    id.Name,
				LogAnalyticsWorkspaceId: workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID(),
			}

			if props := resp.WatchlistProperties; props != nil {
				if props.DisplayName != nil {
					model.DisplayName = *props.DisplayName
				}
				if props.Description != nil {
					model.Description = *props.Description
				}
				if props.Labels != nil {
					model.Labels = *props.Labels
				}
				if props.DefaultDuration != nil {
					model.DefaultDuration = *props.DefaultDuration
				}
				if props.ItemsSearchKey != nil {
					model.ItemSearchKey = *props.ItemsSearchKey
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r WatchlistResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistsClient

			id, err := parse.WatchlistID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
