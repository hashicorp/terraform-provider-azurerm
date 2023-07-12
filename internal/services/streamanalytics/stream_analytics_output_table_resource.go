// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type OutputTableResource struct{}

var (
	_ sdk.ResourceWithCustomImporter = OutputTableResource{}
	_ sdk.ResourceWithStateMigration = OutputTableResource{}
)

type OutputTableResourceModel struct {
	Name               string   `tfschema:"name"`
	StreamAnalyticsJob string   `tfschema:"stream_analytics_job_name"`
	ResourceGroup      string   `tfschema:"resource_group_name"`
	StorageAccount     string   `tfschema:"storage_account_name"`
	StorageAccountKey  string   `tfschema:"storage_account_key"`
	Table              string   `tfschema:"table"`
	PartitionKey       string   `tfschema:"partition_key"`
	RowKey             string   `tfschema:"row_key"`
	BatchSize          int64    `tfschema:"batch_size"`
	ColumnsToRemove    []string `tfschema:"columns_to_remove"`
}

func (r OutputTableResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"stream_analytics_job_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_account_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"table": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"partition_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"row_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"batch_size": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"columns_to_remove": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func (r OutputTableResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r OutputTableResource) ModelObject() interface{} {
	return &OutputTableResourceModel{}
}

func (r OutputTableResource) ResourceType() string {
	return "azurerm_stream_analytics_output_table"
}

func (r OutputTableResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return outputs.ValidateOutputID
}

func (r OutputTableResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model OutputTableResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.OutputsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := outputs.NewOutputID(subscriptionId, model.ResourceGroup, model.StreamAnalyticsJob, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			tableOutputProps := &outputs.AzureTableOutputDataSourceProperties{
				AccountName:  utils.String(model.StorageAccount),
				AccountKey:   utils.String(model.StorageAccountKey),
				Table:        utils.String(model.Table),
				PartitionKey: utils.String(model.PartitionKey),
				RowKey:       utils.String(model.RowKey),
				BatchSize:    utils.Int64(model.BatchSize),
			}

			if v := model.ColumnsToRemove; len(v) > 0 {
				tableOutputProps.ColumnsToRemove = &v
			}

			props := outputs.Output{
				Name: utils.String(model.Name),
				Properties: &outputs.OutputProperties{
					Datasource: &outputs.AzureTableOutputDataSource{
						Properties: tableOutputProps,
					},
				},
			}

			var opts outputs.CreateOrReplaceOperationOptions
			if _, err = client.CreateOrReplace(ctx, id, props, opts); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r OutputTableResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient
			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					output, ok := props.Datasource.(outputs.AzureTableOutputDataSource)
					if !ok {
						return fmt.Errorf("converting %s to a Table Output", *id)
					}

					if output.Properties != nil {
						if output.Properties.AccountName == nil || output.Properties.Table == nil || output.Properties.PartitionKey == nil || output.Properties.RowKey == nil || output.Properties.BatchSize == nil {
							return nil
						}

						state := OutputTableResourceModel{
							Name:               id.OutputName,
							ResourceGroup:      id.ResourceGroupName,
							StreamAnalyticsJob: id.StreamingJobName,
							StorageAccountKey:  metadata.ResourceData.Get("storage_account_key").(string),
						}

						accountName := ""
						if v := output.Properties.AccountName; v != nil {
							accountName = *v
						}
						state.StorageAccount = accountName

						table := ""
						if v := output.Properties.Table; v != nil {
							table = *v
						}
						state.Table = table

						partitonKey := ""
						if v := output.Properties.PartitionKey; v != nil {
							partitonKey = *v
						}
						state.PartitionKey = partitonKey

						rowKey := ""
						if v := output.Properties.RowKey; v != nil {
							rowKey = *v
						}
						state.RowKey = rowKey

						var batchSize int64
						if v := output.Properties.BatchSize; v != nil {
							batchSize = *v
						}
						state.BatchSize = batchSize

						var columnsToRemove []string
						if columns := output.Properties.ColumnsToRemove; columns != nil && len(*columns) > 0 {
							columnsToRemove = *columns
						}
						state.ColumnsToRemove = columnsToRemove

						return metadata.Encode(&state)
					}
				}
			}
			return nil
		},
	}
}

func (r OutputTableResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient
			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state OutputTableResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			props := &outputs.AzureTableOutputDataSourceProperties{
				AccountName:  utils.String(state.StorageAccount),
				AccountKey:   utils.String(state.StorageAccountKey),
				Table:        utils.String(state.Table),
				PartitionKey: utils.String(state.PartitionKey),
				RowKey:       utils.String(state.RowKey),
				BatchSize:    utils.Int64(state.BatchSize),
			}

			if metadata.ResourceData.HasChange("columns_to_remove") {
				props.ColumnsToRemove = &state.ColumnsToRemove
			}

			output := outputs.Output{
				Name: utils.String(state.Name),
				Properties: &outputs.OutputProperties{
					Datasource: &outputs.AzureTableOutputDataSource{
						Properties: props,
					},
				},
			}
			var opts outputs.UpdateOperationOptions
			if _, err = client.Update(ctx, *id, output, opts); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r OutputTableResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient
			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if resp, err := client.Delete(ctx, *id); err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}

func (r OutputTableResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.StreamAnalytics.OutputsClient
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("reading %s: %+v", *id, err)
		}

		props := resp.Model.Properties
		if _, ok := props.Datasource.(outputs.AzureTableOutputDataSource); !ok {
			return fmt.Errorf("specified output is not of type")
		}
		return nil
	}
}

func (r OutputTableResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsOutputTableV0ToV1{},
		},
	}
}
