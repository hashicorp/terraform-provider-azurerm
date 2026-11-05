// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_hybrid_runbook_worker_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/hybridrunbookworkergroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.HybridRunbookWorkerGroup
			id, err := hybridrunbookworkergroup.ParseHybridRunbookWorkerGroupID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model HybridRunbookWorkerGroupModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			var upd hybridrunbookworkergroup.HybridRunbookWorkerGroupCreateOrUpdateParameters
			if meta.ResourceData.HasChange("credential_name") {
				upd.Properties = &hybridrunbookworkergroup.HybridRunbookWorkerGroupCreateOrUpdateProperties{}
				upd.Properties.Credential = &hybridrunbookworkergroup.RunAsCredentialAssociationProperty{
					Name: pointer.To(model.CredentialName),
				}
			}
			if _, err = client.Update(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}
