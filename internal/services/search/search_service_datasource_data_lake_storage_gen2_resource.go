// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search

import (
	"context"
	"errors"
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

var (
	_ sdk.ResourceWithUpdate         = SearchServiceDatasourceDataLakeStorageGen2Resource{}
	_ sdk.ResourceWithCustomImporter = SearchServiceDatasourceDataLakeStorageGen2Resource{}
)

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
				"the `name` must be 2-128 characters, start and end with an alphanumeric character, and contain only lowercase letters, digits, or dashes",
			),
		},

		"search_service_endpoint": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: searchValidate.SearchServiceEndpoint,
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

		"filesystem_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: storageValidate.StorageContainerName,
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
			RequiredWith: []string{"soft_delete_marker_value"},
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

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := datasources.ParseDatasourceID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.Search.SearchDataPlaneClient.DataSources.Clone(id.BaseURI)
		_, err = r.get(ctx, client, *id)
		return err
	}
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SearchServiceDatasourceDataLakeStorageGen2Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := datasources.NewDatasourceID(model.SearchServiceEndpoint, model.Name)
			client := metadata.Client.Search.SearchDataPlaneClient.DataSources.Clone(id.BaseURI)
			existing, err := r.get(ctx, client, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return err
			}

			opts := datasources.DefaultCreateOrUpdateOperationOptions()
			if response.WasNotFound(existing.HttpResponse) {
				opts.IfNoneMatch = pointer.To("*")
			} else {
				if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
				opts.IfMatch = existing.Model.OdataEtag
			}

			connectionString, err := pluginsdk.GetWriteOnly(metadata.ResourceData, "connection_string_wo", cty.String)
			if err != nil {
				return err
			}
			if connectionString.IsNull() || !connectionString.IsKnown() {
				return fmt.Errorf("`connection_string_wo` must be provided when creating %s", id)
			}

			parameters, err := model.expand()
			if err != nil {
				return err
			}
			parameters.Credentials.ConnectionString = pointer.To(connectionString.AsString())

			if _, err := client.CreateOrUpdate(ctx, id, *parameters, opts); err != nil {
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
			id, err := datasources.ParseDatasourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Search.SearchDataPlaneClient.DataSources.Clone(id.BaseURI)
			resp, err := r.get(ctx, client, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return err
			}

			model := resp.Model
			state := SearchServiceDatasourceDataLakeStorageGen2Model{
				Name:                      id.DatasourceName,
				SearchServiceEndpoint:     id.BaseURI,
				FilesystemName:            model.Container.Name,
				ContainerQuery:            pointer.From(model.Container.Query),
				Description:               pointer.From(model.Description),
				ConnectionStringWOVersion: int64(metadata.ResourceData.Get("connection_string_wo_version").(int)),
			}

			if policy, ok := model.DataDeletionDetectionPolicy.(datasources.SoftDeleteColumnDeletionDetectionPolicy); ok {
				state.SoftDeleteColumnName = pointer.From(policy.SoftDeleteColumnName)
				state.SoftDeleteMarkerValue = pointer.From(policy.SoftDeleteMarkerValue)
			}

			encryptionKey, err := schema.FlattenSearchDatasourceEncryptionKey(model.EncryptionKey, metadata.ResourceData)
			if err != nil {
				return fmt.Errorf("flattening `encryption_key`: %+v", err)
			}
			state.EncryptionKey = encryptionKey

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := datasources.ParseDatasourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Search.SearchDataPlaneClient.DataSources.Clone(id.BaseURI)
			resp, err := r.get(ctx, client, *id)
			if err != nil {
				return err
			}

			var model SearchServiceDatasourceDataLakeStorageGen2Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters, err := model.expand()
			if err != nil {
				return err
			}
			parameters.Credentials = resp.Model.Credentials
			if metadata.ResourceData.HasChange("connection_string_wo_version") {
				connectionString, err := pluginsdk.GetWriteOnly(metadata.ResourceData, "connection_string_wo", cty.String)
				if err != nil {
					return err
				}
				if connectionString.IsNull() || !connectionString.IsKnown() {
					return errors.New("`connection_string_wo` must be provided when updating `connection_string_wo_version`")
				}
				parameters.Credentials.ConnectionString = pointer.To(connectionString.AsString())
			}

			opts := datasources.CreateOrUpdateOperationOptions{
				IfMatch: resp.Model.OdataEtag,
				Prefer:  pointer.To(datasources.PreferReturnRepresentation),
			}
			if _, err := client.CreateOrUpdate(ctx, *id, *parameters, opts); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := datasources.ParseDatasourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Search.SearchDataPlaneClient.DataSources.Clone(id.BaseURI)
			resp, err := r.get(ctx, client, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}
				return err
			}

			opts := datasources.DeleteOperationOptions{IfMatch: resp.Model.OdataEtag}
			if resp, err := client.Delete(ctx, *id, opts); err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) get(ctx context.Context, client *datasources.DataSourcesClient, id datasources.DatasourceId) (datasources.GetOperationResponse, error) {
	resp, err := client.Get(ctx, id, datasources.DefaultGetOperationOptions())
	if err != nil {
		return resp, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if resp.Model == nil {
		return resp, fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if resp.Model.Type != datasources.SearchIndexerDataSourceTypeAdlsgenTwo {
		return resp, fmt.Errorf("%s has type `%s`, expected `adlsgen2`", id, resp.Model.Type)
	}
	return resp, nil
}

func (m SearchServiceDatasourceDataLakeStorageGen2Model) expand() (*datasources.SearchIndexerDataSource, error) {
	encryptionKey, err := schema.ExpandSearchDatasourceEncryptionKey(m.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("expanding `encryption_key`: %+v", err)
	}

	parameters := &datasources.SearchIndexerDataSource{
		Name: m.Name,
		Type: datasources.SearchIndexerDataSourceTypeAdlsgenTwo,
		Container: datasources.SearchIndexerDataContainer{
			Name:  m.FilesystemName,
			Query: pointer.To(m.ContainerQuery),
		},
		Description:   pointer.To(m.Description),
		EncryptionKey: encryptionKey,
	}
	if m.SoftDeleteColumnName != "" {
		parameters.DataDeletionDetectionPolicy = datasources.SoftDeleteColumnDeletionDetectionPolicy{
			SoftDeleteColumnName:  pointer.To(m.SoftDeleteColumnName),
			SoftDeleteMarkerValue: pointer.To(m.SoftDeleteMarkerValue),
		}
	}
	return parameters, nil
}
