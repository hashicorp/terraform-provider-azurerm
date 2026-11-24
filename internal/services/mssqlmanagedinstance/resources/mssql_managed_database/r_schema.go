// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_database

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	atLeastOneOf := []string{
		"long_term_retention_policy.0.weekly_retention", "long_term_retention_policy.0.monthly_retention",
		"long_term_retention_policy.0.yearly_retention", "long_term_retention_policy.0.week_of_year",
	}

	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Managed Database to create.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlManagedInstanceDatabaseName,
		},

		"managed_instance_id": {
			Description:  "The ID of the Azure SQL Managed Instance on which to create this Managed Database.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlManagedInstanceIDInsensitively),
		},

		"long_term_retention_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// WeeklyRetention - The weekly retention policy for an LTR backup in an ISO 8601 format.
					"weekly_retention": {
						Description:  "The weekly retention policy for an LTR backup in an ISO 8601 format.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "PT0S",
						ValidateFunc: validation.ISO8601Duration,
						AtLeastOneOf: atLeastOneOf,
					},

					// MonthlyRetention - The monthly retention policy for an LTR backup in an ISO 8601 format.
					"monthly_retention": {
						Description:  "The monthly retention policy for an LTR backup in an ISO 8601 format.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "PT0S",
						ValidateFunc: validation.ISO8601Duration,
						AtLeastOneOf: atLeastOneOf,
					},

					// YearlyRetention - The yearly retention policy for an LTR backup in an ISO 8601 format.
					"yearly_retention": {
						Description:  "The yearly retention policy for an LTR backup in an ISO 8601 format.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "PT0S",
						ValidateFunc: validation.ISO8601Duration,
						AtLeastOneOf: atLeastOneOf,
					},

					// WeekOfYear - The week of year to take the yearly backup.
					"week_of_year": {
						Description:  "The week of year to take the yearly backup. Value has to be between `1` and `52`.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(0, 52),
						AtLeastOneOf: atLeastOneOf,
					},
				},
			},
		},

		"short_term_retention_days": {
			Description:  "The backup retention period in days. This is how many days Point-in-Time Restore will be supported.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 35),
			Default:      7,
		},

		"short_term_retention_immutability_locked": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"point_in_time_restore": {
			Type:          schema.TypeList,
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			ConflictsWith: []string{"recover_database_id"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"restore_point_in_time": {
						Description:      "The point in time for the restore from `source_database_id`.",
						Type:             pluginsdk.TypeString,
						Required:         true,
						ForceNew:         true,
						DiffSuppressFunc: suppress.RFC3339Time,
						ValidateFunc:     validation.IsRFC3339Time,
					},
					"source_database_id": {
						Description:  "The source database ID that will be used to restore from.",
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.Any(validation.AsGeneratedID(commonids.ParseSqlManagedInstanceDatabaseIDInsensitively), validate.RestorableDatabaseID),
					},
				},
			},
		},

		"recover_database_id": {
			Description:   "The ID of the Recoverable Managed Database from which to perform a cross-region / recovery restore.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"point_in_time_restore"},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			if oldVal, newVal := rd.GetChange("short_term_retention_immutability_locked"); oldVal.(bool) && !newVal.(bool) {
				// the service does not support unlock.
				if err := rd.ForceNew("short_term_retention_immutability_locked"); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
