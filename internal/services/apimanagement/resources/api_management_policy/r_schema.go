// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apimanagementservice.ValidateServiceID,
		},

		"xml_content": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because when `xml_link` is provided the API downloads it into `xml_content`
			Computed:         true,
			ConflictsWith:    []string{"xml_link"},
			ExactlyOneOf:     []string{"xml_link", "xml_content"},
			DiffSuppressFunc: helpers.XmlWithDotNetInterpolationsDiffSuppress,
		},

		"xml_link": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"xml_content"},
			ExactlyOneOf:  []string{"xml_link", "xml_content"},
		},
	}
}
