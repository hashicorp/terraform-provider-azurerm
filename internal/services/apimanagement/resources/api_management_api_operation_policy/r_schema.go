// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementApiOperationPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"api_name": schemaz.SchemaApiManagementApiName("The name of the API within the API Management Service where the Operation exists."),

		"operation_id": schemaz.SchemaApiManagementChildName("The operation identifier within an API. Must be unique in the current API Management service instance."),

		"xml_content": {
			Description: "The XML Content for this Policy.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because when `xml_link` is provided the API downloads it into `xml_content`
			Computed:         true,
			DiffSuppressFunc: helpers.XmlWithDotNetInterpolationsDiffSuppress,
			ConflictsWith:    []string{"xml_link"},
		},

		"xml_link": {
			Description:   "A link to a Policy XML Document, which must be publicly available.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"xml_content"},
		},
	}
}
