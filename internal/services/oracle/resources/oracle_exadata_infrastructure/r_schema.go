// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_exadata_infrastructure

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// Azure
		"location": commonschema.Location(),

		"name": {
			Description:  "The name which should be used for this Cloud Exadata Infrastructure.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// Required
		"compute_count": {
			Description:  "The number of compute servers for the Cloud Exadata Infrastructure.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(2, 32),
		},

		"database_server_type": {
			Description:  "The database server model type of the cloud Exadata infrastructure resource.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 255),
		},

		"display_name": {
			Description:  "The user-friendly name for the Cloud Exadata Infrastructure resource. The name does not need to be unique.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CloudVMClusterName,
		},

		"shape": {
			Description: "The shape of the ODB@A infrastructure resource.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"storage_count": {
			Description:  "The number of storage servers for the Cloud Exadata Infrastructure.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(3, 64),
		},

		// Optional
		"customer_contacts": {
			Description: "The email address used by Oracle to send notifications regarding databases and infrastructure.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"maintenance_window": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"days_of_week": {
						Description: "Days during the week when maintenance should be performed.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.DaysOfWeek,
						},
					},

					"hours_of_day": {
						Description: "The window of hours during the day when maintenance should be performed. The window is a 4 hour slot.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validate.HoursOfDay,
						},
					},

					"lead_time_in_weeks": {
						Description:  "Lead time window allows user to set a lead time to prepare for a down time. The lead time is in weeks.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ForceNew:     true,
						ValidateFunc: validation.IntBetween(1, 4),
					},

					"months": {
						Description: "Months during the year when maintenance should be performed.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.Month,
						},
					},

					"patching_mode": {
						Description:  "Cloud Exadata Infrastructure node patching method, either `ROLLING` or `NONROLLING`. IMPORTANT: Non-rolling infrastructure patching involves system down time. See [Oracle-Managed Infrastructure Maintenance Updates](https://docs.cloud.oracle.com/iaas/Content/Database/Concepts/examaintenance.htm#Oracle) for more information.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ForceNew:     true,
						ValidateFunc: validate.PatchingMode,
					},

					"preference": {
						Description:  "The maintenance window scheduling preference.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ForceNew:     true,
						ValidateFunc: validate.Preference,
					},

					"weeks_of_month": {
						Description: "Weeks during the month when maintenance should be performed. Weeks start on the 1st, 8th, 15th, and 22nd days of the month, and have a duration of 7 days. Weeks start and end based on calendar dates, not days of the week. For example, to allow maintenance during the 2nd week of the month (from the 8th day to the 14th day of the month), use the value 2. Maintenance cannot be scheduled for the fifth week of months that contain more than 28 days. This parameter works in conjunction with the daysOfWeek and hoursOfDay parameters to allow you to specify specific days of the week and hours that maintenance will be performed.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(1, 4),
						},
					},
				},
			},
		},

		"storage_server_type": {
			Description:  "The storage server model type of the cloud Exadata infrastructure resource.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 255),
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
