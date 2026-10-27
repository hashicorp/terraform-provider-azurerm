// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool_registration_info

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopHostPoolRegistrationInfoCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	hostPoolId, err := hostpool.ParseHostPoolID(d.Get("hostpool_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(hostPoolId.HostPoolName, hostPoolResourceType)
	defer locks.UnlockByName(hostPoolId.HostPoolName, hostPoolResourceType)

	// This is a virtual resource so the last segment is hardcoded
	id := parse.NewHostPoolRegistrationInfoID(hostPoolId.SubscriptionId, hostPoolId.ResourceGroupName, hostPoolId.HostPoolName, "default")

	existing, err := client.Get(ctx, *hostPoolId)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("%s could not be found: %s", hostPoolId, err)
		}
		return fmt.Errorf("reading %s: %s", hostPoolId, err)
	}

	payload := hostpool.HostPoolPatch{
		Properties: &hostpool.HostPoolPatchProperties{
			RegistrationInfo: &hostpool.RegistrationInfoPatch{
				ExpirationTime:             pointer.To(d.Get("expiration_date").(string)),
				RegistrationTokenOperation: pointer.To(hostpool.RegistrationTokenOperationUpdate),
			},
		},
	}
	if _, err := client.Update(ctx, *hostPoolId, payload); err != nil {
		return fmt.Errorf("updating registration token for %s: %+v", hostPoolId, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceVirtualDesktopHostPoolRegistrationInfoRead(d, meta)
}
