// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custom_ip_prefix

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/customipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 17 * time.Hour,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			r.client = metadata.Client.Network.CustomIPPrefixes

			id, err := customipprefixes.ParseCustomIPPrefixID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state CustomIpPrefixModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			desiredState := customipprefixes.CommissionedStateProvisioned
			if state.CommissioningEnabled {
				if state.InternetAdvertisingDisabled {
					desiredState = customipprefixes.CommissionedStateCommissionedNoInternetAdvertise
				} else {
					desiredState = customipprefixes.CommissionedStateCommissioned
				}
			}

			commissionedState, err := r.updateCommissionedState(ctx, *id, desiredState)
			if err != nil {
				return err
			}
			if commissionedState == nil {
				return fmt.Errorf("waiting for CommissionedState: final commissionedState was nil")
			}

			log.Printf("[DEBUG] Final CommissionedState is %q for %s..", *commissionedState, id)
			return nil
		},
	}
}
