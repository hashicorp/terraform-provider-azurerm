// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_virtual_machine_instance_guest_agent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	systemcentervirtualmachinemanagerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := systemcentervirtualmachinemanagerClient.From(metadata.Client).GuestAgents

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, commonids.NewScopeID(id.Scope))
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentModel{
				ScopedResourceId: id.Scope,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if v := props.Credentials; v != nil {
						state.Username = v.Username
						state.Password = metadata.ResourceData.Get("password").(string)
					}

					state.ProvisioningAction = pointer.FromEnum(props.ProvisioningAction)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
