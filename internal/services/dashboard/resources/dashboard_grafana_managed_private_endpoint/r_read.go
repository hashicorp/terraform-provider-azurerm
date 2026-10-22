// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dashboard_grafana_managed_private_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedgrafanas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedprivateendpointmodels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dashboard.ManagedPrivateEndpointsClient
			id, err := managedprivateendpointmodels.ParseManagedPrivateEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ManagedPrivateEndpointsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			grafanaId := managedgrafanas.NewGrafanaID(id.SubscriptionId, id.ResourceGroupName, id.GrafanaName)
			state := ManagedPrivateEndpointModel{
				Name:      id.ManagedPrivateEndpointName,
				GrafanaId: grafanaId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.GroupIds = pointer.From(props.GroupIds)
					state.PrivateLinkResourceId = pointer.From(props.PrivateLinkResourceId)
					state.PrivateLinkResourceRegion = pointer.From(props.PrivateLinkResourceRegion)
					state.RequestMessage = pointer.From(props.RequestMessage)
					state.PrivateLinkServiceURL = pointer.From(props.PrivateLinkServiceURL)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}
