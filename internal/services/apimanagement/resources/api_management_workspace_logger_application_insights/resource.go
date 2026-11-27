// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_logger_application_insights

// NOTE: this resource and azurerm_api_management_workspace_logger_eventhub share the same Azure resource type
// but are split by LoggerType: "applicationInsights" vs "azureEventHub"

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_api_management_workspace_logger_application_insights"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return logger.ValidateWorkspaceLoggerID
}
