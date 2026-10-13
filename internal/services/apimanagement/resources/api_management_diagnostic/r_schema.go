// Copyright IBM Corp. 2014, 2025
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"applicationinsights",
				"azuremonitor",
			}, false),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"api_management_logger_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: logger.ValidateLoggerID,
		},

		"sampling_percentage": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.FloatBetween(0.0, 100.0),
		},

		"always_log_errors": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"verbosity": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(diagnostic.PossibleValuesForVerbosity(), false),
		},

		"log_client_ip": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"http_correlation_protocol": {
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
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(diagnostic.PossibleValuesForOperationNameFormat(), false),
		},
	}
}
