package loganalytics

import (
	"context"
	"fmt"
	"log"
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
	_ sdk.ResourceWithUpdate        = WorkspaceTableMicrosoftResource{}
	_ sdk.ResourceWithCustomizeDiff = WorkspaceTableMicrosoftResource{}
)

type WorkspaceTableMicrosoftResource struct{}

type WorkspaceTableMicrosoftResourceModel struct {
	Name                 string   `tfschema:"name"`
	WorkspaceId          string   `tfschema:"workspace_id"`
	DisplayName          string   `tfschema:"display_name"`
	Description          string   `tfschema:"description"`
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
	Type             string `tfschema:"type"`
	Hidden           bool   `tfschema:"hidden"`
	DisplayByDefault bool   `tfschema:"display_by_default"`
}

// CustomizeDiff validates the configuration before applying changes
func (r WorkspaceTableMicrosoftResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var table WorkspaceTableMicrosoftResourceModel
			if err := metadata.DecodeDiff(&table); err != nil {
				return fmt.Errorf("decoding diff: %+v", err)
			}

			if table.RetentionInDays > 0 && table.TotalRetentionInDays > 0 {
				if table.TotalRetentionInDays < table.RetentionInDays {
					return fmt.Errorf("`total_retention_in_days` (%d) must be greater than or equal to `retention_in_days` (%d)", table.TotalRetentionInDays, table.RetentionInDays)
				}
			}

			return nil
		},
	}
}

func (r WorkspaceTableMicrosoftResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Alert",
				"AppCenterError",
				"ComputerGroup",
				"InsightsMetrics",
				"Operation",
				"Usage",
			}, false),
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

		"column": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringMatch(regexp.MustCompile(`.*_CF`), "`column.name` must end with '_CF'"),
					},

					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(tables.PossibleValuesForColumnTypeEnum(), false),
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
				},
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
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkspaceTableMicrosoftResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient

			tableName := model.Name
			log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Workspace Table %s create", tableName)

			workspaceId, err := workspaces.ParseWorkspaceID(model.WorkspaceId)
			if err != nil {
				return err
			}

			id := tables.NewTableID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, tableName)
			_, err = client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", tableName, err)
			}
			// Microsoft tables are always automatically provisioned whenever log analytics workspaces are provisioned, so there's no point in returning a 'resource already exists' error

			param := tables.Table{
				Properties: &tables.TableProperties{
					Plan:                 pointer.To(tables.TablePlanEnumAnalytics),
					RetentionInDays:      defaultRetentionInDays,
					TotalRetentionInDays: defaultRetentionInDays,
					Schema: &tables.Schema{
						Columns:     expandWorkspaceTableMicrosoftColumns(model.Columns),
						DisplayName: pointer.To(model.DisplayName),
						Description: pointer.To(model.Description),
						Labels:      pointer.To(model.Labels),
						Name:        pointer.To(tableName),
						TableType:   pointer.To(tables.TableTypeEnumMicrosoft),
					},
				},
			}

			if model.RetentionInDays > 0 {
				param.Properties.RetentionInDays = pointer.To(model.RetentionInDays)
			}

			if model.TotalRetentionInDays > 0 {
				param.Properties.TotalRetentionInDays = pointer.To(model.TotalRetentionInDays)
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

			workspaceId := pointer.To(workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName))

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
						if asDefault := props.RetentionInDaysAsDefault; asDefault != nil && !*asDefault {
							state.RetentionInDays = pointer.From(props.RetentionInDays)
						}
						if asDefault := props.TotalRetentionInDaysAsDefault; asDefault != nil && !*asDefault {
							state.TotalRetentionInDays = pointer.From(props.TotalRetentionInDays)
						}
					}

					if schema := props.Schema; schema != nil {
						state.DisplayName = pointer.From(schema.DisplayName)
						state.Description = pointer.From(schema.Description)
						state.Labels = pointer.From(schema.Labels)
						state.Solutions = pointer.From(schema.Solutions)
						state.Columns = flattenWorkspaceTableMicrosoftColumns(schema.Columns)
						state.StandardColumns = flattenWorkspaceTableMicrosoftColumns(schema.StandardColumns)
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

			var config WorkspaceTableMicrosoftResourceModel
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
				props.Schema.Columns = expandWorkspaceTableMicrosoftColumns(config.Columns)
			}

			if metadata.ResourceData.HasChange("display_name") {
				props.Schema.DisplayName = pointer.To(config.DisplayName)
			}

			if metadata.ResourceData.HasChange("description") {
				props.Schema.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("labels") {
				props.Schema.Labels = pointer.To(config.Labels)
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
				return err
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

func expandWorkspaceTableMicrosoftColumns(columns []Column) *[]tables.Column {
	result := make([]tables.Column, 0, len(columns))
	for _, column := range columns {
		columnToAdd := tables.Column{
			Name:             pointer.To(column.Name),
			Type:             pointer.ToEnum[tables.ColumnTypeEnum](column.Type),
			IsHidden:         pointer.To(column.Hidden),
			IsDefaultDisplay: pointer.To(column.DisplayByDefault),
		}
		if column.DisplayName != "" {
			columnToAdd.DisplayName = pointer.To(column.DisplayName)
		}
		if column.Description != "" {
			columnToAdd.Description = pointer.To(column.Description)
		}
		result = append(result, columnToAdd)
	}
	return pointer.To(result)
}

func flattenWorkspaceTableMicrosoftColumns(columns *[]tables.Column) []Column {
	if columns == nil {
		return nil
	}
	result := make([]Column, 0, len(*columns))
	for _, column := range *columns {
		result = append(result, Column{
			Name:             pointer.From(column.Name),
			DisplayName:      pointer.From(column.DisplayName),
			Description:      pointer.From(column.Description),
			Type:             pointer.FromEnum(column.Type),
			Hidden:           pointer.From(column.IsHidden),
			DisplayByDefault: pointer.From(column.IsDefaultDisplay),
		})
	}
	return result
}
