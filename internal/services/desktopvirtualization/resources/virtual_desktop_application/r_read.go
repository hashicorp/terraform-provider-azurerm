// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_application

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/application"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/applicationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopApplicationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := application.ParseApplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ApplicationName)
	d.Set("application_group_id", applicationgroup.NewApplicationGroupID(id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName).ID())

	if model := resp.Model; model != nil {
		props := model.Properties

		d.Set("friendly_name", props.FriendlyName)
		d.Set("description", props.Description)
		d.Set("path", props.FilePath)
		d.Set("command_line_argument_policy", string(props.CommandLineSetting))
		d.Set("command_line_arguments", props.CommandLineArguments)
		d.Set("show_in_portal", props.ShowInPortal)
		d.Set("icon_path", props.IconPath)
		d.Set("icon_index", props.IconIndex)
	}

	return nil
}

var applicationType = "azurerm_virtual_desktop_application"
