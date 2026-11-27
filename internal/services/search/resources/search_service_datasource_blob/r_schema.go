// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_service_datasource_blob

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/search/schema"
	searchValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Search Service Blob Data Source.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,126}[a-z0-9]$`),
				"The `name` must be 2-128 characters, start and end with an alphanumeric character, and contain only lowercase letters, digits, or dashes",
			),
		},

		"search_service_endpoint": {
			Description:  "The endpoint of the Search Service in which this Blob Data Source should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: searchValidate.SearchServiceEndpoint,
		},

		"container_name": {
			Description:  "The name of the Azure Blob Storage Container from which to read data.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: storageValidate.StorageContainerName,
		},

		"connection_string_wo": {
			Description:  "The connection string to the Azure Blob Storage Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			WriteOnly:    true,
			ValidateFunc: searchValidate.SearchDatasourceStorageConnectionString,
		},

		"connection_string_wo_version": {
			Description:  "An integer value used to trigger an update for `connection_string_wo`.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"container_query": {
			Description:  "A query string that filters the set of blobs in the container.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Description:  "A description for this Search Service Blob Data Source.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"encryption_key": schema.SearchDatasourceEncryptionKeySchema(),

		"soft_delete_column_name": {
			Description:  "The name of the column that indicates soft deletion of a blob.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"soft_delete_marker_value": {
			Description:  "The value in the `soft_delete_column_name` column that indicates a soft-deleted blob.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"soft_delete_column_name"},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
