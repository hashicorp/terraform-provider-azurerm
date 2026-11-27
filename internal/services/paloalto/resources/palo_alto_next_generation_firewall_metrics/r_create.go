// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			var model NextGenerationFirewallMetricsModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := metricsobjectfirewallresources.ParseFirewallID(model.FirewallID)
			if err != nil {
				return err
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.MetricsObjectFirewallGet(ctx, *id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Metrics Object for %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			input := metricsobjectfirewallresources.MetricsObjectFirewallResource{
				Properties: metricsobjectfirewallresources.MetricsObject{
					ApplicationInsightsConnectionString: model.ApplicationInsightsConnectionString,
					ApplicationInsightsResourceId:       model.ApplicationInsightsID,
				},
			}

			if err := client.MetricsObjectFirewallCreateOrUpdateCallbackThenPoll(ctx, *id, input, metadata.SetIDAndIdentityWithTypeCallback(id, pluginsdk.ResourceTypeForIdentityVirtual)); err != nil {
				return fmt.Errorf("creating Metrics Object for %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, id, pluginsdk.ResourceTypeForIdentityVirtual)
		},
	}
}
