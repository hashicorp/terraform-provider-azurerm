// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_slot_custom_hostname_binding

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceSlotCustomHostnameBindingCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	slotId, err := webapps.ParseSlotID(d.Get("app_service_slot_id").(string))
	if err != nil {
		return err
	}

	sslState := d.Get("ssl_state").(string)
	thumbprint := d.Get("thumbprint").(string)

	id := webapps.NewSlotHostNameBindingID(slotId.SubscriptionId, slotId.ResourceGroupName, slotId.SiteName, slotId.SlotName, d.Get("hostname").(string))

	locks.ByName(id.HostNameBindingName, appServiceSlotCustomHostnameBindingResourceName)
	defer locks.UnlockByName(id.HostNameBindingName, appServiceSlotCustomHostnameBindingResourceName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetHostNameBindingSlot(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_app_service_slot_custom_hostname_binding", id.ID())
		}
	}

	payload := webapps.HostNameBinding{
		Properties: &webapps.HostNameBindingProperties{
			SiteName: pointer.To(id.SiteName),
		},
	}

	if sslState != "" {
		payload.Properties.SslState = pointer.ToEnum[webapps.SslState](sslState)
	}

	if thumbprint != "" {
		payload.Properties.Thumbprint = pointer.To(thumbprint)
	}

	if _, err := client.CreateOrUpdateHostNameBindingSlot(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceSlotCustomHostnameBindingRead(d, meta)
}
