// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	desktopvirtualizationClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopHostPoolDelete(d *pluginsdk.ResourceData, meta any) error {
	client := desktopvirtualizationClient.From(meta).HostPoolsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

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
