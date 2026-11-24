// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_hybrid_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.HybridConnectionsClient

			id, err := hybridconnections.ParseHybridConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := RelayHybridConnectionResourceModel{
				Name:               id.HybridConnectionName,
				ResourceGroupName:  id.ResourceGroupName,
				RelayNamespaceName: id.NamespaceName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.RequiresClientAuthorization = pointer.From(props.RequiresClientAuthorization)
					state.UserMetadata = pointer.From(props.UserMetadata)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
