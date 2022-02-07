package streamanalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/streamanalytics/mgmt/2020-03-01-preview/streamanalytics"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type OutputFunctionResource struct {
}

var _ sdk.ResourceWithCustomImporter = OutputFunctionResource{}

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

		"resource_group_name": azure.SchemaResourceGroupName(),

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
	return validate.OutputID
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

			id := parse.NewOutputID(subscriptionId, model.ResourceGroup, model.StreamAnalyticsJob, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.Name)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := streamanalytics.Output{
				Name: utils.String(model.Name),
				OutputProperties: &streamanalytics.OutputProperties{
					Datasource: &streamanalytics.AzureFunctionOutputDataSource{
						Type: streamanalytics.TypeMicrosoftAzureFunction,
						AzureFunctionOutputDataSourceProperties: &streamanalytics.AzureFunctionOutputDataSourceProperties{
							FunctionAppName: utils.String(model.FunctionApp),
							FunctionName:    utils.String(model.FunctionName),
							APIKey:          utils.String(model.ApiKey),
							MaxBatchSize:    utils.Float(float64(model.BatchMaxInBytes)),
							MaxBatchCount:   utils.Float(float64(model.BatchMaxCount)),
						},
					},
				},
			}

			if _, err = client.CreateOrReplace(ctx, props, id.ResourceGroup, id.StreamingjobName, id.Name, "", ""); err != nil {
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
			id, err := parse.OutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if props := resp.OutputProperties; props != nil && props.Datasource != nil {
				v, ok := props.Datasource.AsAzureFunctionOutputDataSource()
				if !ok {
					return fmt.Errorf("converting output data source to a function output: %+v", err)
				}

				if v.FunctionAppName == nil || v.FunctionName == nil || v.MaxBatchCount == nil || v.MaxBatchSize == nil {
					return nil
				}

				state := OutputFunctionResourceModel{
					Name:               id.Name,
					ResourceGroup:      id.ResourceGroup,
					StreamAnalyticsJob: id.StreamingjobName,
					FunctionApp:        *v.FunctionAppName,
					FunctionName:       *v.FunctionName,
					ApiKey:             metadata.ResourceData.Get("api_key").(string),
					BatchMaxInBytes:    int(*v.MaxBatchSize),
					BatchMaxCount:      int(*v.MaxBatchCount),
				}
				return metadata.Encode(&state)
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
			id, err := parse.OutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state OutputFunctionResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			props := streamanalytics.Output{
				Name: utils.String(state.Name),
				OutputProperties: &streamanalytics.OutputProperties{
					Datasource: &streamanalytics.AzureFunctionOutputDataSource{
						Type: streamanalytics.TypeMicrosoftStorageTable,
						AzureFunctionOutputDataSourceProperties: &streamanalytics.AzureFunctionOutputDataSourceProperties{
							FunctionAppName: utils.String(state.FunctionApp),
							FunctionName:    utils.String(state.FunctionName),
							APIKey:          utils.String(state.ApiKey),
							MaxBatchSize:    utils.Float(float64(state.BatchMaxInBytes)),
							MaxBatchCount:   utils.Float(float64(state.BatchMaxCount)),
						},
					},
				},
			}

			if _, err = client.Update(ctx, props, id.ResourceGroup, id.StreamingjobName, id.Name, ""); err != nil {
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
			id, err := parse.OutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if _, err := client.Delete(ctx, id.ResourceGroup, id.StreamingjobName, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r OutputFunctionResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := parse.OutputID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.StreamAnalytics.OutputsClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.Name)
		if err != nil || resp.OutputProperties == nil {
			return fmt.Errorf("reading %s: %+v", *id, err)
		}

		props := resp.OutputProperties
		if _, ok := props.Datasource.AsAzureFunctionOutputDataSource(); !ok {
			return fmt.Errorf("specified output is not of type %s", streamanalytics.TypeMicrosoftAzureFunction)
		}
		return nil
	}
}
