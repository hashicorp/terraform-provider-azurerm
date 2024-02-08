// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

var _ sdk.Resource = DataFactoryDatasetAzureSQLTableResource{}

type DataFactoryDatasetAzureSQLTableResource struct{}

type DataFactoryDatasetAzureSQLTableResourceSchema struct {
	Name                 string                 `tfschema:"name"`
	DataFactoryId        string                 `tfschema:"data_factory_id"`
	LinkedServiceId      string                 `tfschema:"linked_service_id"`
	Schema               string                 `tfschema:"schema"`
	Table                string                 `tfschema:"table"`
	Parameters           map[string]interface{} `tfschema:"parameters"`
	Description          string                 `tfschema:"description"`
	Annotations          []string               `tfschema:"annotations"`
	Folder               string                 `tfschema:"folder"`
	AdditionalProperties map[string]interface{} `tfschema:"additional_properties"`
	SchemaColumn         []DatasetColumn        `tfschema:"schema_column"`
}

func (DataFactoryDatasetAzureSQLTableResource) Arguments() map[string]*pluginsdk.Schema {
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

		"linked_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LinkedServiceID,
		},

		"schema": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"table": {
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

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"annotations": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"folder": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"additional_properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"schema_column": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Byte",
							"Byte[]",
							"Boolean",
							"Date",
							"DateTime",
							"DateTimeOffset",
							"Decimal",
							"Double",
							"Guid",
							"Int16",
							"Int32",
							"Int64",
							"Single",
							"String",
							"TimeSpan",
						}, false),
					},
					"description": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func (DataFactoryDatasetAzureSQLTableResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (DataFactoryDatasetAzureSQLTableResource) ModelObject() interface{} {
	return &DataFactoryDatasetAzureSQLTableResourceSchema{}
}

func (DataFactoryDatasetAzureSQLTableResource) ResourceType() string {
	return "azurerm_data_factory_dataset_azure_sql_table"
}

func (r DataFactoryDatasetAzureSQLTableResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.DatasetClient
			subscriptionId := client.SubscriptionID
			var data DataFactoryDatasetAzureSQLTableResourceSchema
			if err := metadata.Decode(&data); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			dataFactoryId, err := factories.ParseFactoryID(data.DataFactoryId)
			if err != nil {
				return err
			}

			id := parse.NewDataSetID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, data.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			azureSqlDatasetProperties := datafactory.AzureSQLTableDatasetTypeProperties{
				Schema: data.Schema,
				Table:  data.Table,
			}

			linkedServiceId, err := parse.LinkedServiceID(data.LinkedServiceId)
			if err != nil {
				return err
			}
			if linkedServiceId.SubscriptionId != id.SubscriptionId || linkedServiceId.ResourceGroup != id.ResourceGroup || linkedServiceId.FactoryName != id.FactoryName {
				return fmt.Errorf("checking the linked service %s: not within the same data factory as this dataset %s", data.LinkedServiceId, id.ID())
			}
			linkedService := &datafactory.LinkedServiceReference{
				Type:          pointer.To("LinkedServiceReference"),
				ReferenceName: pointer.To(linkedServiceId.Name),
			}

			description := data.Description
			azureSqlTableset := datafactory.AzureSQLTableDataset{
				AzureSQLTableDatasetTypeProperties: &azureSqlDatasetProperties,
				LinkedServiceName:                  linkedService,
				Description:                        &description,
			}

			if data.Folder != "" {
				azureSqlTableset.Folder = &datafactory.DatasetFolder{
					Name: &data.Folder,
				}
			}

			if len(data.Parameters) > 0 {
				azureSqlTableset.Parameters = expandDataSetParameters(data.Parameters)
			}

			if len(data.Annotations) > 0 {
				annotations := make([]interface{}, len(data.Annotations))
				for i, v := range data.Annotations {
					annotations[i] = v
				}
				azureSqlTableset.Annotations = &annotations
			}

			if len(data.AdditionalProperties) > 0 {
				azureSqlTableset.AdditionalProperties = data.AdditionalProperties
			}

			if len(data.SchemaColumn) > 0 {
				azureSqlTableset.Structure = data.SchemaColumn
			}

			datasetType := string(datafactory.TypeBasicDatasetTypeAzureSQLTable)
			dataset := datafactory.DatasetResource{
				Properties: &azureSqlTableset,
				Type:       &datasetType,
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataFactoryDatasetAzureSQLTableResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.DatasetClient
			id, err := parse.DataSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			var data DataFactoryDatasetAzureSQLTableResourceSchema
			if err := metadata.Decode(&data); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			dataset, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", id, err)
			}

			azureSqlTable, ok := dataset.Properties.AsAzureSQLTableDataset()
			if !ok {
				return fmt.Errorf("classifying Data Factory Dataset Azure SQL Table %s: Expected: %q Received: %T", *id, datafactory.TypeBasicDatasetTypeAzureSQLTable, dataset.Properties)
			}

			if metadata.ResourceData.HasChanges("schema", "table") {
				if azureSqlTable.AzureSQLTableDatasetTypeProperties == nil {
					azureSqlTable.AzureSQLTableDatasetTypeProperties = &datafactory.AzureSQLTableDatasetTypeProperties{}
				}

				if metadata.ResourceData.HasChange("schema") {
					azureSqlTable.AzureSQLTableDatasetTypeProperties.Schema = data.Schema
				}

				if metadata.ResourceData.HasChange("table") {
					azureSqlTable.AzureSQLTableDatasetTypeProperties.Table = data.Table
				}
			}

			if metadata.ResourceData.HasChange("linked_service_id") {
				linkedServiceId, err := parse.LinkedServiceID(data.LinkedServiceId)
				if err != nil {
					return err
				}
				if linkedServiceId.SubscriptionId != id.SubscriptionId || linkedServiceId.ResourceGroup != id.ResourceGroup || linkedServiceId.FactoryName != id.FactoryName {
					return fmt.Errorf("checking the linked service %s: not within the same data factory as this dataset %s", data.LinkedServiceId, id.ID())
				}
				azureSqlTable.LinkedServiceName = &datafactory.LinkedServiceReference{
					Type:          pointer.To("LinkedServiceReference"),
					ReferenceName: pointer.To(linkedServiceId.Name),
				}
			}

			if metadata.ResourceData.HasChange("description") {
				azureSqlTable.Description = pointer.To(data.Description)
			}

			if metadata.ResourceData.HasChange("folder") {
				if data.Folder != "" {
					azureSqlTable.Folder = &datafactory.DatasetFolder{
						Name: &data.Folder,
					}
				} else {
					azureSqlTable.Folder = nil
				}
			}

			if metadata.ResourceData.HasChange("parameters") {
				azureSqlTable.Parameters = expandDataSetParameters(data.Parameters)
			}

			if metadata.ResourceData.HasChange("annotations") {
				if len(data.Annotations) > 0 {
					annotations := make([]interface{}, len(data.Annotations))
					for i, v := range data.Annotations {
						annotations[i] = v
					}
					azureSqlTable.Annotations = &annotations
				} else {
					azureSqlTable.Annotations = nil
				}
			}

			if metadata.ResourceData.HasChange("additional_properties") {
				azureSqlTable.AdditionalProperties = data.AdditionalProperties
			}

			if metadata.ResourceData.HasChange("schema_column") {
				azureSqlTable.Structure = data.SchemaColumn
			}

			dataset = datafactory.DatasetResource{
				Type:       pointer.To(string(datafactory.TypeBasicDatasetTypeAzureSQLTable)),
				Properties: azureSqlTable,
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (DataFactoryDatasetAzureSQLTableResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			d := metadata.ResourceData
			client := metadata.Client.DataFactory.DatasetClient
			id, err := parse.DataSetID(d.Id())
			if err != nil {
				return err
			}

			dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)
			var state DataFactoryDatasetAzureSQLTableResourceSchema

			resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.Name = id.Name
			state.DataFactoryId = dataFactoryId.ID()

			azureSqlTable, ok := resp.Properties.AsAzureSQLTableDataset()
			if !ok {
				return fmt.Errorf("classifying Data Factory Dataset Azure SQL Table %s: Expected: %q Received: %T", *id, datafactory.TypeBasicDatasetTypeAzureSQLTable, resp.Properties)
			}

			state.Description = pointer.From(azureSqlTable.Description)
			state.AdditionalProperties = azureSqlTable.AdditionalProperties

			state.Parameters = flattenDataSetParameters(azureSqlTable.Parameters)
			state.Annotations = flattenDataFactoryAnnotations(azureSqlTable.Annotations)

			if linkedService := azureSqlTable.LinkedServiceName; linkedService != nil && linkedService.ReferenceName != nil {
				state.LinkedServiceId = parse.NewLinkedServiceID(id.SubscriptionId, id.ResourceGroup, id.FactoryName, *linkedService.ReferenceName).ID()
			}

			if properties := azureSqlTable.AzureSQLTableDatasetTypeProperties; properties != nil {
				if val, ok := properties.Schema.(string); ok {
					state.Schema = val
				} else {
					state.Schema = ""
					log.Printf("[DEBUG] Skipping `schema` since it's not a string")
				}

				if val, ok := properties.Table.(string); ok {
					state.Table = val
				} else {
					state.Table = ""
					log.Printf("[DEBUG] Skipping `table` since it's not a string")
				}
			}

			state.Folder = ""
			if folder := azureSqlTable.Folder; folder != nil && folder.Name != nil {
				state.Folder = pointer.From(folder.Name)
			}

			state.SchemaColumn = flattenDataFactoryStructureColumnsToDatasetColumn(azureSqlTable.Structure)

			return metadata.Encode(&state)
		},
	}
}

func (DataFactoryDatasetAzureSQLTableResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			d := metadata.ResourceData
			client := metadata.Client.DataFactory.DatasetClient

			id, err := parse.DataSetID(d.Id())
			if err != nil {
				return err
			}

			response, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(response) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (DataFactoryDatasetAzureSQLTableResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataSetID
}
