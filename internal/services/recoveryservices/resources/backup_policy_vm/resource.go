// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name backup_policy_vm -properties "name,resource_group_name,vault_name:recovery_vault_name" -test-name policyTypeDefault

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBackupProtectionPolicyVMCreate,
		Read:   resourceBackupProtectionPolicyVMRead,
		Update: resourceBackupProtectionPolicyVMUpdate,
		Delete: resourceBackupProtectionPolicyVMDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&protectionpolicies.BackupPolicyId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&protectionpolicies.BackupPolicyId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceBackupProtectionPolicyVMSchema(),

		// if daily, we need daily retention
		// if weekly daily cannot be set, and we need weekly
		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
			_, hasDaily := diff.GetOk("retention_daily")
			_, hasWeekly := diff.GetOk("retention_weekly")

			frequency, _ := diff.GetOk("backup.0.frequency")
			switch frequency.(string) {
			case string(protectionpolicies.ScheduleRunTypeHourly):
				if !hasDaily {
					return errors.New("`retention_daily` must be set when backup.0.frequency is hourly")
				}

				if _, ok := diff.GetOk("backup.0.weekdays"); ok {
					return errors.New("`backup.0.weekdays` should be not set when backup.0.frequency is hourly")
				}
			case string(protectionpolicies.ScheduleRunTypeDaily):
				if !hasDaily {
					return errors.New("`retention_daily` must be set when backup.0.frequency is daily")
				}

				if _, ok := diff.GetOk("backup.0.weekdays"); ok {
					return errors.New("`backup.0.weekdays` should be not set when backup.0.frequency is daily")
				}

				if _, ok := diff.GetOk("backup.0.hour_interval"); ok {
					return errors.New("`backup.0.hour_interval` should be not set when backup.0.frequency is daily")
				}

				if _, ok := diff.GetOk("backup.0.hour_duration"); ok {
					return errors.New("`backup.0.hour_duration` should be not set when backup.0.frequency is daily")
				}
			case string(protectionpolicies.ScheduleRunTypeWeekly):
				if hasDaily {
					return errors.New("`retention_daily` must be not set when backup.0.frequency is weekly")
				}
				if !hasWeekly {
					return errors.New("`retention_weekly` must be set when backup.0.frequency is weekly")
				}

				if _, ok := diff.GetOk("backup.0.hour_interval"); ok {
					return errors.New("`backup.0.hour_interval` should be not set when backup.0.frequency is weekly")
				}

				if _, ok := diff.GetOk("backup.0.hour_duration"); ok {
					return errors.New("`backup.0.hour_duration` should be not set when backup.0.frequency is weekly")
				}
			default:
				return errors.New("unrecognized value for backup.0.frequency")
			}

			consistencyType, hasConsistencyType := diff.GetOk("consistency_type")
			if hasConsistencyType && consistencyType.(string) == string(protectionpolicies.IaasVMSnapshotConsistencyTypeOnlyCrashConsistent) {
				if diff.Get("policy_type").(string) != string(protectionpolicies.IAASVMPolicyTypeVTwo) {
					return errors.New("`policy_type` must be `V2` when `consistency_type` is `OnlyCrashConsistent`")
				}
			}

			return nil
		}),
	}
}
