// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_application_group

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/desktop"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopApplicationGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applicationgroup.ParseApplicationGroupID(d.Id())
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

	d.Set("name", id.ApplicationGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		props := model.Properties

		d.Set("friendly_name", props.FriendlyName)
		d.Set("description", props.Description)
		d.Set("type", string(props.ApplicationGroupType))
		defaultDesktopDisplayName := ""
		if props.ApplicationGroupType == applicationgroup.ApplicationGroupTypeDesktop {
			desktopClient := meta.(*clients.Client).DesktopVirtualization.DesktopsClient
			// default desktop name created for Application Group is 'sessionDesktop'
			desktopId := desktop.NewDesktopID(id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName, "sessionDesktop")
			desktopResp, err := desktopClient.Get(ctx, desktopId)
			if err != nil {
				if !response.WasNotFound(desktopResp.HttpResponse) {
					return fmt.Errorf("retrieving default desktop for %s: %+v", *id, err)
				}
			}
			// if the default desktop was found then set the display name attribute
			if desktopModel := desktopResp.Model; desktopModel != nil && desktopModel.Properties != nil && desktopModel.Properties.FriendlyName != nil {
				defaultDesktopDisplayName = *desktopModel.Properties.FriendlyName
			}
		}
		d.Set("default_desktop_display_name", defaultDesktopDisplayName)

		hostPoolId, err := hostpool.ParseHostPoolIDInsensitively(props.HostPoolArmPath)
		if err != nil {
			return fmt.Errorf("parsing Host Pool ID %q: %+v", props.HostPoolArmPath, err)
		}
		d.Set("host_pool_id", hostPoolId.ID())

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

var applicationGroupType = "azurerm_virtual_desktop_application_group"
