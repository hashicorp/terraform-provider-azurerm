// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_file_share

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectionPolicyFileShareCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := protectionpolicies.NewBackupPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), d.Get("name").(string))

	// getting this ready now because its shared between *everything*, time is... complicated for this resource
	// if it's using hourly backup schedule, it passes null in times
	var times []string
	if strings.ToLower(d.Get("backup.0.frequency").(string)) == "daily" {
		timeOfDay := d.Get("backup.0.time").(string)
		dateOfDay, err := time.Parse(time.RFC3339, fmt.Sprintf("2018-07-30T%s:00Z", timeOfDay))
		if err != nil {
			return fmt.Errorf("generating time from %q for %s: %+v", timeOfDay, id, err)
		}
		times = append(make([]string, 0), date.Time{Time: dateOfDay}.String())
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_backup_policy_file_share", id.ID())
			}
		}
	}

	AzureFileShareProtectionPolicyProperties := &protectionpolicies.AzureFileShareProtectionPolicy{
		TimeZone:       pointer.To(d.Get("timezone").(string)),
		WorkLoadType:   pointer.To(protectionpolicies.WorkloadTypeAzureFileShare),
		SchedulePolicy: expandBackupProtectionPolicyFileShareSchedule(d, times),
	}

	if d.Get("backup_tier").(string) == "vault-standard" {
		AzureFileShareProtectionPolicyProperties.VaultRetentionPolicy = &protectionpolicies.VaultRetentionPolicy{
			SnapshotRetentionInDays: int64(d.Get("snapshot_retention_in_days").(int)),
			VaultRetention: &protectionpolicies.LongTermRetentionPolicy{
				DailySchedule:   expandBackupProtectionPolicyFileShareRetentionDaily(d, times),
				WeeklySchedule:  expandBackupProtectionPolicyFileShareRetentionWeekly(d, times),
				MonthlySchedule: expandBackupProtectionPolicyFileShareRetentionMonthly(d, times),
				YearlySchedule:  expandBackupProtectionPolicyFileShareRetentionYearly(d, times),
			},
		}
	} else {
		AzureFileShareProtectionPolicyProperties.RetentionPolicy = &protectionpolicies.LongTermRetentionPolicy{
			DailySchedule:   expandBackupProtectionPolicyFileShareRetentionDaily(d, times),
			WeeklySchedule:  expandBackupProtectionPolicyFileShareRetentionWeekly(d, times),
			MonthlySchedule: expandBackupProtectionPolicyFileShareRetentionMonthly(d, times),
			YearlySchedule:  expandBackupProtectionPolicyFileShareRetentionYearly(d, times),
		}
	}

	policy := protectionpolicies.ProtectionPolicyResource{
		Properties: AzureFileShareProtectionPolicyProperties,
	}

	if _, err := client.CreateOrUpdate(ctx, id, policy, protectionpolicies.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := resourceBackupProtectionPolicyFileShareWaitForUpdate(ctx, client, id, d); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceBackupProtectionPolicyFileShareRead(d, meta)
}
