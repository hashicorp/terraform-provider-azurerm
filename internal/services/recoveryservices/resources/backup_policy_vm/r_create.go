// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectionPolicyVMCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := protectionpolicies.NewBackupPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), d.Get("name").(string))

	// getting this ready now because its shared between *everything*, time is... complicated for this resource
	timeOfDay := fmt.Sprintf("%s:00Z", d.Get("backup.0.time").(string))
	times := append(make([]string, 0), timeOfDay)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_backup_policy_vm", id.ID())
		}
	}

	// Less than 7 daily backups is no longer supported for create/update
	if d.Get("retention_daily.0.count").(int) > 1 && d.Get("retention_daily.0.count").(int) < 7 {
		return fmt.Errorf("the Azure API has recently changed behaviour so that provisioning a `count` for the `retention_daily` field can no longer be less than 7 days for new/updates to existing Backup Policies. Please ensure that `count` is greater than 7, currently %d", d.Get("retention_daily.0.count").(int))
	}

	schedulePolicy, err := expandBackupProtectionPolicyVMSchedule(d, times)
	if err != nil {
		return err
	}

	policyType := protectionpolicies.IAASVMPolicyType(d.Get("policy_type").(string))
	vmProtectionPolicyProperties := &protectionpolicies.AzureIaaSVMProtectionPolicy{
		TimeZone:                pointer.To(d.Get("timezone").(string)),
		SnapshotConsistencyType: pointer.ToEnum[protectionpolicies.IaasVMSnapshotConsistencyType](d.Get("consistency_type").(string)),
		PolicyType:              pointer.To(policyType),
		SchedulePolicy:          schedulePolicy,
		TieringPolicy:           expandBackupProtectionPolicyVMTieringPolicy(d.Get("tiering_policy").([]any)),
		InstantRPDetails:        expandBackupProtectionPolicyVMResourceGroup(d),
		RetentionPolicy: &protectionpolicies.LongTermRetentionPolicy{ // SimpleRetentionPolicy only has duration property ¯\_(ツ)_/¯
			DailySchedule:   expandBackupProtectionPolicyVMRetentionDaily(d, times),
			WeeklySchedule:  expandBackupProtectionPolicyVMRetentionWeekly(d, times),
			MonthlySchedule: expandBackupProtectionPolicyVMRetentionMonthly(d, times),
			YearlySchedule:  expandBackupProtectionPolicyVMRetentionYearly(d, times),
		},
	}

	if d.HasChange("instant_restore_retention_days") {
		days := d.Get("instant_restore_retention_days").(int)
		if protectionpolicies.IAASVMPolicyTypeVOne == policyType && days > 5 {
			return fmt.Errorf("`instant_restore_retention_days` must be less than or equal to `5` when `policy_type` is `V1`")
		}

		vmProtectionPolicyProperties.InstantRpRetentionRangeInDays = pointer.To(int64(days))
	}

	policy := protectionpolicies.ProtectionPolicyResource{
		Properties: vmProtectionPolicyProperties,
	}

	if _, err = client.CreateOrUpdate(ctx, id, policy, protectionpolicies.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = resourceBackupProtectionPolicyVMWaitForUpdate(ctx, client, id, d); err != nil {
		return err
	}

	d.SetId(id.ID())

	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceBackupProtectionPolicyVMRead(d, meta)
}
