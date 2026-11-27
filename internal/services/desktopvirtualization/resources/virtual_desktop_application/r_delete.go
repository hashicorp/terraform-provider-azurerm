// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_application

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/application"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	desktopvirtualizationClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopApplicationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := desktopvirtualizationClient.From(meta).ApplicationsClient

	id, err := application.ParseApplicationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ApplicationName, applicationType)
	defer locks.UnlockByName(id.ApplicationName, applicationType)

	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()
	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
