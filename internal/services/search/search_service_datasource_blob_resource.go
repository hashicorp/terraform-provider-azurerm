// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	searchSchema "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SearchServiceDatasourceBlobResource struct{}

var _ sdk.ResourceWithUpdate = SearchServiceDatasourceBlobResource{}

type SearchServiceDatasourceBlobModel struct {
	Name                  string                                            `tfschema:"name"`
	SearchServiceId       string                                            `tfschema:"search_service_id"`
	ContainerName         string                                            `tfschema:"container_name"`
	ConnectionString      string                                            `tfschema:"connection_string"`
	Description           string                                            `tfschema:"description"`
	ContainerQuery        string                                            `tfschema:"container_query"`
	SoftDeleteColumnName  string                                            `tfschema:"soft_delete_column_name"`
	SoftDeleteMarkerValue string                                            `tfschema:"soft_delete_marker_value"`
	EncryptionKey         []searchSchema.SearchDatasourceEncryptionKeyModel `tfschema:"encryption_key"`
	Etag                  string                                            `tfschema:"etag"`
}

func (r SearchServiceDatasourceBlobResource) Arguments() map[string]*pluginsdk.Schema {
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

		"search_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: services.ValidateSearchServiceID,
		},

		"connection_string": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"container_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.StorageContainerName,
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

		"encryption_key": searchSchema.SearchDatasourceEncryptionKeySchema(),

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

func (r SearchServiceDatasourceBlobResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r SearchServiceDatasourceBlobResource) ResourceType() string {
	return "azurerm_search_service_datasource_blob"
}

func (r SearchServiceDatasourceBlobResource) ModelObject() interface{} {
	return &SearchServiceDatasourceBlobModel{}
}

func (r SearchServiceDatasourceBlobResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datasources.ValidateDatasourceID
}

func (r SearchServiceDatasourceBlobResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SearchServiceDatasourceBlobModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			searchServiceId, err := services.ParseSearchServiceID(model.SearchServiceId)
			if err != nil {
				return fmt.Errorf("parsing `search_service_id` %q: %+v", model.SearchServiceId, err)
			}
			domainSuffix, ok := metadata.Client.Account.Environment.Search.DomainSuffix()
			if !ok {
				return errors.New("could not determine Search domain suffix for the current environment")
			}
			endpoint := fmt.Sprintf("https://%s.%s", searchServiceId.SearchServiceName, *domainSuffix)

			client := metadata.Client.Search.SearchDataPlaneClient.DataSources.Clone(endpoint)

			id := datasources.NewDatasourceID(endpoint, model.Name)

			existing, err := client.Get(ctx, id, datasources.DefaultGetOperationOptions())
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := datasources.SearchIndexerDataSource{
				Name: model.Name,
				Type: datasources.SearchIndexerDataSourceTypeAzureblob,
				Container: datasources.SearchIndexerDataContainer{
					Name: model.ContainerName,
				},
				Credentials: datasources.DataSourceCredentials{},
			}

			if model.ConnectionString != "" {
				parameters.Credentials.ConnectionString = pointer.To(model.ConnectionString)
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

			parameters.EncryptionKey = searchSchema.ExpandSearchDatasourceEncryptionKey(model.EncryptionKey)

			if _, err := client.Create(ctx, parameters, datasources.DefaultCreateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r SearchServiceDatasourceBlobResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceId, err := datasources.ParseDatasourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			endpoint := resourceId.BaseURI

			// Check whether the parent Search Service still exists via ARM before making the
			// data-plane call. If the service has been deleted outside of Terraform the DNS
			// entry is gone and the data-plane call would fail with "no such host" instead of
			// a 404, which would not be caught by the WasNotFound check below.
			searchServiceId, err := services.ParseSearchServiceID(metadata.ResourceData.Get("search_service_id").(string))
			if err != nil {
				return fmt.Errorf("parsing `search_service_id`: %+v", err)
			}
			servicesClient := metadata.Client.Search.ServicesClient
			serviceResp, err := servicesClient.Get(ctx, *searchServiceId, services.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(serviceResp.HttpResponse) {
					return metadata.MarkAsGone(resourceId)
				}
				return fmt.Errorf("checking for existence of Search Service %s: %+v", *searchServiceId, err)
			}

			client := metadata.Client.Search.SearchDataPlaneClient.DataSources.Clone(endpoint)

			resp, err := client.Get(ctx, *resourceId, datasources.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(resourceId)
				}
				return fmt.Errorf("retrieving %s: %+v", *resourceId, err)
			}

			state := SearchServiceDatasourceBlobModel{
				Name:            resourceId.DatasourceName,
				SearchServiceId: searchServiceId.ID(),
			}

			if respModel := resp.Model; respModel != nil {
				state.ContainerName = respModel.Container.Name
				state.ContainerQuery = pointer.From(respModel.Container.Query)
				state.Description = pointer.From(respModel.Description)
				state.Etag = pointer.From(respModel.OdataEtag)

				if policy, ok := respModel.DataDeletionDetectionPolicy.(datasources.SoftDeleteColumnDeletionDetectionPolicy); ok {
					state.SoftDeleteColumnName = pointer.From(policy.SoftDeleteColumnName)
					state.SoftDeleteMarkerValue = pointer.From(policy.SoftDeleteMarkerValue)
				}

				state.EncryptionKey = searchSchema.FlattenSearchDatasourceEncryptionKey(respModel.EncryptionKey, metadata.ResourceData)
			}

			if v := metadata.ResourceData.Get("connection_string").(string); v != "" {
				state.ConnectionString = v
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r SearchServiceDatasourceBlobResource) Update() sdk.ResourceFunc {
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
				return fmt.Errorf("retrieving %s for update: %+v", *resourceId, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s for update: `model` was nil", *resourceId)
			}

			existing := *resp.Model

			if metadata.ResourceData.HasChange("container_name") || metadata.ResourceData.HasChange("container_query") {
				existing.Container.Name = state.ContainerName
				existing.Container.Query = pointer.To(state.ContainerQuery)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				existing.Credentials.ConnectionString = pointer.To(state.ConnectionString)
			}

			if metadata.ResourceData.HasChange("description") {
				existing.Description = pointer.To(state.Description)
			}

			if metadata.ResourceData.HasChange("encryption_key") {
				existing.EncryptionKey = searchSchema.ExpandSearchDatasourceEncryptionKey(state.EncryptionKey)
			}

			// Always nil out DataChangeDetectionPolicy - it is not managed by this resource and the SDK
			// may populate it with a zero-value RawDataChangeDetectionPolicyImpl from the GET response,
			// which serializes to an invalid payload the API rejects.
			existing.DataChangeDetectionPolicy = nil
			existing.DataDeletionDetectionPolicy = nil

			// Always set DataDeletionDetectionPolicy from current state regardless of HasChange.
			// Without this, an unchanged field retains the zero-value RawDataDeletionDetectionPolicyImpl
			// from the GET response which serializes to an invalid abstract-type payload the API rejects.
			if state.SoftDeleteColumnName != "" {
				existing.DataDeletionDetectionPolicy = datasources.SoftDeleteColumnDeletionDetectionPolicy{
					SoftDeleteColumnName:  pointer.To(state.SoftDeleteColumnName),
					SoftDeleteMarkerValue: pointer.To(state.SoftDeleteMarkerValue),
				}
			}

			prefer := datasources.PreferReturnRepresentation
			opts := datasources.CreateOrUpdateOperationOptions{
				Prefer: &prefer,
			}
			if _, err := client.CreateOrUpdate(ctx, *resourceId, existing, opts); err != nil {
				return fmt.Errorf("updating %s: %+v", *resourceId, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r SearchServiceDatasourceBlobResource) Delete() sdk.ResourceFunc {
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
