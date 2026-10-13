// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_policy

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementApiPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"api_name": schemaz.SchemaApiManagementApiName(),

		"xml_content": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because when `xml_link` is provided the API downloads it into `xml_content`
			Computed:         true,
			ConflictsWith:    []string{"xml_link"},
			DiffSuppressFunc: helpers.XmlWithDotNetInterpolationsDiffSuppress,
		},

		"xml_link": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"xml_content"},
		},
	}
}
