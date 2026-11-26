// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_data_lake_storage

import (
	"regexp"

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
			Description: "The name of the Backup Policy for the Azure Backup Policy Data Lake Storage.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{2,149}$"),
				"`name` must be 3 - 150 characters long, contain only letters, numbers and hyphens(-), and cannot start with a number or hyphen.",
			),
		},

		"data_protection_backup_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(pointer.To(basebackuppolicyresources.BackupVaultId{}), "Data Protection Backup Vault"),

		"backup_schedule": {
			Description: "A list of repeating time interval, also known as the backup schedule. It supports daily & weekly backup. It should follow [`ISO 8601` recurring time interval format](https://en.wikipedia.org/wiki/ISO_8601#Recurring_intervals), for example: `R/2021-05-23T02:30:00+00:00/P1W`.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			ForceNew:    true,
			MinItems:    1,
			MaxItems:    5,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.ISO8601RepeatingTime,
			},
		},

		"default_retention_duration": {
			Description:  "The retention duration up to which the backups are to be retained in the data stores. It should follow `ISO 8601` duration format.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"retention_rule": {
			Description: "The priority of each rule is determined by its order in the list, where the first rule has the highest priority.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the retention rule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"duration": {
						Description:  "The retention duration up to which the backups are to be retained in the data stores. It should follow `ISO 8601` duration format.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.ISO8601Duration,
					},

					"absolute_criteria": {
						Description:  "The absolute criteria for the retention rule. These values mean the first successful backup of the day/week/month/year.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForAbsoluteMarker(), false),
					},

					"days_of_week": {
						Description: "A list of days of the week on which the retention rule applies.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						ForceNew:    true,
						MinItems:    1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForDayOfWeek(), false),
						},
					},

					"months_of_year": {
						Description: "A list of months of the year on which the retention rule applies.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						ForceNew:    true,
						MinItems:    1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForMonth(), false),
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
						Description: "A list of weeks of the month on which the retention rule applies.",
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

		"time_zone": {
			Description:  "The Time Zone which should be used by the backup schedule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.BackupPolicyDataLakeStorageTimeZone(),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
