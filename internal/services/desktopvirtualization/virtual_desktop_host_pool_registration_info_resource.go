// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualDesktopHostPoolRegistrationInfo() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopHostPoolRegistrationInfoCreateUpdate,
		Read:   resourceVirtualDesktopHostPoolRegistrationInfoRead,
		Update: resourceVirtualDesktopHostPoolRegistrationInfoCreateUpdate,
		Delete: resourceVirtualDesktopHostPoolRegistrationInfoDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			// TODO: we should refactor this to use the Host Pool ID instead, meaning we can remove this Virtual ID
			_, err := parse.HostPoolRegistrationInfoID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(hostpoolRegistrationInfoCustomDiff),

		Schema: map[string]*pluginsdk.Schema{
			"hostpool_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: hostpool.ValidateHostPoolID,
			},

			"expiration_date": {
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsRFC3339Time,
				Required:     true,
			},

			"token": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func hostpoolRegistrationInfoCustomDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	if d.HasChange("expiration_date") {
		if err := d.SetNewComputed("token"); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func resourceVirtualDesktopHostPoolRegistrationInfoCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	tokenOperation := hostpool.RegistrationTokenOperationUpdate
	payload := hostpool.HostPoolPatch{
		Properties: &hostpool.HostPoolPatchProperties{
			RegistrationInfo: &hostpool.RegistrationInfoPatch{
				ExpirationTime:             utils.String(d.Get("expiration_date").(string)),
				RegistrationTokenOperation: &tokenOperation,
			},
		},
	}
	if _, err := client.Update(ctx, *hostPoolId, payload); err != nil {
		return fmt.Errorf("updating registration token for %s: %+v", hostPoolId, err)
	}

	d.SetId(id.ID())

	return resourceVirtualDesktopHostPoolRegistrationInfoRead(d, meta)
}

func resourceVirtualDesktopHostPoolRegistrationInfoRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceVirtualDesktopHostPoolRegistrationInfoDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

	tokenOperation := hostpool.RegistrationTokenOperationDelete
	payload := hostpool.HostPoolPatch{
		Properties: &hostpool.HostPoolPatchProperties{
			RegistrationInfo: &hostpool.RegistrationInfoPatch{
				RegistrationTokenOperation: &tokenOperation,
			},
		},
	}

	if _, err := client.Update(ctx, hostPoolId, payload); err != nil {
		return fmt.Errorf("removing Registration Token from %s: %+v", hostPoolId, err)
	}

	return nil
}
