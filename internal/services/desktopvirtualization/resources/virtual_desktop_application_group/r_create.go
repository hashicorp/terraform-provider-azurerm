// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_application_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/desktop"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopApplicationGroupCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(name, applicationGroupType)
	defer locks.UnlockByName(name, applicationGroupType)

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := applicationgroup.NewApplicationGroupID(subscriptionId, resourceGroup, name)
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_virtual_desktop_application_group", id.ID())
			}
		}
	}

	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]any)

	payload := applicationgroup.ApplicationGroup{
		Location: location,
		Tags:     tags.Expand(t),
		Properties: applicationgroup.ApplicationGroupProperties{
			ApplicationGroupType: applicationgroup.ApplicationGroupType(d.Get("type").(string)),
			FriendlyName:         pointer.To(d.Get("friendly_name").(string)),
			Description:          pointer.To(d.Get("description").(string)),
			HostPoolArmPath:      d.Get("host_pool_id").(string),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating Virtual Desktop Application Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(id.ID())

	if applicationgroup.ApplicationGroupType(d.Get("type").(string)) == applicationgroup.ApplicationGroupTypeDesktop {
		if desktopFriendlyName := pointer.To(d.Get("default_desktop_display_name").(string)); desktopFriendlyName != nil {
			desktopClient := meta.(*clients.Client).DesktopVirtualization.DesktopsClient
			// default desktop name created for Application Group is 'sessionDesktop'
			desktopId := desktop.NewDesktopID(id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName, "sessionDesktop")
			model, err := desktopClient.Get(ctx, desktopId)
			if err != nil {
				if !response.WasNotFound(model.HttpResponse) {
					return fmt.Errorf("retrieving default desktop for %s: %+v", id, err)
				}
			}

			desktopPatch := desktop.DesktopPatch{
				Properties: &desktop.DesktopPatchProperties{
					FriendlyName: desktopFriendlyName,
				},
			}

			if _, err := desktopClient.Update(ctx, desktopId, desktopPatch); err != nil {
				return fmt.Errorf("setting friendly name for default desktop %s: %+v", id, err)
			}
		}
	}

	return resourceVirtualDesktopApplicationGroupRead(d, meta)
}
