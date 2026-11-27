// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			var model NextGenerationFirewallMetricsModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := metricsobjectfirewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.MetricsObjectFirewallGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving Metrics Object for %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving Metrics Object for %s: `model` was nil", id)
			}

			payload := *existing.Model

			// the API may redact the connection string on GET, so it is always sent from config to keep the PUT payload valid
			payload.Properties.ApplicationInsightsConnectionString = model.ApplicationInsightsConnectionString

			if metadata.ResourceData.HasChange("application_insights_id") {
				payload.Properties.ApplicationInsightsResourceId = model.ApplicationInsightsID
			}

			if err := client.MetricsObjectFirewallCreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating Metrics Object for %s: %+v", id, err)
			}

			return nil
		},
	}
}
