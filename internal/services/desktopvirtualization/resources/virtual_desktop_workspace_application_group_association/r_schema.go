// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_workspace_application_group_association

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func virtualDesktopWorkspaceApplicationGroupAssociationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Description:  "The resource ID for the Virtual Desktop Workspace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspace.ValidateWorkspaceID,
		},

		"application_group_id": {
			Description:  "The resource ID for the Virtual Desktop Application Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: applicationgroup.ValidateApplicationGroupID,
		},
	}
}
