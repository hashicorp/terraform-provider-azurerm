// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_file_share

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectionPolicyFileShareRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protectionpolicies.ParseBackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.Set("name", id.BackupPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	if model := resp.Model; model != nil {
		if properties, ok := model.Properties.(protectionpolicies.AzureFileShareProtectionPolicy); ok {
			d.Set("timezone", properties.TimeZone)

			if schedule, ok := properties.SchedulePolicy.(protectionpolicies.SimpleSchedulePolicy); ok {
				if s, err := flattenBackupProtectionPolicyFileShareSchedule(schedule); err != nil {
					return fmt.Errorf("flattening `backup`: %+v", err)
				} else if err := d.Set("backup", s); err != nil {
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
	}

	return nil
}

func resourceBackupProtectionPolicyFileShareWaitForDeletion(ctx context.Context, client *protectionpolicies.ProtectionPoliciesClient, id protectionpolicies.BackupPolicyId, d *pluginsdk.ResourceData) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"Found"},
		Target:     []string{"NotFound"},
		Refresh:    resourceBackupProtectionPolicyFileShareRefreshFunc(ctx, client, id),
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for delete to finish for %s: %+v", id, err)
	}

	return nil
}

func resourceBackupProtectionPolicyFileShareRefreshFunc(ctx context.Context, client *protectionpolicies.ProtectionPoliciesClient, id protectionpolicies.BackupPolicyId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "NotFound", nil
			}

			return resp, "Error", fmt.Errorf("making Read request on %s: %+v", id, err)
		}

		return resp, "Found", nil
	}
}
