// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ sdk.Resource = DataFactoryDatasetAzureSQLTableResource{}

type DataFactoryDatasetAzureSQLTableResource struct{}

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

		"linked_service_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"table_name": {
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
	return nil
}

func (DataFactoryDatasetAzureSQLTableResource) ResourceType() string {
	return "azurerm_data_factory_dataset_azure_sql_table"
}

func (r DataFactoryDatasetAzureSQLTableResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			d := metadata.ResourceData
			client := metadata.Client.DataFactory.DatasetClient
			subscriptionId := client.SubscriptionID

			dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
			if err != nil {
				return err
			}

			id := parse.NewDataSetID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

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
				TableName: d.Get("table_name").(string),
			}

			linkedServiceName := d.Get("linked_service_name").(string)
			linkedServiceType := "LinkedServiceReference"
			linkedService := &datafactory.LinkedServiceReference{
				ReferenceName: &linkedServiceName,
				Type:          &linkedServiceType,
			}

			description := d.Get("description").(string)
			azureSqlTableset := datafactory.AzureSQLTableDataset{
				AzureSQLTableDatasetTypeProperties: &azureSqlDatasetProperties,
				LinkedServiceName:                  linkedService,
				Description:                        &description,
			}

			if v, ok := d.GetOk("folder"); ok {
				name := v.(string)
				azureSqlTableset.Folder = &datafactory.DatasetFolder{
					Name: &name,
				}
			}

			if v, ok := d.GetOk("parameters"); ok {
				azureSqlTableset.Parameters = expandDataSetParameters(v.(map[string]interface{}))
			}

			if v, ok := d.GetOk("annotations"); ok {
				annotations := v.([]interface{})
				azureSqlTableset.Annotations = &annotations
			}

			if v, ok := d.GetOk("additional_properties"); ok {
				azureSqlTableset.AdditionalProperties = v.(map[string]interface{})
			}

			if v, ok := d.GetOk("schema_column"); ok {
				azureSqlTableset.Structure = expandDataFactoryDatasetStructure(v.([]interface{}))
			}

			datasetType := string(datafactory.TypeBasicDatasetTypeAzureSQLTable)
			dataset := datafactory.DatasetResource{
				Properties: &azureSqlTableset,
				Type:       &datasetType,
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			d.SetId(id.ID())
			return nil
		},
	}
}

func (r DataFactoryDatasetAzureSQLTableResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			d := metadata.ResourceData
			client := metadata.Client.DataFactory.DatasetClient
			id, err := parse.DataSetID(d.Id())
			if err != nil {
				return err
			}

			_, err = client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", id, err)
			}

			// TO-DO: use d.HasChanges
			azureSqlDatasetProperties := datafactory.AzureSQLTableDatasetTypeProperties{
				TableName: d.Get("table_name").(string),
			}

			linkedServiceName := d.Get("linked_service_name").(string)
			linkedServiceType := "LinkedServiceReference"
			linkedService := &datafactory.LinkedServiceReference{
				ReferenceName: &linkedServiceName,
				Type:          &linkedServiceType,
			}

			description := d.Get("description").(string)
			azureSqlTableset := datafactory.AzureSQLTableDataset{
				AzureSQLTableDatasetTypeProperties: &azureSqlDatasetProperties,
				LinkedServiceName:                  linkedService,
				Description:                        &description,
			}

			if v, ok := d.GetOk("folder"); ok {
				name := v.(string)
				azureSqlTableset.Folder = &datafactory.DatasetFolder{
					Name: &name,
				}
			}

			if v, ok := d.GetOk("parameters"); ok {
				azureSqlTableset.Parameters = expandDataSetParameters(v.(map[string]interface{}))
			}

			if v, ok := d.GetOk("annotations"); ok {
				annotations := v.([]interface{})
				azureSqlTableset.Annotations = &annotations
			}

			if v, ok := d.GetOk("additional_properties"); ok {
				azureSqlTableset.AdditionalProperties = v.(map[string]interface{})
			}

			if v, ok := d.GetOk("schema_column"); ok {
				azureSqlTableset.Structure = expandDataFactoryDatasetStructure(v.([]interface{}))
			}

			datasetType := string(datafactory.TypeBasicDatasetTypeAzureSQLTable)
			dataset := datafactory.DatasetResource{
				Properties: &azureSqlTableset,
				Type:       &datasetType,
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			d.SetId(id.ID())
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

			resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			d.Set("name", id.Name)
			d.Set("data_factory_id", dataFactoryId.ID())

			azureSqlTable, ok := resp.Properties.AsAzureSQLTableDataset()
			if !ok {
				return fmt.Errorf("classifying Data Factory Dataset Azure SQL Table %s: Expected: %q Received: %T", *id, datafactory.TypeBasicDatasetTypeAzureSQLTable, resp.Properties)
			}

			d.Set("additional_properties", azureSqlTable.AdditionalProperties)

			if azureSqlTable.Description != nil {
				d.Set("description", azureSqlTable.Description)
			}

			parameters := flattenDataSetParameters(azureSqlTable.Parameters)
			if err := d.Set("parameters", parameters); err != nil {
				return fmt.Errorf("setting `parameters`: %+v", err)
			}

			annotations := flattenDataFactoryAnnotations(azureSqlTable.Annotations)
			if err := d.Set("annotations", annotations); err != nil {
				return fmt.Errorf("setting `annotations`: %+v", err)
			}

			if linkedService := azureSqlTable.LinkedServiceName; linkedService != nil {
				if linkedService.ReferenceName != nil {
					d.Set("linked_service_name", linkedService.ReferenceName)
				}
			}

			if properties := azureSqlTable.AzureSQLTableDatasetTypeProperties; properties != nil {
				val, ok := properties.TableName.(string)
				if !ok {
					log.Printf("[DEBUG] Skipping `table_name` since it's not a string")
				} else {
					d.Set("table_name", val)
				}
			}

			if folder := azureSqlTable.Folder; folder != nil {
				if folder.Name != nil {
					d.Set("folder", folder.Name)
				}
			}

			structureColumns := flattenDataFactoryStructureColumns(azureSqlTable.Structure)
			if err := d.Set("schema_column", structureColumns); err != nil {
				return fmt.Errorf("setting `schema_column`: %+v", err)
			}

			return nil
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
	return func(input interface{}, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		if _, err := parse.DataSetID(v); err != nil {
			errors = append(errors, err)
		}

		return
	}
}
