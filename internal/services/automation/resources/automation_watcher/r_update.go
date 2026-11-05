// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.WatcherClient

			id, err := watcher.ParseWatcherID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WatcherModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			var upd watcher.WatcherUpdateParameters
			upd.Properties = &watcher.WatcherUpdateProperties{}
			if meta.ResourceData.HasChange("execution_frequency_in_seconds") {
				upd.Properties.ExecutionFrequencyInSeconds = pointer.To(model.ExecutionFrequencyInSeconds)
			}
			if _, err = client.Update(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %v", *id, err)
			}

			return nil
		},
	}
}
