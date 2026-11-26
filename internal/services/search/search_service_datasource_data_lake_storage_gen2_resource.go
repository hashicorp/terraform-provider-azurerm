// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/datasources"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/search/schema"
	searchValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SearchServiceDatasourceDataLakeStorageGen2Resource struct{}

var _ sdk.ResourceWithUpdate = SearchServiceDatasourceDataLakeStorageGen2Resource{}

type SearchServiceDatasourceDataLakeStorageGen2Model struct {
	Name                      string                                      `tfschema:"name"`
	SearchServiceEndpoint     string                                      `tfschema:"search_service_endpoint"`
	FilesystemName            string                                      `tfschema:"filesystem_name"`
	ConnectionStringWOVersion int64                                       `tfschema:"connection_string_wo_version"`
	Description               string                                      `tfschema:"description"`
	ContainerQuery            string                                      `tfschema:"container_query"`
	SoftDeleteColumnName      string                                      `tfschema:"soft_delete_column_name"`
	SoftDeleteMarkerValue     string                                      `tfschema:"soft_delete_marker_value"`
	EncryptionKey             []schema.SearchDatasourceEncryptionKeyModel `tfschema:"encryption_key"`
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,126}[a-z0-9]$`),
				"The `name` must be 2-128 characters, start and end with an alphanumeric character, and contain only lowercase letters, digits, or dashes",
			),
		},

		"search_service_endpoint": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: searchValidate.SearchServiceEndpoint,
		},

		"filesystem_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: storageValidate.StorageContainerName,
		},

		"connection_string_wo": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			WriteOnly:    true,
			ValidateFunc: searchValidate.SearchDatasourceStorageConnectionString,
		},

		"connection_string_wo_version": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"container_query": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"encryption_key": schema.SearchDatasourceEncryptionKeySchema(),

		"soft_delete_column_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"soft_delete_marker_value": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"soft_delete_column_name"},
		},
	}
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) ResourceType() string {
	return "azurerm_search_service_datasource_data_lake_storage_gen2"
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) ModelObject() any {
	return &SearchServiceDatasourceDataLakeStorageGen2Model{}
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datasources.ValidateDatasourceID
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SearchServiceDatasourceDataLakeStorageGen2Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			endpoint := model.SearchServiceEndpoint
			client := metadata.Client.Search.SearchDataPlaneClient.DataSources.Clone(endpoint)
			id := datasources.NewDatasourceID(endpoint, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
				Type: datasources.SearchIndexerDataSourceTypeAdlsgenTwo,
				Container: datasources.SearchIndexerDataContainer{
					Name: model.FilesystemName,
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

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) Read() sdk.ResourceFunc {
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

			state := SearchServiceDatasourceDataLakeStorageGen2Model{
				Name:                  resourceId.DatasourceName,
				SearchServiceEndpoint: endpoint,
			}

			if respModel := resp.Model; respModel != nil {
				state.FilesystemName = respModel.Container.Name
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

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceId, err := datasources.ParseDatasourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			endpoint := resourceId.BaseURI

			client := metadata.Client.Search.SearchDataPlaneClient.DataSources.Clone(endpoint)

			var state SearchServiceDatasourceDataLakeStorageGen2Model
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
			if metadata.ResourceData.HasChange("filesystem_name") {
				existing.Container.Name = state.FilesystemName
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

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceId, err := datasources.ParseDatasourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Search.SearchDataPlaneClient.DataSources.Clone(resourceId.BaseURI)

			if _, err := client.Delete(ctx, *resourceId, datasources.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *resourceId, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
