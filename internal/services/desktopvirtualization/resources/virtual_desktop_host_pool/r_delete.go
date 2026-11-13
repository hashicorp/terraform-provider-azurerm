// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/sessionhost"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopHostPoolDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	sessionHostClient := meta.(*clients.Client).DesktopVirtualization.SessionHostsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// Delete session hosts associated with host pool first, or else host pool cannot be deleted
	sessionHostPoolId, err := sessionhost.ParseHostPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := sessionHostClient.ListComplete(ctx, pointer.From(sessionHostPoolId), sessionhost.DefaultListOperationOptions())
	if err != nil {
		return fmt.Errorf("listing Session Hosts for %s: %+v", sessionHostPoolId, err)
	}

	for _, sessionHost := range resp.Items {
		sessionHostId, err := sessionhost.ParseSessionHostIDInsensitively(pointer.From(sessionHost.Id))
		if err != nil {
			return err
		}

		if _, err := sessionHostClient.Delete(ctx, pointer.From(sessionHostId), sessionhost.DefaultDeleteOperationOptions()); err != nil {
			return fmt.Errorf("deleting %s: %+v", *sessionHostId, err)
		}
	}

	id, err := hostpool.ParseHostPoolID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.HostPoolName, hostPoolResourceType)
	defer locks.UnlockByName(id.HostPoolName, hostPoolResourceType)

	options := hostpool.DeleteOperationOptions{
		Force: pointer.To(true),
	}
	if _, err = client.Delete(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
