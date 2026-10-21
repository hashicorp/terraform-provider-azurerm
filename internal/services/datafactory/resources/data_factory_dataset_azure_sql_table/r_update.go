// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_azure_sql_table

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func (r Resource) Update() sdk.ResourceFunc {
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
					azureSqlTable.Table = data.Table
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
				azureSqlTable.Parameters = helpers.ExpandDataSetParameters(data.Parameters)
			}

			if metadata.ResourceData.HasChange("annotations") {
				if len(data.Annotations) > 0 {
					annotations := make([]any, len(data.Annotations))
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
