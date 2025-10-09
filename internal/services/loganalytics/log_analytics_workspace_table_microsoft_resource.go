// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
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
	_ sdk.ResourceWithUpdate        = WorkspaceTableMicrosoftResource{}
	_ sdk.ResourceWithCustomizeDiff = WorkspaceTableMicrosoftResource{}
)

type WorkspaceTableMicrosoftResource struct{}

type WorkspaceTableMicrosoftResourceModel struct {
	Name                 string   `tfschema:"name"`
	WorkspaceId          string   `tfschema:"workspace_id"`
	DisplayName          string   `tfschema:"display_name"`
	Description          string   `tfschema:"description"`
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

func (r WorkspaceTableMicrosoftResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var table WorkspaceTableMicrosoftResourceModel
			if err := metadata.DecodeDiff(&table); err != nil {
				return err
			}

			if strings.HasSuffix(table.Name, "_CL") {
				return errors.New("name must not end with '_CL' for Microsoft tables")
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

func (r WorkspaceTableMicrosoftResource) Arguments() map[string]*pluginsdk.Schema {
	args := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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

		"sub_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForTableSubTypeEnum(), false),
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				// Microsoft tables can flap between "Any" and "DataCollectionRuleBased"
				// This is due to Azure API behavior, not user changes
				if (old == "Any" && new == "DataCollectionRuleBased") ||
					(old == "DataCollectionRuleBased" && new == "Any") {
					return true
				}
				return false
			},
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

	return args
}

func (r WorkspaceTableMicrosoftResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r WorkspaceTableMicrosoftResource) ModelObject() interface{} {
	return &WorkspaceTableMicrosoftResourceModel{}
}

func (r WorkspaceTableMicrosoftResource) ResourceType() string {
	return "azurerm_log_analytics_workspace_table_microsoft"
}

func (r WorkspaceTableMicrosoftResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return tables.ValidateTableID
}

func (r WorkspaceTableMicrosoftResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute, // NOT 30m like other resources, as Microsoft tables are made with the log analytics workspace
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkspaceTableMicrosoftResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient

			tableName := model.Name
			log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Workspace Table %s update.", tableName)

			workspaceId, err := workspaces.ParseWorkspaceID(model.WorkspaceId)
			if err != nil {
				return fmt.Errorf("invalid workspace object ID for table %s: %s", tableName, err)
			}

			id := tables.NewTableID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, tableName)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", tableName, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			if model.SubType == string(tables.TableSubTypeEnumClassic) {
				return errors.New("sub_type 'Classic' tables cannot be created with this resource")
			}

			param := tables.Table{
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
						TableType:    pointer.To(tables.TableTypeEnumMicrosoft),
					},
				},
			}

			if model.Plan == string(tables.TablePlanEnumAnalytics) {
				if model.RetentionInDays == 0 {
					param.Properties.RetentionInDays = defaultRetentionInDays
				} else {
					param.Properties.RetentionInDays = pointer.To(model.RetentionInDays)
				}
				if model.TotalRetentionInDays == 0 {
					param.Properties.TotalRetentionInDays = defaultRetentionInDays
				} else {
					param.Properties.TotalRetentionInDays = pointer.To(model.TotalRetentionInDays)
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

func (r WorkspaceTableMicrosoftResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
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

			client := metadata.Client.LogAnalytics.TablesClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := WorkspaceTableMicrosoftResourceModel{
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
					state.Plan = pointer.FromEnum(props.Plan)

					if props.Schema != nil {
						state.DisplayName = pointer.From(props.Schema.DisplayName)
						state.Description = pointer.From(props.Schema.Description)
						state.SubType = pointer.FromEnum(props.Schema.TableSubType)

						if categories, ok := metadata.ResourceData.GetOk("categories"); ok {
							// Preserve configured categories if API doesn't return them
							categoriesSet := categories.(*pluginsdk.Set)
							categories := make([]string, 0, categoriesSet.Len())
							for _, item := range categoriesSet.List() {
								categories = append(categories, item.(string))
							}
							state.Categories = categories
						}

						state.Labels = pointer.From(props.Schema.Labels)
						state.Solutions = pointer.From(props.Schema.Solutions)

						if props.Schema.Columns != nil {
							state.Columns = flattenColumns(props.Schema.Columns)
						} else {
							state.Columns = nil
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

func (r WorkspaceTableMicrosoftResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.TablesClient
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config WorkspaceTableMicrosoftResourceModel
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

			// Create / Update requests MUST have a nil value for `StandardColumns`
			param.Properties.Schema.StandardColumns = nil

			if metadata.ResourceData.HasChange("plan") {
				param.Properties.Plan = pointer.To(tables.TablePlanEnum(config.Plan))
			}

			if config.Plan == string(tables.TablePlanEnumAnalytics) {
				if metadata.ResourceData.HasChange("retention_in_days") {
					param.Properties.RetentionInDays = defaultRetentionInDays
					if config.RetentionInDays != 0 {
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

			if metadata.ResourceData.HasChange("categories") {
				param.Properties.Schema.Categories = pointer.To(config.Categories)
			}

			if metadata.ResourceData.HasChange("labels") {
				param.Properties.Schema.Labels = pointer.To(config.Labels)
			}

			if metadata.ResourceData.HasChange("column") {
				param.Properties.Schema.Columns = expandColumns(&config.Columns)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, pointer.From(param)); err != nil {
				return fmt.Errorf("updating %s: %+v", id.TableName, err)
			}

			return nil
		},
	}
}

func (r WorkspaceTableMicrosoftResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkspaceTableMicrosoftResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			// We can't delete Microsoft tables, so we'll just set the retention to workspace default
			updateInput := tables.Table{
				Properties: &tables.TableProperties{
					RetentionInDays:      defaultRetentionInDays,
					TotalRetentionInDays: defaultRetentionInDays,
					Schema: &tables.Schema{
						Name:    pointer.To(id.TableName),
						Columns: &[]tables.Column{},
					},
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
