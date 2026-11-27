// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_threat_intelligence_taxii

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Threat Intelligence TAXII Data Connector.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"log_analytics_workspace_id": {
			Description:  "The ID of the Log Analytics Workspace that this Threat Intelligence TAXII Data Connector resides in.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},
		"display_name": {
			Description:  "The friendly name which should be used for this Threat Intelligence TAXII Data Connector.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"api_root_url": {
			Description:  "The API root URI of the TAXII server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"collection_id": {
			Description:  "The collection ID of the TAXII server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"user_name": {
			Description:  "The user name for the TAXII server.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"password": {
			Description:  "The password for the TAXII server.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"polling_frequency": {
			Description: "The polling frequency for the TAXII server.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(dataconnectors.PollingFrequencyOnceAnHour),
			ValidateFunc: validation.StringInSlice(dataconnectors.PossibleValuesForPollingFrequency(),
				false),
		},
		"lookback_date": {
			Description:  "The lookback date for the TAXII server in RFC3339.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "1970-01-01T00:00:00Z",
			ValidateFunc: validation.IsRFC3339Time,
		},
		"tenant_id": {
			Description:  "The ID of the tenant that this Threat Intelligence TAXII Data Connector connects to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
