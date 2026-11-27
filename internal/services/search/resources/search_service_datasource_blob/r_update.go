// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_service_datasource_blob

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/datasources"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/search/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceId, err := datasources.ParseDatasourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			endpoint := resourceId.BaseURI

			client := metadata.Client.Search.SearchDataPlaneClient.DataSources.Clone(endpoint)

			var state SearchServiceDatasourceBlobModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *resourceId, datasources.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *resourceId, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *resourceId)
			}

			existing := *resp.Model
			if metadata.ResourceData.HasChange("container_name") {
				existing.Container.Name = state.ContainerName
			}

			if metadata.ResourceData.HasChange("container_query") {
				existing.Container.Query = pointer.To(state.ContainerQuery)
			}

			if metadata.ResourceData.HasChange("connection_string_wo_version") {
				woConnectionString, err := pluginsdk.GetWriteOnly(metadata.ResourceData, "connection_string_wo", cty.String)
				if err != nil {
					return err
				}
				if !woConnectionString.IsNull() {
					existing.Credentials.ConnectionString = pointer.To(woConnectionString.AsString())
				}
			}

			if metadata.ResourceData.HasChange("description") {
				existing.Description = pointer.To(state.Description)
			}

			if metadata.ResourceData.HasChange("encryption_key") {
				encryptionKey, err := schema.ExpandSearchDatasourceEncryptionKey(state.EncryptionKey)
				if err != nil {
					return fmt.Errorf("expanding `encryption_key`: %+v", err)
				}

				existing.EncryptionKey = encryptionKey
			}

			// The SDK populates zero-value RawData*DetectionPolicyImpl structs from the GET response
			// that serialize to invalid abstract-type payloads the API rejects. Nil them out and
			// rebuild DataDeletionDetectionPolicy from current state.
			existing.DataChangeDetectionPolicy = nil
			existing.DataDeletionDetectionPolicy = nil
			if state.SoftDeleteColumnName != "" {
				existing.DataDeletionDetectionPolicy = datasources.SoftDeleteColumnDeletionDetectionPolicy{
					SoftDeleteColumnName:  pointer.To(state.SoftDeleteColumnName),
					SoftDeleteMarkerValue: pointer.To(state.SoftDeleteMarkerValue),
				}
			}

			opts := datasources.CreateOrUpdateOperationOptions{
				Prefer: pointer.To(datasources.PreferReturnRepresentation),
			}

			if _, err := client.CreateOrUpdate(ctx, *resourceId, existing, opts); err != nil {
				return fmt.Errorf("updating %s: %+v", *resourceId, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
