// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_slot_custom_hostname_binding

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceSlotCustomHostnameBindingDelete(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).WebAppsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := webapps.ParseSlotHostNameBindingID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.HostNameBindingName, appServiceSlotCustomHostnameBindingResourceName)
	defer locks.UnlockByName(id.HostNameBindingName, appServiceSlotCustomHostnameBindingResourceName)

	if _, err := client.DeleteHostNameBindingSlot(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
