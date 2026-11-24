// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_backup_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the NetApp Backup Policy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VolumeQuotaRuleName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"account_name": {
			Description:  "The name of the NetApp account in which the NetApp Policy should be created under.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},

		"tags": commonschema.Tags(),

		"daily_backups_to_keep": {
			Description:  "Provides the number of daily backups to keep.",
			Type:         pluginsdk.TypeInt,
			Default:      2,
			Optional:     true,
			ValidateFunc: validation.IntBetween(2, 1019),
		},

		"weekly_backups_to_keep": {
			Description: "Provides the number of weekly backups to keep. The minimum is `0` and the maximum is `1019`.",
			Type:        pluginsdk.TypeInt,
			// Default in API is 0 but initially set to 1 so keeping to avoid breaking changes
			Default:      1,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 1019),
		},

		"monthly_backups_to_keep": {
			Description: "Provides the number of monthly backups to keep. The minimum is `0` and the maximum is `1019`.",
			Type:        pluginsdk.TypeInt,
			// Default in API is 0 but initially set to 1 so keeping to avoid breaking changes
			Default:      1,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 1019),
		},

		"enabled": {
			Description: "Whether the Backup Policy is enabled.",
			Type:        pluginsdk.TypeBool,
			Default:     true,
			Optional:    true,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
