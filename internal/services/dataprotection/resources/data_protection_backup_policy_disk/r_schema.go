// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_disk

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataProtectionBackupPolicyDiskSchema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "The name which should be used for this Backup Policy Disk.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-a-zA-Z0-9]{3,150}$"),
				"DataProtection BackupPolicy name must be 3 - 150 characters long, contain only letters, numbers and hyphens.",
			),
		},

		"vault_id": {
			Description:  "The ID of the Backup Vault within which the Backup Policy Disk should exist.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: basebackuppolicyresources.ValidateBackupVaultID,
		},

		"backup_repeating_time_intervals": {
			Description: "A list of repeating time interval. It should follow `ISO 8601` repeating time interval.",
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			MinItems:    1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"default_retention_duration": {
			Description:  "The duration of default retention rule. It should follow `ISO 8601` duration format.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"retention_rule": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Description: "The name which should be used for this retention rule.",
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
					},

					"duration": {
						Description:  "Duration of deletion after given timespan. It should follow `ISO 8601` duration format.",
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.ISO8601Duration,
					},

					"criteria": {
						Type:     schema.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"absolute_criteria": {
									Description:  "These values mean the first successful backup of the day/week/month/year.",
									Type:         schema.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForAbsoluteMarker(), false),
								},
							},
						},
					},

					"priority": {
						Description: "Retention Tag priority.",
						Type:        schema.TypeInt,
						Required:    true,
						ForceNew:    true,
					},
				},
			},
		},

		"time_zone": {
			Description:  "The Time Zone which should be used by the backup schedule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
