// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_workspace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopWorkspaceCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workspace.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_virtual_desktop_workspace", id.ID())
			}
		}
	}

	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]any)

	payload := workspace.Workspace{
		Location: location,
		Tags:     tags.Expand(t),
		Properties: &workspace.WorkspaceProperties{
			Description:  pointer.To(d.Get("description").(string)),
			FriendlyName: pointer.To(d.Get("friendly_name").(string)),
		},
	}

	publicNetworkAccess := workspace.PublicNetworkAccessEnabled

	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = workspace.PublicNetworkAccessDisabled
	}

	payload.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualDesktopWorkspaceRead(d, meta)
}
