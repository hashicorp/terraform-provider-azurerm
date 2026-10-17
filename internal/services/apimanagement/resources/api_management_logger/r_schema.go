// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_logger

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementLoggerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"eventhub": {
			Type:          pluginsdk.TypeList,
			MaxItems:      1,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"application_insights"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ValidateEventHubName(),
					},

					"connection_string": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{
							"eventhub.0.connection_string",
							"eventhub.0.endpoint_uri",
						},
						ConflictsWith: []string{
							"eventhub.0.endpoint_uri",
							"eventhub.0.user_assigned_identity_client_id",
						},
					},
					"endpoint_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{
							"eventhub.0.connection_string",
							"eventhub.0.endpoint_uri",
						},
						ConflictsWith: []string{
							"eventhub.0.connection_string",
						},
					},
					"user_assigned_identity_client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
						ConflictsWith: []string{
							"eventhub.0.connection_string",
						},
					},
				},
			},
		},

		"application_insights": {
			Type:          pluginsdk.TypeList,
			MaxItems:      1,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"eventhub"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"connection_string": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{
							"application_insights.0.connection_string",
							"application_insights.0.instrumentation_key",
						},
						ConflictsWith: []string{
							"application_insights.0.instrumentation_key",
						},
					},
					"instrumentation_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{
							"application_insights.0.connection_string",
							"application_insights.0.instrumentation_key",
						},
						ConflictsWith: []string{
							"application_insights.0.connection_string",
						},
					},
					"identity_client_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.Any(
							validation.IsUUID,
							validation.StringInSlice([]string{"SystemAssigned"}, false),
						),
						RequiredWith: []string{
							"application_insights.0.connection_string",
						},
						ConflictsWith: []string{
							"application_insights.0.instrumentation_key",
						},
					},
				},
			},
		},

		"buffered": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}
