// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_blob_storage

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataProtectionBackupPolicyBlobStorageSchema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "The name which should be used for this Backup Policy Blob Storage.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-a-zA-Z0-9]{3,150}$"),
				"DataProtection BackupPolicy name must be 3 - 150 characters long, contain only letters, numbers and hyphens.",
			),
		},

		"vault_id": {
			Description:  "The ID of the Backup Vault within which the Backup Policy Blob Storage should exist.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: basebackuppolicyresources.ValidateBackupVaultID,
		},

		"backup_repeating_time_intervals": {
			Description: "A list of repeating time interval. It should follow `ISO 8601` repeating time interval.",
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MinItems:    1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"operational_default_retention_duration": {
			Description:  "The duration of operational default retention rule. It should follow `ISO 8601` duration format.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.ISO8601Duration,
			AtLeastOneOf: []string{"operational_default_retention_duration", "vault_default_retention_duration"},
		},

		"time_zone": {
			Description:  "The Time Zone which should be used by the backup schedule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"vault_default_retention_duration": {
			Description:  "The duration of vault default retention rule. It should follow `ISO 8601` duration format.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.ISO8601Duration,
			RequiredWith: []string{"backup_repeating_time_intervals"},
			AtLeastOneOf: []string{"operational_default_retention_duration", "vault_default_retention_duration"},
		},

		"retention_rule": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"vault_default_retention_duration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name which should be used for this retention rule.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
					},

					"criteria": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"absolute_criteria": {
									Description: "These values mean the first successful backup of the day/week/month/year.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
									ForceNew:    true,
									ValidateFunc: validation.StringInSlice(
										basebackuppolicyresources.PossibleValuesForAbsoluteMarker(), false,
									),
								},

								"days_of_month": {
									Description: "`0` for last day within the month.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									ForceNew:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeInt,
										ValidateFunc: validation.Any(
											validation.IntBetween(0, 28),
										),
									},
								},

								"days_of_week": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsDayOfTheWeek(false),
									},
								},

								"months_of_year": {
									Description: "When this property is specified, exactly one of the following must also be set: `days_of_month`, `days_of_week`.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									ForceNew:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsMonth(false),
									},
								},

								"scheduled_backup_times": {
									Description: "A list of backup times for backup in the `RFC3339` format.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									ForceNew:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsRFC3339Time,
									},
								},

								"weeks_of_month": {
									Description: "When this property is specified, exactly one of the following must also be set: `days_of_month`, `days_of_week`.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									ForceNew:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForWeekNumber(), false),
									},
								},
							},
						},
					},

					"life_cycle": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_store_type": {
									Description: "The type of data store.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ForceNew:    true,
									ValidateFunc: validation.StringInSlice([]string{
										// confirmed with the service team that currently only `VaultStore` is supported.
										// However, since `ArchiveStore` may be supported in the future, it is open to user specification.
										string(basebackuppolicyresources.DataStoreTypesVaultStore),
									}, false),
								},

								"duration": {
									Description:  "The retention duration up to which the backups are to be retained in the data stores. It should follow `ISO 8601` duration format.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.ISO8601Duration,
								},
							},
						},
					},

					"priority": {
						Description: "The priority of the rule. The priority number must be unique for each rule. The lower the priority number, the higher the priority of the rule.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
						ForceNew:    true,
					},
				},
			},
		},
	}
}
