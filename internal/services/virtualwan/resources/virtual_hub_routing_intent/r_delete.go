// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_routing_intent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.VirtualWan.VirtualWANs

			id, err := virtualwans.ParseRoutingIntentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.RoutingIntentDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
