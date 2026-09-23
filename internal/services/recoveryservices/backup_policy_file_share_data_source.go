// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceBackupPolicyFileShare() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBackupPolicyFileShareRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: dataSourceBackupPolicyFileShareSchema(),
	}
}

func dataSourceBackupPolicyFileShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := protectionpolicies.NewBackupPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), d.Get("name").(string))

	protectionPolicy, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(protectionPolicy.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on%s: %+v", id, err)
	}

	d.SetId(id.ID())

	model := protectionPolicy.Model
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	if properties, ok := model.Properties.(protectionpolicies.AzureFileShareProtectionPolicy); ok {
		d.Set("timezone", properties.TimeZone)

		if schedule, ok := properties.SchedulePolicy.(protectionpolicies.SimpleSchedulePolicy); ok {
			s, err := flattenBackupProtectionPolicyFileShareSchedule(schedule)
			if err != nil {
				return fmt.Errorf("flattening `backup`: %+v", err)
			}
			if err := d.Set("backup", s); err != nil {
				return fmt.Errorf("setting `backup`: %+v", err)
			}
		}

		if retention, ok := properties.RetentionPolicy.(protectionpolicies.LongTermRetentionPolicy); ok {
			if err := d.Set("retention_daily", flattenBackupProtectionPolicyFileShareRetentionDaily(retention.DailySchedule)); err != nil {
				return fmt.Errorf("setting `retention_daily`: %+v", err)
			}

			if err := d.Set("retention_weekly", flattenBackupProtectionPolicyFileShareRetentionWeekly(retention.WeeklySchedule)); err != nil {
				return fmt.Errorf("setting `retention_weekly`: %+v", err)
			}

			if err := d.Set("retention_monthly", flattenBackupProtectionPolicyFileShareRetentionMonthly(retention.MonthlySchedule)); err != nil {
				return fmt.Errorf("setting `retention_monthly`: %+v", err)
			}

			if err := d.Set("retention_yearly", flattenBackupProtectionPolicyFileShareRetentionYearly(retention.YearlySchedule)); err != nil {
				return fmt.Errorf("setting `retention_yearly`: %+v", err)
			}
		}

		if properties.VaultRetentionPolicy != nil {
			d.Set("backup_tier", "vault-standard")
			d.Set("snapshot_retention_in_days", int(properties.VaultRetentionPolicy.SnapshotRetentionInDays))

			if retention, ok := properties.VaultRetentionPolicy.VaultRetention.(protectionpolicies.LongTermRetentionPolicy); ok {
				if err := d.Set("retention_daily", flattenBackupProtectionPolicyFileShareRetentionDaily(retention.DailySchedule)); err != nil {
					return fmt.Errorf("setting `retention_daily`: %+v", err)
				}

				if err := d.Set("retention_weekly", flattenBackupProtectionPolicyFileShareRetentionWeekly(retention.WeeklySchedule)); err != nil {
					return fmt.Errorf("setting `retention_weekly`: %+v", err)
				}

				if err := d.Set("retention_monthly", flattenBackupProtectionPolicyFileShareRetentionMonthly(retention.MonthlySchedule)); err != nil {
					return fmt.Errorf("setting `retention_monthly`: %+v", err)
				}

				if err := d.Set("retention_yearly", flattenBackupProtectionPolicyFileShareRetentionYearly(retention.YearlySchedule)); err != nil {
					return fmt.Errorf("setting `retention_yearly`: %+v", err)
				}
			}
		} else {
			d.Set("backup_tier", "snapshot")
			d.Set("snapshot_retention_in_days", 0)
		}
	}

	return nil
}

func dataSourceBackupPolicyFileShareSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"backup": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"frequency": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"hourly": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"interval": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"start_time": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"window_duration": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
							},
						},
					},

					"time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"backup_tier": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"retention_daily": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"retention_monthly": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"days": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
						},
					},

					"include_last_days": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"weekdays": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"weeks": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"retention_weekly": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"weekdays": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"retention_yearly": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"days": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
						},
					},

					"include_last_days": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"months": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"weekdays": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"weeks": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"snapshot_retention_in_days": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"timezone": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
