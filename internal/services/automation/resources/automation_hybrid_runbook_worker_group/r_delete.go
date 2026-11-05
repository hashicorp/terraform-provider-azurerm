// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_hybrid_runbook_worker_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/hybridrunbookworkergroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/listallhybridrunbookworkergroupinautomationaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := hybridrunbookworkergroup.ParseHybridRunbookWorkerGroupID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			// NOTE: Due to incorrect tags in Azure Swagger, the delete operation is only available
			// in the ListAllHybridRunbookWorkerGroupInAutomationAccount client, not in the
			// main HybridRunbookWorkerGroup client where it logically belongs
			// Issue tracked in: https://github.com/Azure/azure-rest-api-specs/issues/38703
			client := meta.Client.Automation.ListAllHybridRunbookWorkerGroupInAutomationAccount

			// Convert the ID to the listallhybridrunbookworkergroupinautomationaccount package format
			deleteId := listallhybridrunbookworkergroupinautomationaccount.NewHybridRunbookWorkerGroupID(
				id.SubscriptionId,
				id.ResourceGroupName,
				id.AutomationAccountName,
				id.HybridRunbookWorkerGroupName,
			)

			if _, err = client.HybridRunbookWorkerGroupDelete(ctx, deleteId); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}
