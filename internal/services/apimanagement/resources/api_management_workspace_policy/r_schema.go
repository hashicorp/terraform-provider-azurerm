// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_policy

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"xml_content": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C when `xml_link` is provided, the API downloads it into `xml_content`.
			Computed:         true,
			ExactlyOneOf:     []string{"xml_link", "xml_content"},
			DiffSuppressFunc: helpers.XmlWithDotNetInterpolationsDiffSuppress,
		},

		"xml_link": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"xml_link", "xml_content"},
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
