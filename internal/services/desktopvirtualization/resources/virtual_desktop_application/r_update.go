// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_application

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/application"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/applicationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopApplicationCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	applicationGroup, _ := applicationgroup.ParseApplicationGroupID(d.Get("application_group_id").(string))
	id := application.NewApplicationID(subscriptionId, applicationGroup.ResourceGroupName, applicationGroup.ApplicationGroupName, d.Get("name").(string))

	locks.ByName(id.ApplicationName, applicationType)
	defer locks.UnlockByName(id.ApplicationName, applicationType)

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_virtual_desktop_application", id.ID())
			}
		}
	}

	payload := application.Application{
		Properties: application.ApplicationProperties{
			FriendlyName:         pointer.To(d.Get("friendly_name").(string)),
			Description:          pointer.To(d.Get("description").(string)),
			FilePath:             pointer.To(d.Get("path").(string)),
			CommandLineSetting:   application.CommandLineSetting(d.Get("command_line_argument_policy").(string)),
			CommandLineArguments: pointer.To(d.Get("command_line_arguments").(string)),
			ShowInPortal:         pointer.To(d.Get("show_in_portal").(bool)),
			IconPath:             pointer.To(d.Get("icon_path").(string)),
			IconIndex:            pointer.To(int64(d.Get("icon_index").(int))),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}
	return resourceVirtualDesktopApplicationRead(d, meta)
}
