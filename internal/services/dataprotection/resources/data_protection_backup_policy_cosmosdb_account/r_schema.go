// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_cosmosdb_account

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	backupPolicyCosmosdbAccountFrequencyDaily  = "Daily"
	backupPolicyCosmosdbAccountFrequencyWeekly = "Weekly"
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

		"data_protection_backup_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(pointer.To(basebackuppolicyresources.BackupVaultId{})),

		"backup_schedule": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validate.BackupPolicyCosmosdbAccountBackupSchedule(),
				validation.ISO8601RepeatingTime,
			),
		},

		"default_retention_duration": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"backup_frequency": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  backupPolicyCosmosdbAccountFrequencyDaily,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				backupPolicyCosmosdbAccountFrequencyDaily,
				backupPolicyCosmosdbAccountFrequencyWeekly,
			}, false),
		},

		"monthly_retention_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"duration": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.ISO8601Duration,
					},

					"days_of_week": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						MaxItems: 1,
						RequiredWith: []string{
							"monthly_retention_rule.0.days_of_week",
							"monthly_retention_rule.0.weeks_of_month",
						},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForDayOfWeek(), false),
						},
					},

					"weeks_of_month": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						MaxItems: 5,
						RequiredWith: []string{
							"monthly_retention_rule.0.days_of_week",
							"monthly_retention_rule.0.weeks_of_month",
						},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForWeekNumber(), false),
						},
					},
				},
			},
		},

		"weekly_retention_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"duration": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.ISO8601Duration,
					},

					"days_of_week": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						MaxItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForDayOfWeek(), false),
						},
					},
				},
			},
		},

		"yearly_retention_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"duration": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.ISO8601Duration,
					},

					"days_of_week": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						MaxItems: 1,
						RequiredWith: []string{
							"yearly_retention_rule.0.days_of_week",
							"yearly_retention_rule.0.weeks_of_month",
							"yearly_retention_rule.0.months_of_year",
						},
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
						MaxItems: 12,
						RequiredWith: []string{
							"yearly_retention_rule.0.days_of_week",
							"yearly_retention_rule.0.weeks_of_month",
							"yearly_retention_rule.0.months_of_year",
						},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForMonth(), false),
						},
					},

					"weeks_of_month": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						MaxItems: 5,
						RequiredWith: []string{
							"yearly_retention_rule.0.days_of_week",
							"yearly_retention_rule.0.weeks_of_month",
							"yearly_retention_rule.0.months_of_year",
						},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForWeekNumber(), false),
						},
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"incremental_backup_schedules": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
