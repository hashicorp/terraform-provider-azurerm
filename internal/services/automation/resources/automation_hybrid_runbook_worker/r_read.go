// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_hybrid_runbook_worker

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/hybridrunbookworker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := hybridrunbookworker.ParseHybridRunbookWorkerID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			client := meta.Client.Automation.HybridRunbookWorker
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return err
			}
			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}

			var output HybridRunbookWorkerModel

			// the name in response corresponding to work_id in request
			output.WorkerId = pointer.From(result.Model.Name)
			output.AutomationAccountName = id.AutomationAccountName
			output.ResourceGroupName = id.ResourceGroupName
			output.WorkerGroupName = id.HybridRunbookWorkerGroupName
			if prop := result.Model.Properties; prop != nil {
				output.VmResourceId = pointer.From(prop.VMResourceId)
				output.WorkerType = pointer.From((*string)(prop.WorkerType))
				output.LastSeenDateTime = pointer.From(prop.LastSeenDateTime)
				output.RegisteredDateTime = pointer.From(prop.RegisteredDateTime)
				output.Ip = pointer.From(prop.IP)
				output.WorkerName = pointer.From(prop.WorkerName)
			}
			return meta.Encode(&output)
		},
	}
}
