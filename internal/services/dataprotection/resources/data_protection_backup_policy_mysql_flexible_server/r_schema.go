// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_mysql_flexible_server

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BackupPolicyMySQLFlexibleServerName,
		},

		"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(pointer.To(basebackuppolicyresources.BackupVaultId{})),

		"backup_repeating_time_intervals": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.ISO8601RepeatingTime,
			},
		},

		"default_retention_rule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"life_cycle": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_store_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										// Confirmed with the service team that current possible value only support `VaultStore`.
										// However, considering that `ArchiveStore` will be supported in the future, it would be exposed for user specification.
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
				},
			},
		},

		"retention_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"criteria": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"absolute_criteria": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForAbsoluteMarker(), false),
								},

								"days_of_week": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForDayOfWeek(), false),
									},
								},

								"months_of_year": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForMonth(), false),
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
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_store_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										// Confirmed with the service team that currently only `VaultStore` is supported.
										// However, considering that `ArchiveStore` will be supported in the future, it would be exposed for user specification.
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

		"time_zone": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.BackupPolicyMySQLFlexibleServerTimeZone(),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
