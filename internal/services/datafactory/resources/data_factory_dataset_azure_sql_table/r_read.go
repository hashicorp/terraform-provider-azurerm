// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_azure_sql_table

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func (Resource) Read() sdk.ResourceFunc {
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
				if response.WasNotFound(resp.Response.Response) {
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

			state.Parameters = helpers.FlattenDataSetParameters(azureSqlTable.Parameters)
			state.Annotations = helpers.FlattenDataFactoryAnnotations(azureSqlTable.Annotations)

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

			state.SchemaColumn = helpers.FlattenDataFactoryStructureColumnsToDatasetColumn(azureSqlTable.Structure)

			return metadata.Encode(&state)
		},
	}
}

type DataFactoryDatasetAzureSQLTableResourceSchema struct {
	Name                 string                  `tfschema:"name"`
	DataFactoryId        string                  `tfschema:"data_factory_id"`
	LinkedServiceId      string                  `tfschema:"linked_service_id"`
	Schema               string                  `tfschema:"schema"`
	Table                string                  `tfschema:"table"`
	Parameters           map[string]any          `tfschema:"parameters"`
	Description          string                  `tfschema:"description"`
	Annotations          []string                `tfschema:"annotations"`
	Folder               string                  `tfschema:"folder"`
	AdditionalProperties map[string]any          `tfschema:"additional_properties"`
	SchemaColumn         []helpers.DatasetColumn `tfschema:"schema_column"`
}
