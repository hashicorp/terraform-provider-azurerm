// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_workspace_application_group_association

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/resources/virtual_desktop_workspace_application_group_association/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var workspaceResourceType = "azurerm_virtual_desktop_workspace"

var applicationGroupType = "azurerm_virtual_desktop_application_group"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopWorkspaceApplicationGroupAssociationCreate,
		Read:   resourceVirtualDesktopWorkspaceApplicationGroupAssociationRead,
		Delete: resourceVirtualDesktopWorkspaceApplicationGroupAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WorkspaceApplicationGroupAssociationID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WorkspaceApplicationGroupAssociationV0ToV1{},
		}),

		Schema: virtualDesktopWorkspaceApplicationGroupAssociationSchema(),
	}
}
