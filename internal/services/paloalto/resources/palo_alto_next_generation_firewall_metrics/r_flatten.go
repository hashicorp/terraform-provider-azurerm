// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_metrics

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *metricsobjectfirewallresources.FirewallId, model *metricsobjectfirewallresources.MetricsObjectFirewallResource) error {
	state := NextGenerationFirewallMetricsModel{
		FirewallID: id.ID(),
	}

	if model != nil {
		props := model.Properties

		// the API may redact the connection string on GET, so default to the value currently in state
		state.ApplicationInsightsConnectionString = metadata.ResourceData.Get("application_insights_connection_string").(string)
		if props.ApplicationInsightsConnectionString != "" {
			state.ApplicationInsightsConnectionString = props.ApplicationInsightsConnectionString
		}

		if props.ApplicationInsightsResourceId != "" {
			applicationInsightsId, err := componentsapis.ParseComponentIDInsensitively(props.ApplicationInsightsResourceId)
			if err != nil {
				return err
			}
			state.ApplicationInsightsID = applicationInsightsId.ID()
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id, pluginsdk.ResourceTypeForIdentityVirtual); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
