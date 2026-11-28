// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the NetApp Snapshot Policy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SnapshotName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"account_name": {
			Description:  "The name of the NetApp Account in which the NetApp Snapshot Policy should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},

		"enabled": {
			Description: "Whether the NetApp Snapshot Policy is enabled.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},

		"hourly_schedule": {
			Description: "Sets an hourly snapshot schedule.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"snapshots_to_keep": {
						Description:  "How many hourly snapshots to keep, valid range is from 0 to 255.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 255),
					},

					"minute": {
						Description:  "Minute of the hour that the snapshots will be created, valid range is from 0 to 59.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 59),
					},
				},
			},
		},

		"daily_schedule": {
			Description: "Sets a daily snapshot schedule.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"snapshots_to_keep": {
						Description:  "How many hourly snapshots to keep, valid range is from 0 to 255.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 255),
					},

					"hour": {
						Description:  "Hour of the day that the snapshots will be created, valid range is from 0 to 23.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 23),
					},

					"minute": {
						Description:  "Minute of the hour that the snapshots will be created, valid range is from 0 to 59.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 59),
					},
				},
			},
		},

		"weekly_schedule": {
			Description: "Sets a weekly snapshot schedule.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"snapshots_to_keep": {
						Description:  "How many hourly snapshots to keep, valid range is from 0 to 255.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 255),
					},

					"days_of_week": {
						Description: "List of the week days using English names when the snapshots will be created.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						MaxItems:    7,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsDayOfTheWeek(false),
						},
					},

					"hour": {
						Description:  "Hour of the day that the snapshots will be created, valid range is from 0 to 23.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 23),
					},

					"minute": {
						Description:  "Minute of the hour that the snapshots will be created, valid range is from 0 to 59.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 59),
					},
				},
			},
		},

		"monthly_schedule": {
			Description: "Sets a monthly snapshot schedule.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"snapshots_to_keep": {
						Description:  "How many hourly snapshots to keep, valid range is from 0 to 255.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 255),
					},

					"days_of_month": {
						Description: "List of the days of the month when the snapshots will be created, valid range is from 1 to 30.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						MaxItems:    30,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(1, 30),
						},
					},

					"hour": {
						Description:  "Hour of the day that the snapshots will be created, valid range is from 0 to 23.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 23),
					},

					"minute": {
						Description:  "Minute of the hour that the snapshots will be created, valid range is from 0 to 59.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 59),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			for _, schedule := range []string{"hourly_schedule", "daily_schedule", "weekly_schedule", "monthly_schedule"} {
				if oldVal, newVal := metadata.ResourceDiff.GetChange(schedule); len(oldVal.([]any)) > 0 && len(newVal.([]any)) == 0 {
					if err := metadata.ResourceDiff.ForceNew(schedule); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}
