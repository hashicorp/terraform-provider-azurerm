// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool_registration_info

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var hostPoolResourceType = "azurerm_virtual_desktop_host_pool"

func RegisterResource() *pluginsdk.Resource {
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

		Schema: virtualDesktopHostPoolRegistrationInfoSchema(),
	}
}
