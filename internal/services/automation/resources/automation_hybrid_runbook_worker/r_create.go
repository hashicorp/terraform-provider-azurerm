// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_hybrid_runbook_worker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/hybridrunbookworker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.HybridRunbookWorker

			var model HybridRunbookWorkerModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionID := metadata.Client.Account.SubscriptionId
			id := hybridrunbookworker.NewHybridRunbookWorkerID(subscriptionID, model.ResourceGroupName,
				model.AutomationAccountName, model.WorkerGroupName, model.WorkerId)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			req := hybridrunbookworker.HybridRunbookWorkerCreateParameters{
				Properties: &hybridrunbookworker.HybridRunbookWorkerCreateOrUpdateParameters{},
			}
			if model.VmResourceId != "" {
				req.Properties.VMResourceId = pointer.To(model.VmResourceId)
			}

			future, err := client.Create(ctx, id, req)
			if err != nil {
				// Workaround swagger issue https://github.com/Azure/azure-rest-api-specs/issues/19741
				if !response.WasStatusCode(future.HttpResponse, http.StatusCreated) {
					return fmt.Errorf("creating %s: %v", id, err)
				}
			}
			_ = future

			metadata.SetID(id)
			return nil
		},
	}
}
