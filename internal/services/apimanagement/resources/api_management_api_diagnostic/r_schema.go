// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_diagnostic

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apidiagnostic"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementApiDiagnosticSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"identifier": {
			Description: "Identifier of the Diagnostics Logs.",
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

		"api_name": schemaz.SchemaApiManagementApiName("The name of the API on which to configure the Diagnostics Logs."),

		"api_management_logger_id": {
			Description:  "The ID (name) of the Diagnostics Logger.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: logger.ValidateLoggerID,
		},

		"sampling_percentage": {
			Description:  "Sampling (%). For high traffic APIs, please read this [documentation](https://learn.microsoft.com/azure/api-management/api-management-howto-app-insights#performance-implications-and-log-sampling) to understand performance implications and log sampling.",
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
			ValidateFunc: validation.StringInSlice(apidiagnostic.PossibleValuesForVerbosity(), false),
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
			ValidateFunc: validation.StringInSlice(apidiagnostic.PossibleValuesForHTTPCorrelationProtocol(), false),
		},

		"frontend_request": helpers.ResourceApiManagementApiDiagnosticAdditionalContentSchema(),

		"frontend_response": helpers.ResourceApiManagementApiDiagnosticAdditionalContentSchema(),

		"backend_request": helpers.ResourceApiManagementApiDiagnosticAdditionalContentSchema(),

		"backend_response": helpers.ResourceApiManagementApiDiagnosticAdditionalContentSchema(),

		"operation_name_format": {
			Description:  "The format of the Operation Name for Application Insights telemetries.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(apidiagnostic.OperationNameFormatName),
			ValidateFunc: validation.StringInSlice(apidiagnostic.PossibleValuesForOperationNameFormat(), false),
		},
	}
}
