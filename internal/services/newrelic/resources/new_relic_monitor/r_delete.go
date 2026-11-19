// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package new_relic_monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NewRelic.MonitorsClient
			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NewRelicMonitorModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if len(model.UserInfo) == 0 {
				return fmt.Errorf("deleting: `user` not found")
			}

			if err = client.DeleteThenPoll(ctx, *id, monitors.DeleteOperationOptions{UserEmail: &model.UserInfo[0].EmailAddress}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
