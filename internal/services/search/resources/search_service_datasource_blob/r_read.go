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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/search/schema"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceId, err := datasources.ParseDatasourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			endpoint := resourceId.BaseURI
			client := metadata.Client.Search.SearchDataPlaneClient.DataSources.Clone(endpoint)
			resp, err := client.Get(ctx, *resourceId, datasources.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(resourceId)
				}
				return fmt.Errorf("retrieving %s: %+v", *resourceId, err)
			}

			state := SearchServiceDatasourceBlobModel{
				Name:                  resourceId.DatasourceName,
				SearchServiceEndpoint: endpoint,
			}

			if respModel := resp.Model; respModel != nil {
				state.ContainerName = respModel.Container.Name
				state.ContainerQuery = pointer.From(respModel.Container.Query)
				state.Description = pointer.From(respModel.Description)
				if policy, ok := respModel.DataDeletionDetectionPolicy.(datasources.SoftDeleteColumnDeletionDetectionPolicy); ok {
					state.SoftDeleteColumnName = pointer.From(policy.SoftDeleteColumnName)
					state.SoftDeleteMarkerValue = pointer.From(policy.SoftDeleteMarkerValue)
				}

				encryptionKey, err := schema.FlattenSearchDatasourceEncryptionKey(respModel.EncryptionKey, metadata.ResourceData)
				if err != nil {
					return fmt.Errorf("flattening `encryption_key`: %+v", err)
				}

				state.EncryptionKey = encryptionKey
			}

			state.ConnectionStringWOVersion = int64(metadata.ResourceData.Get("connection_string_wo_version").(int))

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}
