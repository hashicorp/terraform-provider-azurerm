// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_azure_sql_table

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
				if err != nil {
					if !response.WasNotFound(existing.Response.Response) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.Response.Response) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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

			azureSqlTableset := datafactory.AzureSQLTableDataset{
				AzureSQLTableDatasetTypeProperties: &azureSqlDatasetProperties,
				LinkedServiceName:                  linkedService,
				Description:                        pointer.To(data.Description),
			}

			if data.Folder != "" {
				azureSqlTableset.Folder = &datafactory.DatasetFolder{
					Name: &data.Folder,
				}
			}

			if len(data.Parameters) > 0 {
				azureSqlTableset.Parameters = helpers.ExpandDataSetParameters(data.Parameters)
			}

			if len(data.Annotations) > 0 {
				annotations := make([]any, len(data.Annotations))
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

			dataset := datafactory.DatasetResource{
				Properties: &azureSqlTableset,
				Type:       pointer.To(string(datafactory.TypeBasicDatasetTypeAzureSQLTable)),
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
