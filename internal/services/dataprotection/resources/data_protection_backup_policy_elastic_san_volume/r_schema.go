// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_elastic_san_volume

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{2,149}$"),
				"`name` must be 3 - 150 characters long, contain only letters, numbers and hyphens, and cannot start with a number or hyphen",
			),
		},

		"data_protection_backup_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(
			pointer.To(basebackuppolicyresources.BackupVaultId{}),
			"Data Protection Backup Vault",
		),

		"daily_backup_schedule": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ExactlyOneOf: []string{
				"daily_backup_schedule",
				"weekly_backup_schedule",
			},
			ValidateFunc: validation.All(
				validation.ISO8601RepeatingTime,
				validation.StringMatch(regexp.MustCompile(`/P1D$`), "must use a daily (`P1D`) interval"),
			),
		},

		"weekly_backup_schedule": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ExactlyOneOf: []string{
				"daily_backup_schedule",
				"weekly_backup_schedule",
			},
			ValidateFunc: validation.All(
				validation.ISO8601RepeatingTime,
				validation.StringMatch(regexp.MustCompile(`/P1W$`), "must use a weekly (`P1W`) interval"),
			),
		},

		"default_retention_duration": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.ISO8601DurationBetween("P7D", "P365D"),
			),
		},

		"daily_retention_duration": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"weekly_backup_schedule"},
			ValidateFunc: validation.All(
				validation.ISO8601DurationBetween("P7D", "P365D"),
			),
		},

		"weekly_retention_duration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.ISO8601DurationBetween("P4W", "P52W"),
			),
		},

		"monthly_retention_duration": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"daily_backup_schedule"},
			ValidateFunc: validation.All(
				validation.ISO8601DurationBetween("P1M", "P11M"),
			),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
