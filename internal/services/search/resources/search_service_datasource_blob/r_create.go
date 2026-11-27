// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_service_datasource_blob

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/datasources"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	searchClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/search/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SearchServiceDatasourceBlobModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			endpoint := model.SearchServiceEndpoint
			client := searchClient.From(metadata.Client).SearchDataPlaneClient.DataSources.Clone(endpoint)
			id := datasources.NewDatasourceID(endpoint, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id, datasources.DefaultGetOperationOptions())
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			woConnectionString, err := pluginsdk.GetWriteOnly(metadata.ResourceData, "connection_string_wo", cty.String)
			if err != nil {
				return err
			}

			encryptionKey, err := schema.ExpandSearchDatasourceEncryptionKey(model.EncryptionKey)
			if err != nil {
				return fmt.Errorf("expanding `encryption_key`: %+v", err)
			}

			parameters := datasources.SearchIndexerDataSource{
				Name: model.Name,
				Type: datasources.SearchIndexerDataSourceTypeAzureblob,
				Container: datasources.SearchIndexerDataContainer{
					Name: model.ContainerName,
				},
				Credentials: datasources.DataSourceCredentials{
					ConnectionString: pointer.To(woConnectionString.AsString()),
				},
				EncryptionKey: encryptionKey,
			}

			if model.Description != "" {
				parameters.Description = pointer.To(model.Description)
			}

			if model.ContainerQuery != "" {
				parameters.Container.Query = pointer.To(model.ContainerQuery)
			}

			if model.SoftDeleteColumnName != "" {
				parameters.DataDeletionDetectionPolicy = datasources.SoftDeleteColumnDeletionDetectionPolicy{
					SoftDeleteColumnName:  pointer.To(model.SoftDeleteColumnName),
					SoftDeleteMarkerValue: pointer.To(model.SoftDeleteMarkerValue),
				}
			}

			if _, err = client.CreateOrUpdate(ctx, id, parameters, datasources.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
