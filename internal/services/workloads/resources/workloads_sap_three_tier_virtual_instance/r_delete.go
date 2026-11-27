// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_three_tier_virtual_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
