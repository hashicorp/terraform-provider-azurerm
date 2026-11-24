// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_version_set

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apiversionset"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementApiVersionSetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName("The name of the API Version Set. May only contain alphanumeric characters and dashes up to 80 characters in length."),

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"display_name": {
			Description:  "The display name of this API Version Set.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"versioning_scheme": {
			Description:  "Where in an Inbound HTTP Request that the API Version should be read from.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(apiversionset.PossibleValuesForVersioningScheme(), false),
		},

		"description": {
			Description:  "The description of API Version Set.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"version_header_name": {
			Description:   "The name of the Header which should be read from Inbound Requests which defines the API Version.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"version_query_name"},
		},

		"version_query_name": {
			Description:   "The name of the Query String which should be read from Inbound Requests which defines the API Version.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"version_header_name"},
		},
	}
}
