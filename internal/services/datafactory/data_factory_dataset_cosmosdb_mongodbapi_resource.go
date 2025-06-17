package datafactory

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/datasets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ sdk.Resource = DataFactoryDatasetCosmosDbMongoDbApiResource{}

type DataFactoryDatasetCosmosDbMongoDbApiResource struct{}

type DataFactoryDatasetCosmosDbMongoDbApiResourceModel struct {
	Annotations       []string          `tfschema:"annotations"`
	CollectionName    string            `tfschema:"collection_name"`
	DataFactoryId     string            `tfschema:"data_factory_id"`
	Description       string            `tfschema:"description"`
	Folder            string            `tfschema:"folder"`
	LinkedServiceName string            `tfschema:"linked_service_name"`
	Name              string            `tfschema:"name"`
	Parameters        map[string]string `tfschema:"parameters"`
}

func (DataFactoryDatasetCosmosDbMongoDbApiResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LinkedServiceDatasetName,
		},
		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},
		"linked_service_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"annotations": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"collection_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"folder": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"parameters": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (DataFactoryDatasetCosmosDbMongoDbApiResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (DataFactoryDatasetCosmosDbMongoDbApiResource) ModelObject() interface{} {
	return &DataFactoryDatasetCosmosDbMongoDbApiResourceModel{}
}

func (DataFactoryDatasetCosmosDbMongoDbApiResource) ResourceType() string {
	return "azurerm_data_factory_dataset_cosmosdb_mongoapi"
}

func (r DataFactoryDatasetCosmosDbMongoDbApiResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.DatasetClientGoAzureSDK
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config DataFactoryDatasetCosmosDbMongoDbApiResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			dataFactoryId, err := factories.ParseFactoryID(config.DataFactoryId)
			if err != nil {
				return err
			}

			id := datasets.NewDatasetID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, config.Name)

			existing, err := client.Get(ctx, id, datasets.DefaultGetOperationOptions())
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			dataset := datasets.CosmosDbMongoDbApiCollectionDataset{
				TypeProperties: datasets.CosmosDbMongoDbApiCollectionDatasetTypeProperties{
					Collection: config.CollectionName,
				},
				LinkedServiceName: datasets.LinkedServiceReference{
					ReferenceName: config.LinkedServiceName,
					Type:          datasets.TypeLinkedServiceReference,
				},
			}

			if config.Annotations != nil {
				dataset.Annotations = pointer.To(utils.FlattenStringSlice(&config.Annotations))
			}

			if config.Description != "" {
				dataset.Description = &config.Description
			}

			if config.Folder != "" {
				dataset.Folder = &datasets.DatasetFolder{
					Name: &config.Folder,
				}
			}

			if config.Parameters != nil {
				dataset.Parameters = expandDataSetParametersGoAzureSdk(&config.Parameters)
			}

			datasetResource := datasets.DatasetResource{
				Properties: &dataset,
				Type:       pointer.To("CosmosDbMongoDbApiCollection"),
			}

			if _, err := client.CreateOrUpdate(ctx, id, datasetResource, datasets.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataFactoryDatasetCosmosDbMongoDbApiResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.DatasetClientGoAzureSDK

			var config DataFactoryDatasetCosmosDbMongoDbApiResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := datasets.ParseDatasetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, datasets.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			dataset, ok := resp.Model.Properties.(datasets.CosmosDbMongoDbApiCollectionDataset)
			if !ok {
				return fmt.Errorf("dataset %s is not a CosmosDbMongoDbAPICollectionDataset", id)
			}

			if metadata.ResourceData.HasChange("annotations") {
				dataset.Annotations = pointer.To(utils.FlattenStringSlice(&config.Annotations))
			}

			if metadata.ResourceData.HasChange("collection_name") {
				dataset.TypeProperties.Collection = config.CollectionName
			}

			if metadata.ResourceData.HasChange("description") {
				dataset.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("folder") {
				dataset.Folder = &datasets.DatasetFolder{
					Name: &config.Folder,
				}
			}

			if metadata.ResourceData.HasChange("linked_service_name") {
				dataset.LinkedServiceName = datasets.LinkedServiceReference{
					ReferenceName: config.LinkedServiceName,
					Type:          datasets.TypeLinkedServiceReference,
				}
			}

			if metadata.ResourceData.HasChange("parameters") {
				dataset.Parameters = expandDataSetParametersGoAzureSdk(&config.Parameters)
			}

			resp.Model.Properties = dataset

			if _, err := client.CreateOrUpdate(ctx, *id, *resp.Model, datasets.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (DataFactoryDatasetCosmosDbMongoDbApiResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.DatasetClientGoAzureSDK

			id, err := datasets.ParseDatasetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, datasets.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			dataset, ok := resp.Model.Properties.(datasets.CosmosDbMongoDbApiCollectionDataset)
			if !ok {
				return fmt.Errorf("dataset %s is not a CosmosDbMongoDbAPICollectionDataset", id)
			}

			state := DataFactoryDatasetCosmosDbMongoDbApiResourceModel{}

			if dataset.Annotations != nil {
				state.Annotations = flattenDataFactoryAnnotations(dataset.Annotations)
			}

			state.CollectionName = dataset.TypeProperties.Collection.(string)

			state.DataFactoryId = datasets.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName).ID()

			if dataset.Description != nil {
				state.Description = *dataset.Description
			}

			if dataset.Folder != nil {
				if dataset.Folder.Name != nil {
					state.Folder = *dataset.Folder.Name
				}
			}

			state.LinkedServiceName = dataset.LinkedServiceName.ReferenceName

			state.Name = id.DatasetName

			if dataset.Parameters != nil {
				state.Parameters = flattenDataSetParametersGoAzureSdk(dataset.Parameters)
			}

			return metadata.Encode(&state)
		},
	}
}

func (DataFactoryDatasetCosmosDbMongoDbApiResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.DatasetClientGoAzureSDK

			id, err := datasets.ParseDatasetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (DataFactoryDatasetCosmosDbMongoDbApiResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.LinkedServiceDatasetName
}
