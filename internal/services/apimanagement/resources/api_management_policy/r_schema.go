// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_policy

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_management_id": {
			Description:  "The ID of the API Management service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apimanagementservice.ValidateServiceID,
		},

		"xml_content": {
			Description: "The XML Content for this Policy as a string. An XML file can be used here with Terraform's [file function](https://developer.hashicorp.com/terraform/language/functions/file) that is similar to Microsoft's `PolicyFilePath` option. To integrate frontend and backend services in Azure API Management, utilize the [`set-backend-service`](https://learn.microsoft.com/azure/api-management/set-backend-service-policy) policy, specifying the `base-url` value. Typically, this value corresponds to the `url` property defined in the [`azurerm_api_management_backend`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/api_management_backend) configuration.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because when `xml_link` is provided the API downloads it into `xml_content`
			Computed:         true,
			ConflictsWith:    []string{"xml_link"},
			ExactlyOneOf:     []string{"xml_link", "xml_content"},
			DiffSuppressFunc: helpers.XmlWithDotNetInterpolationsDiffSuppress,
		},

		"xml_link": {
			Description:   "A link to a Policy XML Document, which must be publicly available.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"xml_content"},
			ExactlyOneOf:  []string{"xml_link", "xml_content"},
		},
	}
}
