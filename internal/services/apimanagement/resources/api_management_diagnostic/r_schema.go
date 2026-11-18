// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_diagnostic

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/diagnostic"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementDiagnosticSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"identifier": {
			Description: "The diagnostic identifier for the API Management Service. At this time the.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"applicationinsights",
				"azuremonitor",
			}, false),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"api_management_logger_id": {
			Description:  "The id of the target API Management Logger where the API Management Diagnostic should be saved.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: logger.ValidateLoggerID,
		},

		"sampling_percentage": {
			Description:  "Sampling (%). For high traffic APIs, please read this [documentation](https://docs.microsoft.com/azure/api-management/api-management-howto-app-insights#performance-implications-and-log-sampling) to understand performance implications and log sampling.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.FloatBetween(0.0, 100.0),
		},

		"always_log_errors": {
			Description: "Always log errors. Send telemetry if there is an erroneous condition, regardless of sampling settings.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"verbosity": {
			Description:  "Logging verbosity.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(diagnostic.PossibleValuesForVerbosity(), false),
		},

		"log_client_ip": {
			Description: "Log client IP address.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"metrics_enabled": {
			Description: "Whether custom metrics emitted via the `emit-metric` policy are sent to Application Insights.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"http_correlation_protocol": {
			Description:  "The HTTP Correlation Protocol to use.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(diagnostic.PossibleValuesForHTTPCorrelationProtocol(), false),
		},

		"frontend_request": helpers.ResourceApiManagementApiDiagnosticAdditionalContentSchema(),

		"frontend_response": helpers.ResourceApiManagementApiDiagnosticAdditionalContentSchema(),

		"backend_request": helpers.ResourceApiManagementApiDiagnosticAdditionalContentSchema(),

		"backend_response": helpers.ResourceApiManagementApiDiagnosticAdditionalContentSchema(),

		"operation_name_format": {
			Description:  "The format of the Operation Name for Application Insights telemetries.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(diagnostic.PossibleValuesForOperationNameFormat(), false),
		},
	}
}
