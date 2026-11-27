// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_office_365

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an Office 365 Data Connector.",
		Create:      resourceSentinelDataConnectorOffice365CreateUpdate,
		Read:        resourceSentinelDataConnectorOffice365Read,
		Update:      resourceSentinelDataConnectorOffice365CreateUpdate,
		Delete:      resourceSentinelDataConnectorOffice365Delete,

		Importer: helpers.ImportDataConnectorUntyped(dataconnectors.DataConnectorKindOfficeThreeSixFive),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: sentinelDataConnectorOffice365Schema(),
	}
}
