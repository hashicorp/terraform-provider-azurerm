// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/custompollers"
)

func (m Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
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

			upd := watcher.WatcherUpdateParameters{
				Name: &id.WatcherName,
			}
			upd.Properties = &watcher.WatcherUpdateProperties{}
			if meta.ResourceData.HasChange("execution_frequency_in_seconds") {
				upd.Properties.ExecutionFrequencyInSeconds = pointer.To(model.ExecutionFrequencyInSeconds)
			}

			pollerType := custompollers.NewAutomationWatcherUpdatePoller(client, *id)
			poller := pollers.NewPoller(pollerType, 1*time.Minute, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return err
			}

			var pollerModel watcher.Watcher
			if err = poller.LatestResponse().Unmarshal(&pollerModel); err != nil {
				return err
			}

			status := pointer.From(pollerModel.Properties.Status)
			if status == "Running" {
				if _, err = client.Stop(ctx, *id); err != nil {
					return fmt.Errorf("stopping %s: %v", *id, err)
				}
			}

			if _, err = client.Update(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %v", *id, err)
			}

			if _, err = client.Start(ctx, *id); err != nil {
				return fmt.Errorf("starting %s: %v", *id, err)
			}

			return nil
		},
	}
}
