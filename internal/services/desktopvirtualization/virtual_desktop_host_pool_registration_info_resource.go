package desktopvirtualization

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/desktopvirtualization/mgmt/2021-09-03-preview/desktopvirtualization"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/validate"
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
			_, err := parse.HostPoolRegistrationInfoID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(hostpoolRegistrationInfoCustomDiff),

		Schema: map[string]*pluginsdk.Schema{
			"hostpool_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.HostPoolID,
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

	hostpoolId, err := parse.HostPoolID(d.Get("hostpool_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(hostpoolId.Name, hostpoolResourceType)
	defer locks.UnlockByName(hostpoolId.Name, hostpoolResourceType)

	// This is a virtual resource so the last segment is hardcoded
	id := parse.NewHostPoolRegistrationInfoID(hostpoolId.SubscriptionId, hostpoolId.ResourceGroup, hostpoolId.Name, "default")

	hostpool, err := client.Get(ctx, id.ResourceGroup, id.HostPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(hostpool.Response) {
			return fmt.Errorf("%s could not be found: %s", hostpoolId, err)
		}
		return fmt.Errorf("reading %s: %s", hostpoolId, err)
	}

	hostpoolPatch := &desktopvirtualization.HostPoolPatch{}
	hostpoolPatch.HostPoolPatchProperties = &desktopvirtualization.HostPoolPatchProperties{}
	hostpoolPatch.HostPoolPatchProperties.RegistrationInfo = &desktopvirtualization.RegistrationInfoPatch{}

	expdt, _ := date.ParseTime(time.RFC3339, d.Get("expiration_date").(string))
	hostpoolPatch.HostPoolPatchProperties.RegistrationInfo.ExpirationTime = &date.Time{
		Time: expdt,
	}
	hostpoolPatch.HostPoolPatchProperties.RegistrationInfo.RegistrationTokenOperation = desktopvirtualization.RegistrationTokenOperationUpdate

	if _, err := client.Update(ctx, id.ResourceGroup, id.HostPoolName, hostpoolPatch); err != nil {
		return fmt.Errorf("Creating Virtual Desktop Host Pool Registration Info %q (Resource Group %q): %+v", id.HostPoolName, id.ResourceGroup, err)
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

	resp, err := client.RetrieveRegistrationToken(ctx, id.ResourceGroup, id.HostPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual Desktop Host Pool %q was not found in Resource Group %q - removing from state!", id.HostPoolName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Making Read request on Virtual Desktop Host Pool %q (Resource Group %q): %+v", id.HostPoolName, id.ResourceGroup, err)
	}

	if resp.ExpirationTime == nil || resp.Token == nil {
		log.Printf("HostPool is missing registration info - marking as gone")
		d.SetId("")
		return nil
	}

	d.Set("hostpool_id", parse.NewHostPoolID(id.SubscriptionId, id.ResourceGroup, id.HostPoolName).ID())
	d.Set("expiration_date", resp.ExpirationTime.Format(time.RFC3339))
	d.Set("token", resp.Token)

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

	hostpoolId := parse.NewHostPoolID(id.SubscriptionId, id.ResourceGroup, id.HostPoolName)

	locks.ByName(hostpoolId.Name, hostpoolResourceType)
	defer locks.UnlockByName(hostpoolId.Name, hostpoolResourceType)

	resp, err := client.Get(ctx, id.ResourceGroup, id.HostPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual Desktop Host Pool %q was not found in Resource Group %q - removing from state!", id.HostPoolName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Making Read request on Virtual Desktop Host Pool %q (Resource Group %q): %+v", id.HostPoolName, id.ResourceGroup, err)
	}

	regInfo, err := client.RetrieveRegistrationToken(ctx, id.ResourceGroup, id.HostPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual Desktop Host Pool %q Registration Info was not found in Resource Group %q - removing from state!", id.HostPoolName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Making Read request on Virtual Desktop Host Pool %q (Resource Group %q): %+v", id.HostPoolName, id.ResourceGroup, err)
	}
	if regInfo.ExpirationTime == nil {
		log.Printf("[INFO] RegistrationInfo for %s was nil, registrationInfo already deleted - removing %s from state", hostpoolId.ID(), id)
		return nil
	}

	hostpoolPatch := &desktopvirtualization.HostPoolPatch{}
	hostpoolPatch.HostPoolPatchProperties = &desktopvirtualization.HostPoolPatchProperties{}
	hostpoolPatch.HostPoolPatchProperties.RegistrationInfo = &desktopvirtualization.RegistrationInfoPatch{}
	hostpoolPatch.HostPoolPatchProperties.RegistrationInfo.RegistrationTokenOperation = desktopvirtualization.RegistrationTokenOperationDelete

	if _, err := client.Update(ctx, id.ResourceGroup, id.HostPoolName, hostpoolPatch); err != nil {
		return fmt.Errorf("deleting Virtual Desktop Host Pool Registration Info %q (Resource Group %q): %+v", id.HostPoolName, id.ResourceGroup, err)
	}

	return nil
}
