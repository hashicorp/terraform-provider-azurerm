// Copyright IBM Corp. 2014, 2025
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
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-a-zA-Z0-9]{3,150}$"),
				"DataProtection BackupPolicy name must be 3 - 150 characters long, contain only letters, numbers and hyphens.",
			),
		},

		"vault_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: basebackuppolicyresources.ValidateBackupVaultID,
		},

		"backup_repeating_time_intervals": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"operational_default_retention_duration": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			AtLeastOneOf: []string{"operational_default_retention_duration", "vault_default_retention_duration"},
			ValidateFunc: validation.ISO8601Duration,
		},

		"time_zone": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"vault_default_retention_duration": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			AtLeastOneOf: []string{"operational_default_retention_duration", "vault_default_retention_duration"},
			RequiredWith: []string{"backup_repeating_time_intervals"},
			ValidateFunc: validation.ISO8601Duration,
		},

		"retention_rule": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"vault_default_retention_duration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"criteria": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"absolute_criteria": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice(
										basebackuppolicyresources.PossibleValuesForAbsoluteMarker(), false,
									),
								},

								"days_of_month": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
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
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsMonth(false),
									},
								},

								"scheduled_backup_times": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsRFC3339Time,
									},
								},

								"weeks_of_month": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
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
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										// confirmed with the service team that currently only `VaultStore` is supported.
										// However, since `ArchiveStore` may be supported in the future, it is open to user specification.
										string(basebackuppolicyresources.DataStoreTypesVaultStore),
									}, false),
								},

								"duration": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.ISO8601Duration,
								},
							},
						},
					},

					"priority": {
						Type:     pluginsdk.TypeInt,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},
	}
}
