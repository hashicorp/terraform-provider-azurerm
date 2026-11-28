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

		"container_name": {
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

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
