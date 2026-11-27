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

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			id, err := metricsobjectfirewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.MetricsObjectFirewallDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting Metrics Object for %s: %+v", id, err)
			}

			return nil
		},
	}
}
