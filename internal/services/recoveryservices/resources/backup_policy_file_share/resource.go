// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_file_share

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBackupProtectionPolicyFileShareCreateUpdate,
		Read:   resourceBackupProtectionPolicyFileShareRead,
		Update: resourceBackupProtectionPolicyFileShareCreateUpdate,
		Delete: resourceBackupProtectionPolicyFileShareDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := protectionpolicies.ParseBackupPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceBackupProtectionPolicyFileShareSchema(),

		// if daily, we need daily retention
		// if weekly daily cannot be set, and we need weekly
		CustomizeDiff: func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
			_, hasDaily := diff.GetOk("retention_daily")
			_, hasWeekly := diff.GetOk("retention_weekly")

			frequencyI, _ := diff.GetOk("backup.0.frequency")
			switch strings.ToLower(frequencyI.(string)) {
			case "daily":
				if !hasDaily {
					return errors.New("`retention_daily` must be set when backup.0.frequency is daily")
				}

				if _, ok := diff.GetOk("backup.0.weekdays"); ok {
					return errors.New("`backup.0.weekdays` should be not set when backup.0.frequency is daily")
				}
			case "weekly":
				if hasDaily {
					return errors.New("`retention_daily` must be not set when backup.0.frequency is weekly")
				}
				if !hasWeekly {
					return errors.New("`retention_weekly` must be set when backup.0.frequency is weekly")
				}
			case "hourly":
				if !hasDaily {
					return errors.New("`retention_daily` must be set when backup.0.frequency is hourly")
				}
			default:
				return errors.New("unrecognized value for backup.0.frequency")
			}

			// validate backup_tier and snapshot_retention_in_days
			if backupTier, ok := diff.GetOk("backup_tier"); ok && strings.ToLower(backupTier.(string)) == "vault-standard" {
				snapshotRetention := 0
				if v, ok := diff.GetOk("snapshot_retention_in_days"); ok {
					snapshotRetention = v.(int)
				}

				if retentionDaily, ok := diff.GetOk("retention_daily"); ok {
					retentionDailyList := retentionDaily.([]any)
					if len(retentionDailyList) > 0 {
						retentionDailyMap := retentionDailyList[0].(map[string]any)
						if count, ok := retentionDailyMap["count"]; ok {
							dailyCount := count.(int)
							if snapshotRetention >= dailyCount {
								return fmt.Errorf("`snapshot_retention_in_days` must be less than `retention_daily` count when `backup_tier` is set to `vault-standard`. Got snapshot_retention_in_days: %d, retention_daily count: %d", snapshotRetention, dailyCount)
							}
						}
					}
				}
			}

			return nil
		},
	}
}
