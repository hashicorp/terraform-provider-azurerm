// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool_registration_info

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopHostPoolRegistrationInfoRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HostPoolRegistrationInfoID(d.Id())
	if err != nil {
		return err
	}

	hostPoolId := hostpool.NewHostPoolID(id.SubscriptionId, id.ResourceGroup, id.HostPoolName)
	resp, err := client.RetrieveRegistrationToken(ctx, hostPoolId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Registration Token was not found for %s - removing from state!", hostPoolId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Registration Token for %s: %+v", hostPoolId, err)
	}

	if resp.Model == nil || resp.Model.ExpirationTime == nil || resp.Model.Token == nil {
		log.Printf("HostPool is missing registration info - marking as gone")
		d.SetId("")
		return nil
	}

	d.Set("hostpool_id", hostPoolId.ID())
	d.Set("expiration_date", resp.Model.ExpirationTime)
	d.Set("token", resp.Model.Token)

	return nil
}

func hostpoolRegistrationInfoCustomDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	if d.HasChange("expiration_date") {
		if err := d.SetNewComputed("token"); err != nil {
			return err
		}
		return nil
	}
	return nil
}
