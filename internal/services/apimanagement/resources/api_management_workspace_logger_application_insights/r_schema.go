// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_logger_application_insights

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName(),

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&logger.WorkspaceId{}),

		"application_insights": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"connection_string": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ExactlyOneOf: []string{"application_insights.0.connection_string", "application_insights.0.instrumentation_key"},
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"instrumentation_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ExactlyOneOf: []string{"application_insights.0.connection_string", "application_insights.0.instrumentation_key"},
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"buffering_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
