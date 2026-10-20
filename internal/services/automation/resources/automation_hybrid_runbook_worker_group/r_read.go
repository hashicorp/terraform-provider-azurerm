// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_hybrid_runbook_worker_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/hybridrunbookworkergroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := hybridrunbookworkergroup.ParseHybridRunbookWorkerGroupID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			client := meta.Client.Automation.HybridRunbookWorkerGroup
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
			var output HybridRunbookWorkerGroupModel

			output.Name = pointer.From(result.Model.Name)
			output.AutomationAccountName = id.AutomationAccountName
			if model := result.Model; model != nil {
				if prop := model.Properties; prop != nil {
					if prop.Credential != nil {
						output.CredentialName = pointer.From(prop.Credential.Name)
					}
				}
			}
			output.ResourceGroupName = id.ResourceGroupName
			return meta.Encode(&output)
		},
	}
}
