// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_managed_private_endpoint

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := streamanalyticsClient.From(metadata.Client).EndpointsClient
			id, err := privateendpoints.ParsePrivateEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if resp.Model.Properties.ManualPrivateLinkServiceConnections == nil {
				return errors.New("no private link service connections available")
			}

			state := ManagedPrivateEndpointModel{
				Name:                   id.PrivateEndpointName,
				ResourceGroup:          id.ResourceGroupName,
				StreamAnalyticsCluster: id.ClusterName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					for _, mplsc := range *props.ManualPrivateLinkServiceConnections {
						state.TargetResourceId = *mplsc.Properties.PrivateLinkServiceId
						state.SubResourceName = strings.Join(*mplsc.Properties.GroupIds, "")
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}
