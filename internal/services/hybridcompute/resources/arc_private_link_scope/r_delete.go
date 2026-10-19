// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_private_link_scope

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privatelinkscopes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (a Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridCompute.PrivateLinkScopesClient

			id, err := privatelinkscopes.ParseProviderPrivateLinkScopeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
