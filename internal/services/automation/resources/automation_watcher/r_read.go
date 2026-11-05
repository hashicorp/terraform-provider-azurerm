// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_watcher

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
)

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := watcher.ParseWatcherID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Automation.WatcherClient
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			var output WatcherModel
			if err := meta.Decode(&output); err != nil {
				return err
			}

			output.Name = id.WatcherName

			if model := resp.Model; model != nil {
				if props := resp.Model.Properties; props != nil {
					output.AutomationAccountID = watcher.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName).ID()
					output.ExecutionFrequencyInSeconds = pointer.From(props.ExecutionFrequencyInSeconds)
					output.ScriptName = pointer.From(props.ScriptName)
					output.ScriptRunOn = pointer.From(props.ScriptRunOn)
					output.Description = pointer.From(props.Description)
					output.Status = pointer.From(props.Status)

					if props.ScriptParameters != nil {
						output.ScriptParameters = helpers.FlattenMap(*props.ScriptParameters)
					}
				}
			}

			// tags, etag and location are not returned by response, so do NOT encode them

			return meta.Encode(&output)
		},
	}
}
