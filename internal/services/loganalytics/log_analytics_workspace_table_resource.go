// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LogAnalyticsWorkspaceTableResource struct{}

var (
	_                            sdk.ResourceWithUpdate        = LogAnalyticsWorkspaceTableResource{}
	_                            sdk.ResourceWithCustomizeDiff = LogAnalyticsWorkspaceTableResource{}
	useWorkspaceDefaultRetention                               = pointer.To(int64(-1))
)

type LogAnalyticsWorkspaceTableResourceModel struct {
	WorkspaceId          string   `tfschema:"workspace_id"`
	Name                 string   `tfschema:"name"`
	DisplayName          string   `tfschema:"display_name"`
	Description          string   `tfschema:"description"`
	Type                 string   `tfschema:"type"`
	SubType              string   `tfschema:"sub_type"`
	Plan                 string   `tfschema:"plan"`
	Categories           []string `tfschema:"categories"`
	Columns              []Column `tfschema:"column"`
	Labels               []string `tfschema:"labels"`
	Solutions            []string `tfschema:"solutions"`
	StandardColumns      []Column `tfschema:"standard_column"`
	RetentionInDays      int64    `tfschema:"retention_in_days"`
	TotalRetentionInDays int64    `tfschema:"total_retention_in_days"`
}

type Column struct {
	Name             string `tfschema:"name"`
	DisplayName      string `tfschema:"display_name"`
	Description      string `tfschema:"description"`
	IsHidden         bool   `tfschema:"hidden"`
	IsDefaultDisplay bool   `tfschema:"display_by_default"`
	Type             string `tfschema:"type"`
	TypeHint         string `tfschema:"type_hint"`
}

func (r LogAnalyticsWorkspaceTableResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			var table LogAnalyticsWorkspaceTableResourceModel
			if err := metadata.DecodeDiff(&table); err != nil {
				return err
			}

			switch table.Type {
			case string(tables.TableTypeEnumMicrosoft):
				if strings.HasSuffix(table.Name, "_CL") {
					return fmt.Errorf("name must not end with '_CL' for Microsoft tables")
				}

			case string(tables.TableTypeEnumCustomLog):
				if !strings.HasSuffix(table.Name, "_CL") {
					return fmt.Errorf("name must end with '_CL' for CustomLog tables")
				}
				if table.SubType == "" {
					return fmt.Errorf("sub_type must be set for CustomLog tables")
				}
				if table.SubType == string(tables.TableSubTypeEnumAny) {
					return fmt.Errorf("sub_type cannot be 'Any' for CustomLog tables")
				}
			}

			for _, column := range table.Columns {
				if column.TypeHint != "" && column.Type != string(tables.ColumnTypeEnumString) {
					return fmt.Errorf("type_hint can only be set for columns of type 'string'")
				}
			}

			if table.Plan == string(tables.TablePlanEnumBasic) {
				if _, ok := metadata.ResourceDiff.GetOk("retention_in_days"); ok {
					return fmt.Errorf("cannot set retention_in_days because the retention is fixed at eight days on Basic plan")
				}
			}

			return nil
		},
	}
}

func (r LogAnalyticsWorkspaceTableResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
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

		"type": {
			Type:         pluginsdk.TypeString,
			Required:     features.FivePointOh(),
			Optional:     !features.FivePointOh(),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForTableTypeEnum(), false),
			Default: func() interface{} {
				if !features.FivePointOh() {
					return string(tables.TableTypeEnumMicrosoft)
				}
				return nil
			}(),
		},

		"sub_type": {
			Type:         pluginsdk.TypeString,
			Required:     features.FivePointOh(),
			Optional:     !features.FivePointOh(),
			Computed:     !features.FivePointOh(),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForTableSubTypeEnum(), false),
		},

		"plan": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(tables.TablePlanEnumAnalytics),
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForTablePlanEnum(), false),
		},

		"categories": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"column": {
			Type:     pluginsdk.TypeList, // Order matters for display in Log Analytics
			Optional: true,
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
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(4, 730),
		},

		"total_retention_in_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.Any(validation.IntBetween(4, 730), validation.IntInSlice([]int{1095, 1460, 1826, 2191, 2556, 2922, 3288, 3653, 4018, 4383})),
		},
	}
}

func (r LogAnalyticsWorkspaceTableResource) Attributes() map[string]*pluginsdk.Schema {
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
				Schema: columnSchema(),
			},
		},
	}
}

func (r LogAnalyticsWorkspaceTableResource) ModelObject() interface{} {
	return &LogAnalyticsWorkspaceTableResourceModel{}
}

func (r LogAnalyticsWorkspaceTableResource) ResourceType() string {
	return "azurerm_log_analytics_workspace_table"
}

func (r LogAnalyticsWorkspaceTableResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return tables.ValidateTableID
}

func (r LogAnalyticsWorkspaceTableResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LogAnalyticsWorkspaceTableResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			tableName := model.Name
			log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Workspace Table %s update.", tableName)

			workspaceId, err := workspaces.ParseWorkspaceID(model.WorkspaceId)
			if err != nil {
				return fmt.Errorf("invalid workspace object ID for table %s: %s", tableName, err)
			}

			id := tables.NewTableID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, tableName)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("error checking for presence of existing table %s: %v", tableName, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			if model.SubType == string(tables.TableSubTypeEnumClassic) {
				return fmt.Errorf("sub_type 'Classic' tables cannot be created with this resource")
			}

			updateInput := tables.Table{
				Properties: &tables.TableProperties{
					Plan: pointer.To(tables.TablePlanEnum(model.Plan)),
					Schema: &tables.Schema{
						Categories:   pointer.To(model.Categories),
						Columns:      expandColumns(&model.Columns),
						DisplayName:  pointer.To(model.DisplayName),
						Description:  pointer.To(model.Description),
						Labels:       pointer.To(model.Labels),
						Name:         pointer.To(tableName),
						TableSubType: pointer.To(tables.TableSubTypeEnum(model.SubType)),
						TableType:    pointer.To(tables.TableTypeEnum(model.Type)),
					},
				},
			}

			if model.Plan == string(tables.TablePlanEnumAnalytics) {
				if model.RetentionInDays == 0 {
					// Set the retention period to the workspace default
					updateInput.Properties.RetentionInDays = useWorkspaceDefaultRetention
				} else {
					updateInput.Properties.RetentionInDays = pointer.To(model.RetentionInDays)
				}
				if model.TotalRetentionInDays == 0 {
					// Set the retention period to the workspace default
					updateInput.Properties.TotalRetentionInDays = useWorkspaceDefaultRetention
				} else {
					updateInput.Properties.TotalRetentionInDays = pointer.To(model.TotalRetentionInDays)
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, updateInput); err != nil {
				return fmt.Errorf("failed to update table %s in workspace %s in resource group %s: %s", tableName, workspaceId.WorkspaceName, workspaceId.ResourceGroupName, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LogAnalyticsWorkspaceTableResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.TablesClient
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LogAnalyticsWorkspaceTableResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Log Analytics Workspace Table %s: %v", id, err)
			}

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					updateInput := tables.Table{
						Properties: &tables.TableProperties{
							Plan: props.Plan,
							Schema: &tables.Schema{
								Categories:  props.Schema.Categories,
								Columns:     props.Schema.Columns,
								DisplayName: props.Schema.DisplayName,
								Description: props.Schema.Description,
								Labels:      props.Schema.Labels,
								Name:        props.Schema.Name,
							},
						},
					}

					if metadata.ResourceData.HasChange("plan") {
						updateInput.Properties.Plan = pointer.To(tables.TablePlanEnum(state.Plan))
					}

					if state.Plan == string(tables.TablePlanEnumAnalytics) {
						if metadata.ResourceData.HasChange("retention_in_days") {
							if state.RetentionInDays == 0 {
								// Set the retention period to the workspace default
								updateInput.Properties.RetentionInDays = useWorkspaceDefaultRetention
							} else {
								// Set the retention period to the workspace default
								updateInput.Properties.RetentionInDays = pointer.To(state.RetentionInDays)
							}
						}

						if metadata.ResourceData.HasChange("total_retention_in_days") {
							if state.TotalRetentionInDays == 0 {
								updateInput.Properties.TotalRetentionInDays = useWorkspaceDefaultRetention
							} else {
								updateInput.Properties.TotalRetentionInDays = pointer.To(state.TotalRetentionInDays)
							}
						}
					}

					if metadata.ResourceData.HasChange("display_name") {
						updateInput.Properties.Schema.DisplayName = pointer.To(state.DisplayName)
					}

					if metadata.ResourceData.HasChange("description") {
						updateInput.Properties.Schema.Description = pointer.To(state.Description)
					}

					if metadata.ResourceData.HasChange("categories") {
						updateInput.Properties.Schema.Categories = pointer.To(state.Categories)
					}

					if metadata.ResourceData.HasChange("labels") {
						updateInput.Properties.Schema.Labels = pointer.To(state.Labels)
					}

					if metadata.ResourceData.HasChange("column") {
						updateInput.Properties.Schema.Columns = expandColumns(&state.Columns)
					}

					if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput); err != nil {
						return fmt.Errorf("failed to update table: %s: %+v", id.TableName, err)
					}
				}
			}

			return nil
		},
	}
}

func (r LogAnalyticsWorkspaceTableResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			var workspaceId *workspaces.WorkspaceId
			if workspaceStateId, ok := metadata.ResourceData.GetOk("workspace_id"); ok {
				workspaceId, err = workspaces.ParseWorkspaceID(workspaceStateId.(string))
				if err != nil {
					return fmt.Errorf("while parsing resource ID: %+v", err)
				}
			} else {
				workspaceId = pointer.To(workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName))
			}

			client := metadata.Client.LogAnalytics.TablesClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving Log Analytics Workspace Table %s: %+v", *id, err)
			}

			state := LogAnalyticsWorkspaceTableResourceModel{
				Name:        id.TableName,
				WorkspaceId: workspaceId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if pointer.From(props.Plan) == tables.TablePlanEnumAnalytics {
						if !pointer.From(props.RetentionInDaysAsDefault) {
							state.RetentionInDays = pointer.From(props.RetentionInDays)
						}
						if !pointer.From(props.TotalRetentionInDaysAsDefault) {
							state.TotalRetentionInDays = pointer.From(props.TotalRetentionInDays)
						}
					}
					state.TotalRetentionInDays = pointer.From(props.TotalRetentionInDays)
					state.Plan = string(pointer.From(props.Plan))

					if props.Schema != nil {
						state.DisplayName = pointer.From(props.Schema.DisplayName)
						state.Description = pointer.From(props.Schema.Description)
						state.Type = string(pointer.From(props.Schema.TableType))
						state.SubType = string(pointer.From(props.Schema.TableSubType))
						state.Categories = pointer.From(props.Schema.Categories)
						state.Labels = pointer.From(props.Schema.Labels)
						state.Solutions = pointer.From(props.Schema.Solutions)

						if props.Schema.Columns != nil {
							state.Columns = flattenColumns(props.Schema.Columns)
						}

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

func (r LogAnalyticsWorkspaceTableResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LogAnalyticsWorkspaceTableResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			if model.Type == string(tables.TableTypeEnumMicrosoft) {
				// We can't delete Microsoft tables, so we'll just set the retention to workspace default
				updateInput := tables.Table{
					Properties: &tables.TableProperties{
						RetentionInDays:      useWorkspaceDefaultRetention,
						TotalRetentionInDays: useWorkspaceDefaultRetention,
						Schema: &tables.Schema{
							Name:    pointer.To(id.TableName),
							Columns: &[]tables.Column{},
						},
					},
				}

				if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput); err != nil {
					return fmt.Errorf("failed to update table %s in workspace %s in resource group %s: %s", id.TableName, id.WorkspaceName, id.ResourceGroupName, err)
				}
			} else {
				if err := client.DeleteThenPoll(ctx, *id); err != nil {
					return fmt.Errorf("deleting Log Analytics Workspace Table %s: %v", id, err)
				}
			}

			return nil
		},
	}
}

func columnSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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

		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForColumnTypeEnum(), false),
		},

		"type_hint": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForColumnDataTypeHintEnum(), false),
		},

		"hidden": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"display_by_default": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
	}
}

func expandColumns(columns *[]Column) *[]tables.Column {
	var result []tables.Column
	for _, column := range *columns {
		result = append(result, tables.Column{
			Name:             pointer.To(column.Name),
			DisplayName:      pointer.To(column.DisplayName),
			Description:      pointer.To(column.Description),
			IsHidden:         pointer.To(column.IsHidden),
			IsDefaultDisplay: pointer.To(column.IsDefaultDisplay),
			Type:             pointer.To(tables.ColumnTypeEnum(column.Type)),
			DataTypeHint:     pointer.To(tables.ColumnDataTypeHintEnum(column.TypeHint)),
		})
	}
	return pointer.To(result)
}

func flattenColumns(columns *[]tables.Column) []Column {
	var result []Column
	for _, column := range *columns {
		result = append(result, Column{
			Name:             pointer.From(column.Name),
			DisplayName:      pointer.From(column.DisplayName),
			Description:      pointer.From(column.Description),
			IsHidden:         pointer.From(column.IsHidden),
			IsDefaultDisplay: pointer.From(column.IsDefaultDisplay),
			Type:             string(pointer.From(column.Type)),
			TypeHint:         string(pointer.From(column.DataTypeHint)),
		})
	}
	return result
}
