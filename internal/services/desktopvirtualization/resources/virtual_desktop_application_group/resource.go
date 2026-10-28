// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_application_group

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/applicationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/resources/virtual_desktop_application_group/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopApplicationGroupCreateUpdate,
		Read:   resourceVirtualDesktopApplicationGroupRead,
		Update: resourceVirtualDesktopApplicationGroupCreateUpdate,
		Delete: resourceVirtualDesktopApplicationGroupDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := applicationgroup.ParseApplicationGroupID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApplicationGroupV0ToV1{},
		}),

		Schema: virtualDesktopApplicationGroupSchema(),
	}
}
