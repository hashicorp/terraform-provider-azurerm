// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectionPolicyVMUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protectionpolicies.ParseBackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	// getting this ready now because its shared between *everything*, time is... complicated for this resource
	timeOfDay := fmt.Sprintf("%s:00Z", d.Get("backup.0.time").(string))
	times := append(make([]string, 0), timeOfDay)

	// Less than 7 daily backups is no longer supported for create/update
	if d.HasChange("retention_daily.0.count") && (d.Get("retention_daily.0.count").(int) > 1 && d.Get("retention_daily.0.count").(int) < 7) {
		return fmt.Errorf("the Azure API has recently changed behaviour so that provisioning a `count` for the `retention_daily` field can no longer be less than 7 days for new/updates to existing Backup Policies. Please ensure that `count` is greater than 7, currently %d", d.Get("retention_daily.0.count").(int))
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return err
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	if _, ok := existing.Model.Properties.(protectionpolicies.AzureIaaSVMProtectionPolicy); !ok {
		return fmt.Errorf("retrieving %s: `properties` was not a AzureIaaSVMProtectionPolicy", id)
	}

	model := *existing.Model
	properties := existing.Model.Properties.(protectionpolicies.AzureIaaSVMProtectionPolicy)

	properties.InstantRpRetentionRangeInDays = nil
	if d.HasChange("instant_restore_retention_days") {
		days := d.Get("instant_restore_retention_days").(int)
		if properties.PolicyType != nil && protectionpolicies.IAASVMPolicyTypeVOne == *properties.PolicyType && days > 5 {
			return fmt.Errorf("`instant_restore_retention_days` must be less than or equal to `5` when `policy_type` is `V1`")
		}

		properties.InstantRpRetentionRangeInDays = pointer.To(int64(days))
	}

	if d.HasChange("consistency_type") {
		properties.SnapshotConsistencyType = pointer.ToEnum[protectionpolicies.IaasVMSnapshotConsistencyType](d.Get("consistency_type").(string))
	}

	if d.HasChange("tiering_policy") {
		properties.TieringPolicy = expandBackupProtectionPolicyVMTieringPolicy(d.Get("tiering_policy").([]any))
	}

	if d.HasChange("timezone") {
		properties.TimeZone = pointer.To(d.Get("timezone").(string))
	}

	if d.HasChange("instant_restore_resource_group") {
		properties.InstantRPDetails = expandBackupProtectionPolicyVMResourceGroup(d)
	}

	// If anything changes inside any of the `retention_*` fields, update the `schedulePolicy` object as the API requires all timestamps match
	// lintignore:R019 // deliberate subset: only the fields feeding schedulePolicy; the remaining attributes are applied in separate branches
	if d.HasChanges("backup", "retention_daily", "retention_weekly", "retention_monthly", "retention_yearly") {
		schedulePolicy, err := expandBackupProtectionPolicyVMSchedule(d, times)
		if err != nil {
			return err
		}
		properties.SchedulePolicy = schedulePolicy
	}

	if d.HasChanges("retention_daily", "backup.0.time") {
		if properties.RetentionPolicy == nil {
			properties.RetentionPolicy = &protectionpolicies.LongTermRetentionPolicy{}
		}
		retentionPolicy, ok := properties.RetentionPolicy.(protectionpolicies.LongTermRetentionPolicy)
		if !ok {
			properties.RetentionPolicy = &protectionpolicies.LongTermRetentionPolicy{}
		}

		retentionPolicy.DailySchedule = expandBackupProtectionPolicyVMRetentionDaily(d, times)
		properties.RetentionPolicy = retentionPolicy
	}

	if d.HasChanges("retention_weekly", "backup.0.time") {
		if properties.RetentionPolicy == nil {
			properties.RetentionPolicy = &protectionpolicies.LongTermRetentionPolicy{}
		}
		retentionPolicy, ok := properties.RetentionPolicy.(protectionpolicies.LongTermRetentionPolicy)
		if !ok {
			properties.RetentionPolicy = &protectionpolicies.LongTermRetentionPolicy{}
		}

		retentionPolicy.WeeklySchedule = expandBackupProtectionPolicyVMRetentionWeekly(d, times)
		properties.RetentionPolicy = retentionPolicy
	}

	if d.HasChanges("retention_monthly", "backup.0.time") {
		if properties.RetentionPolicy == nil {
			properties.RetentionPolicy = &protectionpolicies.LongTermRetentionPolicy{}
		}
		retentionPolicy, ok := properties.RetentionPolicy.(protectionpolicies.LongTermRetentionPolicy)
		if !ok {
			properties.RetentionPolicy = &protectionpolicies.LongTermRetentionPolicy{}
		}

		retentionPolicy.MonthlySchedule = expandBackupProtectionPolicyVMRetentionMonthly(d, times)
		properties.RetentionPolicy = retentionPolicy
	}

	if d.HasChanges("retention_yearly", "backup.0.time") {
		if properties.RetentionPolicy == nil {
			properties.RetentionPolicy = &protectionpolicies.LongTermRetentionPolicy{}
		}
		retentionPolicy, ok := properties.RetentionPolicy.(protectionpolicies.LongTermRetentionPolicy)
		if !ok {
			properties.RetentionPolicy = &protectionpolicies.LongTermRetentionPolicy{}
		}

		retentionPolicy.YearlySchedule = expandBackupProtectionPolicyVMRetentionYearly(d, times)
		properties.RetentionPolicy = retentionPolicy
	}

	model.Properties = properties
	if _, err = client.CreateOrUpdate(ctx, *id, model, protectionpolicies.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = resourceBackupProtectionPolicyVMWaitForUpdate(ctx, client, *id, d); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceBackupProtectionPolicyVMRead(d, meta)
}

func resourceBackupProtectionPolicyVMWaitForUpdate(ctx context.Context, client *protectionpolicies.ProtectionPoliciesClient, id protectionpolicies.BackupPolicyId, d *pluginsdk.ResourceData) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"NotFound"},
		Target:     []string{"Found"},
		Refresh:    resourceBackupProtectionPolicyVMRefreshFunc(ctx, client, id),
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to provision: %+v", id, err)
	}

	return nil
}
