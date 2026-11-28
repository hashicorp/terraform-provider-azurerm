// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package tenant_template_deployment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-02-01/templatespecversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func tenantTemplateDeploymentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Template.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.TemplateDeploymentName,
		},

		"location": commonschema.Location(),

		"template_content": {
			Description: "The contents of the ARM Template which should be deployed into this Resource Group.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			StateFunc:   pluginsdk.NormalizeJson,
			ExactlyOneOf: []string{
				"template_content",
				"template_spec_version_id",
			},
		},

		"template_spec_version_id": {
			Description:  "The ID of the Template Spec Version to deploy.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.AsGeneratedID(templatespecversions.ParseTemplateSpecVersionIDInsensitively),
			ExactlyOneOf: []string{
				"template_content",
				"template_spec_version_id",
			},
		},

		// Optional
		"debug_level": {
			Description:  "The Debug Level which should be used for this Resource Group Template Deployment.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(helpers.TemplateDeploymentDebugLevels, false),
		},

		"parameters_content": {
			Description: "The contents of the ARM Template parameters file - containing a JSON list of parameters.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			StateFunc:   pluginsdk.NormalizeJson,
		},

		"tags": commonschema.Tags(),

		// Computed
		"output_content": {
			Description: "The JSON Content of the Outputs of the ARM Template Deployment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			// NOTE:  outputs can be strings, ints, objects etc - whilst using a nested object was considered
			// parsing the JSON using `jsondecode` allows the users to interact with/map objects as required
		},
	}
}
