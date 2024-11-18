package datafactory

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

var _ sdk.Resource = DataFactoryDatasetCosmosDbMongoDbResource{}

type DataFactoryDatasetCosmosDbMongoDbResource struct{}

type DataFactoryDatasetCosmosDbMongoDbResourceModel struct {
	AdditionalProperties map[string]interface{} `tfschema:"additional_properties"`
	Annotations          []string               `tfschema:"annotations"`
	CollectionName       string                 `tfschema:"collection_name"`
	DataFactoryId        string                 `tfschema:"data_factory_id"`
	Description          string                 `tfschema:"description"`
	Folder               string                 `tfschema:"folder"`
	LinkedServiceName    string                 `tfschema:"linked_service_name"`
	Name                 string                 `tfschema:"name"`
	Parameters           map[string]interface{} `tfschema:"parameters"`
}

func (DataFactoryDatasetCosmosDbMongoDbResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"additional_properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
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
		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
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
		"linked_service_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LinkedServiceDatasetName,
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

func (DataFactoryDatasetCosmosDbMongoDbResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (DataFactoryDatasetCosmosDbMongoDbResource) ModelObject() interface{} {
	return &DataFactoryDatasetCosmosDbMongoDbResourceModel{}
}

func (DataFactoryDatasetCosmosDbMongoDbResource) ResourceType() string {
	return "azurerm_data_factory_dataset_cosmosdb_mongoapi"
}

func (r DataFactoryDatasetCosmosDbMongoDbResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.DatasetClient
			subscriptionId := metadata.Client.DataFactory.DatasetClient.SubscriptionID

			var config DataFactoryDatasetCosmosDbMongoDbResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			dataFactoryId, err := factories.ParseFactoryID(config.DataFactoryId)
			if err != nil {
				return err
			}

			id := parse.NewDataSetID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, config.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			datasetProperties := datafactory.CosmosDbMongoDbAPICollectionDataset{
				CosmosDbMongoDbAPICollectionDatasetTypeProperties: &datafactory.CosmosDbMongoDbAPICollectionDatasetTypeProperties{
					Collection: config.CollectionName,
				},
				LinkedServiceName: &datafactory.LinkedServiceReference{
					ReferenceName: &config.LinkedServiceName,
					Type:          pointer.To("LinkedServiceReference"),
				},
			}

			if config.AdditionalProperties != nil {
				datasetProperties.AdditionalProperties = config.AdditionalProperties
			}

			if config.Annotations != nil {
				datasetProperties.Annotations = expandDataFactoryAnnotations(&config.Annotations)
			}

			if config.Description != "" {
				datasetProperties.Description = &config.Description
			}

			if config.Folder != "" {
				datasetProperties.Folder = &datafactory.DatasetFolder{
					Name: &config.Folder,
				}
			}

			if config.Parameters != nil {
				datasetProperties.Parameters = expandDataSetParameters(config.Parameters)
			}

			dataset := datafactory.DatasetResource{
				Properties: &datasetProperties,
				Type:       pointer.To(string(datafactory.TypeBasicDatasetTypeCosmosDbMongoDbAPICollection)),
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataFactoryDatasetCosmosDbMongoDbResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.DatasetClient

			var config DataFactoryDatasetCosmosDbMongoDbResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := parse.DataSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			dataset, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			datasetProperties, ok := dataset.Properties.AsCosmosDbMongoDbAPICollectionDataset()
			if !ok {
				return fmt.Errorf("dataset %s is not a CosmosDbMongoDbAPICollectionDataset", id)
			}

			if metadata.ResourceData.HasChange("additional_properties") {
				datasetProperties.AdditionalProperties = config.AdditionalProperties
			}

			if metadata.ResourceData.HasChange("annotations") {
				datasetProperties.Annotations = expandDataFactoryAnnotations(&config.Annotations)
			}

			if metadata.ResourceData.HasChange("collection_name") {
				datasetProperties.CosmosDbMongoDbAPICollectionDatasetTypeProperties.Collection = config.CollectionName
			}

			if metadata.ResourceData.HasChange("description") {
				datasetProperties.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("folder") {
				datasetProperties.Folder = &datafactory.DatasetFolder{
					Name: &config.Folder,
				}
			}

			if metadata.ResourceData.HasChange("linked_service_name") {
				datasetProperties.LinkedServiceName = &datafactory.LinkedServiceReference{
					ReferenceName: &config.LinkedServiceName,
					Type:          pointer.To("LinkedServiceReference"),
				}
			}

			if metadata.ResourceData.HasChange("parameters") {
				datasetProperties.Parameters = expandDataSetParameters(config.Parameters)
			}

			dataset.Properties = datasetProperties

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (DataFactoryDatasetCosmosDbMongoDbResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.DatasetClient

			id, err := parse.DataSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			dataset, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if utils.ResponseWasNotFound(dataset.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			datasetProperties, ok := dataset.Properties.AsCosmosDbMongoDbAPICollectionDataset()
			if !ok {
				return fmt.Errorf("dataset %s is not a CosmosDbMongoDbAPICollectionDataset", id)
			}

			state := DataFactoryDatasetCosmosDbMongoDbResourceModel{}

			state.AdditionalProperties = datasetProperties.AdditionalProperties

			if datasetProperties.Annotations != nil {
				state.Annotations = flattenDataFactoryAnnotations(datasetProperties.Annotations)
			}

			if collectionName, ok := datasetProperties.Collection.(string); ok {
				state.CollectionName = collectionName
			}

			state.DataFactoryId = factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName).ID()

			if datasetProperties.Description != nil {
				state.Description = *datasetProperties.Description
			}

			if datasetProperties.Folder != nil {
				if datasetProperties.Folder.Name != nil {
					state.Folder = *datasetProperties.Folder.Name
				}
			}

			if datasetProperties.LinkedServiceName != nil {
				if datasetProperties.LinkedServiceName.ReferenceName != nil {
					state.LinkedServiceName = *datasetProperties.LinkedServiceName.ReferenceName
				}
			}

			state.Name = id.Name

			state.Parameters = flattenDataSetParameters(datasetProperties.Parameters)

			return metadata.Encode(&state)
		},
	}
}

func (DataFactoryDatasetCosmosDbMongoDbResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.DatasetClient

			id, err := parse.DataSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (DataFactoryDatasetCosmosDbMongoDbResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.LinkedServiceDatasetName
}
