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
	Name                 string                 `tfschema:"name"`
	WorkspaceId          string                 `tfschema:"workspace_id"`
	DisplayName          string                 `tfschema:"display_name"`
	Description          string                 `tfschema:"description"`
	Plan                 string                 `tfschema:"plan"`
	Columns              []workspaceTableColumn `tfschema:"column"`
	Solutions            []string               `tfschema:"solutions"`
	StandardColumns      []workspaceTableColumn `tfschema:"standard_column"`
	RetentionInDays      int64                  `tfschema:"retention_in_days"`
	TotalRetentionInDays int64                  `tfschema:"total_retention_in_days"`
}

func (r WorkspaceTableCustomLogResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var table WorkspaceTableCustomLogResourceModel
			if err := metadata.DecodeDiff(&table); err != nil {
				return err
			}

			if table.Plan == string(tables.TablePlanEnumBasic) {
				if _, ok := metadata.ResourceDiff.GetOk("retention_in_days"); ok {
					return errors.New("`retention_in_days` cannot be set when `plan` is set to `Basic`")
				}
			}

			return nil
		},
	}
}

func (r WorkspaceTableCustomLogResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`_CL$`), "must end with '_CL'."),
		},

		"column": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: workspaceTableColumnSchema(),
			},
		},

		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
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
				Schema: workspaceTableColumnSchemaComputed(),
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

			workspaceId, err := workspaces.ParseWorkspaceID(config.WorkspaceId)
			if err != nil {
				return err
			}

			id := tables.NewTableID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, config.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := tables.Table{
				Properties: &tables.TableProperties{
					Plan:                 pointer.To(tables.TablePlanEnum(config.Plan)),
					RetentionInDays:      defaultRetentionInDays,
					TotalRetentionInDays: defaultRetentionInDays,
					Schema: &tables.Schema{
						Columns:      expandWorkspaceTableColumns(config.Columns),
						DisplayName:  pointer.To(config.DisplayName),
						Description:  pointer.To(config.Description),
						Name:         pointer.To(config.Name),
						TableSubType: pointer.To(tables.TableSubTypeEnumDataCollectionRuleBased),
						TableType:    pointer.To(tables.TableTypeEnumCustomLog),
					},
				},
			}

			if config.Plan == string(tables.TablePlanEnumAnalytics) {
				if config.RetentionInDays > 0 {
					param.Properties.RetentionInDays = pointer.To(config.RetentionInDays)
				}
				if config.TotalRetentionInDays > 0 {
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

			state := WorkspaceTableCustomLogResourceModel{
				WorkspaceId: pointer.To(workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)).ID(),
				Name:        id.TableName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Plan = pointer.FromEnum(props.Plan)

					if asDefault := props.RetentionInDaysAsDefault; asDefault != nil && !*asDefault {
						state.RetentionInDays = pointer.From(props.RetentionInDays)
					}
					if asDefault := props.TotalRetentionInDaysAsDefault; asDefault != nil && !*asDefault {
						state.TotalRetentionInDays = pointer.From(props.TotalRetentionInDays)
					}

					if schema := props.Schema; schema != nil {
						state.DisplayName = pointer.From(schema.DisplayName)
						state.Description = pointer.From(schema.Description)
						state.Solutions = pointer.From(schema.Solutions)
						state.Columns = flattenWorkspaceTableColumns(schema.Columns)
						state.StandardColumns = flattenWorkspaceTableColumns(schema.StandardColumns)
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
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			props := existing.Model.Properties
			if props.Schema == nil {
				props.Schema = &tables.Schema{}
			}

			// Create / Update requests MUST have a nil value for `StandardColumns` or they will get a 400 response
			props.Schema.StandardColumns = nil

			// If RetentionInDaysAsDefault is true, RetentionInDays still returns the actual value.
			// Since the `Update` reuses the payload from the GET request, that value is included.
			// If this value is then sent to Azure, `RetentionInDaysAsDefault` no longer returns `true`
			// causing a diff on the subsequent read where we then set this value into state.
			if pointer.From(props.RetentionInDaysAsDefault) {
				props.RetentionInDays = defaultRetentionInDays
			}

			// The comment above applies to `TotalRetentionInDaysAsDefault` as well.
			if pointer.From(props.TotalRetentionInDaysAsDefault) {
				props.TotalRetentionInDays = defaultRetentionInDays
			}

			if metadata.ResourceData.HasChange("column") {
				props.Schema.Columns = expandWorkspaceTableColumns(config.Columns)
			}

			if metadata.ResourceData.HasChange("display_name") {
				props.Schema.DisplayName = pointer.To(config.DisplayName)
			}

			if metadata.ResourceData.HasChange("description") {
				props.Schema.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("plan") {
				props.Plan = pointer.ToEnum[tables.TablePlanEnum](config.Plan)
			}

			if metadata.ResourceData.HasChange("retention_in_days") {
				props.RetentionInDays = pointer.To(config.RetentionInDays)
				if config.RetentionInDays == 0 {
					props.RetentionInDays = defaultRetentionInDays
				}
			}

			if metadata.ResourceData.HasChange("total_retention_in_days") {
				props.TotalRetentionInDays = pointer.To(config.TotalRetentionInDays)
				if config.TotalRetentionInDays == 0 {
					props.TotalRetentionInDays = defaultRetentionInDays
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
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
