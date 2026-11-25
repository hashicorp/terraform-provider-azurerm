// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func managementGroupSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name or UUID for this Management Group, which needs to be unique across your tenant. A new UUID will be generated if not provided.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validate.ManagementGroupName,
		},

		"display_name": {
			Description: "A friendly name for this Management Group. If not specified, this will be the same as the `name`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"tenant_scoped_id": {
			Description: "The Management Group ID with the Tenant ID prefix.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"parent_management_group_id": {
			Description:  "The ID of the Parent Management Group.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: commonids.ValidateManagementGroupID,
		},

		"subscription_ids": {
			Description: "A list of Subscription GUIDs which should be assigned to the Management Group.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
			Set: pluginsdk.HashString,
		},
	}
}
