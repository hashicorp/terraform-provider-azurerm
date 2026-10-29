// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool_registration_info

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopHostPoolRegistrationInfoDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HostPoolRegistrationInfoID(d.Id())
	if err != nil {
		return err
	}

	hostPoolId := hostpool.NewHostPoolID(id.SubscriptionId, id.ResourceGroup, id.HostPoolName)

	locks.ByName(hostPoolId.HostPoolName, hostPoolResourceType)
	defer locks.UnlockByName(hostPoolId.HostPoolName, hostPoolResourceType)

	resp, err := client.Get(ctx, hostPoolId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", hostPoolId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", hostPoolId, err)
	}

	regInfo, err := client.RetrieveRegistrationToken(ctx, hostPoolId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Virtual Desktop Host Pool %q Registration Info was not found in Resource Group %q - removing from state!", id.HostPoolName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Registration Token for %s: %+v", hostPoolId, err)
	}
	if regInfo.Model == nil || regInfo.Model.ExpirationTime == nil {
		log.Printf("[INFO] RegistrationInfo for %s was nil, registrationInfo already deleted - removing from state", hostPoolId)
		return nil
	}

	payload := hostpool.HostPoolPatch{
		Properties: &hostpool.HostPoolPatchProperties{
			RegistrationInfo: &hostpool.RegistrationInfoPatch{
				RegistrationTokenOperation: pointer.To(hostpool.RegistrationTokenOperationDelete),
			},
		},
	}

	if _, err := client.Update(ctx, hostPoolId, payload); err != nil {
		return fmt.Errorf("removing Registration Token from %s: %+v", hostPoolId, err)
	}

	return nil
}
