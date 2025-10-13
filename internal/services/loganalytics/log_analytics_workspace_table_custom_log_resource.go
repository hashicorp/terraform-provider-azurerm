// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithUpdate        = WorkspaceTableCustomLogResource{}
	_ sdk.ResourceWithCustomizeDiff = WorkspaceTableCustomLogResource{}
)

type WorkspaceTableCustomLogResource struct{}

type WorkspaceTableCustomLogResourceModel struct {
	Name                 string   `tfschema:"name"`
	WorkspaceId          string   `tfschema:"workspace_id"`
	DisplayName          string   `tfschema:"display_name"`
	Description          string   `tfschema:"description"`
	Plan                 string   `tfschema:"plan"`
	Columns              []Column `tfschema:"column"`
	Labels               []string `tfschema:"labels"`
	Solutions            []string `tfschema:"solutions"`
	StandardColumns      []Column `tfschema:"standard_column"`
	RetentionInDays      int64    `tfschema:"retention_in_days"`
	TotalRetentionInDays int64    `tfschema:"total_retention_in_days"`
}

func (r WorkspaceTableCustomLogResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var table WorkspaceTableCustomLogResourceModel
			if err := metadata.DecodeDiff(&table); err != nil {
				return err
			}

			for _, column := range table.Columns {
				if column.TypeHint != "" && column.Type != string(tables.ColumnTypeEnumString) {
					return errors.New("type_hint can only be set for columns of type 'string'")
				}
			}

			if table.Plan == string(tables.TablePlanEnumBasic) {
				if _, ok := metadata.ResourceDiff.GetOk("retention_in_days"); ok {
					return errors.New("cannot set retention_in_days because the retention is fixed at eight days on Basic plan")
				}
			}

			return nil
		},
	}
}

func (r WorkspaceTableCustomLogResource) Arguments() map[string]*pluginsdk.Schema {
	args := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`_CL$`), "This must end with '_CL'."),
		},

		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"plan": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(tables.TablePlanEnumAnalytics),
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForTablePlanEnum(), false),
		},

		"column": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: columnSchema(),
			},
		},

		"labels": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"retention_in_days": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C If not specified, defaults to the workspace's retention period
			Computed:     true,
			ValidateFunc: validation.IntBetween(4, 730),
		},

		"total_retention_in_days": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C If not specified, defaults to the workspace's retention period
			Computed:     true,
			ValidateFunc: validation.Any(validation.IntBetween(4, 730), validation.IntInSlice([]int{1095, 1460, 1826, 2191, 2556, 2922, 3288, 3653, 4018, 4383})),
		},
	}

	return args
}

func (r WorkspaceTableCustomLogResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"solutions": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"standard_column": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"display_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"description": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"type_hint": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"hidden": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"display_by_default": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r WorkspaceTableCustomLogResource) ModelObject() interface{} {
	return &WorkspaceTableCustomLogResourceModel{}
}

func (r WorkspaceTableCustomLogResource) ResourceType() string {
	return "azurerm_log_analytics_workspace_table_custom_log"
}

func (r WorkspaceTableCustomLogResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return tables.ValidateTableID
}

func (r WorkspaceTableCustomLogResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config WorkspaceTableCustomLogResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient

			tableName := config.Name

			workspaceId, err := workspaces.ParseWorkspaceID(config.WorkspaceId)
			if err != nil {
				return fmt.Errorf("invalid workspace object ID for table %s: %s", tableName, err)
			}

			id := tables.NewTableID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, tableName)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := tables.Table{
				Properties: &tables.TableProperties{
					Plan: pointer.To(tables.TablePlanEnum(config.Plan)),
					Schema: &tables.Schema{
						Columns:      expandColumns(config.Columns),
						DisplayName:  pointer.To(config.DisplayName),
						Description:  pointer.To(config.Description),
						Labels:       pointer.To(config.Labels),
						Name:         pointer.To(tableName),
						TableSubType: pointer.To(tables.TableSubTypeEnumDataCollectionRuleBased),
						TableType:    pointer.To(tables.TableTypeEnumCustomLog),
					},
				},
			}

			if config.Plan == string(tables.TablePlanEnumAnalytics) {
				if config.RetentionInDays == 0 {
					param.Properties.RetentionInDays = defaultRetentionInDays
				} else {
					param.Properties.RetentionInDays = pointer.To(config.RetentionInDays)
				}
				if config.TotalRetentionInDays == 0 {
					param.Properties.TotalRetentionInDays = defaultRetentionInDays
				} else {
					param.Properties.TotalRetentionInDays = pointer.To(config.TotalRetentionInDays)
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkspaceTableCustomLogResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.TablesClient

			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			var workspaceId *workspaces.WorkspaceId
			if workspaceStateId, ok := metadata.ResourceData.GetOk("workspace_id"); ok {
				workspaceId, err = workspaces.ParseWorkspaceID(workspaceStateId.(string))
				if err != nil {
					return err
				}
			} else {
				workspaceId = pointer.To(workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName))
			}

			state := WorkspaceTableCustomLogResourceModel{
				WorkspaceId: workspaceId.ID(),
				Name:        id.TableName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					planAllowsCustomRetention := props.Plan != nil && *props.Plan == tables.TablePlanEnumAnalytics
					if planAllowsCustomRetention {
						if !pointer.From(props.RetentionInDaysAsDefault) {
							state.RetentionInDays = pointer.From(props.RetentionInDays)
						}
						if !pointer.From(props.TotalRetentionInDaysAsDefault) {
							state.TotalRetentionInDays = pointer.From(props.TotalRetentionInDays)
						}
					}
					state.Plan = pointer.FromEnum(props.Plan)

					if props.Schema != nil {
						state.DisplayName = pointer.From(props.Schema.DisplayName)
						state.Description = pointer.From(props.Schema.Description)
						state.Labels = pointer.From(props.Schema.Labels)
						state.Solutions = pointer.From(props.Schema.Solutions)

						state.Columns = flattenColumns(props.Schema.Columns)
						if props.Schema.StandardColumns != nil {
							state.StandardColumns = flattenColumns(props.Schema.StandardColumns)
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkspaceTableCustomLogResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.TablesClient

			var config WorkspaceTableCustomLogResourceModel
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Log Analytics Workspace Table %s: %v", id, err)
			}

			if model := existing.Model; model == nil {
				return fmt.Errorf("model is nil: %+v", existing)
			}

			props := existing.Model.Properties

			if props == nil {
				return fmt.Errorf("properties is nil: %+v", existing)
			}

			param := existing.Model

			// Create / Update requests MUST have a nil value for `StandardColumns` or they will get a 400 response
			param.Properties.Schema.StandardColumns = nil

			if metadata.ResourceData.HasChange("plan") {
				param.Properties.Plan = pointer.To(tables.TablePlanEnum(config.Plan))
			}

			if config.Plan == string(tables.TablePlanEnumAnalytics) {
				if metadata.ResourceData.HasChange("retention_in_days") {
					if config.RetentionInDays == 0 {
						param.Properties.RetentionInDays = defaultRetentionInDays
					} else {
						param.Properties.RetentionInDays = pointer.To(config.RetentionInDays)
					}
				}

				if metadata.ResourceData.HasChange("total_retention_in_days") {
					if config.TotalRetentionInDays == 0 {
						param.Properties.TotalRetentionInDays = defaultRetentionInDays
					} else {
						param.Properties.TotalRetentionInDays = pointer.To(config.TotalRetentionInDays)
					}
				}
			}

			if metadata.ResourceData.HasChange("display_name") {
				param.Properties.Schema.DisplayName = pointer.To(config.DisplayName)
			}

			if metadata.ResourceData.HasChange("description") {
				param.Properties.Schema.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("labels") {
				param.Properties.Schema.Labels = pointer.To(config.Labels)
			}

			if metadata.ResourceData.HasChange("column") {
				param.Properties.Schema.Columns = expandColumns(config.Columns)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, pointer.From(param)); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r WorkspaceTableCustomLogResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkspaceTableCustomLogResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
