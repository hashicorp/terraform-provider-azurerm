// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custom_ip_prefix

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/customipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 17 * time.Hour,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			r.client = metadata.Client.Network.CustomIPPrefixes

			id, err := customipprefixes.ParseCustomIPPrefixID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Must be de-provisioned before deleting
			if _, err = r.updateCommissionedState(ctx, *id, customipprefixes.CommissionedStateDeprovisioned); err != nil {
				return err
			}

			if err := r.client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
