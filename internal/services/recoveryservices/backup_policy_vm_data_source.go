// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceBackupPolicyVm() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBackupPolicyVmRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: dataSourceBackupPolicyVmSchema(),
	}
}

func dataSourceBackupPolicyVmRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	subscriptionid := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := protectionpolicies.NewBackupPolicyID(subscriptionid, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), d.Get("name").(string))

	protectionPolicy, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(protectionPolicy.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	if protectionPolicy.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	d.SetId(id.ID())

	if properties, ok := protectionPolicy.Model.Properties.(protectionpolicies.AzureIaaSVMProtectionPolicy); ok {
		d.Set("timezone", properties.TimeZone)
		d.Set("instant_restore_retention_days", properties.InstantRpRetentionRangeInDays)
		d.Set("tiering_policy", flattenBackupProtectionPolicyVMTieringPolicy(properties.TieringPolicy))

		if schedule, ok := properties.SchedulePolicy.(protectionpolicies.SimpleSchedulePolicy); ok {
			if err := d.Set("backup", flattenBackupProtectionPolicyVMSchedule(schedule)); err != nil {
				return fmt.Errorf("setting `backup`: %+v", err)
			}
		}

		if schedule, ok := properties.SchedulePolicy.(protectionpolicies.SimpleSchedulePolicyV2); ok {
			if err := d.Set("backup", flattenBackupProtectionPolicyVMScheduleV2(schedule)); err != nil {
				return fmt.Errorf("setting `backup`: %+v", err)
			}
		}

		policyType := string(protectionpolicies.IAASVMPolicyTypeVOne)
		if pointer.From(properties.PolicyType) != "" {
			policyType = string(pointer.From(properties.PolicyType))
		}
		d.Set("policy_type", policyType)

		d.Set("consistency_type", pointer.FromEnum(properties.SnapshotConsistencyType))

		if retention, ok := properties.RetentionPolicy.(protectionpolicies.LongTermRetentionPolicy); ok {
			if s := retention.DailySchedule; s != nil {
				if err := d.Set("retention_daily", flattenBackupProtectionPolicyVMRetentionDaily(s)); err != nil {
					return fmt.Errorf("setting `retention_daily`: %+v", err)
				}
			}

			if s := retention.WeeklySchedule; s != nil {
				if err := d.Set("retention_weekly", flattenBackupProtectionPolicyVMRetentionWeekly(s)); err != nil {
					return fmt.Errorf("setting `retention_weekly`: %+v", err)
				}
			}

			if s := retention.MonthlySchedule; s != nil {
				if err := d.Set("retention_monthly", flattenBackupProtectionPolicyVMRetentionMonthly(s)); err != nil {
					return fmt.Errorf("setting `retention_monthly`: %+v", err)
				}
			}

			if s := retention.YearlySchedule; s != nil {
				if err := d.Set("retention_yearly", flattenBackupProtectionPolicyVMRetentionYearly(s)); err != nil {
					return fmt.Errorf("setting `retention_yearly`: %+v", err)
				}
			}
		}

		if instantRPDetail := properties.InstantRPDetails; instantRPDetail != nil {
			d.Set("instant_restore_resource_group", flattenBackupProtectionPolicyVMResourceGroup(*instantRPDetail))
		}
	}

	return nil
}

func dataSourceBackupPolicyVmSchema() map[string]*pluginsdk.Schema {
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

					"hour_duration": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"hour_interval": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"time": {
						Type:     pluginsdk.TypeString,
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

		"consistency_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"instant_restore_resource_group": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"prefix": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"suffix": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"instant_restore_retention_days": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"policy_type": {
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

		"tiering_policy": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"archived_restore_point": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"duration": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"duration_type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"mode": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"timezone": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
