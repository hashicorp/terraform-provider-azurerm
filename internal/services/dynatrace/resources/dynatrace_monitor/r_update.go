// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dynatrace_monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.MonitorsClient
			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MonitorsResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") {
				props := monitors.MonitorResourceUpdate{
					Tags: pointer.To(state.Tags),
				}

				if _, err := client.Update(ctx, *id, props); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
