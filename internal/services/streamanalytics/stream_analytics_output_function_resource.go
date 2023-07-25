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

type OutputFunctionResource struct {
}

var (
	_ sdk.ResourceWithCustomImporter = OutputFunctionResource{}
	_ sdk.ResourceWithStateMigration = OutputFunctionResource{}
)

type OutputFunctionResourceModel struct {
	Name               string `tfschema:"name"`
	StreamAnalyticsJob string `tfschema:"stream_analytics_job_name"`
	ResourceGroup      string `tfschema:"resource_group_name"`
	FunctionApp        string `tfschema:"function_app"`
	FunctionName       string `tfschema:"function_name"`
	ApiKey             string `tfschema:"api_key"`
	BatchMaxInBytes    int    `tfschema:"batch_max_in_bytes"`
	BatchMaxCount      int    `tfschema:"batch_max_count"`
}

func (r OutputFunctionResource) Arguments() map[string]*pluginsdk.Schema {
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

		"function_app": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"function_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"api_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"batch_max_in_bytes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  262144, // 256kB
		},

		"batch_max_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  100,
		},
	}
}

func (r OutputFunctionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r OutputFunctionResource) ModelObject() interface{} {
	return &OutputFunctionResourceModel{}
}

func (r OutputFunctionResource) ResourceType() string {
	return "azurerm_stream_analytics_output_function"
}

func (r OutputFunctionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return outputs.ValidateOutputID
}

func (r OutputFunctionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model OutputFunctionResourceModel
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

			props := outputs.Output{
				Name: utils.String(model.Name),
				Properties: &outputs.OutputProperties{
					Datasource: &outputs.AzureFunctionOutputDataSource{
						Properties: &outputs.AzureFunctionOutputDataSourceProperties{
							FunctionAppName: utils.String(model.FunctionApp),
							FunctionName:    utils.String(model.FunctionName),
							ApiKey:          utils.String(model.ApiKey),
							MaxBatchSize:    utils.Float(float64(model.BatchMaxInBytes)),
							MaxBatchCount:   utils.Float(float64(model.BatchMaxCount)),
						},
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

func (r OutputFunctionResource) Read() sdk.ResourceFunc {
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
					output, ok := props.Datasource.(outputs.AzureFunctionOutputDataSource)
					if !ok {
						return fmt.Errorf("converting %s to a Function Output", *id)
					}

					if output.Properties != nil {
						if output.Properties.FunctionAppName == nil || output.Properties.FunctionName == nil || output.Properties.MaxBatchCount == nil || output.Properties.MaxBatchSize == nil {
							return nil
						}

						state := OutputFunctionResourceModel{
							Name:               id.OutputName,
							ResourceGroup:      id.ResourceGroupName,
							StreamAnalyticsJob: id.StreamingJobName,
							ApiKey:             metadata.ResourceData.Get("api_key").(string),
						}

						functionApp := ""
						if v := output.Properties.FunctionAppName; v != nil {
							functionApp = *v
						}
						state.FunctionApp = functionApp

						functionName := ""
						if v := output.Properties.FunctionName; v != nil {
							functionName = *v
						}
						state.FunctionName = functionName

						batchMaxInBytes := 0
						if v := output.Properties.MaxBatchSize; v != nil {
							batchMaxInBytes = int(*v)
						}
						state.BatchMaxInBytes = batchMaxInBytes

						batchMaxCount := 0
						if v := output.Properties.MaxBatchCount; v != nil {
							batchMaxCount = int(*v)
						}
						state.BatchMaxCount = batchMaxCount

						return metadata.Encode(&state)
					}
				}
			}
			return nil
		},
	}
}

func (r OutputFunctionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient
			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state OutputFunctionResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			props := outputs.Output{
				Name: utils.String(state.Name),
				Properties: &outputs.OutputProperties{
					Datasource: &outputs.AzureFunctionOutputDataSource{
						Properties: &outputs.AzureFunctionOutputDataSourceProperties{
							FunctionAppName: utils.String(state.FunctionApp),
							FunctionName:    utils.String(state.FunctionName),
							ApiKey:          utils.String(state.ApiKey),
							MaxBatchSize:    utils.Float(float64(state.BatchMaxInBytes)),
							MaxBatchCount:   utils.Float(float64(state.BatchMaxCount)),
						},
					},
				},
			}

			var opts outputs.UpdateOperationOptions
			if _, err = client.Update(ctx, *id, props, opts); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r OutputFunctionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient
			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r OutputFunctionResource) CustomImporter() sdk.ResourceRunFunc {
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
		if _, ok := props.Datasource.(outputs.AzureFunctionOutputDataSource); !ok {
			return fmt.Errorf("specified output is not of type")
		}
		return nil
	}
}

func (r OutputFunctionResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsOutputFunctionV0ToV1{},
		},
	}
}
