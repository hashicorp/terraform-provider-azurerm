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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},
		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"api_root_url": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"collection_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"user_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"polling_frequency": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(dataconnectors.PollingFrequencyOnceAnHour),
			ValidateFunc: validation.StringInSlice(dataconnectors.PossibleValuesForPollingFrequency(),
				false),
		},
		"lookback_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "1970-01-01T00:00:00Z",
			ValidateFunc: validation.IsRFC3339Time,
		},
		"tenant_id": {
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
