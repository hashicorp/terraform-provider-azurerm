// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm_workload

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectionPoliciesClient

			id, err := protectionpolicies.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := BackupProtectionPolicyVMWorkloadModel{
				Name:              id.BackupPolicyName,
				ResourceGroupName: id.ResourceGroupName,
				RecoveryVaultName: id.VaultName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					vmWorkload, _ := props.(protectionpolicies.AzureVMWorkloadProtectionPolicy)
					state.WorkloadType = pointer.FromEnum(vmWorkload.WorkLoadType)
					state.Settings = flattenBackupProtectionPolicyVMWorkloadSettings(vmWorkload.Settings)
					state.ProtectionPolicies = flattenBackupProtectionPolicyVMWorkloadProtectionPolicies(vmWorkload.SubProtectionPolicy)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

type ProtectionPolicy struct {
	Backup           []Backup           `tfschema:"backup"`
	PolicyType       string             `tfschema:"policy_type"`
	RetentionDaily   []RetentionDaily   `tfschema:"retention_daily"`
	RetentionWeekly  []RetentionWeekly  `tfschema:"retention_weekly"`
	RetentionMonthly []RetentionMonthly `tfschema:"retention_monthly"`
	RetentionYearly  []RetentionYearly  `tfschema:"retention_yearly"`
	SimpleRetention  []SimpleRetention  `tfschema:"simple_retention"`
}

type Backup struct {
	Frequency          string   `tfschema:"frequency"`
	FrequencyInMinutes int64    `tfschema:"frequency_in_minutes"`
	Time               string   `tfschema:"time"`
	Weekdays           []string `tfschema:"weekdays"`
}

type RetentionDaily struct {
	Count int64 `tfschema:"count"`
}

type RetentionWeekly struct {
	Count    int64    `tfschema:"count"`
	Weekdays []string `tfschema:"weekdays"`
}

type RetentionMonthly struct {
	Count      int64    `tfschema:"count"`
	FormatType string   `tfschema:"format_type"`
	Monthdays  []int64  `tfschema:"monthdays"`
	Weeks      []string `tfschema:"weeks"`
	Weekdays   []string `tfschema:"weekdays"`
}

type RetentionYearly struct {
	Count      int64    `tfschema:"count"`
	FormatType string   `tfschema:"format_type"`
	Months     []string `tfschema:"months"`
	Monthdays  []int64  `tfschema:"monthdays"`
	Weeks      []string `tfschema:"weeks"`
	Weekdays   []string `tfschema:"weekdays"`
}

type SimpleRetention struct {
	Count int64 `tfschema:"count"`
}

type Settings struct {
	CompressionEnabled bool   `tfschema:"compression_enabled"`
	TimeZone           string `tfschema:"time_zone"`
}
