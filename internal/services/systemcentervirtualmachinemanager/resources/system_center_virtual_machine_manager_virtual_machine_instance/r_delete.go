// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_virtual_machine_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachineinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineInstances

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, commonids.NewScopeID(id.Scope), virtualmachineinstances.DeleteOperationOptions{
				DeleteFromHost: pointer.To(virtualmachineinstances.DeleteFromHostTrue),
				Force:          pointer.To(virtualmachineinstances.ForceDeleteTrue),
			}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
