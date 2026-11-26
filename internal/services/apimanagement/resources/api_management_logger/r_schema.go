// Copyright IBM Corp. 2014, 2026
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
		"name": schemaz.SchemaApiManagementChildName("The name of this Logger, which must be unique within the API Management Service."),

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_id": {
			Description:  "The target resource ID which will be linked in the API-Management portal page.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"eventhub": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			ConflictsWith: []string{"application_insights"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of an EventHub.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ValidateEventHubName(),
					},

					"connection_string": {
						Description:  "The connection string of an EventHub Namespace.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ConflictsWith: []string{
							"eventhub.0.endpoint_uri",
							"eventhub.0.user_assigned_identity_client_id",
						},
						AtLeastOneOf: []string{
							"eventhub.0.connection_string",
							"eventhub.0.endpoint_uri",
						},
					},
					"endpoint_uri": {
						Description:  "The endpoint address of an EventHub Namespace.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						ConflictsWith: []string{
							"eventhub.0.connection_string",
						},
						AtLeastOneOf: []string{
							"eventhub.0.connection_string",
							"eventhub.0.endpoint_uri",
						},
					},
					"user_assigned_identity_client_id": {
						Description:  "The Client ID of the User Assigned Identity with the \"Azure Event Hubs Data Sender\" role to the target EventHub Namespace. If not specified the System Assigned Identity will be used.",
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
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			ConflictsWith: []string{"eventhub"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"connection_string": {
						Description:  "The connection string of Application Insights.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ConflictsWith: []string{
							"application_insights.0.instrumentation_key",
						},
						AtLeastOneOf: []string{
							"application_insights.0.connection_string",
							"application_insights.0.instrumentation_key",
						},
					},
					"instrumentation_key": {
						Description:  "The instrumentation key used to push data to Application Insights.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ConflictsWith: []string{
							"application_insights.0.connection_string",
						},
						AtLeastOneOf: []string{
							"application_insights.0.connection_string",
							"application_insights.0.instrumentation_key",
						},
					},
					"identity_client_id": {
						Description: "The Client ID of the User Assigned Identity, or `SystemAssigned` to use the System Assigned Identity, that has the \"Monitoring Metrics Publisher\" role on the target Application Insights resource.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
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
			Description: "Whether records should be buffered in the Logger before publishing.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"description": {
			Description: "A description of this Logger.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},
	}
}
