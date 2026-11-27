// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_microsoft_threat_protection

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithCustomImporter = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_sentinel_data_connector_microsoft_threat_protection"
}

func (r Resource) Description() string {
	return "Manages a Microsoft Threat Protection Data Connector."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dataconnectors.ValidateDataConnectorID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return helpers.ImportDataConnectorTyped(dataconnectors.DataConnectorKindMicrosoftThreatProtection)
}
