// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_azure_security_center

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an Azure Security Center Data Connector.",
		Create:      resourceSentinelDataConnectorAzureSecurityCenterCreate,
		Read:        resourceSentinelDataConnectorAzureSecurityCenterRead,
		Delete:      resourceSentinelDataConnectorAzureSecurityCenterDelete,

		Importer: helpers.ImportDataConnectorUntyped(dataconnectors.DataConnectorKindAzureSecurityCenter),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: sentinelDataConnectorAzureSecurityCenterSchema(),
	}
}
